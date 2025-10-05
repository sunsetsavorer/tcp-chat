package main

import "github.com/sunsetsavorer/tcp-chat.git/app"

func main() {

	app := app.NewApp()

	err := app.Run()
	if err != nil {
		panic(err)
	}
}
