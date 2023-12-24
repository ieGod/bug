package elements

import (
	"bug/constants"
	"bug/coordinates"
	"bug/definitions"
	"bug/resources/images"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Bug struct {
	Sprite          *ebiten.Image
	action          definitions.BugAction
	location        coordinates.Vector
	targetlocation  coordinates.Vector
	animationframes int
	cycle           int
}

func NewBug() *Bug {
	return &Bug{
		Sprite:          ebiten.NewImage(constants.BugWidth, constants.BugHeight),
		action:          definitions.BugActionIdle,
		location:        coordinates.Vector{},
		targetlocation:  coordinates.Vector{},
		animationframes: constants.AnimationFrames,
	}
}

func (bug *Bug) Animate() {
	bug.Sprite.Clear()

	//todo: cycle offset depends on target number of frames
	ox := (bug.cycle % bug.animationframes) * constants.BugWidth
	oy := int(bug.action) * constants.BugHeight
	bug.Sprite.DrawImage(images.BugImages[images.IMGSHADOW], nil)
	bug.Sprite.DrawImage(images.BugImages[images.IMGBUG].SubImage(image.Rect(ox, oy, ox+constants.BugWidth, oy+constants.BugHeight)).(*ebiten.Image), nil)
	bug.cycle++
}

func (bug *Bug) SetLocation(pos coordinates.Vector) {
	bug.location = pos
}

func (bug *Bug) SetTargetLocation(pos coordinates.Vector) {
	bug.targetlocation = pos
}

func (bug *Bug) GetLocation() coordinates.Vector {
	return bug.location
}

func (bug *Bug) GetTargetLocation() coordinates.Vector {
	return bug.targetlocation
}

func (bug *Bug) SetTargetFrameCycles(cycles int) {
	bug.animationframes = cycles
}

func (bug *Bug) SetRole(action definitions.BugAction) {
	var targetframes int
	bug.action = action
	switch bug.action {
	case definitions.BugActionIdle:
		targetframes = constants.BugIdleFramecount
	case definitions.BugActionForwardRun:
		targetframes = constants.BugForwardRunFramecount
	}
	bug.SetTargetFrameCycles(targetframes)
}
