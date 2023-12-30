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
	bug  *elements.Bug
	hand *elements.Handy

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

	scene.DrawScene(img)
	scene.DrawGlitchIndicator(img)

	if scene.gameover {
		text.Draw(img, constants.Strings.GameOverReset, fonts.Bugger.ArcadeLarge, constants.OffsetGameOverResetX, constants.OffsetGameOverResetY, color.White)
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

		scene.bug.SetRole(definitions.BugActionDeath, coordinates.Direction{})
		scene.hand.SetRole(definitions.BugActionDeath, coordinates.Direction{})
	} else if scene.gameover {
		scene.SetBugIdle()
	}

	if scene.tick%7 == 0 {
		scene.bug.Animate()
		scene.hand.Animate()

	}
	scene.hand.CloseTargets()

	//have maurice re-initiate target acquisition periodically, and often >:]
	var chasecheck bool = scene.tick%30 == 0
	if scene.glitching {
		chasecheck = scene.tick%(60*2) == 0
	}

	if chasecheck {
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
		Easing: constants.AnimationCameraEasing,
	})

	scene.bug.ForceAllPositionsGrid(constants.LocationBugStart)
	scene.hand.ForceAllPositionsGrid(constants.LocationMauriceStart)
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
		x0 = (4 - rand.Intn(4))
		y0 = (3 - rand.Intn(3))
	case definitions.NodeTileWallBottom:
		x0 = 1
		y0 = 4
	case definitions.NodeTileWallLeft,
		definitions.NodeTileWallTopLeft:
		x0 = 0
		y0 = 1
	case definitions.NodeTileWallRight,
		definitions.NodeTileWallTopRight:
		x0 = 5
		y0 = 1
	case definitions.NodeTileWallTop:
		x0 = 1
		y0 = 0
	case definitions.NodeTileWallBottomLeft:
		x0 = 0
		y0 = 4
	case definitions.NodeTileWallBottomRight:
		x0 = 5
		y0 = 4
	case definitions.NodeTileBlank:
		x0 = 8
		y0 = 7
	case definitions.NodeTileInsideTopLeft:
		x0 = 5
		y0 = 5
	case definitions.NodeTileInsideTopRight:
		x0 = 0
		y0 = 5
	}

	x0 = constants.BugHeight * x0
	y0 = constants.BugHeight * y0
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

			ox := float64(w * constants.TileWidth)
			oy := float64(h * constants.TileHeight)

			op.GeoM.Reset()
			op.GeoM.Translate(ox, oy)
			scene.scratch.DrawImage(srcimg, op)
		}
	}
}

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
		params.TargetLocation.X -= constants.TileWidth * params.Scale.X
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
		params.TargetLocation.X += constants.TileHeight * params.Scale.X
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
			newpos.X = newpos.X + 1
			scene.bug.SetRole(definitions.BugActionGlitch, direction)

			scene.glitchcooldown = constants.TimerGlitchCooldown
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

	//find grid indeces for the npc and player
	npcloc := scene.hand.GetLocation()
	gidx0 := npcloc.Y*scene.bugmap.Dimensions.Width + npcloc.X
	start := scene.bugmap.Nodes[gidx0]

	bugloc := scene.bug.GetLocation()
	var gidx1 int
	var goal *bugmap.BugNode

	if scene.glitching {
		//send maurice to a random walkable location on the map
		nodes := scene.bugmap.GetWalkableNodes()
		goal = nodes[rand.Intn(len(nodes))]

	} else {
		gidx1 = bugloc.Y*scene.bugmap.Dimensions.Width + bugloc.X
		goal = scene.bugmap.Nodes[gidx1]

	}

	wp := bugmap.AStar(start, goal, scene.bugmap)

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

// render our map, player, npc, and associated views
func (scene *AsphodelScene) DrawScene(img *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	//draw the npc
	scene.scene.Clear()
	scene.scene.DrawImage(scene.scratch, nil)
	npcloc := scene.hand.GetLoc64()
	op.GeoM.Reset()
	op.GeoM.Translate(npcloc.X, npcloc.Y)
	scene.scene.DrawImage(scene.hand.Sprite, op)

	//draw visible scene portion based on camera settings
	paramecas := scene.bugcam.GetParams()
	mx := float64(paramecas.Location.X)
	my := float64(paramecas.Location.Y)
	sx := paramecas.Scale.X
	sy := paramecas.Scale.Y
	op.GeoM.Reset()
	op.GeoM.Scale(sx, sy)

	if scene.glitching {
		mx = mx + (rand.Float64()*3 - 1.5)
		my = my + (rand.Float64()*3 - 1.5)
	}

	op.GeoM.Translate(-mx, -my)
	img.DrawImage(scene.scene, op)

	//draw the bug
	op.GeoM.Reset()
	op.GeoM.Scale(constants.CameraScale, constants.CameraScale)
	op.GeoM.Translate(constants.CameraScale*constants.BugWidth*float64(constants.LocationBugStart.X), constants.CameraScale*constants.BugHeight*float64(constants.LocationBugStart.Y))
	img.DrawImage(scene.bug.Sprite, op)

	//clone the scene onto our maurice cam scratchpad and add the player bug
	scene.mcscratch.DrawImage(scene.scene, nil)
	op.GeoM.Reset()
	mcx := mx/constants.CameraScale + constants.BugWidth*float64(constants.LocationBugStart.X)
	mcy := my/constants.CameraScale + constants.BugHeight*float64(constants.LocationBugStart.Y)
	op.GeoM.Translate(mcx, mcy)
	scene.mcscratch.DrawImage(scene.bug.Sprite, op)

	//draw MauriceCam™
	ox := int(npcloc.X) - 3*constants.BugWidth
	oy := int(npcloc.Y) - 2*constants.BugHeight

	if scene.glitching {

		ox = rand.Intn(1280)
		oy = rand.Intn(720)
	}

	ox1 := ox + constants.BugWidth*7
	oy1 := oy + constants.BugHeight*5
	op.GeoM.Reset()
	op.GeoM.Translate(50, 50)
	vector.DrawFilledRect(img, 45, 45, 234, 170, fx.HexToRGBA(0x4c2f49, 0xff), true)
	img.DrawImage(scene.mcscratch.SubImage(image.Rect(ox, oy, ox1, oy1)).(*ebiten.Image), op)
}

func (scene *AsphodelScene) DrawGlitchIndicator(img *ebiten.Image) {
	//glitch cooldown indicator
	vector.DrawFilledRect(img, 1280-140, 720-35, 140, 35, fx.HexToRGBA(0x4c2f49, 0xff), true)
	vector.DrawFilledCircle(img, 1280-140, 720-17.5, 17.5, fx.HexToRGBA(0x4c2f49, 0xff), true)
	var gclr color.RGBA
	var gtxt string = "GLITCH!"
	if scene.canglitch {
		gclr = fx.HexToRGBA(0x00ff84, 0xff)
	} else {
		gclr = fx.HexToRGBA(0x888888, 0xff)
		gtxt = fmt.Sprintf("GLITCH %d", scene.glitchcooldown/60+1)
	}

	text.Draw(img, gtxt, fonts.Bugger.Arcade, 1280-130, 720-12, gclr)
}
