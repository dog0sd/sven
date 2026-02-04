//go:build windows || darwin

package audio

import "fmt"

type PulsePlayer struct{}

func (p *PulsePlayer) Play(mp3Data []byte) error {
	return fmt.Errorf("pulse backend is not supported on this OS, use oto instead")
}
