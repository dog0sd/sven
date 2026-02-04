package audio

import (
	"fmt"
	"runtime"
)

type Player interface {
	Play(mp3Data []byte) error
}

// DefaultBackend returns the default audio backend for the current OS.
func DefaultBackend() string {
	if runtime.GOOS == "linux" {
		return "pulse"
	}
	return "oto"
}

func NewPlayer(backend string) (Player, error) {
	switch backend {
	case "pulse":
		return &PulsePlayer{}, nil
	case "oto":
		return &OtoPlayer{}, nil
	default:
		return nil, fmt.Errorf("unknown audio backend: %q (supported: pulse, oto)", backend)
	}
}
