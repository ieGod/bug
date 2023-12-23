package scenes

import (
	"bug/defaults"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type IntroScene struct {
	loaded   bool
	complete bool
	cycle    int
	tick     int
}

func NewIntroScene() *IntroScene {
	var scene *IntroScene = &IntroScene{
		cycle:    0,
		tick:     0,
		loaded:   false,
		complete: false,
	}
	return scene
}

func (scene *IntroScene) Draw(img *ebiten.Image) {
	scene.cycle++
}

func (scene *IntroScene) Update() error {

	scene.handleInput()
	scene.tick++
	return nil
}

func (scene *IntroScene) IsLoaded() bool {
	return scene.loaded
}

func (scene *IntroScene) Load() {

}

func (scene *IntroScene) GetName() string {
	return defaults.Strings.IntroName
}

func (scene *IntroScene) IsComplete() bool {
	return scene.complete
}

func (scene *IntroScene) handleInput() {
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		scene.complete = true
	}
}
