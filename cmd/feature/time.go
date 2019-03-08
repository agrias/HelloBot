package feature

import (
	"strings"
	"YmirBot/proto"
	"time"
)

const TIME_CMD = "!time"

type TimeFeatureHandler struct {
}

func (h *TimeFeatureHandler) Filter(text string) bool {
	if strings.HasPrefix(text, TIME_CMD) {
		return true;
	}

	return false;
}

func (h *TimeFeatureHandler) Handle(request *proto.BotRequest) *proto.BotResponse {

	return &proto.BotResponse{Id: request.Id, Text: time.Now().String()}
}

func (h *TimeFeatureHandler) Type() string {
	return TIME_CMD
}