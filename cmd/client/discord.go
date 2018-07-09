package client

import (
	"github.com/bwmarrin/discordgo"
	"os"
	"syscall"
	"os/signal"
	"YmirBot/proto"
	"google.golang.org/grpc"
	"context"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

type discordBotClient struct {
	Session *discordgo.Session
	Client proto.BotClient
	BotState *BotState
}

func NewDiscordBot() YmirClient {

	log.Println("Creating client...")

	discord, err := discordgo.New("Bot "+"")
	if (err != nil) {
		log.Fatalln(err)
		panic(err)
	}

	client := GetClient()

	return &discordBotClient{discord, client, &BotState{client}}
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
	err := b.Session.Open()
	if err != nil {
		log.Error("Error opening discord session: %s", err)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	b.Session.Close()
}

type BotState struct {
	Client proto.BotClient
}

func (b *BotState) onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {

	log.Info("Message received: %s %s %s\n", m.ChannelID, m.Content, m.Author)
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
		return
	}

	resp, err := b.Client.GetResponse(context.TODO(), &proto.BotRequest{Id: uuid.NewV1().String(),Text: m.Content})
	if err != nil {
		log.Errorln("Problem getting response from Ymir Server")
	}

	s.ChannelMessageSend(m.ChannelID, resp.Text)
}