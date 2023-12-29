package coordinates

type Dimension struct {
	Height int `json:"Height"`
	Width  int `json:"Width"`
}

type Vector struct {
	X int `json:"X"`
	Y int `json:"Y"`
	Z int `json:"Z"`
}

type Vector64 struct {
	X float64
	Y float64
	Z float64
}

type Direction struct {
	Straight bool
	Right    bool
	Forward  bool
}

func (v Vector) GetManhattanDist(target Vector) int {
	deltax := target.X - v.X
	deltay := target.Y - v.Y
	if deltax < 0 {
		deltax = -deltax
	}
	if deltay < 0 {
		deltay = -deltay
	}
	return deltax + deltay
}
