package fonts

import (
	"bug/constants"
	"log"

	_ "embed"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
)

type BugFont struct {
	Standard  font.Face
	Large     font.Face
	Glitch    font.Face
	GlitchBig font.Face
}

var (

	//go:embed agencyb.ttf
	agency_ttf []byte

	Bugger BugFont
)

func LoadFontFatal(src []byte) *sfnt.Font {
	tt, err := opentype.Parse(src)
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
	Bugger.Standard = GetFaceFatal(fnt, constants.DPI, constants.FontSizeStandard)
	Bugger.Large = GetFaceFatal(fnt, constants.DPI, constants.FontSizeLarge)

	fnt2 := LoadFontFatal(agency_ttf)
	Bugger.Glitch = GetFaceFatal(fnt2, constants.DPI, constants.FontSizeLarge)
	Bugger.GlitchBig = GetFaceFatal(fnt2, constants.DPI, constants.FontSizeBig)
}
