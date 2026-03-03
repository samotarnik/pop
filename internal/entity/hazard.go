package entity

import "pop/internal/physics"

type Hazard struct {
	Type  string
	X     float64
	Y     float64
	W     float64
	H     float64
	Speed float64
}

func NewHazard(kind string, x, y, speed float64) *Hazard {
	h := &Hazard{Type: kind, X: x, Y: y, Speed: speed}
	switch kind {
	case "cat":
		h.W = 36
		h.H = 34
	case "hay_bale":
		h.W = 40
		h.H = 30
	default:
		h.W = 30
		h.H = 30
	}
	return h
}

func (h *Hazard) Update(dt, worldWidth float64) {
	if h.Type != "hay_bale" {
		return
	}
	h.X += h.Speed * dt
	if h.X < 0 {
		h.X = 0
		h.Speed *= -1
	}
	if h.X+h.W > worldWidth {
		h.X = worldWidth - h.W
		h.Speed *= -1
	}
}

func (h *Hazard) Rect() physics.Rect {
	return physics.Rect{X: h.X, Y: h.Y, W: h.W, H: h.H}
}
