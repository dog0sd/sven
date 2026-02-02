package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/dog0sd/sven/internal/audio"
	"github.com/dog0sd/sven/internal/config"
	"github.com/dog0sd/sven/internal/tts"
)

func StartServer(port string, config config.Config, player audio.Player) error {
	http.HandleFunc("/tts", func(w http.ResponseWriter, r *http.Request) {
		handleTTS(w, r, config, player)
	})
	return http.ListenAndServe(port, nil)
}

type elevenLabsConfig struct {
	Model           string  `json:"model"`
	SimilarityBoost float32 `json:"similarity_boost"`
	Stability       float32 `json:"stability"`
	Style           float32 `json:"style"`
	Speed           float32 `json:"speed"`
}

type TTSRequest struct {
	Text       string           `json:"text"`
	Elevenlabs elevenLabsConfig `json:"voice_settings"`
	PText      string           `json:"ptext"` // Previous text
}

func handleTTS(w http.ResponseWriter, r *http.Request, config config.Config, player audio.Player) {
	reqBody := TTSRequest{}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error decoding request")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	reqConfig := mergeElevenLabsTTSSettings(reqBody, &config)
	mp3Data, err := tts.Synthesize(reqConfig.Elevenlabs, reqBody.Text, reqBody.PText)
	if err != nil {
		fmt.Fprintf(os.Stderr, "elevenlabs error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if err := player.Play(mp3Data); err != nil {
		fmt.Fprintf(os.Stderr, "playback error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("OK"))
}

func mergeElevenLabsTTSSettings(req TTSRequest, c *config.Config) *config.Config {
	reqConfig := *c
	if req.Elevenlabs.Model != "" {
		reqConfig.Elevenlabs.Model = req.Elevenlabs.Model
	}
	if req.Elevenlabs.SimilarityBoost != 0.0 {
		reqConfig.Elevenlabs.Settings.SimilarityBoost = req.Elevenlabs.SimilarityBoost
	}
	if req.Elevenlabs.Stability != 0.0 {
		reqConfig.Elevenlabs.Settings.Stability = req.Elevenlabs.Stability
	}
	if req.Elevenlabs.Style != 0.0 {
		reqConfig.Elevenlabs.Settings.Style = req.Elevenlabs.Style
	}
	if req.Elevenlabs.Speed != 1.0 {
		reqConfig.Elevenlabs.Settings.Speed = req.Elevenlabs.Speed
	}
	return &reqConfig
}
