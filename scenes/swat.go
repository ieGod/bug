package scenes

import (
	"bug/coordinates"
	"bug/defaults"
	"bug/elements"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type SwatScene struct {
	dimensions coordinates.Dimension
	loaded     bool
	complete   bool
	cycle      int
	tick       int

	//scene elements
	bug *elements.Bug
}

func NewSwatScene(dimensions coordinates.Dimension) *SwatScene {
	var scene *SwatScene = &SwatScene{
		bug:        elements.NewBug(),
		cycle:      0,
		tick:       0,
		loaded:     false,
		complete:   false,
		dimensions: dimensions,
	}
	return scene
}

func (scene *SwatScene) Draw(img *ebiten.Image) {
	img.Clear()
	img.Fill(color.White)

	scene.RenderBug(img)
	scene.cycle++
}

func (scene *SwatScene) Update() error {
	scene.handleInputs()

	if scene.tick%7 == 0 {
		scene.bug.Animate()
	}

	scene.tick++
	return nil
}

func (scene *SwatScene) IsLoaded() bool {
	return scene.loaded
}

func (scene *SwatScene) Load() {
	scene.loaded = true
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

func (scene *SwatScene) RenderBug(img *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(100, 100)
	img.DrawImage(scene.bug.Sprite, op)
}
