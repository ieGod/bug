package scenes

import "github.com/hajimehoshi/ebiten/v2"

type BugScene interface {
	Draw(img *ebiten.Image)
	Update() error
	Load()
	IsLoaded() bool
	IsComplete() bool
	GetName() string //textual description of scene implementation
}
