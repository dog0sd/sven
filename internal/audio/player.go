package audio

import "fmt"

type Player interface {
	Play(mp3Data []byte) error
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
