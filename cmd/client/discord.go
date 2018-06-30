package client

import (
	"google.golang.org/grpc"
	"context"
	"YmirBot/proto"
	"github.com/bwmarrin/discordgo"
)

/*
type BotClient interface {
	GetResponse(ctx context.Context, in *BotRequest, opts ...grpc.CallOption) (*BotResponse, error)
}
*/
type DiscordBotClient struct {
	Session *discordgo.Session
}

func (b *DiscordBotClient) GetResponse(ctx context.Context, in *proto.BotRequest, opts ...grpc.CallOption) (*proto.BotResponse, error) {

	return nil, nil
}

func NewDiscordBot() proto.BotClient {

	discord, err := discordgo.New("Bot "+"")

	if (err != nil) {

	}

	return &DiscordBotClient{discord}
}