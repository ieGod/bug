package main

import (
	"bug/coordinates"
	"bug/defaults"
	"bug/manager"
	"bug/scenes"
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	fmt.Println(defaults.Strings.VersionInfo)

	bugmanager := manager.NewBugManager()
	bugmanager.LoadScene(scenes.NewIntroScene(coordinates.Dimension{Width: defaults.ScreenWidth, Height: defaults.ScreenHeight}))
	bugmanager.LoadScene(scenes.NewSwatScene(coordinates.Dimension{Width: defaults.ScreenWidth, Height: defaults.ScreenHeight}))

	ebiten.SetWindowSize(defaults.ScreenWidth, defaults.ScreenHeight)
	ebiten.SetWindowTitle(defaults.Strings.Title)
	if err := ebiten.RunGame(bugmanager); err != nil {
		log.Fatal(err)
	}

}
