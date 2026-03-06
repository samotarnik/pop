package sprite

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Animator struct {
	sheet    *Sheet
	clipName string
	framePos float64
}

func NewAnimator(sheet *Sheet, clip string) *Animator {
	a := &Animator{sheet: sheet}
	a.SetClip(clip)
	return a
}

func (a *Animator) SetClip(clip string) {
	if a.sheet == nil {
		return
	}
	if _, ok := a.sheet.Clips[clip]; !ok {
		clip = "idle"
	}
	if a.clipName == clip {
		return
	}
	a.clipName = clip
	a.framePos = 0
}

func (a *Animator) Update(dt float64) {
	if a.sheet == nil {
		return
	}
	clip, ok := a.sheet.Clips[a.clipName]
	if !ok || len(clip.Frames) == 0 {
		return
	}
	a.framePos += dt * clip.FPS
	if clip.Loop {
		a.framePos = math.Mod(a.framePos, float64(len(clip.Frames)))
		if a.framePos < 0 {
			a.framePos += float64(len(clip.Frames))
		}
		return
	}
	maxFrame := float64(len(clip.Frames) - 1)
	if a.framePos > maxFrame {
		a.framePos = maxFrame
	}
}

func (a *Animator) Draw(screen *ebiten.Image, footX, footY, cameraX float64, facing int) {
	if a.sheet == nil {
		return
	}
	clip, ok := a.sheet.Clips[a.clipName]
	if !ok || len(clip.Frames) == 0 {
		return
	}
	idx := int(a.framePos)
	if idx < 0 {
		idx = 0
	}
	if idx >= len(clip.Frames) {
		idx = len(clip.Frames) - 1
	}
	a.sheet.drawFrame(screen, clip.Frames[idx], footX, footY, cameraX, facing)
}
