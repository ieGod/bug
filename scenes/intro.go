package scenes

import (
	"bug/constants"
	"bug/coordinates"
	"bug/definitions"
	"bug/elements"
	"bug/fonts"
	"bug/fx"
	"bug/resources/images"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type IntroScene struct {
	dimensions coordinates.Dimension
	bug        *elements.Bug
	loaded     bool
	complete   bool
	cycle      int
	tick       int

	//scene elements

}

func NewIntroScene(dimensions coordinates.Dimension) *IntroScene {
	var scene *IntroScene = &IntroScene{
		bug:        elements.NewBug(),
		cycle:      0,
		tick:       0,
		loaded:     false,
		complete:   false,
		dimensions: dimensions,
	}
	return scene
}

func (scene *IntroScene) Draw(img *ebiten.Image) {
	img.Clear()
	img.Fill(fx.HexToRGBA(0x8c0002, 0xff))
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(5, 5)
	op.GeoM.Translate(float64(scene.dimensions.Width/2-32*5/2), float64(scene.dimensions.Height/2-200))

	//img.DrawImage(images.BugImages[images.IMGTITLE], op)
	text.Draw(img, "B", fonts.Bugger.ArcadeHuge, 500, 250, color.White)
	img.DrawImage(scene.bug.Sprite, op)
	text.Draw(img, "G", fonts.Bugger.ArcadeHuge, scene.dimensions.Width/2+32*5/2, 250, color.White)

	text.Draw(img, constants.Strings.PressEnter, fonts.Bugger.Arcade, scene.dimensions.Width/2-100, scene.dimensions.Height/2+200, color.White)

	scene.cycle++
}

func (scene *IntroScene) Update() error {

	scene.handleInput()

	if scene.tick%7 == 0 {
		scene.bug.Animate()
	}
	scene.tick++
	return nil
}

func (scene *IntroScene) IsLoaded() bool {
	return scene.loaded
}

func (scene *IntroScene) Load() {
	images.LoadImageAssets()
	dir := coordinates.Direction{
		Straight: true,
		Forward:  true,
		Right:    false,
	}
	scene.bug.SetRole(definitions.BugActionGlitch, dir)
	scene.loaded = true
}

func (scene *IntroScene) GetName() string {
	return constants.Strings.IntroName
}

func (scene *IntroScene) IsComplete() bool {
	return scene.complete
}

func (scene *IntroScene) handleInput() {
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) ||
		inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		scene.complete = true
	}
}
