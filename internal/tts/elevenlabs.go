package tts

import (
	"context"
	"fmt"
	"time"

	"github.com/Mliviu79/elevenlabs-go"

	"github.com/dog0sd/sven/internal/config"
)

func Synthesize(elConfig config.ElevenLabsConfig, text string, previousText string) ([]byte, error) {
	client := elevenlabs.NewClient(context.Background(), elConfig.Token, 30*time.Second)
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
