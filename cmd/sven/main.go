package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dog0sd/sven/internal/config"
	"github.com/dog0sd/sven/internal/server"
	"github.com/dog0sd/sven/internal/tts"
)

func main() {
	flag.Usage = func() {
		fmt.Printf("Usage: %s [text]\n", os.Args[0])
		fmt.Printf("CLI mode: %s 'Hello, world!'\n", os.Args[0])
		fmt.Println("OR")
		fmt.Printf("Server mode: %s\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(0)
	}
	flag.Parse()
	var config, err = config.LoadConfig()
	if err != nil {
		log.Fatalf("error loading configuration: %v", err)
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
