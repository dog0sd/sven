package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/dog0sd/sven/internal/config"
	"github.com/dog0sd/sven/internal/tts"
)

func StartServer(port string, config config.Config) error {
	http.HandleFunc("/tts", func(w http.ResponseWriter, r *http.Request) {
		handleTTS(w, r, config)
	})
	return http.ListenAndServe(port, nil)
}

type elevenLabsConfig struct {
	Model           string  `json:"model"`
	SimilarityBoost float32 `json:"similarity_boost"`
	Stability       float32 `json:"stability"`
	Style           float32 `json:"style"`
}

type TTSRequest struct {
	Text       string           `json:"text"`
	Elevenlabs elevenLabsConfig `json:"11labs"`
}

func handleTTS(w http.ResponseWriter, r *http.Request, config config.Config) {
	reqBody := TTSRequest{}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error decoding request")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	mergeElevenLabsTTSSettings(reqBody, &config)
	err = tts.ElevenlabsTTS(reqBody.Text, config.Elevenlabs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "elevelabs error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("OK"))
}

func mergeElevenLabsTTSSettings(req TTSRequest, config *config.Config) {
	if req.Elevenlabs.Model != "" {
		config.Elevenlabs.Model = req.Elevenlabs.Model
	}
	if req.Elevenlabs.SimilarityBoost != 0.0 {
		config.Elevenlabs.Settings.SimilarityBoost = req.Elevenlabs.SimilarityBoost
	}
	if req.Elevenlabs.Stability != 0.0 {
		config.Elevenlabs.Settings.Stability = req.Elevenlabs.Stability
	}
	if req.Elevenlabs.Style != 0.0 {
		config.Elevenlabs.Settings.Style = req.Elevenlabs.Style
	}
}
