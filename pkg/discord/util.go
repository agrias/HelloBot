package discord

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"errors"
)

// often returns a 404
func GetChannel(s *discordgo.Session, id string) (*discordgo.Channel, error) {

	for retry := 0; retry <= 3; retry++ {
		channel, err := s.Channel(id)
		if err != nil && retry == 3 {
			log.Errorf("Issue finding channel for %s: %s", id, err.Error())
			return nil, err
		} else {
			return channel, nil
		}
	}

	return nil, errors.New("This should never happen")
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
			channel, err := GetChannel(s, v.ChannelID)
			if err != nil {
				return nil, err
			}

			return channel, nil
		}
	}

	return nil, errors.New("User not found in VoiceChannel")
}