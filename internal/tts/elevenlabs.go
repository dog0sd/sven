package tts

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Mliviu79/elevenlabs-go"

	"github.com/dog0sd/sven/internal/config"
)

type VoiceInfo struct {
	Name        string
	Description string
	VoiceId     string
}

func timeoutDuration(elConfig config.ElevenLabsConfig) time.Duration {
	if elConfig.Timeout <= 0 {
		return 30 * time.Second
	}
	return time.Duration(elConfig.Timeout) * time.Second
}

func GetVoices(elConfig config.ElevenLabsConfig) ([]VoiceInfo, error) {
	client := elevenlabs.NewClient(context.Background(), elConfig.Token, timeoutDuration(elConfig))
	voices, err := client.GetVoices()
	if err != nil {
		return nil, fmt.Errorf("request error: %v", err)
	}
	result := make([]VoiceInfo, len(voices))
	for i, v := range voices {
		result[i] = VoiceInfo{
			Name:        v.Name,
			Description: v.Description,
			VoiceId:     v.VoiceId,
		}
	}
	return result, nil
}

// ResolveVoiceName finds a voice ID by case-insensitive name match.
func ResolveVoiceName(elConfig config.ElevenLabsConfig) (string, error) {
	voices, err := GetVoices(elConfig)
	if err != nil {
		return "", fmt.Errorf("fetching voices: %v", err)
	}
	for _, v := range voices {
		name := v.Name
		if idx := strings.Index(name, " - "); idx >= 0 {
			name = name[:idx]
		}
		if strings.EqualFold(name, elConfig.VoiceName) {
			return v.VoiceId, nil
		}
	}
	return "", fmt.Errorf("voice %q not found", elConfig.VoiceName)
}

type ModelInfo struct {
	Name        string
	Description string
	ModelId     string
}

func GetModels(elConfig config.ElevenLabsConfig) ([]ModelInfo, error) {
	client := elevenlabs.NewClient(context.Background(), elConfig.Token, timeoutDuration(elConfig))
	models, err := client.GetModels()
	if err != nil {
		return nil, fmt.Errorf("request error: %v", err)
	}
	result := make([]ModelInfo, len(models))
	for i, m := range models {
		result[i] = ModelInfo{
			Name:        m.Name,
			Description: m.Description,
			ModelId:     m.ModelId,
		}
	}
	return result, nil
}

func Synthesize(elConfig config.ElevenLabsConfig, text string, previousText string) ([]byte, error) {
	client := elevenlabs.NewClient(context.Background(), elConfig.Token, timeoutDuration(elConfig))
	var voiceSettings elevenlabs.VoiceSettings
	voiceSettings.SimilarityBoost = elConfig.Settings.SimilarityBoost
	voiceSettings.Stability = elConfig.Settings.Stability
	voiceSettings.Style = elConfig.Settings.Style
	voiceSettings.SpeakerBoost = elConfig.Settings.SpeakerBoost
	voiceSettings.Speed = elConfig.Settings.Speed

	ttsReq := elevenlabs.TextToSpeechRequest{
		Text:          text,
		ModelID:       elConfig.Model,
		VoiceSettings: &voiceSettings,
		PreviousText:  previousText,
	}
	audio, err := client.TextToSpeech(elConfig.VoiceId, ttsReq)
	if err != nil {
		return nil, fmt.Errorf("request error: %v", err)
	}
	return audio, nil
}
