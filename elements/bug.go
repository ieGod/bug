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
	direction       coordinates.Direction
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
		direction:       coordinates.Direction{Straight: true, Right: true, Forward: true},
	}
}

func (elem *Bug) Animate() {
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

func (elem *Bug) SetLocation(pos coordinates.Vector) {
	elem.location = pos
}

func (elem *Bug) SetTargetLocation(pos coordinates.Vector) {
	elem.targetlocation = pos
}

func (elem *Bug) ForceAllPositionsGrid(pos coordinates.Vector) {
	elem.SetLocation(pos)
	elem.SetTargetLocation(pos)
}

func (elem *Bug) GetLocation() coordinates.Vector {
	return elem.location
}

func (elem *Bug) GetTargetLocation() coordinates.Vector {
	return elem.targetlocation
}

func (elem *Bug) SetTargetFrameCycles(cycles int) {
	elem.animationframes = cycles
}

func (elem *Bug) SetRole(action definitions.BugAction, direction coordinates.Direction) {
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

func (elem *Bug) GetAction() definitions.BugAction {
	return elem.action
}
