package main

import (
	"context"
	"ichat/internal/app"
)

func main() {
	ctx := context.Background()
	app := app.New()

	app.Run(ctx)
}
