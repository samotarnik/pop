package scene

import (
	"pop/internal/game"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type MenuScene struct{}

func NewMenuScene() *MenuScene {
	return &MenuScene{}
}

func (m *MenuScene) Update(g *game.Game) error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.State.ResetRun()
		g.SetScene(NewStoryScene("intro", func(_ *game.Game) game.Scene {
			return NewLevelScene(0)
		}))
	}
	return nil
}

func (m *MenuScene) Draw(_ *game.Game, screen *ebiten.Image) {
	text := "Prince of Prekmurje\n\nPress Enter to Start\n\nControls:\nJ/L or Left/Right: Move\nI or Up: Jump\nK or Down: Duck\nSpace: Pitchfork attack (while standing)\nR: Restart level\nEsc: Pause\nM: Mute toggle"
	ebitenutil.DebugPrintAt(screen, text, 40, 50)
}
