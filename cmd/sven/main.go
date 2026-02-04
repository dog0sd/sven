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
	stability := flag.Float64("stability", -1, "voice stability 0.0-1.0 (overrides config)")
	similarity := flag.Float64("similarity", -1, "voice similarity boost 0.0-1.0 (overrides config)")
	style := flag.Float64("style", -1, "voice style 0.0-1.0 (overrides config)")
	speed := flag.Float64("speed", -1, "voice speed 0.7-1.2 (overrides config)")
	flag.Usage = func() {
		fmt.Printf("Usage: %s [flags] [text]\n\n", os.Args[0])
		fmt.Printf("CLI mode: %s 'Hello, world!'\n", os.Args[0])
		fmt.Printf("Server mode: %s (no arguments)\n\n", os.Args[0])
		fmt.Println("Flags:")
		fmt.Println("  -backend string")
		fmt.Println("        audio backend: pulse or oto")
		fmt.Println("  -stability float")
		fmt.Println("        voice stability (0.0-1.0)")
		fmt.Println("  -similarity float")
		fmt.Println("        voice similarity boost (0.0-1.0)")
		fmt.Println("  -style float")
		fmt.Println("        voice style (0.0-1.0)")
		fmt.Println("  -speed float")
		fmt.Println("        voice speed (0.7-1.2)")
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
	if *stability >= 0 {
		cfg.Elevenlabs.Settings.Stability = float32(*stability)
	}
	if *similarity >= 0 {
		cfg.Elevenlabs.Settings.SimilarityBoost = float32(*similarity)
	}
	if *style >= 0 {
		cfg.Elevenlabs.Settings.Style = float32(*style)
	}
	if *speed >= 0 {
		cfg.Elevenlabs.Settings.Speed = float32(*speed)
	}

	player, err := audio.NewPlayer(cfg.AudioBackend)
	if err != nil {
		log.Fatalf("error creating audio player: %v", err)
	}

	if flag.NArg() == 0 {
		config.LogStartupInfo(cfg)
		if err := server.StartServer(cfg.Port, cfg, player); err != nil {
			log.Fatal("HTTP server failed: ", err)
		}
	} else {
		text := strings.TrimSpace(strings.Join(flag.Args(), " "))
		if text == "" {
			log.Fatal("no text provided")
		}
		mp3Data, err := tts.Synthesize(cfg.Elevenlabs, text, "")
		if err != nil {
			log.Fatal("elevenlabs synthesis error: ", err)
		}
		if err := player.Play(mp3Data); err != nil {
			log.Fatal("playback error: ", err)
		}
	}
}
