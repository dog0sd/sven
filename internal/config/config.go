package config

import (
	"fmt"

	"github.com/jinzhu/configor"
)

func LoadConfig() (Config, error) {
	var config Config

	err := configor.Load(&config, "/etc/sven.yml", "sven.yml")
	if err != nil {
		return config, err
	}
	if err = validateElevenLabsSettings(config); err != nil {
		return config, err
	}
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
	fmt.Printf("[startup] voice id: %s\n", config.Elevenlabs.VoiceId)
	fmt.Printf("[startup] model id: %s\n", config.Elevenlabs.Model)
	fmt.Printf("[startup] default similarity boost: %.2f\n", config.Elevenlabs.Settings.SimilarityBoost)
	fmt.Printf("[startup] default stability: %.2f\n", config.Elevenlabs.Settings.Stability)
	fmt.Printf("[startup] default style: %.2f\n", config.Elevenlabs.Settings.Style)
	fmt.Printf("[startup] default speaker boost: %v\n", config.Elevenlabs.Settings.SpeakerBoost)
	fmt.Printf("[startup] default speed: %.2f\n", config.Elevenlabs.Settings.Speed)
	return nil
}
