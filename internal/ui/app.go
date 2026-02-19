package ui

import (
	"context"
	"fmt"
	chatdomain "ichat/internal/domain/chat"
	"ichat/internal/service"
	"log"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Option func(*UI)

type UI struct {
	app      fyne.App
	srv      service.ChatService
	window   fyne.Window
	ctx      context.Context
	cancel   context.CancelFunc
	statusMu sync.RWMutex
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

func (a *UI) Start(ctx context.Context) {
	a.ctx, a.cancel = context.WithCancel(ctx)
	a.window = a.app.NewWindow("LIZZARD")
	a.window.CenterOnScreen()
	a.window.Resize(fyne.NewSize(800, 600))

	// Start connection in background and show loading screen
	a.showLoadingScreen(a.window)

	// Start connection
	err := a.srv.Connect(a.ctx)
	if err != nil {
		log.Printf("Connection init error: %v", err)
	}

	// Watch connection status
	go a.watchConnectionStatus()

	a.window.ShowAndRun()
}

func (a *UI) watchConnectionStatus() {
	lastStatus := service.StatusDisconnected
	for {
		select {
		case <-a.ctx.Done():
			return
		default:
			status := a.srv.GetStatus()
			if status != lastStatus {
				lastStatus = status
				a.updateConnectionStatus(status)
			}
		}
	}
}

func (a *UI) updateConnectionStatus(status service.ConnStatus) {
	fyne.Do(func() {
		switch status {
		case service.StatusConnected:
			a.showLoginScreen(a.window)
		case service.StatusDisconnected:
			// Show error dialog if we were connected and now disconnected
			if a.ctx.Err() == nil {
				dialog.ShowError(fmt.Errorf("connection lost, attempting to reconnect..."), a.window)
			}
		}
	})
}

func (a *UI) Close(w fyne.Window) {
	if a.cancel != nil {
		a.cancel()
	}
	a.srv.Close()
	a.app.Quit()
}

func (a *UI) showEnterScreen(w fyne.Window) {
	a.showLoginScreen(w)
}

func (a *UI) showLoadingScreen(w fyne.Window) {
	progress := widget.NewProgressBarInfinite()
	progress.Start()

	statusLabel := widget.NewLabel("Connecting to server...")
	statusLabel.Alignment = fyne.TextAlignCenter

	icon := widget.NewIcon(theme.ContentClearIcon())

	content := container.NewVBox(
		container.NewCenter(icon),
		container.NewCenter(statusLabel),
		container.NewCenter(progress),
	)

	w.SetContent(content)
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
	lgnLog.SetPlaceHolder("Username")
	pswLog := widget.NewEntry()
	pswLog.SetPlaceHolder("Password")
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

	chatTypes := []string{"Private", "Group", "Public"}
	selectedType := "Private"

	selectWidget := widget.NewSelect(chatTypes, func(value string) {
		selectedType = value
	})

	btn := widget.NewButton("Create", func() {
		name := entry.Text
		if name == "" {
			dialog.ShowInformation("Error", "Enter chat name", w)
			return
		}

		var chatType int
		switch selectedType {
		case "Private":
			chatType = 0
		case "Group":
			chatType = 2
		case "Public":
			chatType = 2
		}

		a.srv.CreateChat(chatType, name)
		a.showChatsListScreen(w)
	})

	content := container.NewVBox(
		entry,
		widget.NewLabel("Chat type:"),
		selectWidget,
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
		dialog.ShowInformation("Error", "Error joining chat", w)
		a.showChatsListScreen(w)
		return
	}

	// Top info
	msgList := container.NewVBox()
	// Example initial message
	msgList.Add(widget.NewLabel("System: Welcome to the chat"))

	msgCtx, msgCancel := context.WithCancel(a.ctx)
	defer msgCancel()

	msgCh, err := a.srv.ReceiveMessages(chat.ID)
	if err != nil {
		fmt.Printf("Error receiving messages: %v\n", err)
	}

	go func() {
		for {
			select {
			case <-msgCtx.Done():
				return
			case msg, ok := <-msgCh:
				if !ok {
					return
				}
				fyne.Do(func() {
					msgList.Add(widget.NewLabel(fmt.Sprintf("%s: %s", msg.SenderID, msg.Content)))
				})
			}
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
			dialog.ShowInformation("Error", "Error leaving chat", w)
		}
		msgCancel()
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
				err := a.srv.SendMessage(chatdomain.Message{
					SenderID: "current_user_id", // replace with actual user ID
					Content:  t,
					ChatID:   chat.ID,
					Action:   string(chatdomain.ActionSendText),
				})
				if err != nil {
					fyne.Do(func() {
						dialog.ShowInformation("Error", fmt.Sprintf("Failed to send: %v", err), w)
					})
				}
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
