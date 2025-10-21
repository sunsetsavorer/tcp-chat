package app

import (
	"github.com/sunsetsavorer/tcp-chat-server/config"
	"github.com/sunsetsavorer/tcp-chat-server/server"
)

type App struct {
	Config *config.Config
	Server *server.ChatServer
}

func NewApp() *App {

	return &App{}
}

func (a *App) Run() error {

	config := config.New()

	err := config.LoadConfig()
	if err != nil {
		return err
	}

	a.Config = config
	a.Server = server.New(a.Config)

	err = a.Server.Run()
	if err != nil {
		return err
	}

	return nil
}
