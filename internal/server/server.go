package server

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dog0sd/sven/internal/audio"
	"github.com/dog0sd/sven/internal/config"
	"github.com/dog0sd/sven/internal/tts"
)

const shutdownTimeout = 5 * time.Second

func StartServer(listen string, cfg config.Config, player audio.Player) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/tts", func(w http.ResponseWriter, r *http.Request) {
		handleTTS(w, r, cfg, player)
	})

	server := &http.Server{
		Addr:    listen,
		Handler: mux,
	}

	// Channel to listen for OS signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		slog.Info("starting HTTP server", "listen", listen)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("HTTP server error", "error", err)
		}
	}()

	// Wait for shutdown signal
	<-stop
	slog.Info("shutting down server...")

	// Create a deadline for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return err
	}
	slog.Info("server stopped")
	return nil
}

type elevenLabsConfig struct {
	Model           string   `json:"model"`
	SimilarityBoost *float32 `json:"similarity_boost"`
	Stability       *float32 `json:"stability"`
	Style           *float32 `json:"style"`
	Speed           *float32 `json:"speed"`
}

type TTSRequest struct {
	Text       string           `json:"text"`
	Elevenlabs elevenLabsConfig `json:"voice_settings"`
	PText      string           `json:"ptext"` // Previous text
}

func handleTTS(w http.ResponseWriter, r *http.Request, cfg config.Config, player audio.Player) {
	reqBody := TTSRequest{}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		slog.Error("request decode failed", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	reqConfig := mergeElevenLabsTTSSettings(reqBody, &cfg)
	mp3Data, err := tts.Synthesize(reqConfig.Elevenlabs, reqBody.Text, reqBody.PText)
	if err != nil {
		slog.Error("elevenlabs error", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if err := player.Play(mp3Data); err != nil {
		slog.Error("playback error", "error", err)
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
	if req.Elevenlabs.SimilarityBoost != nil {
		reqConfig.Elevenlabs.Settings.SimilarityBoost = *req.Elevenlabs.SimilarityBoost
	}
	if req.Elevenlabs.Stability != nil {
		reqConfig.Elevenlabs.Settings.Stability = *req.Elevenlabs.Stability
	}
	if req.Elevenlabs.Style != nil {
		reqConfig.Elevenlabs.Settings.Style = *req.Elevenlabs.Style
	}
	if req.Elevenlabs.Speed != nil {
		reqConfig.Elevenlabs.Settings.Speed = *req.Elevenlabs.Speed
	}
	return &reqConfig
}
