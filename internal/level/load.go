package level

import (
	"encoding/json"
	"fmt"
	"io/fs"
)

func LoadLevel(fsys fs.FS, path string) (*LevelData, error) {
	b, err := fs.ReadFile(fsys, path)
	if err != nil {
		return nil, fmt.Errorf("read level %s: %w", path, err)
	}
	var out LevelData
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, fmt.Errorf("decode level %s: %w", path, err)
	}
	if out.ID == "" || out.Width <= 0 || out.GroundY <= 0 || out.GoalX <= 0 {
		return nil, fmt.Errorf("invalid level %s", path)
	}
	return &out, nil
}

func LoadStoryCards(fsys fs.FS, path string) (map[string]StoryCard, error) {
	b, err := fs.ReadFile(fsys, path)
	if err != nil {
		return nil, fmt.Errorf("read story cards %s: %w", path, err)
	}
	out := map[string]StoryCard{}
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, fmt.Errorf("decode story cards %s: %w", path, err)
	}
	return out, nil
}
