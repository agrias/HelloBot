package server

import (
	"YmirBot/proto"
	"context"
	"net"
	"fmt"
	"google.golang.org/grpc"
	log "github.com/sirupsen/logrus"
	"YmirBot/cmd/db"
	"strings"
	"YmirBot/cmd/feature"
)

/*
type BotServer interface {
	GetResponse(context.Context, *BotRequest) (*BotResponse, error)
}
*/

type botServer struct {
	cache db.Database
}

func (s *botServer) GetResponse(context context.Context, req *proto.BotRequest) (*proto.BotResponse, error) {

	log.Info("GRPC request received...")

	response := &proto.BotResponse{Id: req.Id, Text: "Hello "+req.Name}

	if strings.HasPrefix(req.Text, "!roll") {
		num_dice, sides, modifier := feature.ParseDiceString(req.Text)

		result := feature.RollDiceModifierWithHistory(num_dice, sides, modifier, s.cache, req.Name)

		response.Text = "<@"+req.Name+"> "+result.String()
	} else if strings.HasPrefix(req.Text, "!group") {
		response.Text = feature.ProcessGroupCommand(req, s.cache)
	} else if strings.HasPrefix(req.Text, "!stats") {
		response.Text = feature.GetDiceHistoryStats(s.cache, req.Name)
	}

	return response, nil
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

	cache := db.NewDiskvCache("/app-data")
	server := &botServer{cache}

	go Start(server)

	return server
}
