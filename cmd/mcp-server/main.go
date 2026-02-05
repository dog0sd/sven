package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/dog0sd/sven/internal/audio"
	"github.com/dog0sd/sven/internal/config"
	"github.com/dog0sd/sven/internal/tts"
)

func main() {

	s := server.NewMCPServer("SVEN", "0.2.1", server.WithToolCapabilities(false))

	voiceIt := mcp.NewTool("voice_it",
		mcp.WithDescription("Use this tool to voice out some information to the user."),
		mcp.WithString("text", mcp.Required(), mcp.Description("The text that needs to be spoken")),
	)

	s.AddTool(voiceIt, elevenlabsSynthetic)

	if err := server.ServeStdio(s); err != nil {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}

func elevenlabsSynthetic(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	cfg, err := config.LoadConfigFromEnv()
	if err != nil {
		return mcp.NewToolResultErrorFromErr("", err), nil
	}

	text := request.GetString("text", "")
	if text == "" {
		return mcp.NewToolResultError("no text provided"), nil
	}

	mp3Data, err := tts.Synthesize(cfg.Elevenlabs, text, "")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("synthesis error: ", err), nil
	}

	player, err := audio.NewPlayer(cfg.AudioBackend)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("audio backend error: ", err), nil
	}

	if err := player.Play(mp3Data); err != nil {
		return mcp.NewToolResultErrorFromErr("playback error: ", err), nil
	}

	return mcp.NewToolResultText("done"), nil
}
