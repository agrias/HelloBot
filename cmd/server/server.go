package server

import (
	"YmirBot/proto"
	"context"
	"net"
	"fmt"
	"google.golang.org/grpc"
	log "github.com/sirupsen/logrus"
)

/*
type BotServer interface {
	GetResponse(context.Context, *BotRequest) (*BotResponse, error)
}
*/

type botServer struct {

}

func (s *botServer) GetResponse(context context.Context, req *proto.BotRequest) (*proto.BotResponse, error) {

	log.Info("GRPC request received...")

	return &proto.BotResponse{Id: req.Id, Text: "Hello World"}, nil
}

func Start(server *botServer) {

	log.Info("Starting GRPC server...")

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 9095))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	proto.RegisterBotServer(grpcServer, server)
	grpcServer.Serve(lis)
}

func NewBotServer() proto.BotServer{

	server := &botServer{}

	go Start(server)

	return server
}