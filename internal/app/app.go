package app

import (
	"context"
	"ichat/internal/service"
	chatsrv "ichat/internal/service/chat"
	"ichat/internal/ui"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	ui      *ui.UI
	chatSrv service.ChatService
}

func New() *App {
	a := &App{}
	return a
}

func (a *App) Run(ctx context.Context) error {
	// Setup signal handling for graceful shutdown
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
	}()

	a.getUI().Start(ctx)
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
