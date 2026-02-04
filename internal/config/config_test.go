package config

import (
	"testing"
)

func TestValidateElevenLabsSettings(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid settings",
			config: Config{
				Elevenlabs: ElevenLabsConfig{
					Settings: ElevenlabsVoiceSettings{
						SimilarityBoost: 0.5,
						Stability:       0.8,
						Style:           0.3,
						Speed:           1.0,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "similarity_boost too high",
			config: Config{
				Elevenlabs: ElevenLabsConfig{
					Settings: ElevenlabsVoiceSettings{
						SimilarityBoost: 1.5,
						Speed:           1.0,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "stability too high",
			config: Config{
				Elevenlabs: ElevenLabsConfig{
					Settings: ElevenlabsVoiceSettings{
						Stability: 1.5,
						Speed:     1.0,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "style too high",
			config: Config{
				Elevenlabs: ElevenLabsConfig{
					Settings: ElevenlabsVoiceSettings{
						Style: 1.5,
						Speed: 1.0,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "speed too high",
			config: Config{
				Elevenlabs: ElevenLabsConfig{
					Settings: ElevenlabsVoiceSettings{
						Speed: 1.5,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "speed too low",
			config: Config{
				Elevenlabs: ElevenLabsConfig{
					Settings: ElevenlabsVoiceSettings{
						Speed: 0.5,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "zero values are valid",
			config: Config{
				Elevenlabs: ElevenLabsConfig{
					Settings: ElevenlabsVoiceSettings{
						SimilarityBoost: 0.0,
						Stability:       0.0,
						Style:           0.0,
						Speed:           0.7,
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateElevenLabsSettings(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateElevenLabsSettings() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateAudioBackend(t *testing.T) {
	tests := []struct {
		name            string
		backend         string
		wantErr         bool
		expectedBackend string
	}{
		{
			name:            "pulse backend",
			backend:         "pulse",
			wantErr:         false,
			expectedBackend: "pulse",
		},
		{
			name:            "oto backend",
			backend:         "oto",
			wantErr:         false,
			expectedBackend: "oto",
		},
		{
			name:            "empty defaults to OS default",
			backend:         "",
			wantErr:         false,
			expectedBackend: "", // will be set by the function
		},
		{
			name:    "invalid backend",
			backend: "alsa",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{AudioBackend: tt.backend}
			err := validateAudioBackend(cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateAudioBackend() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tt.expectedBackend != "" && cfg.AudioBackend != tt.expectedBackend {
				t.Errorf("validateAudioBackend() backend = %v, want %v", cfg.AudioBackend, tt.expectedBackend)
			}
			if !tt.wantErr && tt.backend == "" && cfg.AudioBackend == "" {
				t.Errorf("validateAudioBackend() should set default backend for empty input")
			}
		})
	}
}
