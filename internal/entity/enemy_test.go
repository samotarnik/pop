package entity

import "testing"

func TestBossPressureFactorAffectsTuning(t *testing.T) {
	boss := NewEnemy(300, 420, "boss", 3)
	baseSpeed := boss.moveSpeed()
	baseCD := boss.attackCooldown()
	baseReach := boss.attackReach()

	boss.Health = 1
	if boss.moveSpeed() <= baseSpeed {
		t.Fatalf("expected boss speed to increase as health drops")
	}
	if boss.attackCooldown() >= baseCD {
		t.Fatalf("expected boss cooldown to decrease as health drops")
	}
	if boss.attackReach() <= baseReach {
		t.Fatalf("expected boss reach to increase as health drops")
	}
}
