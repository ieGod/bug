package elements

import "bug/coordinates"

type MappedEntity interface {
	GetLocation() coordinates.Vector
	GetTargetLocation() coordinates.Vector
	SetLocation(coordinates.Vector)
	SetTargetLocation(coordinates.Vector)
}
