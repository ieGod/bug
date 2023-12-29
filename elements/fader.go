package elements

import (
	"bug/coordinates"
	"bug/definitions"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Fader struct {
	Sprite     *ebiten.Image
	src        *ebiten.Image
	dimensions coordinates.Dimension
	fadetype   definitions.FadeType
	color      color.RGBA
	alpha      float32
	cycle      int
	duration   int
	complete   bool
}

func NewFader(dim coordinates.Dimension, fadetype definitions.FadeType, color color.RGBA, tickduration int) *Fader {
	fader := &Fader{
		Sprite:     ebiten.NewImage(dim.Width, dim.Height),
		src:        ebiten.NewImage(dim.Width, dim.Height),
		dimensions: dim,
		fadetype:   fadetype,
		color:      color,
		duration:   tickduration,
		complete:   false,
	}

	switch fadetype {
	case definitions.FadeTypeIn:
		fader.alpha = 0
	case definitions.FadeTypeOut:
		fader.alpha = 1
	}

	fader.src.Fill(color)

	return fader
}

func (fader *Fader) Animate() {

	fader.Sprite.Clear()
	op := &ebiten.DrawImageOptions{}
	op.ColorScale.ScaleAlpha(fader.alpha)
	fader.Sprite.DrawImage(fader.src, op)

	switch fader.fadetype {
	case definitions.FadeTypeIn:
		fader.alpha = float32(fader.cycle) / float32(fader.duration)
		if fader.alpha > 1 {
			fader.complete = true
		}
	case definitions.FadeTypeOut:
		fader.alpha = 1 - float32(fader.cycle)/float32(fader.duration)
		if fader.alpha < 0 {
			fader.complete = true
		}
	}

	fader.cycle++
}

func (fader *Fader) IsComplete() bool {
	return fader.complete
}
