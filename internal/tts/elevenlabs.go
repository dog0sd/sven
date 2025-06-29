package tts

import (
	"bytes"
	"context"
	"io"
	"time"

	"github.com/haguro/elevenlabs-go"
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
	ttsReq := elevenlabs.TextToSpeechRequest{
		Text:    text,
		ModelID: elConfig.Model,
		VoiceSettings: &voiceSettings,
	}
	audio, err := client.TextToSpeech(elConfig.VoiceId, ttsReq)
	if err != nil {
		return err
	}

	decoder, err := mp3.NewDecoder(bytes.NewReader(audio))
	if err != nil {
		return err
	}

	ctx, err := oto.NewContext(decoder.SampleRate(), 2, 2, 8192)
	if err != nil {
		return err
	}
	defer ctx.Close()

	p := ctx.NewPlayer()
	defer p.Close()
	if _, err := io.Copy(p, decoder); err != nil {
		return err
	}
	return nil
}