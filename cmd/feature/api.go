package feature

import "HelloBot/proto"

type FeatureHandler interface {
	Filter(string) bool
	Handle(*proto.BotRequest) *proto.BotResponse
	Type() string
}