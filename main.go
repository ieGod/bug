package main

import (
	"bug/constants"
	"bug/coordinates"
	"bug/manager"
	"bug/scenes"
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	fmt.Println(constants.Strings.VersionInfo)

	bugmanager := manager.NewBugManager()
	bugmanager.LoadScene(scenes.NewIntroScene(coordinates.Dimension{Width: constants.ScreenWidth, Height: constants.ScreenHeight}))
	bugmanager.LoadScene(scenes.NewSwatScene(coordinates.Dimension{Width: constants.ScreenWidth, Height: constants.ScreenHeight}))

	ebiten.SetWindowSize(constants.ScreenWidth, constants.ScreenHeight)
	ebiten.SetWindowTitle(constants.Strings.Title)
	if err := ebiten.RunGame(bugmanager); err != nil {
		log.Fatal(err)
	}

}
