package entity

import "pop/internal/physics"

// ResolveStab applies the prince attack to the first intersecting alive enemy.
func ResolveStab(player *Player, enemies []*Enemy) bool {
	atkRect, ok := player.AttackRect()
	if !ok {
		return false
	}
	for _, e := range enemies {
		if !e.Alive() {
			continue
		}
		if physics.Intersects(atkRect, e.Rect()) {
			e.Health--
			player.RegisterAttackHit()
			return true
		}
	}
	return false
}
