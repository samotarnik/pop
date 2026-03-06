package audio

import (
	"bytes"
	"encoding/binary"
	"math"

	ebiaudio "github.com/hajimehoshi/ebiten/v2/audio"
)

const sampleRate = 44100

type Note struct {
	// Freq is in Hz. Use 0 for a rest (silence).
	Freq float64
	// Beats uses musical length in beats:
	// 0.5=eighth note, 1=quarter note, 2=half note.
	Beats float64
}

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
	const bpm = 110.0
	secondsPerBeat := 60.0 / bpm
	melody := []Note{
		// MARKO SKACE!
		// https://www.youtube.com/watch?v=5q3Uc8Jk3KQ
		{Freq: 329.63, Beats: 0.5}, // E4, eigth note
		{Freq: 392.00, Beats: 0.5},
		{Freq: 392.00, Beats: 0.5},
		{Freq: 349.23, Beats: 0.5},
		{Freq: 329.63, Beats: 0.5},
		{Freq: 392.00, Beats: 0.5},
		{Freq: 392.00, Beats: 0.5},
		{Freq: 349.23, Beats: 0.5},
		// --
		{Freq: 329.63, Beats: 0.5},
		{Freq: 329.63, Beats: 0.5},
		{Freq: 293.66, Beats: 0.5},
		{Freq: 293.66, Beats: 0.5},
		{Freq: 261.63, Beats: 0.5},
		{Freq: 246.94, Beats: 0.5},
		{Freq: 220.00, Beats: 0.5},
		{Freq: 196.00, Beats: 0.5},
		// --
		{Freq: 261.63, Beats: 0.5},
		{Freq: 196.00, Beats: 0.5},
		{Freq: 293.66, Beats: 0.5},
		{Freq: 196.00, Beats: 0.5},
		{Freq: 329.63, Beats: 0.5},
		{Freq: 392.00, Beats: 0.5},
		{Freq: 392.00, Beats: 0.5},
		{Freq: 349.23, Beats: 0.5},
		{Freq: 329.63, Beats: 0.5},
		{Freq: 329.63, Beats: 0.5},
		{Freq: 293.66, Beats: 0.5},
		{Freq: 293.66, Beats: 0.5},
		{Freq: 261.63, Beats: 0.5},
		{Freq: 246.94, Beats: 0.5},
		{Freq: 220.00, Beats: 0.5},
		{Freq: 196.00, Beats: 0.5},
		// --
		{Freq: 261.63, Beats: 0.5},
		{Freq: 196.00, Beats: 0.5},
		{Freq: 293.66, Beats: 0.5},
		{Freq: 196.00, Beats: 0.5},
		{Freq: 329.63, Beats: 0.5},
		{Freq: 392.00, Beats: 0.5},
		{Freq: 392.00, Beats: 1.0},
		{Freq: 329.63, Beats: 0.5},
		{Freq: 329.63, Beats: 0.5},
		{Freq: 293.66, Beats: 0.5},
		{Freq: 293.66, Beats: 0.5},
		{Freq: 261.63, Beats: 1.0},
		{Freq: 261.63, Beats: 1.0},
	}
	amp := 0.22

	totalSamples := 0
	for _, n := range melody {
		totalSamples += int(float64(sampleRate) * secondsPerBeat * n.Beats)
	}
	out := make([]byte, 0, totalSamples*4)

	for _, note := range melody {
		n := int(float64(sampleRate) * secondsPerBeat * note.Beats)
		for i := 0; i < n; i++ {
			v := 0.0
			if note.Freq > 0 {
				v = square(note.Freq, i, sampleRate) * amp
			}
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
