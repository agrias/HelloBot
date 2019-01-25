package client

import (
	"testing"
)

func TestNewDiscordBot(t *testing.T) {
	bot := NewDiscordBot()
	go bot.Run()
}
