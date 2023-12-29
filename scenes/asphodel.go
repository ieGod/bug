package scenes

import (
	"bug/bugmap"
	"bug/constants"
	"bug/coordinates"
	"bug/definitions"
	"bug/elements"
	"bug/fonts"
	"bug/fx"
	"bug/resources/images"
	_ "embed"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	//go:embed asphodel.json
	mapjson []byte
)

type AsphodelScene struct {
	dimensions coordinates.Dimension
	loaded     bool
	tick       int

	//scene components
	bugmap    *bugmap.Level
	scratch   *ebiten.Image //our pre-generated map, we don't update this
	mcscratch *ebiten.Image //maurice-cam scratchpad
	scene     *ebiten.Image //our updated scene, drawn overtop the pre-generated map
	bugcam    *elements.BugCam

	//scene elements
	ground *ebiten.Image
	wall   *ebiten.Image
	bug    *elements.Bug
	hand   *elements.Handy

	//logical states
	canglitch      bool
	glitching      bool
	gameover       bool
	glitchcooldown int
}

func NewAsphodelScene(dimensions coordinates.Dimension) *AsphodelScene {
	asphodel := &AsphodelScene{
		dimensions:     dimensions,
		loaded:         false,
		bugcam:         elements.NewBugCam(),
		bug:            elements.NewBug(),
		hand:           elements.NewHandy(),
		canglitch:      true,
		glitching:      false,
		glitchcooldown: 0,
		gameover:       false,
	}

	return asphodel
}

func (scene *AsphodelScene) Draw(img *ebiten.Image) {

	img.Clear()
	img.Fill(fx.HexToRGBA(0x25131a, 0xff))

	op := &ebiten.DrawImageOptions{}
	//op.GeoM.Scale(2, 2)

	//draw the npc
	scene.scene.Clear()
	scene.scene.DrawImage(scene.scratch, nil)
	npcloc := scene.hand.GetLoc64()
	op.GeoM.Reset()
	op.GeoM.Translate(npcloc.X, npcloc.Y)
	scene.scene.DrawImage(scene.hand.Sprite, op)

	//draw visible scene portion based on camera settings
	paramecas := scene.bugcam.GetParams()
	mx := paramecas.Location.X
	my := paramecas.Location.Y
	sx := paramecas.Scale.X
	sy := paramecas.Scale.Y
	op.GeoM.Reset()
	op.GeoM.Scale(sx, sy)
	op.GeoM.Translate(-float64(mx), -float64(my))
	img.DrawImage(scene.scene, op)

	//draw the bug
	op.GeoM.Reset()
	op.GeoM.Scale(constants.CameraScale, constants.CameraScale)
	op.GeoM.Translate(constants.CameraScale*constants.BugWidth*6, constants.CameraScale*constants.BugHeight*3)
	img.DrawImage(scene.bug.Sprite, op)

	//clone the scene onto our maurice cam scratchpad and add the player bug
	scene.mcscratch.DrawImage(scene.scene, nil)
	op.GeoM.Reset()
	op.GeoM.Translate(mx/constants.CameraScale+constants.BugWidth*6, my/constants.CameraScale+constants.BugHeight*3)
	scene.mcscratch.DrawImage(scene.bug.Sprite, op)

	//draw MauriceCam™
	ox := int(npcloc.X) - 3*constants.BugWidth
	oy := int(npcloc.Y) - 2*constants.BugHeight
	ox1 := ox + constants.BugWidth*7
	oy1 := oy + constants.BugHeight*5
	op.GeoM.Reset()
	op.GeoM.Translate(50, 50)
	vector.DrawFilledRect(img, 45, 45, 234, 170, fx.HexToRGBA(0x4c2f49, 0xff), true)
	img.DrawImage(scene.mcscratch.SubImage(image.Rect(ox, oy, ox1, oy1)).(*ebiten.Image), op)

	//glitch cooldown indicator
	vector.DrawFilledRect(img, 1280-140, 720-35, 140, 35, fx.HexToRGBA(0x4c2f49, 0xff), true)
	var gclr color.RGBA
	var gtxt string = "GLITCH!"
	if scene.canglitch {
		gclr = fx.HexToRGBA(0x00ff84, 0xff)
	} else {
		gclr = fx.HexToRGBA(0x888888, 0xff)
		gtxt = fmt.Sprintf("GLITCH %d", scene.glitchcooldown/60+1)
	}

	text.Draw(img, gtxt, fonts.Bugger.Arcade, 1280-120, 720-12, gclr)

	if scene.gameover {
		text.Draw(img, "GAME OVER F5 TO RESTART", fonts.Bugger.ArcadeLarge, 175, 325, color.White)
	}
}

func (scene *AsphodelScene) Update() error {
	if !scene.gameover {
		scene.handleInputs()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF5) {
		scene.Load()
	}

	scene.bugcam.CloseTargets()

	if scene.hand.GetLocation().GetManhattanDist(scene.bug.GetLocation()) == 0 && !scene.glitching {
		scene.gameover = true
	}

	if scene.tick%7 == 0 {
		scene.bug.Animate()
		scene.hand.Animate()

	}
	scene.hand.CloseTargets()

	//have maurice re-initiate target acquisition periodically, and often >:]
	if scene.tick%30 == 0 {
		scene.ChasePlayer()
	}

	if scene.glitchcooldown/60 < 5 && scene.glitching {
		scene.SetBugIdle()
		scene.glitching = false

	}

	if scene.glitchcooldown > 0 {
		scene.glitchcooldown -= 1
	} else {
		scene.canglitch = true
	}

	scene.tick++
	return nil
}

func (scene *AsphodelScene) Load() {
	var err error

	//load up and generate map but only if we haven't already, otherwise
	//we can just reset the scene's logical states
	if !scene.loaded {
		//path := "bugmap.json"

		scene.bugmap = &bugmap.Level{}

		err = json.Unmarshal(mapjson, scene.bugmap)

		if err != nil {
			log.Fatal("invalid map.")
		}

		scenedimensions := coordinates.Dimension{
			Width:  scene.bugmap.Dimensions.Width * constants.BugWidth,
			Height: scene.bugmap.Dimensions.Height * constants.BugHeight,
		}

		scene.scene = ebiten.NewImage(scenedimensions.Width, scenedimensions.Height)
		scene.scratch = ebiten.NewImage(scenedimensions.Width, scenedimensions.Height)
		scene.mcscratch = ebiten.NewImage(scenedimensions.Width, scenedimensions.Height)
		scene.ground = ebiten.NewImage(32, 32)
		scene.wall = ebiten.NewImage(32, 32)
		scene.ground.Fill(fx.HexToRGBA(0x44FF44, 0xFF))
		scene.wall.Fill(fx.HexToRGBA(0x000044, 0xFF))

		scene.GenerateMap()
	}

	scene.bugcam.SetParams(definitions.Paramecas{
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
			X: constants.CameraScale,
			Y: constants.CameraScale,
			Z: 0,
		},
		Easing: 8.,
	})

	scene.bug.SetLocation(coordinates.Vector{X: 6, Y: 3})
	scene.bug.SetTargetLocation(coordinates.Vector{X: 6, Y: 3})
	scene.hand.ForceAllPositionsGrid(coordinates.Vector{X: 94, Y: 3})
	scene.SetBugIdle()

	scene.canglitch = true
	scene.glitching = false
	scene.glitchcooldown = 0
	scene.gameover = false

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
		x0 = 32 * (4 - rand.Intn(4))
		y0 = 32 * (3 - rand.Intn(3))
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
	case definitions.NodeTileBlank:
		x0 = 32 * 8
		y0 = 32 * 7
	case definitions.NodeTileInsideTopLeft:
		x0 = 32 * 5
		y0 = 32 * 5
	case definitions.NodeTileInsideTopRight:
		x0 = 32 * 0
		y0 = 32 * 5
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
			scene.scratch.DrawImage(srcimg, op)
		}
	}
}

/*
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
*/

func (scene *AsphodelScene) handleInputs() {

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

		if scene.canglitch && scene.glitchcooldown == 0 {
			direction := coordinates.Direction{
				Straight: true,
				Right:    false,
				Forward:  false,
			}
			newpos.X = newpos.X + 1 //constants.BugSpeed
			scene.bug.SetRole(definitions.BugActionGlitch, direction)

			scene.glitchcooldown = 60 * 10 //5 second cooldown
			scene.canglitch = false
			scene.glitching = true
		}

	}

	if inpututil.IsKeyJustReleased(ebiten.KeyW) ||
		inpututil.IsKeyJustReleased(ebiten.KeyA) ||
		inpututil.IsKeyJustReleased(ebiten.KeyS) ||
		inpututil.IsKeyJustReleased(ebiten.KeyD) ||
		inpututil.IsKeyJustReleased(ebiten.KeyArrowUp) ||
		inpututil.IsKeyJustReleased(ebiten.KeyArrowLeft) ||
		inpututil.IsKeyJustReleased(ebiten.KeyArrowDown) ||
		inpututil.IsKeyJustReleased(ebiten.KeyArrowRight) {
		scene.SetBugIdle()
	}

	if update {

		//fmt.Printf("%d, %d\n", newpos.X, newpos.Y)

		//first criteria passed, we're within map bounds and an update is required, next we
		//check whether or not the new position is a legal move
		gidx := newpos.Y*scene.bugmap.Dimensions.Width + newpos.X

		if scene.bugmap.Nodes[gidx].Nodetype != definitions.NodeTypeWall && scene.tick%7 == 0 {
			scene.bug.SetLocation(newpos)
			scene.bug.SetTargetLocation(newpos)
			scene.bugcam.SetParams(params)
			scene.glitching = false
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		scene.hand.GenWaypoints()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyN) {
		scene.ChasePlayer()
	}
}

// we're going to use a* to find a path from the npc to the player
func (scene *AsphodelScene) ChasePlayer() {

	var waypoints []coordinates.Vector

	//find gidxs for the npc and player
	npcloc := scene.hand.GetLocation()
	gidx0 := npcloc.Y*scene.bugmap.Dimensions.Width + npcloc.X
	start := scene.bugmap.Nodes[gidx0]

	bugloc := scene.bug.GetLocation()
	var gidx1 int
	var goal *bugmap.BugNode

	if scene.glitching {
		gidx1 = 2*scene.bugmap.Dimensions.Width + 2
		goal = scene.bugmap.Nodes[gidx1]

	} else {
		gidx1 = bugloc.Y*scene.bugmap.Dimensions.Width + bugloc.X
		goal = scene.bugmap.Nodes[gidx1]

	}

	wp := bugmap.AStar(start, goal, scene.bugmap)

	//fmt.Println(wp)

	for _, point := range wp {
		waypoints = append(waypoints, point.Location)
	}

	scene.hand.SetWaypoints(waypoints)
}

func (scene *AsphodelScene) SetBugIdle() {
	direction := coordinates.Direction{
		Straight: true,
		Right:    false,
		Forward:  true,
	}
	scene.bug.SetRole(definitions.BugActionIdle, direction)
}
