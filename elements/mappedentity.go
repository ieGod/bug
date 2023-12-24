package elements

import "bug/coordinates"

type MappedEntity interface {
	GetLocation() coordinates.Vector
	SetCoordinates(coordinates.Vector)
}
