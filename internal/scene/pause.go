package scene

import (
	"image/color"

	"pop/internal/config"
	"pop/internal/game"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type PauseScene struct {
	resume game.Scene
}

func NewPauseScene(resume game.Scene) *PauseScene {
	return &PauseScene{resume: resume}
}

func (p *PauseScene) Update(g *game.Game) error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.SetScene(p.resume)
	}
	return nil
}

func (p *PauseScene) Draw(g *game.Game, screen *ebiten.Image) {
	p.resume.Draw(g, screen)
	ebitenutil.DrawRect(screen, 0, 0, float64(config.InternalWidth), float64(config.InternalHeight), color.RGBA{A: 150})
	ebitenutil.DebugPrintAt(screen, "Paused\nEnter/Esc: Resume", config.InternalWidth/2-70, config.InternalHeight/2-20)
}
