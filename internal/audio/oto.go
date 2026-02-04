//go:build windows || darwin

package audio

import (
	"bytes"
	"fmt"
	"io"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
)

type OtoPlayer struct{}

func (p *OtoPlayer) Play(mp3Data []byte) error {
	decoder, err := mp3.NewDecoder(bytes.NewReader(mp3Data))
	if err != nil {
		return fmt.Errorf("mp3 decode error: %v", err)
	}

	ctx, err := oto.NewContext(decoder.SampleRate(), 2, 2, 8192)
	if err != nil {
		return fmt.Errorf("oto context error: %v", err)
	}
	defer ctx.Close()

	player := ctx.NewPlayer()
	defer player.Close()

	if _, err := io.Copy(player, decoder); err != nil {
		return fmt.Errorf("playing audio error: %v", err)
	}
	return nil
}
