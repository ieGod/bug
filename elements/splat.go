package elements

import (
	"bug/constants"
	"bug/coordinates"
	"bug/resources/images"

	"github.com/hajimehoshi/ebiten/v2"
)

type Splat struct {
	Sprite         *ebiten.Image
	location       coordinates.Vector
	targetlocation coordinates.Vector
	cycle          int
}

func NewSplat() *Splat {
	return &Splat{
		Sprite: ebiten.NewImage(constants.SwatWidth, constants.SwatHeight),
		cycle:  0,
	}
}

func (elem *Splat) GetLocation() coordinates.Vector {
	return elem.location
}

func (elem *Splat) GetTargetLocation() coordinates.Vector {
	return elem.targetlocation
}

func (elem *Splat) SetLocation(pos coordinates.Vector) {
	elem.location = pos
}

func (elem *Splat) SetTargetLocation(pos coordinates.Vector) {
	elem.targetlocation = pos
}

func (elem *Splat) Animate() {
	elem.Sprite.Clear()
	elem.Sprite.DrawImage(images.BugImages[images.IMGSPLAT], nil)
	elem.cycle++
}
