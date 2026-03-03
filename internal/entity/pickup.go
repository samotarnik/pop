package entity

import "pop/internal/physics"

type Pickup struct {
	Type      string
	X         float64
	Y         float64
	W         float64
	H         float64
	Collected bool
}

func NewPickup(kind string, x, y float64) *Pickup {
	return &Pickup{Type: kind, X: x, Y: y, W: 24, H: 24}
}

func (p *Pickup) Rect() physics.Rect {
	return physics.Rect{X: p.X, Y: p.Y, W: p.W, H: p.H}
}
