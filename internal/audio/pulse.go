//go:build linux

package audio

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"

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

	sampleRate := decoder.SampleRate()
	totalFrames := len(samples) / 2 // stereo
	audioDuration := time.Duration(int64(totalFrames)*1000/int64(sampleRate)) * time.Millisecond

	client, err := pulse.NewClient()
	if err != nil {
		return fmt.Errorf("pulse client error: %v", err)
	}
	defer client.Close()

	// We must NOT return EndOfData before Start() completes, because the
	// pulse library's Start() waits on <-p.started, but the "started"
	// callback only fires when state==running. If the reader returns
	// EndOfData during the first buffer fill (before "Started" arrives),
	// the state transitions to idle and Start() deadlocks.
	//
	// Solution: once all real audio data is consumed, return silence instead
	// of EndOfData. We stop the stream ourselves after the audio duration.
	offset := 0
	totalSamples := len(floatSamples)

	playback, err := client.NewPlayback(
		pulse.Float32Reader(func(buf []float32) (int, error) {
			if offset >= totalSamples {
				for i := range buf {
					buf[i] = 0
				}
				return len(buf), nil
			}
			remaining := floatSamples[offset:]
			n := copy(buf, remaining)
			offset += n
			return n, nil
		}),
		pulse.PlaybackStereo,
		pulse.PlaybackSampleRate(sampleRate),
		pulse.PlaybackMediaName("sven"),
	)
	if err != nil {
		return fmt.Errorf("pulse playback error: %v", err)
	}

	playback.Start()

	// Wait for the audio to finish playing, plus buffer for PulseAudio latency
	time.Sleep(audioDuration + 300*time.Millisecond)

	playback.Stop()
	playback.Close()

	return nil
}
