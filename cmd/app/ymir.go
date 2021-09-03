package app

import (
	"HelloBot/cmd/client"
	"HelloBot/cmd/server"
	"HelloBot/proto"
)

type Bot interface {
	Start()
}

type ymirBot struct {
	Server proto.BotServer
	Clients []client.YmirClient
}

func (b *ymirBot) Start() {
	b.Server = server.NewBotServer()
	b.Clients = append(b.Clients, client.NewDiscordBot())

	for _, client := range b.Clients {
		client.Run()
	}
}

func NewHelloBot() Bot {
	return &ymirBot{}
}