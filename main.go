package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"encoding/json"

	"github.com/haguro/elevenlabs-go"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
	"github.com/jinzhu/configor"
)

type Config struct {
	Elevenlabs struct {
		Enabled bool `default:"false"`
		VoiceId string `required:"true"`
		Model string `default:"eleven_multilingual_v2"`
		Token string `required:"true"`
	}
	Port string `default:"8080"`
}

type TTSRequest struct {
	Text string
}

var (
	config Config
)

func main() {
	err := configor.Load(&config, "sven.yml")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	elevenlabsToken := config.Elevenlabs.Token
	voiceId := config.Elevenlabs.VoiceId
	model := config.Elevenlabs.Model

	var text string
	if len(os.Args) == 1 {
		fmt.Println("Running as server on " + config.Port + "...")
		http.HandleFunc("/tts", handleTTS)
		log.Fatal(http.ListenAndServe(":" + config.Port, nil))
	} else {
		text = strings.Join(os.Args[1:], " ")
	}
	err = elevenlabsTTS(text, model, voiceId, elevenlabsToken)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
}

func handleTTS(w http.ResponseWriter, r *http.Request) {
	req := TTSRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error decoding request")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	err = elevenlabsTTS(req.Text, config.Elevenlabs.Model, config.Elevenlabs.VoiceId, config.Elevenlabs.Token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "elevelabs error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("OK"))
}

func elevenlabsTTS(text string, model string, voiceId string, token string) error {

	client := elevenlabs.NewClient(context.Background(), token, 30*time.Second)
	ttsReq := elevenlabs.TextToSpeechRequest{
		Text:    text,
		ModelID: model,
	}
	audio, err := client.TextToSpeech(voiceId, ttsReq)
	if err != nil {
		return err
	}

	if err := os.WriteFile("audio.mp3", audio, 0644); err != nil {
		return err
	}

	f, err := os.Open("audio.mp3")
	if err != nil {
		return err
	}
	defer f.Close()

	d, err := mp3.NewDecoder(f)
	if err != nil {
		return err
	}

	c, err := oto.NewContext(d.SampleRate(), 2, 2, 8192)
	if err != nil {
		return err
	}
	defer c.Close()

	p := c.NewPlayer()
	defer p.Close()
	if _, err := io.Copy(p, d); err != nil {
		return err
	}
	return nil
}