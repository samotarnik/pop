package config

const (
	InternalWidth  = 640
	InternalHeight = 480
	WindowScale    = 2
	TPS            = 60
)

const (
	PlayerWidth       = 36
	PlayerHeight      = 64
	PlayerDuckHeight  = 36
	PlayerRunSpeed    = 180.0
	PlayerJumpSpeed   = -460.0
	Gravity           = 1020.0
	AttackDurationSec = 0.14
	HitIFramesSec     = 0.50
	JumpCutGravityMul = 1.55
)

const (
	StartHealth = 3
	MaxHealth   = 5
	StartLives  = 3
)

const (
	EnemyAttackDamage = 1
)

const (
	FarmerMoveSpeed = 76.0
	BossMoveSpeed   = 106.0

	FarmerAttackCooldownSec = 1.05
	BossAttackCooldownSec   = 0.58
	BossCooldownStepSec     = 0.08
	BossMinCooldownSec      = 0.34

	FarmerAttackReach = 38.0
	BossAttackReach   = 52.0
	BossReachStep     = 3.0
	BossMoveStep      = 10.0

	EnemyNearSlowdownFactor = 0.58
	EnemyNearSlowdownRadius = 1.7
)
