package feature

import "YmirBot/proto"

type FeatureHandler interface {
	Filter(string) bool
	Handle(*proto.BotRequest) *proto.BotResponse
	Type() string
}