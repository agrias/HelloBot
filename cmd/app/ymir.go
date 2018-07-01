package app

import (
	"YmirBot/proto"
	"YmirBot/cmd/server"
	"YmirBot/cmd/client"
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

func NewYmirBot() Bot {
	return &ymirBot{}
}