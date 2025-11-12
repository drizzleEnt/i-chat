package app

import (
	"ichat/internal/service"
	chatsrv "ichat/internal/service/chat"
	"ichat/internal/ui"
)

type App struct {
	ui *ui.UI
	chatSrv service.ChatService
}

func New() *App {
	a := &App{}
	return a
}

func (a *App) Run() error {
	a.getUI().Start()
	return nil
}

func (a *App) getUI() *ui.UI {
	if a.ui == nil {
		a.ui = ui.New(ui.WithChatService(a.ChatService()))
	}

	return a.ui
}

func (a *App) ChatService() service.ChatService {
	if a.chatSrv == nil {
		a.chatSrv = chatsrv.NewConnAdapter()
	}

	return a.chatSrv
}