package elements

import (
	"bug/coordinates"
	"bug/definitions"
	"math"
)

type BugCam struct {
	params definitions.Paramecas
}

func NewBugCam() *BugCam {
	return &BugCam{
		params: definitions.Paramecas{},
	}
}

func (elem *BugCam) GetParams() definitions.Paramecas {
	return elem.params
}

func (elem *BugCam) SetParams(params definitions.Paramecas) {
	elem.params = params
}

func (elem *BugCam) GetLocation() coordinates.Vector {
	p := elem.params.Location
	return coordinates.Vector{X: int(p.X), Y: int(p.Y), Z: int(p.Z)}
}

func (elem *BugCam) GetTargetLocation() coordinates.Vector {
	p := elem.params.TargetLocation
	return coordinates.Vector{X: int(p.X), Y: int(p.Y), Z: int(p.Z)}
}

func (elem *BugCam) SetLocation(loc coordinates.Vector) {
	elem.params.Location.X = float64(loc.X)
	elem.params.Location.Y = float64(loc.Y)
	elem.params.Location.Z = float64(loc.Z)
}

func (elem *BugCam) SetTargetLocation(loc coordinates.Vector) {
	elem.params.TargetLocation.X = float64(loc.X)
	elem.params.TargetLocation.Y = float64(loc.Y)
	elem.params.TargetLocation.Z = float64(loc.Z)
}

func (elem *BugCam) CloseTargets() {

	dx := float64(elem.params.TargetLocation.X - elem.params.Location.X)
	dy := float64(elem.params.TargetLocation.Y - elem.params.Location.Y)

	dist := math.Sqrt(dx*dx + dy*dy)

	if dist > 0.005 {
		angle := math.Atan2(dy, dx)

		xadjust := dist * math.Cos(angle) / 4
		yadjust := dist * math.Sin(angle) / 4

		elem.params.Location.X += xadjust
		elem.params.Location.Y += yadjust
	}

}
