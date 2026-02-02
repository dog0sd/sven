package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dog0sd/sven/internal/audio"
	"github.com/dog0sd/sven/internal/config"
	"github.com/dog0sd/sven/internal/server"
	"github.com/dog0sd/sven/internal/tts"
)

func main() {
	backend := flag.String("backend", "", "audio backend: pulse or oto (overrides config)")
	flag.Usage = func() {
		fmt.Printf("Usage: %s [flags] [text]\n", os.Args[0])
		fmt.Printf("CLI mode: %s 'Hello, world!'\n", os.Args[0])
		fmt.Println("OR")
		fmt.Printf("Server mode: %s\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(0)
	}
	flag.Parse()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("error loading configuration: %v", err)
	}

	if *backend != "" {
		cfg.AudioBackend = *backend
	}

	player, err := audio.NewPlayer(cfg.AudioBackend)
	if err != nil {
		log.Fatalf("error creating audio player: %v", err)
	}

	if flag.NArg() == 0 {
		if err := server.StartServer(cfg.Port, cfg, player); err != nil {
			log.Fatal("HTTP server failed: ", err)
		}
	} else {
		text := strings.TrimSpace(strings.Join(flag.Args(), " "))
		if text != "" {
			mp3Data, err := tts.Synthesize(cfg.Elevenlabs, text, "")
			if err != nil {
				log.Fatal("elevenlabs synthesis error: ", err)
			}
			if err := player.Play(mp3Data); err != nil {
				log.Fatal("playback error: ", err)
			}
		}
	}
}
