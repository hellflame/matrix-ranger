package styles

import "image/color"

type Theme struct {
	Block      *Block
	Background *Background
	Trophy     *Trophy

	Shapes []*Block
}

type Block struct {
	Color color.NRGBA
	Image string
}

type Background struct {
	Image string
}

type Trophy struct {
	Image string
}

func loadDefaultTheme() *Theme {
	return &Theme{
		Block: &Block{
			Color: color.NRGBA{R: 148, G: 153, B: 160, A: 0x66},
		},
		Shapes: []*Block{{
			Color: color.NRGBA{R: 167, G: 143, B: 211, A: 0xff},
		}, {
			Color: color.NRGBA{R: 103, G: 58, B: 183, A: 0xff},
		}, {
			Color: color.NRGBA{R: 33, G: 150, B: 243, A: 0xff},
		}, {
			Color: color.NRGBA{R: 0, G: 188, B: 212, A: 0xff},
		}, {
			Color: color.NRGBA{R: 0, G: 150, B: 136, A: 0xff},
		}, {
			Color: color.NRGBA{R: 76, G: 175, B: 80, A: 0xff},
		}, {
			Color: color.NRGBA{R: 255, G: 193, B: 7, A: 0xff},
		}, {
			Color: color.NRGBA{R: 255, G: 87, B: 34, A: 0xff},
		}, {
			Color: color.NRGBA{R: 96, G: 125, B: 139, A: 0xff},
		}, {
			Color: color.NRGBA{R: 255, G: 161, B: 169, A: 0xff},
		}, {
			Color: color.NRGBA{R: 187, G: 193, B: 238, A: 0xff},
		},
		},
	}
}
