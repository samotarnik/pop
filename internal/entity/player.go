package entity

import (
	"math"
	"pop/internal/config"
	"pop/internal/physics"
)

type Player struct {
	X float64
	Y float64

	vx float64
	vy float64

	Facing int

	W float64
	H float64

	Ducking  bool
	OnGround bool

	Attacking   bool
	attackTimer float64
	attackHit   bool
	invulnTimer float64
}

func NewPlayer(spawnX, groundY float64) *Player {
	p := &Player{
		X:      spawnX,
		Y:      groundY - config.PlayerHeight,
		W:      config.PlayerWidth,
		H:      config.PlayerHeight,
		Facing: 1,
	}
	p.OnGround = true
	return p
}

func (p *Player) Update(dt float64, in Input, groundY float64) {
	if p.invulnTimer > 0 {
		p.invulnTimer -= dt
		if p.invulnTimer < 0 {
			p.invulnTimer = 0
		}
	}

	if p.Attacking {
		p.attackTimer -= dt
		if p.attackTimer <= 0 {
			p.Attacking = false
			p.attackHit = false
		}
	}

	if !p.Attacking {
		if in.Duck && p.OnGround {
			if !p.Ducking {
				p.Y += config.PlayerHeight - config.PlayerDuckHeight
			}
			p.Ducking = true
			p.H = config.PlayerDuckHeight
		} else {
			if p.Ducking {
				p.Y -= config.PlayerHeight - config.PlayerDuckHeight
			}
			p.Ducking = false
			p.H = config.PlayerHeight
		}

		p.vx = 0
		if !p.Ducking {
			if in.Left {
				p.vx = -config.PlayerRunSpeed
				p.Facing = -1
			}
			if in.Right {
				p.vx = config.PlayerRunSpeed
				p.Facing = 1
			}
		}

		if in.Jump && p.OnGround {
			p.vy = config.PlayerJumpSpeed
			p.OnGround = false
		}

		if in.Attack && p.CanAttack() {
			p.Attacking = true
			p.attackTimer = config.AttackDurationSec
			p.attackHit = false
			p.vx = 0
		}
	} else {
		p.vx = 0
	}

	if !p.OnGround {
		p.vy += config.Gravity * dt
	}

	p.X += p.vx * dt
	p.Y += p.vy * dt

	if p.Y+p.H >= groundY {
		p.Y = groundY - p.H
		p.vy = 0
		p.OnGround = true
	}
}

func (p *Player) CanAttack() bool {
	return p.OnGround && !p.Ducking && math.Abs(p.vx) < 0.5 && !p.Attacking
}

func (p *Player) Rect() physics.Rect {
	return physics.Rect{X: p.X, Y: p.Y, W: p.W, H: p.H}
}

func (p *Player) AttackRect() (physics.Rect, bool) {
	if !p.Attacking || p.attackHit {
		return physics.Rect{}, false
	}
	if p.Facing >= 0 {
		return physics.Rect{X: p.X + p.W, Y: p.Y + 14, W: 28, H: p.H - 18}, true
	}
	return physics.Rect{X: p.X - 28, Y: p.Y + 14, W: 28, H: p.H - 18}, true
}

func (p *Player) RegisterAttackHit() {
	p.attackHit = true
}

func (p *Player) CanBeHit() bool {
	return p.invulnTimer <= 0
}

func (p *Player) RegisterHit() {
	p.invulnTimer = config.HitIFramesSec
}
