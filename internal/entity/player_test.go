package entity

import "testing"

func TestCanAttackOnlyStanding(t *testing.T) {
	p := NewPlayer(0, 420)
	if !p.CanAttack() {
		t.Fatalf("expected new player to be able to attack while standing")
	}

	p.Update(1.0/60.0, Input{Right: true}, 420)
	if p.CanAttack() {
		t.Fatalf("expected moving player to be unable to attack")
	}
}

func TestResolveStabReducesEnemyHealth(t *testing.T) {
	p := NewPlayer(100, 420)
	enemy := NewEnemy(136, 420, "farmer", 2)

	p.Update(1.0/60.0, Input{Attack: true}, 420)
	hit := ResolveStab(p, []*Enemy{enemy})
	if !hit {
		t.Fatalf("expected stab to hit enemy")
	}
	if enemy.Health != 1 {
		t.Fatalf("expected enemy health to become 1, got %d", enemy.Health)
	}
}
