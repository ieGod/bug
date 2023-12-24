package images

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type BugAssetName string

const (
	IMGTITLE  BugAssetName = "TITLE"
	IMGBUG    BugAssetName = "BUG"
	IMGSHADOW BugAssetName = "SHADOW"
)

var (
	BugImages map[BugAssetName]*ebiten.Image

	//go:embed bugtitle.png
	bugtitle_img []byte

	//go:embed bug.png
	bug_img []byte

	//go:embed shadow.png
	shadow_img []byte
)

func init() {
	BugImages = make(map[BugAssetName]*ebiten.Image)
}

func LoadImageAssets() {
	BugImages[IMGTITLE] = LoadImagesFatal(bugtitle_img)
	BugImages[IMGBUG] = LoadImagesFatal(bug_img)
	BugImages[IMGSHADOW] = LoadImagesFatal(shadow_img)

}

func LoadImagesFatal(b []byte) *ebiten.Image {
	img, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		log.Fatal(err)
	}
	return ebiten.NewImageFromImage(img)
}
