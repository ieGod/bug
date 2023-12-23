package manager

import (
	"bug/defaults"
	"bug/fonts"
	"bug/scenes"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type Manager struct {
	scenes     []scenes.BugScene
	loadcalled bool
	cycle      int
	tick       int
}

func NewBugManager() *Manager {
	return &Manager{
		loadcalled: false,
		cycle:      0,
		tick:       0,
	}
}

func (m *Manager) LoadScene(scene scenes.BugScene) {
	m.scenes = append(m.scenes, scene)
}

func (m *Manager) ClearScenes() {
	m.scenes = m.scenes[:0]
}

func (m *Manager) Draw(screen *ebiten.Image) {

	//nothing to do if we have no scenes
	if len(m.scenes) == 0 {
		return
	}

	if m.scenes[0].IsLoaded() {
		m.scenes[0].Draw(screen)
	} else {
		m.DrawLoading(screen)
	}

	m.cycle++
}

func (m *Manager) Update() error {

	//we're pretty much done if we have no scenes remaining
	if len(m.scenes) == 0 {
		fmt.Println(defaults.Strings.ThanksTxt)
		return ebiten.Termination
	} else if m.scenes[0].IsLoaded() {
		//let's call our respective update
		m.scenes[0].Update()

		//check for scene completion, if so remove from queue and bounce to next one
		if m.scenes[0].IsComplete() {
			m.scenes = m.scenes[1:]
			m.loadcalled = false
			m.tick = 0
		}

	} else if !m.loadcalled {
		m.scenes[0].Load()
		m.loadcalled = true
	}

	m.tick++
	return nil
}

func (m *Manager) Layout(width, height int) (int, int) {
	return width, height
}

func (m *Manager) DrawLoading(img *ebiten.Image) {
	img.Clear()

	var dots string
	//every half second for five counts
	for i := 0; i < (m.tick/30)%5; i++ {
		dots = dots + " ."
	}

	loadtext := fmt.Sprintf("%s %s %s", defaults.Strings.Loading, m.scenes[0].GetName(), dots)
	text.Draw(img, loadtext, fonts.Bugger.Standard, defaults.OffsetLoadingX, defaults.OffsetLoadingY, color.White)
}
