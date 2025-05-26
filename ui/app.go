package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type UI struct {
}

func New() *UI {
	return &UI{}
}

func (a *UI) Start() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Chat App")
	myWindow.CenterOnScreen()
	myWindow.Resize(fyne.NewSize(800, 600))

	a.showEnterScreen(myWindow)

	myWindow.ShowAndRun()
}

func (a *UI) showEnterScreen(w fyne.Window) {
	a.showLoginScreen(w)
}

func (a *UI) showMainMenu(w fyne.Window) {

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
	})

	loginContent.Add(lgnLog)
	loginContent.Add(pswLog)
	loginContent.Add(loginBtn)

	w.SetContent(loginContent)
}

func (a *UI) showRegisterScreen(w fyne.Window) {

}

func (a *UI) showChatsListScreen(w fyne.Window) {

}

func (a *UI) showChatScreen(w fyne.Window) {

}
