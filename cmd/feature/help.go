package feature

import (
	"YmirBot/proto"
	"strings"
)

/*
type featureHandler interface {
	Filter(string) bool
	Handle(*proto.BotRequest) *proto.BotResponse
}*/

const HELP_CMD = "!help"

type HelpFeatureHandler struct {
	Handlers []FeatureHandler
}

func (h *HelpFeatureHandler) Filter(text string) bool {
	if strings.HasPrefix(text, HELP_CMD) {
		return true;
	}

	return false;
}

func (h *HelpFeatureHandler) Handle(request *proto.BotRequest) *proto.BotResponse {

	response := "Available Commands:\n!help\n"

	for _, handler := range h.Handlers {
		response = response + handler.Type() + "\n"
	}

	return &proto.BotResponse{Id: request.Id, Text: response}
}

func (h *HelpFeatureHandler) Type() string {
	return HELP_CMD
}