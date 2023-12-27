package scenes

import (
	"bug/bugmap"
	"bug/constants"
	"bug/coordinates"
	"bug/definitions"
	"bug/elements"
	"bug/fx"
	"bug/resources/images"
	"encoding/json"
	"image"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type AsphodelScene struct {
	dimensions coordinates.Dimension
	loaded     bool
	tick       int

	//scene components
	bugmap *bugmap.Level
	scene  *ebiten.Image
	bugcam *elements.BugCam

	//scene elements
	ground *ebiten.Image
	wall   *ebiten.Image
	bug    *elements.Bug
}

func NewAsphodelScene(dimensions coordinates.Dimension) *AsphodelScene {
	asphodel := &AsphodelScene{
		dimensions: dimensions,
		loaded:     false,
		bugcam:     elements.NewBugCam(),
		bug:        elements.NewBug(),
	}

	asphodel.bugcam.SetParams(definitions.Paramecas{
		Location: coordinates.Vector64{
			X: 0,
			Y: 0,
			Z: 0,
		},
		TargetLocation: coordinates.Vector64{
			X: 0,
			Y: 0,
			Z: 0,
		},
		Scale: coordinates.Vector64{
			X: constants.Scale,
			Y: constants.Scale,
			Z: 0,
		},
		Easing: 8.,
	})

	asphodel.bug.SetLocation(coordinates.Vector{X: 6, Y: 3})
	asphodel.bug.SetTargetLocation(coordinates.Vector{X: 6, Y: 3})
	return asphodel
}

func (scene *AsphodelScene) Draw(img *ebiten.Image) {

	img.Clear()

	op := &ebiten.DrawImageOptions{}
	//op.GeoM.Scale(2, 2)

	paramecas := scene.bugcam.GetParams()

	mx := paramecas.Location.X
	my := paramecas.Location.Y
	sx := paramecas.Scale.X
	sy := paramecas.Scale.Y

	op.GeoM.Scale(sx, sy)
	op.GeoM.Translate(-float64(mx), -float64(my))

	img.DrawImage(scene.scene, op)

	//todo: retrieve bug position from location, decouple from camera
	op.GeoM.Reset()
	op.GeoM.Scale(constants.Scale, constants.Scale)
	op.GeoM.Translate(constants.Scale*constants.BugWidth*6, constants.Scale*constants.BugHeight*3)
	img.DrawImage(scene.bug.Sprite, op)
}

func (scene *AsphodelScene) Update() error {
	scene.handleInputs2()
	scene.bugcam.CloseTargets()

	if scene.tick%7 == 0 {
		scene.bug.Animate()
	}

	scene.tick++
	return nil
}

func (scene *AsphodelScene) Load() {
	var err error
	path := "bugmap.json"

	scene.bugmap = &bugmap.Level{}

	//check file exists, load, deserialize
	_, err = os.Stat(path)
	if err == nil {
		rawbytes, err := os.ReadFile(path)

		if err == nil {
			err = json.Unmarshal(rawbytes, scene.bugmap)
		}

		if err != nil {
			log.Fatal("invalid map.")
		}
	}

	scenedimensions := coordinates.Dimension{
		Width:  scene.bugmap.Dimensions.Width * constants.BugWidth,
		Height: scene.bugmap.Dimensions.Height * constants.BugHeight,
	}

	scene.scene = ebiten.NewImage(scenedimensions.Width, scenedimensions.Height)
	scene.ground = ebiten.NewImage(32, 32)
	scene.wall = ebiten.NewImage(32, 32)
	scene.ground.Fill(fx.HexToRGBA(0x44FF44, 0xFF))
	scene.wall.Fill(fx.HexToRGBA(0x000044, 0xFF))

	scene.GenerateMap()

	scene.loaded = true

}

func (scene *AsphodelScene) IsLoaded() bool {
	return scene.loaded
}

func (scene *AsphodelScene) IsComplete() bool {
	return false
}

func (scene *AsphodelScene) GetName() string {
	return constants.Strings.AsphodelName
}

func (scene *AsphodelScene) GetImageFromNodeTile(tiletype definitions.NodeTile) *ebiten.Image {

	var img *ebiten.Image
	var x0, y0 int
	switch tiletype {
	case definitions.NodeTileGround:
		x0 = 32
		y0 = 32
	case definitions.NodeTileWallBottom:
		x0 = 32
		y0 = 32 * 4
	case definitions.NodeTileWallLeft,
		definitions.NodeTileWallTopLeft:
		x0 = 0
		y0 = 32
	case definitions.NodeTileWallRight,
		definitions.NodeTileWallTopRight:
		x0 = 32 * 5
		y0 = 32
	case definitions.NodeTileWallTop:
		x0 = 32
		y0 = 0
	case definitions.NodeTileWallBottomLeft:
		x0 = 0
		y0 = 32 * 4
	case definitions.NodeTileWallBottomRight:
		x0 = 32 * 5
		y0 = 32 * 4
	}

	img = images.BugImages[images.IMGTILESET].SubImage(image.Rect(x0, y0, x0+32, y0+32)).(*ebiten.Image)
	return img

}

func (scene *AsphodelScene) GenerateMap() {
	op := &ebiten.DrawImageOptions{}
	//draw the map/tiles
	for w := 0; w < scene.bugmap.Dimensions.Width; w++ {
		for h := 0; h < scene.bugmap.Dimensions.Height; h++ {
			grididx := h*scene.bugmap.Dimensions.Width + w

			var srcimg *ebiten.Image = scene.GetImageFromNodeTile(scene.bugmap.Nodes[grididx].Nodetile)

			ox := float64(w * 32)
			oy := float64(h * 32)

			op.GeoM.Reset()
			op.GeoM.Translate(ox, oy)
			scene.scene.DrawImage(srcimg, op)
		}
	}
}

func (scene *AsphodelScene) handleInputs() {

	var update bool = false

	params := scene.bugcam.GetParams()
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		params.TargetLocation.X += 32 * params.Scale.X
		update = true
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		params.TargetLocation.X -= 32 * params.Scale.X
		update = true
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		params.TargetLocation.Y -= 32 * params.Scale.Y
		update = true
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		params.TargetLocation.Y += 32 * params.Scale.Y
		update = true
	}

	if update {
		scene.bugcam.SetParams(params)
	}
}

func (scene *AsphodelScene) handleInputs2() {

	var newpos coordinates.Vector = scene.bug.GetLocation()

	var update bool = false

	params := scene.bugcam.GetParams()

	if ebiten.IsKeyPressed(ebiten.KeyW) ||
		ebiten.IsKeyPressed(ebiten.KeyArrowUp) {

		direction := coordinates.Direction{
			Straight: true,
			Right:    false,
			Forward:  false,
		}
		newpos.Y = newpos.Y - 1 //constants.BugSpeed
		scene.bug.SetRole(definitions.BugActionReverseRun, direction)
		params.TargetLocation.Y -= 32 * params.Scale.Y
		update = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) ||
		ebiten.IsKeyPressed(ebiten.KeyArrowDown) {

		direction := coordinates.Direction{
			Straight: true,
			Right:    false,
			Forward:  true,
		}
		newpos.Y = newpos.Y + 1 //constants.BugSpeed
		scene.bug.SetRole(definitions.BugActionForwardRun, direction)
		params.TargetLocation.Y += 32 * params.Scale.Y
		update = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) ||
		ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		direction := coordinates.Direction{
			Straight: false,
			Right:    false,
			Forward:  true,
		}
		newpos.X = newpos.X - 1 //constants.BugSpeed
		scene.bug.SetRole(definitions.BugActionSideRun, direction)
		params.TargetLocation.X -= 32 * params.Scale.X
		update = true

	}
	if ebiten.IsKeyPressed(ebiten.KeyD) ||
		ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		direction := coordinates.Direction{
			Straight: false,
			Right:    true,
			Forward:  true,
		}
		newpos.X = newpos.X + 1 //constants.BugSpeed
		scene.bug.SetRole(definitions.BugActionSideRun, direction)
		params.TargetLocation.X += 32 * params.Scale.X
		update = true
	}

	//GLITCH TIME
	if ebiten.IsKeyPressed(ebiten.KeyF) {
		direction := coordinates.Direction{
			Straight: true,
			Right:    false,
			Forward:  false,
		}
		newpos.X = newpos.X + 1 //constants.BugSpeed
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

	if update && 0 <= newpos.X && newpos.X < scene.bugmap.Dimensions.Width &&
		0 <= newpos.Y && newpos.Y < scene.bugmap.Dimensions.Height {

		//fmt.Printf("%d, %d\n", newpos.X, newpos.Y)

		//gidx := newpos.Y*scene.bugmap.Dimensions.Width + newpos.X

		if scene.tick%7 == 0 {
			scene.bug.SetLocation(newpos)
			scene.bug.SetTargetLocation(newpos)
			scene.bugcam.SetParams(params)
		}
	}
}
