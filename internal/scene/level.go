package scene

import (
	"fmt"
	"image/color"

	"pop/internal/config"
	"pop/internal/entity"
	"pop/internal/game"
	"pop/internal/level"
	"pop/internal/physics"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type LevelScene struct {
	levelIndex int

	player  *entity.Player
	enemies []*entity.Enemy
	hazards []*entity.Hazard
	pickups []*entity.Pickup

	cameraX float64
}

func NewLevelScene(levelIndex int) *LevelScene {
	return &LevelScene{levelIndex: levelIndex}
}

func (l *LevelScene) setup(g *game.Game) {
	data := g.Levels[l.levelIndex]
	g.Audio.PlayBGM(data.BGM)
	l.player = entity.NewPlayer(data.PlayerSpawn, data.GroundY)
	l.enemies = make([]*entity.Enemy, 0, len(data.Enemies))
	for _, e := range data.Enemies {
		l.enemies = append(l.enemies, entity.NewEnemy(e.X, data.GroundY, e.Kind, e.Health))
	}
	l.hazards = make([]*entity.Hazard, 0, len(data.Hazards))
	for _, h := range data.Hazards {
		l.hazards = append(l.hazards, entity.NewHazard(h.Type, h.X, h.Y, h.Speed))
	}
	l.pickups = make([]*entity.Pickup, 0, len(data.Pickups))
	for _, p := range data.Pickups {
		l.pickups = append(l.pickups, entity.NewPickup(p.Type, p.X, p.Y))
	}
	g.State.CurrentLevel = l.levelIndex
}

func (l *LevelScene) Update(g *game.Game) error {
	if l.player == nil {
		l.setup(g)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.SetScene(NewPauseScene(l))
		return nil
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		l.setup(g)
		return nil
	}

	data := g.Levels[l.levelIndex]
	dt := 1.0 / config.TPS

	in := entity.Input{
		Left:   ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyJ),
		Right:  ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyL),
		Jump:   inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyI),
		Duck:   ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyK),
		Attack: inpututil.IsKeyJustPressed(ebiten.KeySpace),
	}

	l.player.Update(dt, in, data.GroundY)
	if l.player.X < 0 {
		l.player.X = 0
	}
	if l.player.X+l.player.W > data.Width {
		l.player.X = data.Width - l.player.W
	}

	for _, h := range l.hazards {
		h.Update(dt, data.Width)
	}
	for _, e := range l.enemies {
		e.Update(dt, l.player.X)
	}

	playerRect := l.player.Rect()
	if entity.ResolveStab(l.player, l.enemies) {
		g.State.Heal(1)
	}

	for _, e := range l.enemies {
		if !e.CanAttack(playerRect) {
			continue
		}
		e.RegisterAttack()
		if l.player.CanBeHit() {
			dead := g.State.Damage(config.EnemyAttackDamage)
			l.player.RegisterHit()
			if dead {
				l.handleDeath(g)
				return nil
			}
		}
	}

	for _, h := range l.hazards {
		if !physics.Intersects(h.Rect(), playerRect) {
			continue
		}
		switch h.Type {
		case "cat":
			g.State.Kill()
			l.handleDeath(g)
			return nil
		case "hay_bale":
			if l.player.CanBeHit() {
				dead := g.State.Damage(1)
				l.player.RegisterHit()
				if dead {
					l.handleDeath(g)
					return nil
				}
			}
		}
	}

	for _, p := range l.pickups {
		if p.Collected {
			continue
		}
		if physics.Intersects(p.Rect(), playerRect) {
			p.Collected = true
			if p.Type == "beer" {
				g.State.Heal(1)
			}
		}
	}

	if l.player.X+l.player.W >= data.GoalX {
		l.handleLevelComplete(g)
		return nil
	}

	l.cameraX = l.player.X - float64(config.InternalWidth)/2
	if l.cameraX < 0 {
		l.cameraX = 0
	}
	maxCam := data.Width - float64(config.InternalWidth)
	if maxCam < 0 {
		maxCam = 0
	}
	if l.cameraX > maxCam {
		l.cameraX = maxCam
	}

	return nil
}

func (l *LevelScene) handleDeath(g *game.Game) {
	if g.State.LoseLife() {
		g.SetScene(NewGameOverScene())
		return
	}
	g.State.Respawn()
	l.setup(g)
}

func (l *LevelScene) handleLevelComplete(g *game.Game) {
	data := g.Levels[l.levelIndex]
	if l.levelIndex == len(g.Levels)-1 {
		g.SetScene(NewStoryScene("outro", func(_ *game.Game) game.Scene {
			return NewVictoryScene()
		}))
		return
	}
	next := l.levelIndex + 1
	nextStory := data.NextStoryID
	g.SetScene(NewStoryScene(nextStory, func(_ *game.Game) game.Scene {
		return NewLevelScene(next)
	}))
}

func (l *LevelScene) Draw(g *game.Game, screen *ebiten.Image) {
	if l.player == nil {
		ebitenutil.DebugPrint(screen, "Loading level...")
		return
	}

	data := g.Levels[l.levelIndex]
	l.drawBackground(screen, data)

	ebitenutil.DrawRect(screen, 0, data.GroundY, float64(config.InternalWidth), float64(config.InternalHeight)-data.GroundY, color.RGBA{R: 0x5a, G: 0xa0, B: 0x47, A: 0xff})

	for _, h := range l.hazards {
		x := h.X - l.cameraX
		if !onScreen(x, h.W) {
			continue
		}
		c := color.RGBA{R: 0x33, G: 0x33, B: 0x33, A: 0xff}
		if h.Type == "hay_bale" {
			c = color.RGBA{R: 0xcc, G: 0x9a, B: 0x32, A: 0xff}
		}
		ebitenutil.DrawRect(screen, x, h.Y, h.W, h.H, c)
	}

	for _, p := range l.pickups {
		if p.Collected {
			continue
		}
		x := p.X - l.cameraX
		if !onScreen(x, p.W) {
			continue
		}
		ebitenutil.DrawRect(screen, x, p.Y, p.W, p.H, color.RGBA{R: 0xff, G: 0xd4, B: 0x3b, A: 0xff})
	}

	for _, e := range l.enemies {
		if !e.Alive() {
			continue
		}
		x := e.X - l.cameraX
		if !onScreen(x, e.W) {
			continue
		}
		c := color.RGBA{R: 0xb1, G: 0x37, B: 0x2e, A: 0xff}
		if e.Kind == "boss" {
			c = color.RGBA{R: 0x72, G: 0x1f, B: 0x17, A: 0xff}
		}
		ebitenutil.DrawRect(screen, x, e.Y, e.W, e.H, c)
		l.drawEnemyHealthBar(screen, x, e.Y-8, e.W, e.Health, e.MaxHealth)
	}

	goalX := data.GoalX - l.cameraX
	ebitenutil.DrawRect(screen, goalX, data.GroundY-100, 8, 100, color.RGBA{R: 0x14, G: 0x5a, B: 0x26, A: 0xff})
	ebitenutil.DrawRect(screen, goalX+8, data.GroundY-100, 20, 12, color.RGBA{R: 0xe9, G: 0xdf, B: 0x58, A: 0xff})

	pc := color.RGBA{R: 0x2f, G: 0x3f, B: 0xd1, A: 0xff}
	if l.player.Attacking {
		pc = color.RGBA{R: 0x1e, G: 0x2f, B: 0x8f, A: 0xff}
	}
	ebitenutil.DrawRect(screen, l.player.X-l.cameraX, l.player.Y, l.player.W, l.player.H, pc)
	if atkRect, ok := l.player.AttackRect(); ok {
		ebitenutil.DrawRect(screen, atkRect.X-l.cameraX, atkRect.Y, atkRect.W, atkRect.H, color.RGBA{R: 0xe8, G: 0xe8, B: 0xe8, A: 0xff})
	}
	l.drawHUD(g, screen, data)
}

func (l *LevelScene) drawBackground(screen *ebiten.Image, data *level.LevelData) {
	ebitenutil.DrawRect(screen, 0, 0, float64(config.InternalWidth), float64(config.InternalHeight), color.RGBA{R: 0x70, G: 0xc7, B: 0xf5, A: 0xff})
	ebitenutil.DrawRect(screen, 0, 0, float64(config.InternalWidth), 120, color.RGBA{R: 0x87, G: 0xd5, B: 0xf8, A: 0xff})
	ebitenutil.DrawRect(screen, 0, 120, float64(config.InternalWidth), 90, color.RGBA{R: 0x74, G: 0xc7, B: 0xea, A: 0xff})

	baseHills := []float64{40, 65, 52, 74, 58, 68}
	offset := -l.cameraX * 0.15
	for i := 0; i < len(baseHills); i++ {
		x := float64(i)*130 + offset
		for x > float64(config.InternalWidth)+140 {
			x -= 900
		}
		for x < -160 {
			x += 900
		}
		ebitenutil.DrawRect(screen, x, data.GroundY-baseHills[i]-16, 150, baseHills[i], color.RGBA{R: 0x73, G: 0xb5, B: 0x5e, A: 0xff})
	}
}

func (l *LevelScene) drawEnemyHealthBar(screen *ebiten.Image, x, y, w float64, health, maxHealth int) {
	if maxHealth <= 1 {
		return
	}
	barH := 5.0
	ebitenutil.DrawRect(screen, x, y, w, barH, color.RGBA{R: 0x2a, G: 0x2a, B: 0x2a, A: 0xff})
	fillW := w * clamp01(float64(health)/float64(maxHealth))
	ebitenutil.DrawRect(screen, x, y, fillW, barH, color.RGBA{R: 0xdb, G: 0x45, B: 0x39, A: 0xff})
}

func (l *LevelScene) drawHUD(g *game.Game, screen *ebiten.Image, data *level.LevelData) {
	panelH := 58.0
	ebitenutil.DrawRect(screen, 0, 0, float64(config.InternalWidth), panelH, color.RGBA{R: 0x1d, G: 0x25, B: 0x1d, A: 0xe8})
	ebitenutil.DrawRect(screen, 0, panelH-2, float64(config.InternalWidth), 2, color.RGBA{R: 0x8f, G: 0xb3, B: 0x7d, A: 0xff})

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Field %d: %s", l.levelIndex+1, data.Name), 12, 8)

	for i := 0; i < g.State.MaxHealth; i++ {
		x := 12.0 + float64(i)*24
		y := 30.0
		ebitenutil.DrawRect(screen, x, y, 18, 20, color.RGBA{R: 0x38, G: 0x2d, B: 0x1d, A: 0xff})
		if i < g.State.Health {
			ebitenutil.DrawRect(screen, x+2, y+2, 14, 16, color.RGBA{R: 0xf2, G: 0xc2, B: 0x3f, A: 0xff})
		} else {
			ebitenutil.DrawRect(screen, x+2, y+2, 14, 16, color.RGBA{R: 0x6f, G: 0x64, B: 0x55, A: 0xff})
		}
	}

	livesX := 142.0
	for i := 0; i < config.StartLives; i++ {
		x := livesX + float64(i)*18
		c := color.RGBA{R: 0x44, G: 0x4b, B: 0x44, A: 0xff}
		if i < g.State.Lives {
			c = color.RGBA{R: 0x85, G: 0xe0, B: 0x8b, A: 0xff}
		}
		ebitenutil.DrawRect(screen, x, 34, 12, 12, c)
	}
	ebitenutil.DebugPrintAt(screen, "Lives", 142, 22)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Enemies %d", l.aliveEnemies()), 220, 30)

	progress := clamp01((l.player.X + l.player.W) / data.GoalX)
	barX := 332.0
	barY := 33.0
	barW := 280.0
	barH := 14.0
	ebitenutil.DrawRect(screen, barX, barY, barW, barH, color.RGBA{R: 0x3d, G: 0x43, B: 0x3d, A: 0xff})
	ebitenutil.DrawRect(screen, barX+1, barY+1, (barW-2)*progress, barH-2, color.RGBA{R: 0xa6, G: 0xd9, B: 0x75, A: 0xff})
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Progress %d%%", int(progress*100)), 332, 20)

	audioStatus := "BGM ON"
	if g.State.Muted {
		audioStatus = "MUTED"
	}
	ebitenutil.DebugPrintAt(screen, audioStatus, 564, 8)
}

func (l *LevelScene) aliveEnemies() int {
	count := 0
	for _, e := range l.enemies {
		if e.Alive() {
			count++
		}
	}
	return count
}

func onScreen(x, w float64) bool {
	return x+w >= 0 && x <= float64(config.InternalWidth)
}

func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}
