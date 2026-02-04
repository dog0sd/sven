//go:build linux

package audio

import "fmt"

type OtoPlayer struct{}

func (p *OtoPlayer) Play(mp3Data []byte) error {
	return fmt.Errorf("oto backend is not supported on Linux, use pulse instead")
}
