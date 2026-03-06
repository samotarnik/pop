package sprite

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/png"
	"io/fs"
	"path"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

type Clip struct {
	Frames []int
	FPS    float64
	Loop   bool
	Events map[string]int
}

type Sheet struct {
	Image   *ebiten.Image
	FrameW  int
	FrameH  int
	AnchorX float64
	AnchorY float64
	Cols    int
	Clips   map[string]Clip
}

type manifest struct {
	Image     string                  `json:"image"`
	FrameSize frameSize               `json:"frameSize"`
	Anchor    anchor                  `json:"anchor"`
	Anims     map[string]animManifest `json:"animations"`
}

type frameSize struct {
	W int `json:"w"`
	H int `json:"h"`
}

type anchor struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type animManifest struct {
	Frames []int          `json:"frames"`
	FPS    float64        `json:"fps"`
	Loop   bool           `json:"loop"`
	Events map[string]int `json:"events,omitempty"`
}

func LoadSheet(fsys fs.FS, metaPath string) (*Sheet, error) {
	metaBytes, err := fs.ReadFile(fsys, metaPath)
	if err != nil {
		return nil, fmt.Errorf("read sprite meta %s: %w", metaPath, err)
	}

	var m manifest
	if err := json.Unmarshal(metaBytes, &m); err != nil {
		return nil, fmt.Errorf("decode sprite meta %s: %w", metaPath, err)
	}
	if m.FrameSize.W <= 0 || m.FrameSize.H <= 0 {
		return nil, fmt.Errorf("invalid frame size in %s", metaPath)
	}
	if len(m.Anims) == 0 {
		return nil, fmt.Errorf("no animations in %s", metaPath)
	}

	imgPath := resolveImagePath(metaPath, m.Image)
	imgFile, err := fsys.Open(imgPath)
	if err != nil {
		return nil, fmt.Errorf("open sprite image %s: %w", imgPath, err)
	}
	defer imgFile.Close()

	decoded, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, fmt.Errorf("decode sprite image %s: %w", imgPath, err)
	}

	ebImg := ebiten.NewImageFromImage(decoded)
	b := ebImg.Bounds()
	if b.Dx() < m.FrameSize.W || b.Dy() < m.FrameSize.H {
		return nil, fmt.Errorf("sprite image %s is smaller than one frame", imgPath)
	}

	cols := b.Dx() / m.FrameSize.W
	if cols <= 0 {
		return nil, fmt.Errorf("invalid sprite columns for %s", imgPath)
	}

	clips := make(map[string]Clip, len(m.Anims))
	for name, anim := range m.Anims {
		if len(anim.Frames) == 0 {
			continue
		}
		fps := anim.FPS
		if fps <= 0 {
			fps = 8
		}
		clips[name] = Clip{Frames: anim.Frames, FPS: fps, Loop: anim.Loop, Events: anim.Events}
	}
	if len(clips) == 0 {
		return nil, fmt.Errorf("no valid animations in %s", metaPath)
	}

	return &Sheet{
		Image:   ebImg,
		FrameW:  m.FrameSize.W,
		FrameH:  m.FrameSize.H,
		AnchorX: float64(m.Anchor.X),
		AnchorY: float64(m.Anchor.Y),
		Cols:    cols,
		Clips:   clips,
	}, nil
}

func resolveImagePath(metaPath, imagePath string) string {
	if strings.Contains(imagePath, "/") {
		return imagePath
	}
	return path.Join(path.Dir(metaPath), imagePath)
}

func (s *Sheet) frameRect(frame int) image.Rectangle {
	col := frame % s.Cols
	row := frame / s.Cols
	return image.Rect(col*s.FrameW, row*s.FrameH, (col+1)*s.FrameW, (row+1)*s.FrameH)
}

func (s *Sheet) drawFrame(screen *ebiten.Image, frame int, footX, footY, cameraX float64, facing int) {
	if s == nil || s.Image == nil {
		return
	}
	src := s.Image.SubImage(s.frameRect(frame)).(*ebiten.Image)
	drawX := footX - s.AnchorX - cameraX
	drawY := footY - s.AnchorY
	op := &ebiten.DrawImageOptions{}
	if facing >= 0 {
		op.GeoM.Translate(drawX, drawY)
	} else {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(drawX+float64(s.FrameW), drawY)
	}
	screen.DrawImage(src, op)
}
