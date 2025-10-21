package main

import "github.com/sunsetsavorer/tcp-chat-server/app"

func main() {

	app := app.NewApp()

	err := app.Run()
	if err != nil {
		panic(err)
	}
}
