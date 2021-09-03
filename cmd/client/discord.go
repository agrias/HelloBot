package client

import (
	"HelloBot/pkg/discord"
	"HelloBot/pkg/nicehash"
	"HelloBot/pkg/youtube"
	"HelloBot/proto"
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const DND_CHANNEL = 328682308054024193
const YMIR_ID = "<@462698059701420052>"

type discordBotClient struct {
	Session *discordgo.Session
	Client proto.BotClient
	BotState *BotState
}

func NewDiscordBot() YmirClient {

	log.Println("Creating client...")

	client := GetClient()
	context := &discordBotContext{}
	context.InitializeEnv()

	discord, err := discordgo.New("Bot ")
	if (err != nil) {
		log.Fatalln(err)
		panic(err)
	}

	return &discordBotClient{discord, client, &BotState{Client: client, Context: context, Youtube: youtube.NewYoutubeService()}}
}

func GetClient() proto.BotClient {
	serverAddr := "localhost:9095"
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		panic(err)
	}

	return proto.NewBotClient(conn)
}

func (b *discordBotClient) Run() {

	b.Session.AddHandler(b.BotState.onMessage)
	b.Session.AddHandler(b.BotState.onVoiceStateUpdate)

	err := b.Session.Open()
	if err != nil {
		log.Error("Error opening discord session: %s", err)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	b.Session.Close()
}

type discordBotContext struct {

}

func (context *discordBotContext) InitializeEnv() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("discord")
}

type BotState struct {
	Client proto.BotClient
	Youtube *youtube.YoutubeService
	Context *discordBotContext
}

func (b *BotState) onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {

	log.Infof("Message received: %s %s %s\n", m.ChannelID, m.Content, m.Author)
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
		return
	}

	if strings.HasPrefix(m.Content, "!balance") {
		old_balance := .00100562
		curr_balance := .007361

		now_balance := nicehash.GetBalance()
		output := now_balance - curr_balance - old_balance

		historical := now_balance - old_balance

		log.Infof("Nicehash: %f\n", output)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("BTC Balance: %f ($%f)\nHistorical Balance: %f ($%f)", output, output*nicehash.GetPrice(), historical, historical*nicehash.GetPrice()))
		return
	}

	if strings.HasPrefix(m.Content, "!play") {

		log.Infof("Playing song...")

		//user_meta := discord.GetUser(s, m.Author.ID)
		channel_meta, err := discord.GetChannel(s, m.ChannelID)
		if err != nil {
			log.Infof("Could not find channel...", m.ChannelID)
		}

		// guild_meta := discord.GetGuild(s, channel_meta.GuildID)

		channel_id := "883469626850312253"

		/*
		for k, v := range s.VoiceConnections {
			log.Infof("Voice information: %s, %s", k, v)
		}
		*/

		strings.TrimPrefix(m.Content, "!play")

		log.Infof("Joining channel: %s", channel_id)
		voicechannel, err := s.ChannelVoiceJoin(channel_meta.GuildID, channel_id, false, false)
		if err != nil {
			log.Errorf("Error joining channel... %s", err.Error())
			s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{Description: "This command must be done from a voice channel!", Author: &discordgo.MessageEmbedAuthor{Name: "@"+m.Author.Username}})
			return
		}

		//openchannels := s.VoiceConnections

		//youtube_url := "https://www.youtube.com%s"

		//results := youtube.GetVideosFromSearch(query)

		//s.ChannelMessageSend(m.ChannelID, "Playing... "+discord.FormatHyperlink(results[0].Url, results[0].Title)+"...")
		//s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{Description: "Playing... "+discord.FormatHyperlink(fmt.Sprintf(youtube_url, results[0].Url), results[0].Title)+"...", Author: &discordgo.MessageEmbedAuthor{Name: "@"+m.Author.Username}})
		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{Description: "Playing song", Author: &discordgo.MessageEmbedAuthor{Name: "@"+m.Author.Username}})
/*
		voicechannel := openchannels[guild_meta.ID]

		if voicechannel == nil {
			s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{Description: err.Error(), Author: &discordgo.MessageEmbedAuthor{Name: "@"+m.Author.Username}})
			return
		}
*/
		log.Infof("Starting song...")
		err = b.Youtube.PlayYoutubeVideo(voicechannel, "https://www.youtube.com/watch?v=rvrZJ5C_Nwg", "Blah")
		if err != nil {
			s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{Description: err.Error(), Author: &discordgo.MessageEmbedAuthor{Name: "@"+m.Author.Username}})
			return
		}
		log.Infof("Finishing song...")

		return
	}

	if strings.HasPrefix(m.Content, "!pause") {

		log.Infof("Attempting to stop song...")

		channel_meta, err := discord.GetChannel(s, m.ChannelID)
		if err != nil {
			log.Infof("Could not find channel...", m.ChannelID)
		}

		guild_meta := discord.GetGuild(s, channel_meta.GuildID)

		var channel_id string

		for _, people := range guild_meta.VoiceStates {
			if m.Author.ID == people.UserID {
				channel_id = people.ChannelID
			}
		}

		b.Youtube.OpenStreams[channel_id].SetPaused(true)
		return
	}

	if strings.HasPrefix(m.Content, "!clear") {

		channel_meta, err := discord.GetChannel(s, m.ChannelID)
		if err != nil {
			log.Infof("Could not find channel...", m.ChannelID)
		}

		b.Youtube.ClearQueue(channel_meta.GuildID)

		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{Description: "Queue cleared.", Author: &discordgo.MessageEmbedAuthor{Name: "@"+m.Author.Username}})
		return
	}

	if strings.HasPrefix(m.Content, "!unpause") {

		log.Infof("Attempting to restart song...")

		channel_meta, err := discord.GetChannel(s, m.ChannelID)
		if err != nil {
			log.Infof("Could not find channel...", m.ChannelID)
		}

		guild_meta := discord.GetGuild(s, channel_meta.GuildID)

		var channel_id string

		for _, people := range guild_meta.VoiceStates {
			if m.Author.ID == people.UserID {
				channel_id = people.ChannelID
			}
		}

		b.Youtube.OpenStreams[channel_id].SetPaused(false)
		return
	}

	if strings.HasPrefix(m.Content, "!next") {

		log.Infof("Attempting to skip song...")

		channel_meta, err := discord.GetChannel(s, m.ChannelID)
		if err != nil {
			log.Infof("Could not find channel...", m.ChannelID)
		}

		/*
		guild_meta := discord.GetGuild(s, channel_meta.GuildID)


		var channel_id string

		for _, people := range guild_meta.VoiceStates {
			if m.Author.ID == people.UserID {
				channel_id = people.ChannelID
			}
		}
		*/

		stop := b.Youtube.OpenStreams[channel_meta.GuildID]
		if stop != nil {
			//stop.End()
		}

		return
	}

	if strings.HasPrefix(m.Content, "!queue") {

		openchannels := s.VoiceConnections

		channel_meta, err := discord.GetChannel(s, m.ChannelID)
		if err != nil {
			log.Infof("Could not find channel...", m.ChannelID)
		}

		voicechannel := openchannels[channel_meta.GuildID]

		if voicechannel == nil {
			s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{Description: "Nothing has been queued.", Author: &discordgo.MessageEmbedAuthor{Name: "@"+m.Author.Username}})
			return
		}

		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{Description: b.Youtube.GetQueue(channel_meta.GuildID), Author: &discordgo.MessageEmbedAuthor{Name: "@"+m.Author.Username}})
		return
	}

	if strings.HasPrefix(m.Content, "!") {
		user_meta := discord.GetUser(s, m.Author.ID)

		log.Infof("User meta: %s, %s, %s, %s\n", user_meta.ID, user_meta.Username, user_meta.Email, user_meta.Discriminator)

		resp, err := b.Client.GetResponse(context.TODO(), &proto.BotRequest{Id: uuid.NewV1().String(), Text: m.Content, Name: user_meta.Username})

		if err != nil {
			log.Errorln("Problem getting response from Ymir Server")
		}
		//s.ChannelMessageSend(m.ChannelID, resp.Text)
		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{Description: resp.Text, Author: &discordgo.MessageEmbedAuthor{Name: "@"+m.Author.Username}})
	}

	return
}

func (b *BotState) onVoiceStateUpdate(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {

	log.Infof("VoiceStateUpdate received: %s, %s, %s\n", v.ChannelID, v.UserID, v.SessionID)

	channel_meta, err := discord.GetChannel(s, v.ChannelID)
	if err != nil || channel_meta == nil {
		return
	}

	log.Infof("Channel information: %s, %s, %s\n", channel_meta.ID, channel_meta.Name, channel_meta.ParentID)
}

