package main

import (
	"context"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/dog0sd/sven/internal/config"
	"github.com/dog0sd/sven/internal/tts"
)


func main() {

	s := server.NewMCPServer("SVEN", "0.2.1", server.WithToolCapabilities(false))
	
	voiceIt := mcp.NewTool("voice_it",
	    mcp.WithDescription("Use this tool to voice out some information to the user."),
		mcp.WithString("text", mcp.Required(), mcp.Description("The text that needs to be spoken"),),
	)
	
	s.AddTool(voiceIt, elevenlabsSynthetic)

	if err := server.ServeStdio(s); err != nil {
		log.Fatal("server error: ", err)
	}
}

func elevenlabsSynthetic(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var err error
	var cfg config.Config
	if cfg, err = config.LoadConfigFromEnv(); err != nil {
		return mcp.NewToolResultErrorFromErr("", err), nil
	}
	
	if text := request.GetString("text", ""); text == "" {
		return mcp.NewToolResultError("no text provided"), nil
	}
	
	if err = tts.ElevenlabsTTS(cfg.Elevenlabs, request.GetString("text", ""), ""); err != nil {
		return mcp.NewToolResultErrorFromErr("playback error: ", err), nil
	}

    return mcp.NewToolResultText("done"), nil
}