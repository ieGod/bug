package scenes

import (
	"bug/constants"
	"bug/coordinates"
	"bug/definitions"
	"bug/elements"
	"bug/resources/images"
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

	scene.RenderSurface(img)
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
		scene.bug.SetRole(definitions.BugActionForwardRun)
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) ||
		ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		newpos.X = newpos.X + constants.BugSpeed
	}
	scene.bug.SetLocation(newpos)

	if inpututil.IsKeyJustReleased(ebiten.KeyW) ||
		inpututil.IsKeyJustReleased(ebiten.KeyA) ||
		inpututil.IsKeyJustReleased(ebiten.KeyS) ||
		inpututil.IsKeyJustReleased(ebiten.KeyD) ||
		inpututil.IsKeyJustReleased(ebiten.KeyArrowUp) ||
		inpututil.IsKeyJustReleased(ebiten.KeyArrowLeft) ||
		inpututil.IsKeyJustReleased(ebiten.KeyArrowDown) ||
		inpututil.IsKeyJustReleased(ebiten.KeyArrowRight) {
		scene.bug.SetRole(definitions.BugActionIdle)
	}

}

func (scene *SwatScene) RenderBug(img *ebiten.Image) {
	offset := scene.bug.GetLocation()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(3, 3)
	op.GeoM.Translate(float64(offset.X), float64(offset.Y))
	img.DrawImage(scene.bug.Sprite, op)
}

func (scene *SwatScene) RenderSurface(img *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	tableimg := images.BugImages[images.IMGTABLE]
	sx := float64(scene.dimensions.Width) / float64(tableimg.Bounds().Dx())
	sy := float64(scene.dimensions.Height) / float64(tableimg.Bounds().Dy())
	op.GeoM.Scale(sx, sy)
	img.DrawImage(tableimg, op)
}
