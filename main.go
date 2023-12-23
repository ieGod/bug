package main

import (
	"bug/defaults"
	"bug/manager"
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	fmt.Println(defaults.Strings.VersionInfo)

	bugmanager := manager.NewBugManager()
	bugmanager.LoadScenes()

	ebiten.SetWindowSize(defaults.ScreenWidth, defaults.ScreenHeight)
	ebiten.SetWindowTitle(defaults.Strings.Title)
	if err := ebiten.RunGame(bugmanager); err != nil {
		log.Fatal(err)
	}

}
