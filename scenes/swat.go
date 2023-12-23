package scenes

import (
	"bug/defaults"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type SwatScene struct {
	loaded   bool
	complete bool
	cycle    int
	tick     int
}

func NewSwatScene() *SwatScene {
	var scene *SwatScene = &SwatScene{
		cycle:    0,
		tick:     0,
		loaded:   false,
		complete: false,
	}
	return scene
}

func (scene *SwatScene) Draw(img *ebiten.Image) {
	scene.cycle++
}

func (scene *SwatScene) Update() error {
	scene.handleInputs()
	scene.tick++
	return nil
}

func (scene *SwatScene) IsLoaded() bool {
	return scene.loaded
}

func (scene *SwatScene) Load() {

}

func (scene *SwatScene) GetName() string {
	return defaults.Strings.SwatName
}

func (scene *SwatScene) IsComplete() bool {
	return scene.complete
}

func (scene *SwatScene) handleInputs() {
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		scene.complete = true
	}
}
