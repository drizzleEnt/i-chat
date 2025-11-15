package ui

import (
	"fmt"
	chatdomain "ichat/internal/domain/chat"
	"ichat/internal/service"
	"log"

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

	err := a.srv.Connect()
	if err != nil {
		log.Fatal(err)
	}

	myWindow.ShowAndRun()
}

func (a *UI) Close(w fyne.Window) {
	a.srv.Close()
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

	createChatButton := widget.NewButton("Create Chat", func() {
		a.showCreateChatScreen(w)
	})

	chats, err := a.srv.GetChats()
	if err != nil {
		fmt.Printf("Error Get Chats: %v\n", err.Error())
		dialog.ShowInformation("Error", "Unable to load chats", w)
		return
	}

	chatList := container.NewVBox()

	for _, chat := range chats {
		chatItem := container.NewHBox(
			widget.NewButton(chat.Name, func() {
				// Logic to open the chat screen can be added here
				a.showChatScreen(w, chat)
			}),
		)
		chatList.Add(chatItem)
	}

	scrollContainer := container.NewScroll(chatList)
	content := container.NewBorder(
		container.NewStack(
			createChatButton,
		), // top
		mainMenuBtn, // bottom
		nil,         // left
		nil,         // right
		container.NewStack(
			widget.NewCard("Chats", "", scrollContainer),
		),
	)
	w.SetContent(content)
}

func (a *UI) showCreateChatScreen(w fyne.Window) {
	entry := widget.NewEntry()
	entry.SetPlaceHolder("Enter chat name")

	btn := widget.NewButton("Create", func() {
		name := entry.Text
		if name == "" {
			dialog.ShowInformation("Error", "Enter chat name", w)
			return
		}
		a.srv.CreateChat(name)
		a.showChatsListScreen(w)
	})

	content := container.NewVBox(
		entry,
		btn,
	)

	w.SetContent(content)
}

func (a *UI) showChatScreen(w fyne.Window, chat *chatdomain.Chat) {
	joinMsg := chatdomain.Message{
		SenderID: "current_user_id",
		ChatID:   chat.ID,
		Action:   string(chatdomain.ActionJoinChat),
	}
	err := a.srv.SendMessage(joinMsg)
	if err != nil {
		dialog.ShowInformation("Error", "Enter connecting chat", w)
		a.showChatsListScreen(w)
		return
	}

	// Top info
	msgList := container.NewVBox()
	// Example initial message
	msgList.Add(widget.NewLabel("System: Welcome to the chat"))
	go func() {
		msg, err := a.srv.ReceiveMessages(chat.ID)
		if err != nil {
			fmt.Printf("Error receiving messages: %v\n", err)
			return
		}
		for msg := range msg {
			msgList.Add(widget.NewLabel(fmt.Sprintf("%s: %s", msg.SenderID, msg.Content)))
		}
	}()

	chatTitle := widget.NewLabel(chat.Name)
	chatTitle.Alignment = fyne.TextAlignCenter
	chatTitle.TextStyle = fyne.TextStyle{Bold: true}

	chatInfo := widget.NewLabel("Participants: Alice, Bob")
	chatInfo.Alignment = fyne.TextAlignCenter

	backBtn := widget.NewButton("← Back", func() {
		leaveMsg := chatdomain.Message{
			SenderID: "current_user_id",
			ChatID:   chat.ID,
			Action:   string(chatdomain.ActionLeaveChat),
		}
		err := a.srv.SendMessage(leaveMsg)
		if err != nil {
			dialog.ShowInformation("Error", "Enter connecting chat", w)
			a.showChatsListScreen(w)
			return
		}
		a.showChatsListScreen(w)
	})

	top := container.NewBorder(nil, nil, backBtn, nil, container.NewVBox(chatTitle, chatInfo))

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
				a.srv.SendMessage(chatdomain.Message{
					SenderID: "current_user_id", // replace with actual user ID
					Content:  t,
					ChatID:   chat.ID,
					Action:   string(chatdomain.ActionSendText),
				})
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
