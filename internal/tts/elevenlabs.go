package tts

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/Mliviu79/elevenlabs-go"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"

	"github.com/dog0sd/sven/internal/config"
)


func ElevenlabsTTS(text string, elConfig config.ElevenLabsConfig) error {
	client := elevenlabs.NewClient(context.Background(), elConfig.Token, 30*time.Second)
	var voiceSettings elevenlabs.VoiceSettings
	voiceSettings.SimilarityBoost = elConfig.Settings.SimilarityBoost
	voiceSettings.Stability = elConfig.Settings.Stability
	voiceSettings.Style = elConfig.Settings.Style
	voiceSettings.SpeakerBoost = elConfig.Settings.SpeakerBoost
	voiceSettings.Speed = elConfig.Settings.Speed

	ttsReq := elevenlabs.TextToSpeechRequest{
		Text:    text,
		ModelID: elConfig.Model,
		VoiceSettings: &voiceSettings,
	}
	audio, err := client.TextToSpeech(elConfig.VoiceId, ttsReq)
	if err != nil {
		return fmt.Errorf("request error: %v", err)
	}
	reader := bytes.NewReader(audio)
	decoder, err := mp3.NewDecoder(reader)
		
	if err != nil {
		return fmt.Errorf("decode error: %v", err)
	}

	ctx, err := oto.NewContext(decoder.SampleRate(), 2, 2, 8192)
	if err != nil {
		return fmt.Errorf("oto context error: %v", err)
	}
	defer ctx.Close()

	p := ctx.NewPlayer()
	defer p.Close()
	if _, err := io.Copy(p, decoder); err != nil {
		return fmt.Errorf("playing audio error: %v", err)
	}
	return nil
}
