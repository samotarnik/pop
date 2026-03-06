package game

import (
	"image/color"

	"pop/internal/audio"
	"pop/internal/config"
	"pop/internal/level"
	"pop/internal/sprite"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Scene interface {
	Update(*Game) error
	Draw(*Game, *ebiten.Image)
}

type Game struct {
	State  *State
	Levels []*level.LevelData
	Cards  map[string]level.StoryCard
	Audio  *audio.Manager
	Prince *sprite.Sheet
	scene  Scene
}

func New(levels []*level.LevelData, cards map[string]level.StoryCard, prince *sprite.Sheet) *Game {
	return &Game{
		State:  NewState(),
		Levels: levels,
		Cards:  cards,
		Audio:  audio.NewManager(),
		Prince: prince,
	}
}

func (g *Game) SetScene(scene Scene) {
	g.scene = scene
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		g.Audio.ToggleMute()
		g.State.Muted = g.Audio.IsMuted()
	}
	if g.scene == nil {
		return nil
	}
	return g.scene.Update(g)
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 0x6a, G: 0xc7, B: 0xf2, A: 0xff})
	if g.scene != nil {
		g.scene.Draw(g, screen)
	}
}

func (g *Game) Layout(_, _ int) (int, int) {
	return config.InternalWidth, config.InternalHeight
}
