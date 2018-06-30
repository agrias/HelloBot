package app

import (
	"YmirBot/proto"
	"YmirBot/cmd/server"
)

type Bot interface {
	Start()
}

type ymirBot struct {
	Server proto.BotServer
}

func (b *ymirBot) Start() {
	b.Server = server.NewBotServer()
}

func NewYmirBot() Bot {
	return &ymirBot{}
}