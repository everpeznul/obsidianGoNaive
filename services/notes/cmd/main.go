package main

import (
	"obsidianGoNaive/services/notes/cmd/app"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	app := app.NewApp()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := app.GRPCServer.Serve(app.Listener); err != nil {
			app.Log.Error("server error", "error", err)
		}
	}()
	<-sigChan

	app.GRPCServer.GracefulStop()
	app.Db.Close()

}
