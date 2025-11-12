package ui

import (
	"fmt"
	"ichat/internal/service"

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

	loginContent := container.NewVBox()
	loginBtn := widget.NewButton("Enter", func() {
		lgn := lgnLog.Text
		if lgn == "" {
			pswLog.SetText("")
			dialog.ShowInformation("Error", "Enter password", w)
		}

		pswd := pswLog.Text
		if pswd == "" {
			dialog.ShowInformation("Error", "Enter password", w)
		}
		pswLog.SetText("")

		fmt.Println(lgn)
		fmt.Println(pswd)

		a.showMainMenu(w)
	})

	loginContent.Add(lgnLog)
	loginContent.Add(pswLog)
	loginContent.Add(loginBtn)

	w.SetContent(loginContent)
}

func (a *UI) showRegisterScreen(w fyne.Window) {

}

func (a *UI) showChatsListScreen(w fyne.Window) {
	chats := []string{"Chat 1", "Chat 2", "Chat 3"} // Example chat list
	chatList := container.NewVBox()

	for _, chat := range chats {
		chatItem := container.NewHBox(
			widget.NewButton(chat, func() {
				// Logic to open the chat screen can be added here
				a.showChatScreen(w)
			}),
		)
		chatList.Add(chatItem)
	}

	scrollContainer := container.NewScroll(chatList)
	w.SetContent(scrollContainer)
}

func (a *UI) showChatScreen(w fyne.Window) {

}
