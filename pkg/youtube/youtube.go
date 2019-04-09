package youtube

import (
	"github.com/jonas747/dca"
	"github.com/rylio/ytdl"
	"github.com/bwmarrin/discordgo"
	"io"
	"log"
	"errors"
	"fmt"
	"os"
	"time"
	"github.com/spf13/viper"
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
	videoInfo, err := ytdl.GetVideoInfo(url)
	if err != nil {
		// Handle the error
		return errors.New("Issue playing found video, please try again...")
	}

	if videoInfo.Duration > time.Hour {

		return errors.New("I currently do not support music longer than one hour.")
	}

	log.Printf("Parse download URL... %s\n", url)
	formats := videoInfo.Formats.Extremes(ytdl.FormatAudioBitrateKey, true)

	if len(formats) < 1 {
		return errors.New("Issue with format of video.")
	}

	//downloadURL, err := videoInfo.GetDownloadURL(formats[0])
	if err != nil {
		// Handle the error
	}

	filepath := viper.GetString("tmpdir")
	filename := filepath+connection.GuildID+".mp4"

	file, _ := os.Create(filename)
	defer file.Close()
	videoInfo.Download(videoInfo.Formats[0], file)

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