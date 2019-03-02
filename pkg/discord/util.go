package discord

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"errors"
)

func GetChannel(s *discordgo.Session, id string) *discordgo.Channel {
	channel, err := s.Channel(id)

	if err != nil {
		log.Errorln(err.Error())
	}

	return channel
}

func GetUser(s *discordgo.Session, id string) *discordgo.User {
	user, err := s.User(id)

	if err != nil {
		log.Errorln(err.Error())
	}

	return user
}

func GetGuild(s *discordgo.Session, id string) *discordgo.Guild {
	guild, err := s.Guild(id)

	if err != nil {
		log.Errorln(err.Error())
	}

	return guild
}

func GetUserVoiceChannelInGuild(s *discordgo.Session, user_id string, guild_id string) (*discordgo.Channel, error) {
	guild := GetGuild(s, guild_id)

	for _, v := range guild.VoiceStates {
		if v.UserID == user_id {
			return GetChannel(s, v.ChannelID), nil
		}
	}

	return nil, errors.New("User not found in VoiceChannel")
}