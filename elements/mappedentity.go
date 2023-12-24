package elements

import "bug/coordinates"

type MappedEntity interface {
	GetLocation() coordinates.Vector
	GetTargetLocation() coordinates.Vector
	SetCoordinates(coordinates.Vector)
	SetTargetCoordinates(coordinates.Vector)
}
