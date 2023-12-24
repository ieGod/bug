package coordinates

type Dimension struct {
	Height int
	Width  int
}

type Vector struct {
	X int
	Y int
	Z int
}

type Direction struct {
	Straight bool
	Right    bool
	Forward  bool
}
