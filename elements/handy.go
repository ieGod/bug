package elements

import (
	"bug/constants"
	"bug/coordinates"
	"bug/definitions"
	"bug/resources/images"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Handy struct {
	Sprite          *ebiten.Image
	action          definitions.BugAction
	direction       coordinates.Direction
	location        coordinates.Vector
	targetlocation  coordinates.Vector
	animationframes int
	cycle           int
}

func NewHandy() *Handy {
	return &Handy{
		Sprite:          ebiten.NewImage(constants.BugWidth, constants.BugHeight),
		action:          definitions.BugActionIdle,
		location:        coordinates.Vector{},
		targetlocation:  coordinates.Vector{},
		animationframes: constants.AnimationFrames,
		direction:       coordinates.Direction{Straight: true, Right: true, Forward: true},
	}
}

func (elem *Handy) Animate() {
	elem.Sprite.Clear()

	ox := (elem.cycle % elem.animationframes) * constants.BugWidth
	oy := int(elem.action) * constants.BugHeight

	op := &ebiten.DrawImageOptions{}
	if !elem.direction.Straight && !elem.direction.Right {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(constants.BugWidth, 0)
	}
	elem.Sprite.DrawImage(images.BugImages[images.IMGSHADOW], op)
	elem.Sprite.DrawImage(images.BugImages[images.IMGBUG].SubImage(image.Rect(ox, oy, ox+constants.BugWidth, oy+constants.BugHeight)).(*ebiten.Image), op)
	elem.cycle++
}

func (elem *Handy) SetLocation(pos coordinates.Vector) {
	elem.location = pos
}

func (elem *Handy) SetTargetLocation(pos coordinates.Vector) {
	elem.targetlocation = pos
}

func (elem *Handy) GetLocation() coordinates.Vector {
	return elem.location
}

func (elem *Handy) GetTargetLocation() coordinates.Vector {
	return elem.targetlocation
}

func (elem *Handy) SetTargetFrameCycles(cycles int) {
	elem.animationframes = cycles
}

func (elem *Handy) SetRole(action definitions.BugAction, direction coordinates.Direction) {
	var targetframes int = constants.AnimationFrames
	elem.action = action
	elem.direction = direction
	switch elem.action {
	case definitions.BugActionIdle:
		targetframes = constants.BugIdleFramecount
	case definitions.BugActionForwardRun,
		definitions.BugActionSideRun,
		definitions.BugActionReverseRun,
		definitions.BugActionGlitch:
		targetframes = constants.BugForwardRunFramecount
	}
	elem.SetTargetFrameCycles(targetframes)
}

func (elem *Handy) GetAction() definitions.BugAction {
	return elem.action
}
