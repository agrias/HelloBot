package music

import (
	"os"
	"github.com/hajimehoshi/go-mp3"
	"fmt"
	"io"
	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	"log"
	"time"
)

func GetAudioBytes() (*mp3.Decoder, error) {
	f, err := os.Open("G:\\dev\\GOPATH\\src\\HelloBot\\pkg\\music\\hero.mp3")
	if err != nil {
		return nil, err
	}

	d, err := mp3.NewDecoder(f)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Length: %d[bytes]\n", d.Length())

	return d, nil

	//defer d.Close()
}

func PlayAudioFile(v *discordgo.VoiceConnection, filename string) {
	// Send "speaking" packet over the voice websocket
	err := v.Speaking(true)
	if err != nil {
		log.Fatal("Failed setting speaking", err)
	}

	// Send not "speaking" packet over the websocket when we finish
	defer v.Speaking(false)

	opts := dca.StdEncodeOptions
	opts.RawOutput = true
	opts.Bitrate = 120

	encodeSession, err := dca.EncodeFile(filename, opts)
	if err != nil {
		log.Fatal("Failed creating an encoding session: ", err)
	}

	done := make(chan error)
	stream := dca.NewStream(encodeSession, v, done)

	ticker := time.NewTicker(time.Second)

	for {
		select {
		case err := <-done:
			if err != nil && err != io.EOF {
				log.Fatal("An error occured", err)
			}

			// Clean up incase something happened and ffmpeg is still running
			encodeSession.Truncate()
			return
		case <-ticker.C:
			stats := encodeSession.Stats()
			playbackPosition := stream.PlaybackPosition()

			fmt.Printf("Playback: %10s, Transcode Stats: Time: %5s, Size: %5dkB, Bitrate: %6.2fkB, Speed: %5.1fx\r", playbackPosition, stats.Duration.String(), stats.Size, stats.Bitrate, stats.Speed)
		}
	}
}