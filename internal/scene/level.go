package scene

import (
	"fmt"
	"image/color"

	"pop/internal/config"
	"pop/internal/entity"
	"pop/internal/game"
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
			dead := g.State.Damage(1)
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

	ebitenutil.DrawRect(screen, 0, data.GroundY, float64(config.InternalWidth), float64(config.InternalHeight)-data.GroundY, color.RGBA{R: 0x58, G: 0x9b, B: 0x42, A: 0xff})

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
	}

	goalX := data.GoalX - l.cameraX
	ebitenutil.DrawRect(screen, goalX, data.GroundY-86, 8, 86, color.RGBA{R: 0x14, G: 0x5a, B: 0x26, A: 0xff})

	pc := color.RGBA{R: 0x2f, G: 0x3f, B: 0xd1, A: 0xff}
	if l.player.Attacking {
		pc = color.RGBA{R: 0x1e, G: 0x2f, B: 0x8f, A: 0xff}
	}
	ebitenutil.DrawRect(screen, l.player.X-l.cameraX, l.player.Y, l.player.W, l.player.H, pc)
	if atkRect, ok := l.player.AttackRect(); ok {
		ebitenutil.DrawRect(screen, atkRect.X-l.cameraX, atkRect.Y, atkRect.W, atkRect.H, color.RGBA{R: 0xe8, G: 0xe8, B: 0xe8, A: 0xff})
	}

	ebitenutil.DebugPrintAt(
		screen,
		fmt.Sprintf("Level %d  Health %d/%d  Lives %d  Mute %v", l.levelIndex+1, g.State.Health, g.State.MaxHealth, g.State.Lives, g.State.Muted),
		10,
		8,
	)
}

func onScreen(x, w float64) bool {
	return x+w >= 0 && x <= float64(config.InternalWidth)
}
