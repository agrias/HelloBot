package youtube

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	"github.com/kkdai/youtube/v2"
	"github.com/spf13/viper"
	"io"
	"log"
	"os"
)

type YoutubeService struct {
	RunningConnections map[string]bool
	OpenStreams map[string]*dca.StreamingSession
	Queues map[string][]*VideoQueueData
	Volume int
}

type VideoQueueData struct {
	Url string
	Title string
}

const YOUTUBE_URL = "https://www.youtube.com%s"

func (svc *YoutubeService) PlayYoutubeVideo(connection *discordgo.VoiceConnection, url string, title string) error {

	// init map of queues
	if svc.Queues[connection.GuildID] == nil {
		svc.Queues[connection.GuildID] = make([]*VideoQueueData, 0)
	}

	// init running connections map
	if svc.RunningConnections[connection.GuildID] == true {
		svc.Queues[connection.GuildID] = append(svc.Queues[connection.GuildID], &VideoQueueData{url, title})
		return errors.New("Alreadying playing a video, queuing "+url)
	} else {
		svc.RunningConnections[connection.GuildID] = true
	}

	options := dca.StdEncodeOptions
	options.RawOutput = true
	options.Bitrate = 96
	options.Application = "lowdelay"
	options.Volume = svc.Volume

	log.Println("Getting video information...")

	filepath := viper.GetString("tmpdir")
	filename := filepath+connection.GuildID+".mp4"

	file, _ := os.Create(filename)
	defer file.Close()
	//videoInfo.Download(videoInfo.Formats[0], file)

	client := youtube.Client{}

	video, err := client.GetVideo(url)
	if err != nil {
		panic(err)
	}

	stream, _, err := client.GetStream(video, &video.Formats[0])
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(file, stream)
	if err != nil {
		panic(err)
	}

	log.Println("Encoding file/url")
	encodingSession, err := dca.EncodeFile(filename, options)
	if err != nil {
		// Handle the error
	}
	defer encodingSession.Cleanup()

	done := make(chan error)

	log.Println("Starting youtube stream...")
	instance := dca.NewStream(encodingSession, connection, done)
	svc.OpenStreams[connection.GuildID] = instance

	instance.Finished()

	err = <- done
	if err != nil && err != io.EOF {
		// Handle the error
	}

	svc.RunningConnections[connection.GuildID] = false

	if len(svc.Queues[connection.GuildID]) > 0 {
		pop := svc.Queues[connection.GuildID][0]
		svc.Queues[connection.GuildID] = svc.Queues[connection.GuildID][1:]
		svc.PlayYoutubeVideo(connection, pop.Url, pop.Title)
	}

	return nil
}

func (svc *YoutubeService) GetQueue(guild string) string {

	output := "\nQueued Songs:\n"

	for index, item := range svc.Queues[guild] {
		output = output + fmt.Sprintf("%d. %s\n", index+1, item.Title)
	}

	return output
}

func (svc *YoutubeService) ClearQueue(guild string) {

	svc.Queues[guild] = make([]*VideoQueueData, 0)
}

func (svc *YoutubeService) SetVolume(volume int) {
	svc.Volume = volume
}

func NewYoutubeService() *YoutubeService {
	return &YoutubeService{make(map[string]bool, 5), make(map[string]*dca.StreamingSession, 5), make(map[string][]*VideoQueueData, 0), 100}
}