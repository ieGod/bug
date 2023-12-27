package scenes

import (
	"bug/bugmap"
	"bug/constants"
	"bug/coordinates"
	"bug/definitions"
	"bug/fx"
	"bug/resources/images"
	"encoding/json"
	"image"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

type AsphodelScene struct {
	dimensions coordinates.Dimension
	loaded     bool

	//scene components
	bugmap *bugmap.Level
	scene  *ebiten.Image

	//scene elements
	ground *ebiten.Image
	wall   *ebiten.Image
}

func NewAsphodelScene(dimensions coordinates.Dimension) *AsphodelScene {
	return &AsphodelScene{
		dimensions: dimensions,
		loaded:     false,
	}
}

func (scene *AsphodelScene) Draw(img *ebiten.Image) {

	img.Clear()

	img.DrawImage(scene.scene, nil)
}

func (scene *AsphodelScene) Update() error {
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
		Width:  scene.bugmap.Dimensions.Width * 32,
		Height: scene.bugmap.Dimensions.Height * 32,
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
