package scenes

import (
	"bug/constants"
	"bug/coordinates"
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
	return constants.Strings.SwatName
}

func (scene *SwatScene) IsComplete() bool {
	return scene.complete
}

func (scene *SwatScene) handleInputs() {
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		scene.complete = true
	}

	var newpos coordinates.Vector = scene.bug.GetLocation()

	if ebiten.IsKeyPressed(ebiten.KeyW) ||
		ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		newpos.Y = newpos.Y - constants.BugSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) ||
		ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		newpos.X = newpos.X - constants.BugSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) ||
		ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		newpos.Y = newpos.Y + constants.BugSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) ||
		ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		newpos.X = newpos.X + constants.BugSpeed
	}
	scene.bug.SetLocation(newpos)

}

func (scene *SwatScene) RenderBug(img *ebiten.Image) {
	offset := scene.bug.GetLocation()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(offset.X), float64(offset.Y))
	img.DrawImage(scene.bug.Sprite, op)
}
