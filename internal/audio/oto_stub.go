//go:build !oto

package audio

import "fmt"

type OtoPlayer struct{}

func (p *OtoPlayer) Play(mp3Data []byte) error {
	return fmt.Errorf("oto backend not available: build with -tags oto")
}
