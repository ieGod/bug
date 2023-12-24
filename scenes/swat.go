package scenes

import (
	"bug/constants"
	"bug/coordinates"
	"bug/definitions"
	"bug/elements"
	"bug/fonts"
	"bug/resources/images"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type SwatScene struct {
	dimensions coordinates.Dimension
	loaded     bool
	complete   bool
	cycle      int
	tick       int

	//scene elements
	splatmask    *ebiten.Image
	collidermask *ebiten.Image
	bug          *elements.Bug
	splat        *elements.Splat
	bugcollision bool
}

func NewSwatScene(dimensions coordinates.Dimension) *SwatScene {
	var scene *SwatScene = &SwatScene{
		bug:          elements.NewBug(),
		splat:        elements.NewSplat(),
		splatmask:    ebiten.NewImage(constants.SplatWidth+96*2, constants.SplatHeight+96*2),
		collidermask: ebiten.NewImage(constants.BugWidth*3, constants.BugHeight*3),
		cycle:        0,
		tick:         0,
		loaded:       false,
		complete:     false,
		dimensions:   dimensions,
		bugcollision: false,
	}
	return scene
}

func (scene *SwatScene) Draw(img *ebiten.Image) {
	img.Clear()
	img.Fill(color.White)

	scene.RenderSurface(img)
	scene.RenderBug(img)
	scene.RenderSplat(img)

	if scene.bugcollision {
		text.Draw(img, constants.Strings.Targeted, fonts.Bugger.Standard, 50, 150, color.White)
	}

	scene.cycle++
}

func (scene *SwatScene) Update() error {
	scene.handleInputs()

	if scene.tick%7 == 0 {
		scene.bug.Animate()
	}

	scene.splat.Animate()

	scene.CheckCollisions()

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
		direction := coordinates.Direction{
			Straight: false,
			Right:    false,
			Forward:  true,
		}
		newpos.X = newpos.X - constants.BugSpeed
		scene.bug.SetRole(definitions.BugActionSideRun, direction)

	}
	if ebiten.IsKeyPressed(ebiten.KeyS) ||
		ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		newpos.Y = newpos.Y + constants.BugSpeed
		direction := coordinates.Direction{
			Straight: true,
			Right:    false,
			Forward:  true,
		}
		scene.bug.SetRole(definitions.BugActionForwardRun, direction)
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) ||
		ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		direction := coordinates.Direction{
			Straight: false,
			Right:    true,
			Forward:  true,
		}
		newpos.X = newpos.X + constants.BugSpeed
		scene.bug.SetRole(definitions.BugActionSideRun, direction)
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
		direction := coordinates.Direction{
			Straight: true,
			Right:    false,
			Forward:  true,
		}
		scene.bug.SetRole(definitions.BugActionIdle, direction)
	}

	var mpos coordinates.Vector
	mpos.X, mpos.Y = ebiten.CursorPosition()
	scene.splat.SetLocation(mpos)

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

func (scene *SwatScene) RenderSplat(img *ebiten.Image) {
	offset := scene.splat.GetLocation()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(constants.SplatWidth)/2.-40, -float64(constants.SplatHeight)/2.)
	//op.GeoM.Scale(2, 2)
	op.GeoM.Translate(float64(offset.X), float64(offset.Y))
	img.DrawImage(scene.splat.Sprite, op)

	img.DrawImage(scene.collidermask, nil)
}

func (scene *SwatScene) CheckCollisions() {

	bugloc := scene.bug.GetLocation()
	splatloc := scene.splat.GetLocation()
	splatloc.X = splatloc.X - constants.SplatWidth/2 - 40
	splatloc.Y = splatloc.Y - constants.SplatHeight/2
	var c bool = false
	//check bounding boxes first, if we're in range we need a more precise check
	if bugloc.X >= splatloc.X-96 && bugloc.X < splatloc.X+constants.SplatWidth &&
		bugloc.Y >= splatloc.Y-96 && bugloc.Y < splatloc.Y+constants.SplatHeight {
		//fmt.Println("collision requires more precise check")

		scene.collidermask.Clear()
		ox := -float64(bugloc.X - splatloc.X)
		oy := -float64(bugloc.Y - splatloc.Y)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(3, 3)
		scene.collidermask.DrawImage(scene.bug.Sprite, op)

		op.GeoM.Reset()
		op.GeoM.Translate(96, 96) //this provides some margin for the mask that matches the bounds of the triple scaled 32x32 bug sprite
		scene.splatmask.DrawImage(scene.splat.Sprite, op)

		op.GeoM.Reset()
		op.Blend = ebiten.BlendSourceIn
		op.GeoM.Translate(ox-96, oy-96)
		scene.collidermask.DrawImage(scene.splatmask, op)

		var pixels []byte = make([]byte, constants.BugWidth*3*constants.BugWidth*3*4)
		scene.collidermask.ReadPixels(pixels)
		for i := 0; i < len(pixels); i = i + 4 {
			if pixels[i+3] != 0 {
				//fmt.Println("pixel collision")
				c = true
			}
		}

	}
	scene.bugcollision = c

}
