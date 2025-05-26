package app

import "ichat/ui"

type App struct {
	ui *ui.UI
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
		a.ui = ui.New()
	}

	return a.ui
}
