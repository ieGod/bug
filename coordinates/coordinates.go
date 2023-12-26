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

type Direction struct {
	Straight bool
	Right    bool
	Forward  bool
}
