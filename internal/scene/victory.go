package scene

import (
	"pop/internal/game"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type VictoryScene struct{}

func NewVictoryScene() *VictoryScene {
	return &VictoryScene{}
}

func (s *VictoryScene) Update(g *game.Game) error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.SetScene(NewMenuScene())
	}
	return nil
}

func (s *VictoryScene) Draw(_ *game.Game, screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, "You Won\n\nPress Enter to return to menu", 120, 180)
}
