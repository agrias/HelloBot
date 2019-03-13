package youtube

import (
	"github.com/jonas747/dca"
	"github.com/rylio/ytdl"
	"github.com/bwmarrin/discordgo"
	"io"
	"log"
	"errors"
	"fmt"
)

type YoutubeService struct {
	RunningConnections map[string]bool
	OpenStreams map[string]*dca.StreamingSession
	Queue []string
}

const YOUTUBE_URL = "https://www.youtube.com%s"

func (svc *YoutubeService) PlayYoutubeVideo(connection *discordgo.VoiceConnection, url string) error {

	if svc.RunningConnections[connection.ChannelID] == true {
		svc.Queue = append(svc.Queue, url)
		return errors.New("Alreadying playing a video, queuing "+url)
	} else {
		svc.RunningConnections[connection.ChannelID] = true
	}

	options := dca.StdEncodeOptions
	options.RawOutput = true
	options.Bitrate = 96
	options.Application = "lowdelay"

	log.Println("Getting video information...")
	videoInfo, err := ytdl.GetVideoInfo(url)
	if err != nil {
		// Handle the error
	}

	log.Printf("Parse download URL... %s\n", url)
	formats := videoInfo.Formats.Extremes(ytdl.FormatAudioBitrateKey, true)

	if len(formats) < 1 {
		return errors.New("Issue with format of video.")
	}

	downloadURL, err := videoInfo.GetDownloadURL(formats[0])
	if err != nil {
		// Handle the error
	}

	log.Println("Encoding file/url")
	encodingSession, err := dca.EncodeFile(downloadURL.String(), options)
	if err != nil {
		// Handle the error
	}
	defer encodingSession.Cleanup()

	done := make(chan error)

	log.Println("Starting youtube stream...")
	instance := dca.NewStream(encodingSession, connection, done)
	svc.OpenStreams[connection.ChannelID] = instance

	instance.Finished()

	err = <- done
	if err != nil && err != io.EOF {
		// Handle the error
	}

	svc.RunningConnections[connection.ChannelID] = false

	if len(svc.Queue) > 0 {
		pop := svc.Queue[0]
		svc.Queue = svc.Queue[1:]
		svc.PlayYoutubeVideo(connection, pop)
	}

	return nil
}

func (svc *YoutubeService) GetQueue() string {

	output := "\nQueued Songs:\n"

	for index, item := range svc.Queue {
		output = output + fmt.Sprintf("%d. %s\n", index, item)
	}

	return output
}

func NewYoutubeService() *YoutubeService {
	return &YoutubeService{make(map[string]bool, 5), make(map[string]*dca.StreamingSession, 5), make([]string, 0)}
}