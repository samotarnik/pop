package main

import (
	"fmt"
	"log"

	"pop/assets"
	"pop/internal/config"
	"pop/internal/game"
	"pop/internal/level"
	"pop/internal/scene"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	levels, err := loadLevels()
	if err != nil {
		log.Fatal(err)
	}
	cards, err := level.LoadStoryCards(assets.FS, "story/cards.json")
	if err != nil {
		log.Fatal(err)
	}

	g := game.New(levels, cards)
	g.SetScene(scene.NewMenuScene())

	ebiten.SetWindowSize(config.InternalWidth*config.WindowScale, config.InternalHeight*config.WindowScale)
	ebiten.SetWindowTitle("Prince of Prekmurje")
	ebiten.SetTPS(config.TPS)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func loadLevels() ([]*level.LevelData, error) {
	paths := []string{
		"levels/level1.json",
		"levels/level2.json",
		"levels/level3.json",
	}
	out := make([]*level.LevelData, 0, len(paths))
	for _, p := range paths {
		lvl, err := level.LoadLevel(assets.FS, p)
		if err != nil {
			return nil, fmt.Errorf("load campaign levels: %w", err)
		}
		out = append(out, lvl)
	}
	return out, nil
}
