package sfx

import (
	_ "embed"
	"io"
)

const (
	SamepleRate = 44100
)

var (
	//BugSfx map[BugAssetName]*ebiten.Image

	//go:embed swat.wav
	SwatWav []byte

	//go:embed caught.wav
	CaughtWav []byte

	//go:embed asphodel.mp3
	AsphodelMp3 []byte

	//go:embed glitch.mp3
	GlitchMp3 []byte
)

type AudioStream interface {
	io.ReadSeeker
	Length() int64
}

func init() {
	//BugImages = make(map[BugAssetName]*ebiten.Image)
}
