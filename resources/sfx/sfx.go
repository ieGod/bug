package sfx

import (
	_ "embed"
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
)

func init() {
	//BugImages = make(map[BugAssetName]*ebiten.Image)
}
