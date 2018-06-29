package app

import "YmirBot/proto"

type Bot interface {
	Start()
}

type ymirBot struct {
	Server proto.BotServer
	Client proto.BotClient
}

