package scenes

import (
	"bug/bugmap"
	"bug/constants"
	"bug/coordinates"
	"bug/definitions"
	"bug/fx"
	"encoding/json"
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

	op := &ebiten.DrawImageOptions{}
	//draw the map/tiles
	for w := 0; w < scene.bugmap.Dimensions.Width; w++ {
		for h := 0; h < scene.bugmap.Dimensions.Height; h++ {
			grididx := h*scene.bugmap.Dimensions.Width + w

			var srcimg *ebiten.Image
			switch scene.bugmap.Nodes[grididx].Nodetype {
			case definitions.NodeTypeGround:
				srcimg = scene.ground
			case definitions.NodeTypeWall:
				srcimg = scene.wall
			}

			ox := float64(w * 32)
			oy := float64(h * 32)

			op.GeoM.Reset()
			op.GeoM.Translate(ox, oy)
			scene.scene.DrawImage(srcimg, op)
		}
	}

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
