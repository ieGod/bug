package fonts

import (
	"bug/defaults"
	"log"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
)

type BugFont struct {
	Standard font.Face
}

var (
	Bugger BugFont
)

func LoadFontFatal(src []byte) *sfnt.Font {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}
	return tt
}

func GetFaceFatal(fnt *sfnt.Font, dpi, size float64) font.Face {
	var face font.Face
	var err error

	if dpi > 0 && size > 0 && fnt != nil {
		face, err = opentype.NewFace(fnt, &opentype.FaceOptions{
			Size:    size,
			DPI:     dpi,
			Hinting: font.HintingVertical,
		})
		if err != nil {
			log.Fatal(err)
		}
	}
	return face
}

func init() {
	Bugger = BugFont{}

	fnt := LoadFontFatal(fonts.MPlus1pRegular_ttf)
	Bugger.Standard = GetFaceFatal(fnt, defaults.DPI, defaults.FontSizeStandard)
}
