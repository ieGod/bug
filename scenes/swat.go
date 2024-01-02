package scenes

import (
	"bug/constants"
	"bug/coordinates"
	"bug/definitions"
	"bug/elements"
	"bug/fonts"
	"bug/fx"
	"bug/resources/images"
	"bug/resources/sfx"
	"bytes"
	"image/color"
	"io"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type SwatScene struct {
	dimensions coordinates.Dimension
	loaded     bool
	complete   bool
	cycle      int
	tick       int

	//scene elements
	fader        *elements.Fader
	splatmask    *ebiten.Image
	collidermask *ebiten.Image
	bug          *elements.Bug
	swatter      *elements.Splat
	bugcollision bool //does our splat mask cover part of the bug?
	whack        bool //are we now whacking?
	gameover     bool //did we get hit?

	seCh         chan []byte
	audioContext *audio.Context
	audioPlayer  *audio.Player
	seBytes      []byte
	firstflag    bool
	musicstarted bool
}

func NewSwatScene(dimensions coordinates.Dimension) *SwatScene {
	type audioStream interface {
		io.ReadSeeker
		Length() int64
	}

	var stream audioStream
	stream, err := mp3.DecodeWithoutResampling(bytes.NewReader(sfx.AsphodelMp3))
	if err != nil {
		log.Fatal(err)
	}

	ac := audio.NewContext(sfx.SamepleRate)
	ap, err := ac.NewPlayer(stream)
	if err != nil {
		log.Fatal(err)
	}

	var scene *SwatScene = &SwatScene{
		bug:          elements.NewBug(),
		swatter:      elements.NewSplat(),
		splatmask:    ebiten.NewImage(constants.SwatWidth+(constants.BugWidth*constants.SwatBugScale)*2, constants.SwatHeight+(constants.BugHeight*constants.SwatBugScale)*2),
		collidermask: ebiten.NewImage(constants.BugWidth*constants.SwatBugScale, constants.BugHeight*constants.SwatBugScale),
		fader:        elements.NewFader(dimensions, definitions.FadeTypeIn, fx.HexToRGBA(0xFFFFFF, 0xff), 4*60),
		cycle:        0,
		tick:         0,
		loaded:       false,
		complete:     false,
		dimensions:   dimensions,
		bugcollision: false,
		whack:        false,
		gameover:     false,
		audioContext: ac,
		audioPlayer:  ap,
		seBytes:      sfx.SwatWav,
		firstflag:    false,
		musicstarted: false,
	}

	return scene
}

func (scene *SwatScene) Draw(img *ebiten.Image) {
	img.Clear()
	img.Fill(color.White)

	scene.RenderSurface(img)
	scene.RenderBug(img)
	scene.RenderSplat(img)
	scene.RenderSwatter(img)
	scene.RenderSwatCam(img)

	if scene.bug.GetAction() != definitions.BugActionGlitch {
		if scene.bugcollision && !scene.whack && !scene.gameover {
			text.Draw(img, constants.Strings.Targeted, fonts.Bugger.Arcade, 1280-150, 150, color.White)
		} else if scene.bugcollision && scene.whack || scene.gameover {
			scene.gameover = true
			text.Draw(img, constants.Strings.Splat, fonts.Bugger.Arcade, 1280-150, 150, color.White)
		}
	}

	if scene.gameover {
		img.DrawImage(scene.fader.Sprite, nil)
		text.Draw(img, constants.Strings.VileScum, fonts.Bugger.Arcade, 500, 300, color.Black)
	}

	scene.cycle++
}

func (scene *SwatScene) Update() error {

	if !scene.gameover {
		scene.handleBugInputs()

		if scene.tick%7 == 0 {
			scene.bug.Animate()
		}
	} else {

		scene.fader.Animate()
		if scene.fader.IsComplete() {
			scene.complete = true
		}
	}

	select {
	case scene.seBytes = <-scene.seCh:
		close(scene.seCh)
		scene.seCh = nil
	default:
	}

	scene.handleOtherInputs()
	scene.swatter.Animate()
	scene.CheckCollisions()

	if scene.gameover && !scene.firstflag {
		scene.firstflag = true
		sePlayer := scene.audioContext.NewPlayerFromBytes(scene.seBytes)
		sePlayer.Play()
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

func (scene *SwatScene) handleBugInputs() {

	var newpos coordinates.Vector = scene.bug.GetLocation()

	if ebiten.IsKeyPressed(ebiten.KeyW) ||
		ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		newpos.Y = newpos.Y - constants.BugSpeed
		direction := coordinates.Direction{
			Straight: true,
			Right:    false,
			Forward:  false,
		}
		scene.bug.SetRole(definitions.BugActionReverseRun, direction)
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

	//GLITCH TIME
	if ebiten.IsKeyPressed(ebiten.KeyF) {
		direction := coordinates.Direction{
			Straight: true,
			Right:    false,
			Forward:  false,
		}
		newpos.X = newpos.X + constants.BugSpeed
		scene.bug.SetRole(definitions.BugActionGlitch, direction)
	}

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

}

func (scene *SwatScene) handleOtherInputs() {

	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		scene.complete = true
	}

	scene.whack = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	var mpos coordinates.Vector
	mpos.X, mpos.Y = ebiten.CursorPosition()
	scene.swatter.SetLocation(mpos)

}

func (scene *SwatScene) RenderBug(img *ebiten.Image) {
	offset := scene.bug.GetLocation()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(constants.SwatBugScale, constants.SwatBugScale)
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

func (scene *SwatScene) RenderSwatter(img *ebiten.Image) {
	offset := scene.swatter.GetLocation()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(constants.SwatWidth)/2.-constants.OffsetSplatX, -float64(constants.SwatHeight)/2.)
	op.GeoM.Translate(float64(offset.X), float64(offset.Y))
	img.DrawImage(scene.swatter.Sprite, op)
}

func (scene *SwatScene) CheckCollisions() {

	bugloc := scene.bug.GetLocation()
	splatloc := scene.swatter.GetLocation()
	splatloc.X = splatloc.X - constants.SwatWidth/2 - constants.OffsetSplatX
	splatloc.Y = splatloc.Y - constants.SwatHeight/2
	var c bool = false
	//check bounding boxes first, if we're in range we need a more precise check
	if bugloc.X >= splatloc.X-(constants.BugWidth*constants.SwatBugScale) && bugloc.X < splatloc.X+constants.SwatWidth &&
		bugloc.Y >= splatloc.Y-(constants.BugHeight*constants.SwatBugScale) && bugloc.Y < splatloc.Y+constants.SwatHeight {
		//fmt.Println("collision requires more precise check")

		scene.collidermask.Clear()
		ox := -float64(bugloc.X - splatloc.X)
		oy := -float64(bugloc.Y - splatloc.Y)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(constants.SwatBugScale, constants.SwatBugScale)
		scene.collidermask.DrawImage(scene.bug.Sprite, op)

		op.GeoM.Reset()
		op.GeoM.Translate((constants.BugWidth * constants.SwatBugScale), (constants.BugHeight * constants.SwatBugScale)) //this provides some margin for the mask that matches the bounds of the triple scaled 32x32 bug sprite
		scene.splatmask.DrawImage(scene.swatter.Sprite, op)

		op.GeoM.Reset()
		op.Blend = ebiten.BlendSourceIn
		op.GeoM.Translate(ox-constants.BugWidth*constants.SwatBugScale, oy-constants.BugHeight*constants.SwatBugScale)
		scene.collidermask.DrawImage(scene.splatmask, op)

		//every fourth byte is our alpha channel
		var pixels []byte = make([]byte, constants.BugWidth*constants.SwatBugScale*constants.BugWidth*constants.SwatBugScale*4)
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

func (scene *SwatScene) RenderSplat(img *ebiten.Image) {

	if scene.bugcollision && scene.whack && scene.bug.GetAction() != definitions.BugActionGlitch || scene.gameover {
		loc := scene.bug.GetLocation()

		loc.X = loc.X - constants.BugWidth/2*constants.SwatBugScale - constants.SplatWidth/2
		loc.Y = loc.Y - constants.BugHeight/2*constants.SwatBugScale - constants.SplatHeight/2

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(loc.X), float64(loc.Y))
		img.DrawImage(images.BugImages[images.IMGBLOOD], op)
	}
}

// let's visualize what's going on with the collision mask and call it the swat cam
func (scene *SwatScene) RenderSwatCam(img *ebiten.Image) {

	ox := float64(img.Bounds().Bounds().Dx() - constants.SplatCamWidth - constants.OffsetSplatCamRightMargin)
	oy := float64(constants.OffsetSplatCamTopMargin)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(ox, oy)
	img.DrawImage(images.BugImages[images.IMGSWATCAM], op)

	if (scene.tick/30)%2 == 0 {
		vector.DrawFilledCircle(img, float32(ox+constants.OffsetSplatCamTopMargin), float32(oy+constants.OffsetSplatCamLeftMargin), 5, fx.HexToRGBA(0xFF0000, 0xFF), true)
	}

	text.Draw(img, "SWAT CAM", fonts.Bugger.Arcade, int(ox)+20, int(oy)+20, color.Black)

	//include an inner margin before we draw the collider mask
	if scene.bugcollision {
		op.GeoM.Translate(constants.OffsetSplatCamLeftMargin, oy)
		img.DrawImage(scene.collidermask, op)
	}
}
