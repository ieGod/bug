package elements

import (
	"bug/constants"
	"bug/coordinates"
	"bug/resources/images"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Bug struct {
	Sprite          *ebiten.Image
	location        coordinates.Vector
	animationframes int
	cycle           int
}

func NewBug() *Bug {
	return &Bug{
		Sprite:          ebiten.NewImage(constants.BugWidth, constants.BugHeight),
		location:        coordinates.Vector{},
		animationframes: constants.AnimationFrames,
	}
}

func (bug *Bug) Animate() {
	bug.Sprite.Clear()

	//todo: cycle offset depends on target number of frames
	ox := (bug.cycle % bug.animationframes) * constants.BugWidth
	bug.Sprite.DrawImage(images.BugImages[images.IMGBUG].SubImage(image.Rect(ox, 0, ox+constants.BugWidth, constants.BugHeight)).(*ebiten.Image), nil)
	bug.cycle++
}

func (bug *Bug) SetTargetFrameCycles(cycles int) {
	bug.animationframes = cycles
}
