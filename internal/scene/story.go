package scene

import (
	"fmt"

	"pop/internal/game"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type StoryScene struct {
	cardID string
	next   func(*game.Game) game.Scene
}

func NewStoryScene(cardID string, next func(*game.Game) game.Scene) *StoryScene {
	return &StoryScene{cardID: cardID, next: next}
}

func (s *StoryScene) Update(g *game.Game) error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.SetScene(s.next(g))
	}
	return nil
}

func (s *StoryScene) Draw(g *game.Game, screen *ebiten.Image) {
	card, ok := g.Cards[s.cardID]
	if !ok {
		ebitenutil.DebugPrintAt(screen, "Missing story card: "+s.cardID+"\nPress Enter", 40, 60)
		return
	}
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%s\n\n%s\n\nPress Enter", card.Title, card.Body), 40, 60)
}
