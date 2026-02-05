package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/jinzhu/configor"
	"github.com/dog0sd/sven/internal/audio"
)

func findConfigFiles() []string {
	candidates := []string{"sven.yml", "sven.yaml"}
	if home, err := os.UserHomeDir(); err == nil {
		candidates = append(candidates,
			home+"/.config/sven.yml",
			home+"/.config/sven.yaml",
		)
	}
	candidates = append(candidates, "/etc/sven.yml", "/etc/sven.yaml")
	var paths []string
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			paths = append(paths, p)
		}
	}
	return paths
}

// LoadTokenConfig loads only the API token from config file and environment.
// Use this when full config validation is not needed (e.g. listing voices/models).
func LoadTokenConfig() (ElevenLabsConfig, error) {
	var config Config
	if paths := findConfigFiles(); len(paths) > 0 {
		configor.New(&configor.Config{Silent: true}).Load(&config, paths...)
	}
	if envToken := os.Getenv("ELEVENLABS_API_KEY"); envToken != "" {
		config.Elevenlabs.Token = envToken
	}
	if config.Elevenlabs.Token == "" {
		return config.Elevenlabs, fmt.Errorf("elevenlabs token is required (set in config or ELEVENLABS_API_KEY env)")
	}
	if config.Elevenlabs.Timeout <= 0 {
		config.Elevenlabs.Timeout = 30
	}
	return config.Elevenlabs, nil
}

func LoadConfig() (Config, error) {
	var config Config

	err := configor.Load(&config, findConfigFiles()...)
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
	config.Elevenlabs.Timeout = 30
	config.Elevenlabs.Settings.Stability = 1.0
	config.Elevenlabs.Settings.Speed = 1.0
	config.AudioBackend = "pulse"

	if err := validateElevenLabsSettings(config); err != nil {
		return config, err
	}
	if err := validateAudioBackend(&config); err != nil {
		return config, err
	}
	return config, nil
}

func validateElevenLabsSettings(config Config) error {
	if config.Elevenlabs.VoiceId == "" && config.Elevenlabs.VoiceName == "" {
		return fmt.Errorf("either voiceid or voicename must be set")
	}
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
	slog.Info("configuration loaded",
		"voice_id", config.Elevenlabs.VoiceId,
		"voice_name", config.Elevenlabs.VoiceName,
		"model", config.Elevenlabs.Model,
		"similarity_boost", config.Elevenlabs.Settings.SimilarityBoost,
		"stability", config.Elevenlabs.Settings.Stability,
		"style", config.Elevenlabs.Settings.Style,
		"speaker_boost", config.Elevenlabs.Settings.SpeakerBoost,
		"speed", config.Elevenlabs.Settings.Speed,
		"audio_backend", config.AudioBackend,
		"listen", config.Listen,
		"timeout", config.Elevenlabs.Timeout,
	)
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
