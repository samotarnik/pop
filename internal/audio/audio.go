package audio

import (
	"bytes"
	"encoding/binary"
	"math"

	ebiaudio "github.com/hajimehoshi/ebiten/v2/audio"
)

const sampleRate = 44100

type Manager struct {
	ctx        *ebiaudio.Context
	player     *ebiaudio.Player
	muted      bool
	baseVolume float64
	currentBGM string
}

func NewManager() *Manager {
	return &Manager{
		ctx:        ebiaudio.NewContext(sampleRate),
		baseVolume: 0.16,
	}
}

func (m *Manager) ToggleMute() {
	m.muted = !m.muted
	m.applyVolume()
}

func (m *Manager) IsMuted() bool {
	return m.muted
}

func (m *Manager) PlayBGM(track string) {
	if m.player != nil && m.currentBGM == track {
		if !m.player.IsPlaying() {
			m.player.Play()
		}
		return
	}
	if m.player != nil {
		_ = m.player.Close()
		m.player = nil
	}

	pcm := buildLoopPCM()
	loop := ebiaudio.NewInfiniteLoop(bytes.NewReader(pcm), int64(len(pcm)))
	player, err := m.ctx.NewPlayer(loop)
	if err != nil {
		return
	}
	m.player = player
	m.currentBGM = track
	m.applyVolume()
	m.player.Play()
}

func (m *Manager) applyVolume() {
	if m.player == nil {
		return
	}
	if m.muted {
		m.player.SetVolume(0)
		return
	}
	m.player.SetVolume(m.baseVolume)
}

func buildLoopPCM() []byte {
	notes := []float64{
		261.63, 329.63, 392.00, 523.25,
		392.00, 329.63, 293.66, 329.63,
		261.63, 196.00, 261.63, 329.63,
		392.00, 523.25, 659.25, 523.25,
	}
	noteSeconds := 0.18
	totalSamples := int(float64(sampleRate) * noteSeconds * float64(len(notes)))
	out := make([]byte, 0, totalSamples*4)
	amp := 0.22

	for _, freq := range notes {
		n := int(float64(sampleRate) * noteSeconds)
		for i := 0; i < n; i++ {
			v := square(freq, i, sampleRate) * amp
			sample := int16(v * math.MaxInt16)
			out = binary.LittleEndian.AppendUint16(out, uint16(sample))
			out = binary.LittleEndian.AppendUint16(out, uint16(sample))
		}
	}
	return out
}

func square(freq float64, sampleIndex, rate int) float64 {
	period := float64(rate) / freq
	if math.Mod(float64(sampleIndex), period) < period/2 {
		return 1
	}
	return -1
}
