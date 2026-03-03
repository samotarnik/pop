package game

import "pop/internal/config"

type State struct {
	CurrentLevel int
	Health       int
	MaxHealth    int
	Lives        int
	Score        int
	Muted        bool
}

func NewState() *State {
	s := &State{}
	s.ResetRun()
	return s
}

func (s *State) ResetRun() {
	s.CurrentLevel = 0
	s.Health = config.StartHealth
	s.MaxHealth = config.MaxHealth
	s.Lives = config.StartLives
	s.Score = 0
}

func (s *State) Heal(n int) {
	s.Health += n
	if s.Health > s.MaxHealth {
		s.Health = s.MaxHealth
	}
}

func (s *State) Damage(n int) bool {
	s.Health -= n
	if s.Health <= 0 {
		s.Health = 0
		return true
	}
	return false
}

func (s *State) Kill() bool {
	s.Health = 0
	return true
}

func (s *State) LoseLife() bool {
	s.Lives--
	return s.Lives <= 0
}

func (s *State) Respawn() {
	s.Health = config.StartHealth
}
