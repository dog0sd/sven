package config

import (
	"fmt"

	"github.com/jinzhu/configor"
)


func LoadConfig() (Config, error) {
	var config Config

	err := configor.Load(&config, "/etc/sven.yml", "sven.yml")
	fmt.Printf("Voice id: %s\n", config.Elevenlabs.VoiceId)
	fmt.Printf("Default Model id: %s\n", config.Elevenlabs.Model)
	fmt.Printf("Default SimilarityBoost: %.2f\n", config.Elevenlabs.Settings.SimilarityBoost)
	fmt.Printf("Default Stability: %.2f\n", config.Elevenlabs.Settings.Stability)
	fmt.Printf("Default Style: %.2f\n", config.Elevenlabs.Settings.Style)
	fmt.Printf("SpeakerBoost: %v\n", config.Elevenlabs.Settings.SpeakerBoost)
	if err != nil {
		return config, err
	}
	return config, nil
}
