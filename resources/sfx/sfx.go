package sfx

import (
	_ "embed"
)

const ()

var (
	//BugSfx map[BugAssetName]*ebiten.Image

	//go:embed caught.wav
	Caught_mp3 []byte
)

func init() {
	//BugImages = make(map[BugAssetName]*ebiten.Image)
}
