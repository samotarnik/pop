package scene

import (
	"pop/internal/game"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type GameOverScene struct{}

func NewGameOverScene() *GameOverScene {
	return &GameOverScene{}
}

func (s *GameOverScene) Update(g *game.Game) error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.SetScene(NewMenuScene())
	}
	return nil
}

func (s *GameOverScene) Draw(_ *game.Game, screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, "Game Over\n\nPress Enter to return to menu", 120, 180)
}
