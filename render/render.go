package render

import "github.com/matttproud/go-quake/bsp"

type Texture struct {
	Name             string
	Width, Height    uint
	AnimTotal        int
	AnimMin, AnimMax int
	Next             *Texture
	Alternate        *Texture
	Offsets          [bsp.MipLevels]uint
	Data             []byte
}

var NoTexture *Texture

func init() {
	NoTexture = &Texture{
		Width:  16,
		Height: 16,
	}
	NoTexture.Offsets[0] = 0
	NoTexture.Offsets[1] = NoTexture.Offsets[0] + 16*16
	NoTexture.Offsets[2] = NoTexture.Offsets[1] + 8*8
	NoTexture.Offsets[3] = NoTexture.Offsets[2] + 4*4
	for m := uint(0); m < uint(len(NoTexture.Offsets)); m++ {
		d := NoTexture.Offsets[m]

		for y := uint(0); y < (uint(16) >> m); y++ {
			for x := uint(0); x < (uint(16) >> m); x++ {
				l := (y < (uint(8) >> m))
				r := (x < (uint(8) >> m))
				if (r || l) && !(r && l) {
					NoTexture.Data[d] = 0
				} else {
					NoTexture.Data[d] = 0xff
				}
				d++

			}
		}
	}
}
