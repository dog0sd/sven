package server

import (
	"testing"

	"github.com/dog0sd/sven/internal/config"
)

func floatPtr(f float32) *float32 {
	return &f
}

func TestMergeElevenLabsTTSSettings(t *testing.T) {
	baseConfig := config.Config{
		Elevenlabs: config.ElevenLabsConfig{
			Model:   "eleven_turbo_v2_5",
			VoiceId: "test-voice",
			Settings: config.ElevenlabsVoiceSettings{
				SimilarityBoost: 0.8,
				Stability:       0.7,
				Style:           0.5,
				Speed:           1.0,
			},
		},
	}

	tests := []struct {
		name     string
		request  TTSRequest
		expected config.ElevenlabsVoiceSettings
	}{
		{
			name: "no overrides - keep defaults",
			request: TTSRequest{
				Text: "test",
			},
			expected: config.ElevenlabsVoiceSettings{
				SimilarityBoost: 0.8,
				Stability:       0.7,
				Style:           0.5,
				Speed:           1.0,
			},
		},
		{
			name: "override all settings",
			request: TTSRequest{
				Text: "test",
				Elevenlabs: elevenLabsConfig{
					SimilarityBoost: floatPtr(0.3),
					Stability:       floatPtr(0.4),
					Style:           floatPtr(0.1),
					Speed:           floatPtr(1.2),
				},
			},
			expected: config.ElevenlabsVoiceSettings{
				SimilarityBoost: 0.3,
				Stability:       0.4,
				Style:           0.1,
				Speed:           1.2,
			},
		},
		{
			name: "override with zero values",
			request: TTSRequest{
				Text: "test",
				Elevenlabs: elevenLabsConfig{
					SimilarityBoost: floatPtr(0.0),
					Stability:       floatPtr(0.0),
					Style:           floatPtr(0.0),
					Speed:           floatPtr(0.7),
				},
			},
			expected: config.ElevenlabsVoiceSettings{
				SimilarityBoost: 0.0,
				Stability:       0.0,
				Style:           0.0,
				Speed:           0.7,
			},
		},
		{
			name: "partial override",
			request: TTSRequest{
				Text: "test",
				Elevenlabs: elevenLabsConfig{
					Speed: floatPtr(1.1),
				},
			},
			expected: config.ElevenlabsVoiceSettings{
				SimilarityBoost: 0.8,
				Stability:       0.7,
				Style:           0.5,
				Speed:           1.1,
			},
		},
		{
			name: "override model",
			request: TTSRequest{
				Text: "test",
				Elevenlabs: elevenLabsConfig{
					Model: "eleven_monolingual_v1",
				},
			},
			expected: config.ElevenlabsVoiceSettings{
				SimilarityBoost: 0.8,
				Stability:       0.7,
				Style:           0.5,
				Speed:           1.0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mergeElevenLabsTTSSettings(tt.request, &baseConfig)

			if result.Elevenlabs.Settings.SimilarityBoost != tt.expected.SimilarityBoost {
				t.Errorf("SimilarityBoost = %v, want %v", result.Elevenlabs.Settings.SimilarityBoost, tt.expected.SimilarityBoost)
			}
			if result.Elevenlabs.Settings.Stability != tt.expected.Stability {
				t.Errorf("Stability = %v, want %v", result.Elevenlabs.Settings.Stability, tt.expected.Stability)
			}
			if result.Elevenlabs.Settings.Style != tt.expected.Style {
				t.Errorf("Style = %v, want %v", result.Elevenlabs.Settings.Style, tt.expected.Style)
			}
			if result.Elevenlabs.Settings.Speed != tt.expected.Speed {
				t.Errorf("Speed = %v, want %v", result.Elevenlabs.Settings.Speed, tt.expected.Speed)
			}

			if tt.request.Elevenlabs.Model != "" && result.Elevenlabs.Model != tt.request.Elevenlabs.Model {
				t.Errorf("Model = %v, want %v", result.Elevenlabs.Model, tt.request.Elevenlabs.Model)
			}
		})
	}
}
