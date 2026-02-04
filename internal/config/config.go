package config

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/configor"
	"github.com/dog0sd/sven/internal/audio"
)

func LoadConfig() (Config, error) {
	var config Config

	// Build config paths in priority order (first found wins)
	configPaths := []string{"sven.yml", "sven.yaml"}

	if home, err := os.UserHomeDir(); err == nil {
		configPaths = append(configPaths,
			home+"/.config/sven.yml",
			home+"/.config/sven.yaml",
		)
	}
	configPaths = append(configPaths, "/etc/sven.yml", "/etc/sven.yaml")

	err := configor.Load(&config, configPaths...)
	if err != nil {
		return config, err
	}

	// Environment variables override config file (priority: env > config)
	if envToken := os.Getenv("ELEVENLABS_API_KEY"); envToken != "" {
		config.Elevenlabs.Token = envToken
	}

	if err = validateElevenLabsSettings(config); err != nil {
		return config, err
	}
	if err = validateAudioBackend(&config); err != nil {
		return config, err
	}
	return config, nil
}

func LoadConfigFromEnv() (Config, error) {
	var config Config
	config.Elevenlabs.Token = os.Getenv("ELEVENLABS_API_KEY")
	config.Elevenlabs.VoiceId = os.Getenv("ELEVENLABS_VOICE_ID")
	config.Elevenlabs.Model = os.Getenv("ELEVENLABS_MODEL")
	if config.Elevenlabs.Token == "" || config.Elevenlabs.VoiceId == "" {
		return config, fmt.Errorf("ELEVENLABS_API_KEY or ELEVENLABS_VOICE_ID is empty")
	}
	if config.Elevenlabs.Model == "" {
		config.Elevenlabs.Model = "eleven_turbo_v2_5"
	}
	config.Elevenlabs.Settings.Stability = 1.0
	config.Elevenlabs.Settings.Speed = 1.0
	config.AudioBackend = "pulse"
	return config, nil
}

func validateElevenLabsSettings(config Config) error {
	if config.Elevenlabs.Settings.SimilarityBoost > 1.0 {
		return fmt.Errorf("similarityboost must be 0.0 <= x <= 1.0")
	}
	if config.Elevenlabs.Settings.Stability > 1.0 {
		return fmt.Errorf("stability must be 0.0 <= x <= 1.0")
	}
	if config.Elevenlabs.Settings.Style > 1.0 {
		return fmt.Errorf("style must be 0.0 <= x <= 1.0")
	}
	if config.Elevenlabs.Settings.Speed > 1.2 || config.Elevenlabs.Settings.Speed < 0.7 {
		return fmt.Errorf("speed must be 0.7 <= x <= 1.2")
	}
	return nil
}

// LogStartupInfo logs configuration details at startup.
func LogStartupInfo(config Config) {
	log.Printf("voice id: %s", config.Elevenlabs.VoiceId)
	log.Printf("model id: %s", config.Elevenlabs.Model)
	log.Printf("default similarity boost: %.2f", config.Elevenlabs.Settings.SimilarityBoost)
	log.Printf("default stability: %.2f", config.Elevenlabs.Settings.Stability)
	log.Printf("default style: %.2f", config.Elevenlabs.Settings.Style)
	log.Printf("default speaker boost: %v", config.Elevenlabs.Settings.SpeakerBoost)
	log.Printf("default speed: %.2f", config.Elevenlabs.Settings.Speed)
	log.Printf("audio backend: %s", config.AudioBackend)
}

func validateAudioBackend(config *Config) error {
	switch config.AudioBackend {
	case "pulse", "oto":
		// valid
	case "":
		config.AudioBackend = audio.DefaultBackend()
	default:
		return fmt.Errorf("invalid audiobackend %q (supported: pulse, oto)", config.AudioBackend)
	}
	return nil
}
