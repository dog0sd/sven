package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/haguro/elevenlabs-go"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
	"github.com/joho/godotenv"
)

const (
	movel    = "eleven_monolinguak_v1"
	voice_id = "pNInz6obpgDQGcFmaJgB"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	elevenlabs_token := os.Getenv("ELEVEN_LABS_API_KEY")
	var text string
	if len(os.Args) == 1 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			text += scanner.Text() + " "
		}
	} else {
		text = strings.Join(os.Args[1:], " ")
	}
	client := elevenlabs.NewClient(context.Background(), elevenlabs_token, 30*time.Second)
	ttsReq := elevenlabs.TextToSpeechRequest{
		Text:    text,
		ModelID: "eleven_monolingual_v1",
	}
	audio, err := client.TextToSpeech(voice_id, ttsReq)
	if err != nil {
		log.Fatal(err)
	}

	// Write the audio file bytes to disk
	if err := os.WriteFile("adam.mp3", audio, 0644); err != nil {
		log.Fatal(err)
	}

	// test
	f, err := os.Open("adam.mp3")
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	defer f.Close()

	d, err := mp3.NewDecoder(f)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	c, err := oto.NewContext(d.SampleRate(), 2, 2, 8192)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	defer c.Close()

	p := c.NewPlayer()
	defer p.Close()
	if _, err := io.Copy(p, d); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
