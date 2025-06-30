package main

import (
	"log"
	"os"
	"strings"

	"github.com/dog0sd/sven/internal/config"
	"github.com/dog0sd/sven/internal/server"
	"github.com/dog0sd/sven/internal/tts"
)

func main() {
	var config, err = config.LoadConfig()
	if err != nil {
		log.Fatal("error loading configuration: err")
	}

	if len(os.Args) == 1 {
		if err := server.StartServer(config.Port, config); err != nil {
			log.Fatal("HTTP server failed: ", err)
		}
	} else {
		var text = strings.TrimSpace(strings.Join(os.Args[1:], " "))
		if text != "" {
			if err = tts.ElevenlabsTTS(config.Elevenlabs, text, ""); err != nil {
				log.Fatal("error transcripting using elevenlabs: ", err)
			}
		}
	}
}
