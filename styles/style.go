package styles

type Style struct {
	PileStyle
	BlockCount int

	StageInnerOffset int

	OffsetTop  int
	OffsetLeft int

	MaxLevel int

	CurrentTheme *Theme

	themes map[string]*Theme
}

type PileStyle struct {
	BlockSize  int
	BlockRound int
	BlockSpace int
}

func CreateStyle(size, count, round, space int) *Style {
	s := &Style{
		BlockCount:       count,
		StageInnerOffset: size,
		OffsetTop:        size * 2,
		OffsetLeft:       size,
		PileStyle: PileStyle{
			BlockSize:  size,
			BlockRound: round,
			BlockSpace: space,
		},

		themes: make(map[string]*Theme),
	}
	s.themes["default"] = loadDefaultTheme()

	return s
}

func (s *Style) SetTheme(theme string) {
	s.CurrentTheme = s.themes[theme]
}
