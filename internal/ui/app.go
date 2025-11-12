package ui

import (
	"fmt"
	"ichat/internal/service"
	chatdomain "ichat/internal/service/domain/chat"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type Option func(*UI)

type UI struct {
	app fyne.App
	srv service.ChatService
}

func WithChatService(srv service.ChatService) Option {
	return func(ui *UI) {
		ui.srv = srv
	}
}

func New(opts ...Option) *UI {
	ui := &UI{
		app: app.New(),
	}

	for _, opt := range opts {
		opt(ui)
	}

	return ui
}

func (a *UI) Start() {
	myWindow := a.app.NewWindow("Chat App")
	myWindow.CenterOnScreen()
	myWindow.Resize(fyne.NewSize(800, 600))
	a.showEnterScreen(myWindow)

	myWindow.ShowAndRun()
}

func (a *UI) Close(w fyne.Window) {
	a.app.Quit()
}

func (a *UI) showEnterScreen(w fyne.Window) {
	a.showLoginScreen(w)
}

func (a *UI) showMainMenu(w fyne.Window) {
	menu := container.NewVBox()

	title := widget.NewLabel("Main Menu")
	title.Alignment = fyne.TextAlignCenter

	btnChats := widget.NewButton("Chats", func() {
		a.showChatsListScreen(w)
	})

	btnRegister := widget.NewButton("Register", func() {
		a.showRegisterScreen(w)
	})

	btnLogout := widget.NewButton("Logout", func() {
		a.showLoginScreen(w)
	})

	btnQuit := widget.NewButton("Quit", func() {
		a.app.Quit()
	})

	menu.Add(title)
	menu.Add(widget.NewSeparator())
	menu.Add(btnChats)
	menu.Add(btnRegister)
	menu.Add(btnLogout)
	menu.Add(widget.NewSeparator())
	menu.Add(btnQuit)

	w.SetContent(menu)
}

func (a *UI) showLoginScreen(w fyne.Window) {
	lgnLog := widget.NewEntry()
	pswLog := widget.NewEntry()
	pswLog.Password = true
	registerBtn := widget.NewButton("Register", func() {
		a.showRegisterScreen(w)
	})

	loginContent := container.NewVBox()
	loginBtn := widget.NewButton("Enter", func() {
		lgn := lgnLog.Text
		if lgn == "" {
			pswLog.SetText("")
			dialog.ShowInformation("Error", "Enter login", w)
			return
		}

		pswd := pswLog.Text
		if pswd == "" {
			dialog.ShowInformation("Error", "Enter password", w)
			return
		}
		pswLog.SetText("")

		// Add login logic here

		fmt.Println(lgn)
		fmt.Println(pswd)

		a.showChatsListScreen(w)
	})

	loginContent.Add(lgnLog)
	loginContent.Add(pswLog)
	loginContent.Add(loginBtn)
	loginContent.Add(registerBtn)

	w.SetContent(loginContent)
}

func (a *UI) showRegisterScreen(w fyne.Window) {

}

func (a *UI) showChatsListScreen(w fyne.Window) {
	mainMenuBtn := widget.NewButton("Main Menu", func() {
		a.showMainMenu(w)
	})
	chats := []chatdomain.Chat{
		{ID: "chat1", Name: "Chat 1"},
		{ID: "chat2", Name: "Chat 2"},
		{ID: "chat3", Name: "Chat 3"},
	} // Example chat list
	chatList := container.NewVBox()

	for _, chat := range chats {
		chatItem := container.NewHBox(
			widget.NewButton(chat.Name, func() {
				// Logic to open the chat screen can be added here
				a.showChatScreen(w, chat.ID)
			}),
		)
		chatList.Add(chatItem)
	}

	scrollContainer := container.NewScroll(chatList)
	content := container.NewBorder(
		nil,         // top
		mainMenuBtn, // bottom
		nil,         // left
		nil,         // right
		container.NewStack(
			widget.NewCard("Chats", "", scrollContainer),
		),
	)
	w.SetContent(content)
}

func (a *UI) showChatScreen(w fyne.Window, chatID string) {
	// Top info
	chat := chatdomain.Chat{ // Example chat data
		ID:      "chat1",
		Name:    "Chat 1",
		Members: []chatdomain.Member{{Name: "Alice"}, {Name: "Bob"}},
	}

	chatTitle := widget.NewLabel(chat.Name)
	chatTitle.Alignment = fyne.TextAlignCenter
	chatTitle.TextStyle = fyne.TextStyle{Bold: true}

	chatInfo := widget.NewLabel("Participants: Alice, Bob")
	chatInfo.Alignment = fyne.TextAlignCenter

	backBtn := widget.NewButton("← Back", func() {
		a.showChatsListScreen(w)
	})

	top := container.NewBorder(nil, nil, backBtn, nil, container.NewVBox(chatTitle, chatInfo))

	// Messages list
	msgList := container.NewVBox()
	// Example initial message
	msgList.Add(widget.NewLabel("System: Welcome to the chat"))

	scroll := container.NewVScroll(msgList)
	scroll.SetMinSize(fyne.NewSize(400, 300))

	// Input + send
	input := widget.NewEntry()
	input.SetPlaceHolder("Type a message...")
	input.MultiLine = true
	input.Wrapping = fyne.TextWrapWord
	input.SetMinRowsVisible(3)

	sendFunc := func() {
		text := input.Text
		if text == "" {
			return
		}
		// Add user's message
		lbl := widget.NewLabel(text)
		lbl.Alignment = fyne.TextAlignTrailing
		lbl.Wrapping = fyne.TextWrapWord
		msgList.Add(lbl)
		input.SetText("")
		// Scroll to bottom so newest message is visible
		scroll.ScrollToBottom()

		// Optionally send via service if configured
		if a.srv != nil {
			// non-blocking send; adjust per real service API
			go func(t string) {
				_ = a.srv // placeholder: integrate with a.srv.SendMessage(...) if available
				_ = t
			}(text)
		}
	}

	sendBtn := widget.NewButton("Send", func() { sendFunc() })
	sendBtn.Importance = widget.HighImportance
	// allow pressing Enter to send
	input.OnSubmitted = func(s string) { sendFunc() }

	bottom := container.NewBorder(nil, nil, nil, sendBtn,
		container.NewStack(input),
	)

	// Assemble the screen
	content := container.NewBorder(top, bottom, nil, nil, scroll)
	w.SetContent(content)
}
