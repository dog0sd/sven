package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/dog0sd/sven/internal/audio"
	"github.com/dog0sd/sven/internal/config"
	"github.com/dog0sd/sven/internal/server"
	"github.com/dog0sd/sven/internal/tts"
)

func initLogger(cfg config.Config) {
	var level slog.Level
	switch strings.ToLower(cfg.LogLevel) {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{Level: level}
	var handler slog.Handler
	if strings.ToLower(cfg.LogFormat) == "json" {
		handler = slog.NewJSONHandler(os.Stderr, opts)
	} else {
		handler = slog.NewTextHandler(os.Stderr, opts)
	}
	slog.SetDefault(slog.New(handler))
}

func main() {
	backend := flag.String("backend", "", "audio backend: pulse or oto (overrides config)")
	stability := flag.Float64("stability", -1, "voice stability 0.0-1.0 (overrides config)")
	similarity := flag.Float64("similarity", -1, "voice similarity boost 0.0-1.0 (overrides config)")
	style := flag.Float64("style", -1, "voice style 0.0-1.0 (overrides config)")
	speed := flag.Float64("speed", -1, "voice speed 0.7-1.2 (overrides config)")

	flag.Usage = func() {
		fmt.Printf("Usage: %s [flags] [text]\n\n", os.Args[0])
		fmt.Printf("CLI mode: `%s 'Hello, world!'`\n", os.Args[0])
		fmt.Printf("Server mode: `%s`\n\n", os.Args[0])
		fmt.Println("Commands:")
		fmt.Println("  voices: prints available voices names and ids")
		fmt.Println("  models: prints available models names and ids")

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

	if flag.NArg() > 0 && (flag.Arg(0) == "models" || flag.Arg(0) == "voices") {
		elCfg, err := config.LoadTokenConfig()
		if err != nil {
			slog.Error("config error", "error", err)
			os.Exit(1)
		}
		switch flag.Arg(0) {
		case "models":
			models, err := tts.GetModels(elCfg)
			if err != nil {
				slog.Error("error fetching models", "error", err)
				os.Exit(1)
			}
			for i, m := range models {
				fmt.Printf("Name: %s\n", m.Name)
				fmt.Printf("ID:   %s\n", m.ModelId)
				if m.Description != "" {
					fmt.Printf("Desc: %s\n", m.Description)
				}
				if i < len(models)-1 {
					fmt.Println()
				}
			}
		case "voices":
			voices, err := tts.GetVoices(elCfg)
			if err != nil {
				slog.Error("error fetching voices", "error", err)
				os.Exit(1)
			}
			for i, v := range voices {
				fmt.Printf("Name: %s\n", v.Name)
				fmt.Printf("ID:   %s\n", v.VoiceId)
				if v.Description != "" {
					fmt.Printf("Desc: %s\n", v.Description)
				}
				if i < len(voices)-1 {
					fmt.Println()
				}
			}
		}
		return
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("error loading configuration", "error", err)
		os.Exit(1)
	}

	initLogger(cfg)

	// Resolve voice name to ID if needed
	if cfg.Elevenlabs.VoiceId == "" && cfg.Elevenlabs.VoiceName != "" {
		voiceId, err := tts.ResolveVoiceName(cfg.Elevenlabs)
		if err != nil {
			slog.Error("voice name resolution failed", "error", err)
			os.Exit(1)
		}
		cfg.Elevenlabs.VoiceId = voiceId
		slog.Info("resolved voice name", "name", cfg.Elevenlabs.VoiceName, "voice_id", voiceId)
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
		slog.Error("error creating audio player", "error", err)
		os.Exit(1)
	}

	if flag.NArg() == 0 {
		config.LogStartupInfo(cfg)
		if err := server.StartServer(cfg.Listen, cfg, player); err != nil {
			slog.Error("HTTP server failed", "error", err)
			os.Exit(1)
		}
	} else {
		text := strings.TrimSpace(strings.Join(flag.Args(), " "))
		if text == "" {
			slog.Error("no text provided")
			os.Exit(1)
		}
		mp3Data, err := tts.Synthesize(cfg.Elevenlabs, text, "")
		if err != nil {
			slog.Error("elevenlabs synthesis error", "error", err)
			os.Exit(1)
		}
		if err := player.Play(mp3Data); err != nil {
			slog.Error("playback error", "error", err)
			os.Exit(1)
		}
	}
}
