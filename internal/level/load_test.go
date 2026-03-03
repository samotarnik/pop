package level

import (
	"strings"
	"testing"
	"testing/fstest"
)

func TestLoadLevelSuccess(t *testing.T) {
	fsys := fstest.MapFS{
		"levels/ok.json": &fstest.MapFile{Data: []byte(`{"id":"l1","width":1000,"height":480,"ground_y":400,"player_spawn_x":10,"goal_x":900}`)},
	}
	lvl, err := LoadLevel(fsys, "levels/ok.json")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if lvl.ID != "l1" {
		t.Fatalf("expected id l1, got %s", lvl.ID)
	}
}

func TestLoadLevelInvalid(t *testing.T) {
	fsys := fstest.MapFS{
		"levels/bad.json": &fstest.MapFile{Data: []byte(`{"id":"","width":0}`)},
	}
	_, err := LoadLevel(fsys, "levels/bad.json")
	if err == nil {
		t.Fatalf("expected validation error")
	}
	if !strings.Contains(err.Error(), "invalid level") {
		t.Fatalf("expected invalid level error, got %v", err)
	}
}
