package fx

import "image/color"

func HexToRGBA(hexcolor int, alpha uint8) color.RGBA {
	return color.RGBA{
		R: uint8(hexcolor >> 0x10 & 0xff),
		G: uint8(hexcolor >> 0x08 & 0xff),
		B: uint8(hexcolor >> 0x00 & 0xff),
		A: alpha,
	}
}
