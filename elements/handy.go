package elements

import (
	"bug/constants"
	"bug/coordinates"
	"bug/definitions"
	"bug/resources/images"
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Handy struct {
	Sprite          *ebiten.Image
	action          definitions.BugAction
	direction       coordinates.Direction
	location        coordinates.Vector   //grid
	targetlocation  coordinates.Vector   //grid
	loc64           coordinates.Vector64 //scene position
	targetloc64     coordinates.Vector64 //scene position
	waypoints       []coordinates.Vector
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
		loc64:           coordinates.Vector64{},
		targetloc64:     coordinates.Vector64{},
		waypoints:       make([]coordinates.Vector, 1),
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

func (elem *Handy) SetLoc64(loc coordinates.Vector64) {
	elem.loc64 = loc
}

func (elem *Handy) SetTargetLoc64(loc coordinates.Vector64) {
	elem.targetloc64 = loc
}

func (elem *Handy) GetLoc64() coordinates.Vector64 {
	return elem.loc64
}

func (elem *Handy) GetTargetLoc64() coordinates.Vector64 {
	return elem.targetloc64
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

func (elem *Handy) CloseTargets() {
	dx := elem.targetloc64.X - elem.loc64.X
	dy := elem.targetloc64.Y - elem.loc64.Y

	dist := math.Sqrt(dx*dx + dy*dy)

	if dist > 0.5 {
		angle := math.Atan2(dy, dx)

		xadjust := dist * math.Cos(angle) / 16
		yadjust := dist * math.Sin(angle) / 16

		elem.loc64.X += xadjust
		elem.loc64.Y += yadjust
		elem.location = elem.targetlocation
	} else {
		if len(elem.waypoints) > 0 {
			elem.targetlocation = elem.waypoints[0]
			elem.waypoints = elem.waypoints[1:]

			elem.targetloc64.X = float64(elem.targetlocation.X * constants.BugWidth)
			elem.targetloc64.Y = float64(elem.targetlocation.Y * constants.BugHeight)

		}
	}
}

func (elem *Handy) ForceAllPositionsGrid(loc coordinates.Vector) {

	elem.location = loc
	elem.targetlocation = loc

	elem.loc64.X = float64(elem.targetlocation.X * constants.BugWidth)
	elem.loc64.Y = float64(elem.targetlocation.Y * constants.BugHeight)

	elem.targetloc64.X = float64(elem.targetlocation.X * constants.BugWidth)
	elem.targetloc64.Y = float64(elem.targetlocation.Y * constants.BugHeight)

	elem.waypoints[0] = loc
}

func (elem *Handy) GenWaypoints() {
	var wp []coordinates.Vector

	wp = append(wp, coordinates.Vector{X: 4, Y: 2})
	wp = append(wp, coordinates.Vector{X: 4, Y: 3})
	wp = append(wp, coordinates.Vector{X: 3, Y: 3})
	wp = append(wp, coordinates.Vector{X: 3, Y: 4})
	wp = append(wp, coordinates.Vector{X: 4, Y: 4})
	wp = append(wp, coordinates.Vector{X: 5, Y: 4})

	elem.waypoints = elem.waypoints[:0]
	elem.waypoints = append(elem.waypoints, wp...)

	elem.targetloc64.X = float64(elem.waypoints[0].X * constants.BugWidth)
	elem.targetloc64.Y = float64(elem.waypoints[0].Y * constants.BugHeight)

}
