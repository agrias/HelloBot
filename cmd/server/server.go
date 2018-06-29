package server

import (
	"YmirBot/proto"
	"context"
)

/*
type BotServer interface {
	GetResponse(context.Context, *BotRequest) (*BotResponse, error)
}
*/

type botServer struct {

}

func (s *botServer) GetResponse(context context.Context, req *proto.BotRequest) (*proto.BotResponse, error) {

	return nil, nil
}

func NewBotServer() proto.BotServer{
	return &botServer{}
}