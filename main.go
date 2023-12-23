package main

import (
	"bug/defaults"
	"bug/manager"
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 1280
	screenHeight = 720
)

func main() {
	fmt.Println(defaults.Strings.VersionInfo)

	bugmanager := manager.NewBugManager()
	bugmanager.LoadScenes()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle(defaults.Strings.Title)
	if err := ebiten.RunGame(bugmanager); err != nil {
		log.Fatal(err)
	}

}
