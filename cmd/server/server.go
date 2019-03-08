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
	"YmirBot/pkg/dnd"
	"YmirBot/cmd/feature"
)

/*
type BotServer interface {
	GetResponse(context.Context, *BotRequest) (*BotResponse, error)
}
*/


type botServer struct {
	cache db.Database
	handlers []feature.FeatureHandler
}

func (s *botServer) GetResponse(context context.Context, req *proto.BotRequest) (*proto.BotResponse, error) {

	log.Info("GRPC request received...")

	response := &proto.BotResponse{Id: req.Id, Text: "Hello "+req.Name}

	for _, handler := range s.handlers {

		if handler.Filter(req.Text) {
			return handler.Handle(req), nil
		}
	}

	if strings.HasPrefix(req.Text, "!roll") {
		num_dice, sides, modifier := dnd.ParseDiceString(req.Text)
		rollResults := dnd.RollDiceModifierWithHistory(num_dice, sides, modifier, s.cache, req.Name)
		
		response.Text = dnd.FormatRollResults(rollResults, req.Name, req.Text)
	} else if strings.HasPrefix(req.Text, "!group") {
		response.Text = dnd.ProcessGroupCommand(req, s.cache)
	} else if strings.HasPrefix(req.Text, "!stats") {
		response.Text = dnd.GetDiceHistoryStats(s.cache, req.Name)
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
	handlers := []feature.FeatureHandler{}

	handlers = append(handlers, &feature.TimeFeatureHandler{})
	handlers = append(handlers, &feature.HelpFeatureHandler{handlers})

	server := &botServer{cache: cache, handlers: handlers}

	go Start(server)

	return server
}
