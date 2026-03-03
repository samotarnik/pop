package game

import "testing"

func TestStateHealClamp(t *testing.T) {
	s := NewState()
	s.Heal(10)
	if s.Health != s.MaxHealth {
		t.Fatalf("expected health %d, got %d", s.MaxHealth, s.Health)
	}
}

func TestStateDamageAndDeath(t *testing.T) {
	s := NewState()
	dead := s.Damage(2)
	if dead {
		t.Fatalf("should not be dead after 2 damage")
	}
	dead = s.Damage(1)
	if !dead {
		t.Fatalf("expected death after reaching 0 health")
	}
}

func TestLoseLifeGameOver(t *testing.T) {
	s := NewState()
	s.Lives = 1
	if !s.LoseLife() {
		t.Fatalf("expected game over when last life is consumed")
	}
}
