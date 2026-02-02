package audio

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/hajimehoshi/go-mp3"
	"github.com/jfreymuth/pulse"
)

type PulsePlayer struct{}

func (p *PulsePlayer) Play(mp3Data []byte) error {
	decoder, err := mp3.NewDecoder(bytes.NewReader(mp3Data))
	if err != nil {
		return fmt.Errorf("mp3 decode error: %v", err)
	}

	// Read all PCM data (16-bit signed, stereo)
	var pcmBuf bytes.Buffer
	buf := make([]byte, 8192)
	for {
		n, err := decoder.Read(buf)
		if n > 0 {
			pcmBuf.Write(buf[:n])
		}
		if err != nil {
			break
		}
	}

	// Convert bytes to int16 samples
	raw := pcmBuf.Bytes()
	samples := make([]int16, len(raw)/2)
	for i := range samples {
		samples[i] = int16(binary.LittleEndian.Uint16(raw[i*2 : i*2+2]))
	}

	// Convert int16 to float32 for pulse
	floatSamples := make([]float32, len(samples))
	for i, s := range samples {
		floatSamples[i] = float32(s) / 32768.0
	}

	client, err := pulse.NewClient()
	if err != nil {
		return fmt.Errorf("pulse client error: %v", err)
	}
	defer client.Close()

	playback, err := client.NewPlayback(
		pulse.Float32Reader(func(buf []float32) (int, error) {
			if len(floatSamples) == 0 {
				return 0, pulse.EndOfData
			}
			n := copy(buf, floatSamples)
			floatSamples = floatSamples[n:]
			return n, nil
		}),
		pulse.PlaybackStereo,
		pulse.PlaybackSampleRate(decoder.SampleRate()),
		pulse.PlaybackMediaName("sven"),
	)
	if err != nil {
		return fmt.Errorf("pulse playback error: %v", err)
	}

	playback.Start()
	playback.Drain()
	playback.Stop()

	return nil
}
