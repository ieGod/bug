package scenes

import (
	"bug/constants"
	"bug/coordinates"
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
	loaded     bool
	complete   bool
	cycle      int
	tick       int

	//scene elements

}

func NewIntroScene(dimensions coordinates.Dimension) *IntroScene {
	var scene *IntroScene = &IntroScene{
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
	op.GeoM.Translate(float64(scene.dimensions.Width/2-380), float64(scene.dimensions.Height/2-200))
	img.DrawImage(images.BugImages[images.IMGTITLE], op)
	text.Draw(img, constants.Strings.PressEnter, fonts.Bugger.Standard, scene.dimensions.Width/2-200, scene.dimensions.Height/2+200, color.White)

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
	images.LoadImageAssets()
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
