package entity

import (
	"math"
	"pop/internal/config"
	"pop/internal/physics"
)

type Enemy struct {
	X float64
	Y float64
	W float64
	H float64

	Kind      string
	Health    int
	MaxHealth int
	Facing    int

	attackCD float64
}

func NewEnemy(x, groundY float64, kind string, health int) *Enemy {
	h := 62.0
	w := 34.0
	if kind == "boss" {
		h = 74
		w = 42
	}
	return &Enemy{
		X:         x,
		Y:         groundY - h,
		W:         w,
		H:         h,
		Kind:      kind,
		Health:    health,
		MaxHealth: health,
		Facing:    -1,
	}
}

func (e *Enemy) Alive() bool {
	return e.Health > 0
}

func (e *Enemy) Update(dt, playerX float64) {
	if !e.Alive() {
		return
	}
	if e.attackCD > 0 {
		e.attackCD -= dt
	}

	dist := playerX - e.X
	if dist < 0 {
		e.Facing = -1
	} else {
		e.Facing = 1
	}

	reach := e.attackReach()
	if math.Abs(dist) <= reach {
		return
	}

	speed := e.moveSpeed()
	if math.Abs(dist) <= reach*config.EnemyNearSlowdownRadius {
		speed *= config.EnemyNearSlowdownFactor
	}
	if dist < 0 {
		e.X -= speed * dt
	} else {
		e.X += speed * dt
	}
}

func (e *Enemy) CanAttack(playerRect physics.Rect) bool {
	if !e.Alive() || e.attackCD > 0 {
		return false
	}
	return physics.Intersects(e.AttackRect(), playerRect)
}

func (e *Enemy) RegisterAttack() {
	e.attackCD = e.attackCooldown()
}

func (e *Enemy) Rect() physics.Rect {
	return physics.Rect{X: e.X, Y: e.Y, W: e.W, H: e.H}
}

func (e *Enemy) AttackRect() physics.Rect {
	reach := e.attackReach() * 0.52
	if e.Facing >= 0 {
		return physics.Rect{X: e.X + e.W, Y: e.Y + 12, W: reach, H: e.H - 16}
	}
	return physics.Rect{X: e.X - reach, Y: e.Y + 12, W: reach, H: e.H - 16}
}

func (e *Enemy) moveSpeed() float64 {
	if e.Kind == "boss" {
		return config.BossMoveSpeed + config.BossMoveStep*e.pressureFactor()
	}
	return config.FarmerMoveSpeed
}

func (e *Enemy) attackCooldown() float64 {
	if e.Kind == "boss" {
		cd := config.BossAttackCooldownSec - config.BossCooldownStepSec*e.pressureFactor()
		if cd < config.BossMinCooldownSec {
			return config.BossMinCooldownSec
		}
		return cd
	}
	return config.FarmerAttackCooldownSec
}

func (e *Enemy) attackReach() float64 {
	if e.Kind == "boss" {
		return config.BossAttackReach + config.BossReachStep*e.pressureFactor()
	}
	return config.FarmerAttackReach
}

func (e *Enemy) pressureFactor() float64 {
	if e.MaxHealth <= 1 {
		return 0
	}
	missing := e.MaxHealth - e.Health
	if missing < 0 {
		missing = 0
	}
	return float64(missing) / float64(e.MaxHealth-1)
}
