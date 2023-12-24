package elements

import (
	"bug/coordinates"
	"bug/defaults"
	"bug/resources/images"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Bug struct {
	Sprite   *ebiten.Image
	location coordinates.Vector
	cycle    int
}

func NewBug() *Bug {
	return &Bug{
		Sprite:   ebiten.NewImage(defaults.BugWidth, defaults.BugHeight),
		location: coordinates.Vector{},
	}
}

func (bug *Bug) Animate() {
	bug.Sprite.Clear()
	ox := (bug.cycle % 4) * defaults.BugWidth
	bug.Sprite.DrawImage(images.BugImages[images.IMGBUG].SubImage(image.Rect(ox, 0, ox+32, 32)).(*ebiten.Image), nil)
	bug.cycle++
}
