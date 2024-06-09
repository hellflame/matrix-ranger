package stage

import (
	"image"
	"math/rand"

	"gioui.org/f32"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"github.com/hellflame/matrix-ranger/blocks"
	"github.com/hellflame/matrix-ranger/styles"
)

// type candidateState int

// const (
// 	calm candidateState = iota
// 	dragging
// )

type candidate struct {
	shape blocks.Shape
	theme *styles.Block

	consumed bool
	chosen   bool

	cache op.CallOp

	state  string
	styles map[string]styles.PileStyle

	currentStyle styles.PileStyle

	defaultPos f32.Point
	presentPos f32.Point
}

func NewCandidate(shape blocks.Shape, pos image.Point, calm, drag styles.PileStyle, theme *styles.Block) *candidate {
	c := &candidate{
		shape: shape, theme: theme,
		state: "default",
		styles: map[string]styles.PileStyle{
			"default":  calm,
			"dragging": drag,
		},
	}
	c.currentStyle = c.styles[c.state]

	left, top := c.GetInnerOffset()
	p := f32.Point{X: float32(pos.X + left), Y: float32(pos.Y + top)}
	c.defaultPos = p
	c.presentPos = p
	c.updateCache()
	return c
}

func (c *candidate) ToggleDrag(dragging bool) {
	if dragging {
		c.state = "dragging"
	} else {
		c.state = "default"
	}
	c.currentStyle = c.styles[c.state]
	c.updateCache()
}

func (c *candidate) Interest() event.Filter {
	return pointer.Filter{
		Target: c,
		Kinds:  pointer.Press,
	}
}
func (c *candidate) UpdatePosition(p f32.Point) {
	c.presentPos = p
}

func (c *candidate) OnEvent(ev event.Event) {
	x, ok := ev.(pointer.Event)
	if !ok {
		return
	}
	switch x.Kind {
	case pointer.Press:
		c.ToggleChosen(true)
	}
	println(c.shape.Desc())
}

func (c *candidate) ToggleChosen(chosen bool) {
	c.chosen = chosen
	c.updateCache()
}

func (c *candidate) GetMaxWidth() int {
	return c.currentStyle.BlockSize*5 + 4*c.currentStyle.BlockSpace
}

func (c *candidate) GetInnerOffset() (int, int) {
	maxWidth := c.GetMaxWidth()
	leftOffset := (maxWidth - c.GetWidth()) / 2
	topOffset := (maxWidth - c.GetHeight()) / 2
	return leftOffset, topOffset
}

// offset to top-left corner point
func (c *candidate) GetCenterOffset() (int, int) {
	leftOffset := c.GetWidth() / 2
	topOffset := c.GetHeight() / 2
	return leftOffset, topOffset
}

func (c *candidate) GetWidth() int {
	space := c.currentStyle.BlockSpace
	return len(c.shape[0])*(space+c.currentStyle.BlockSize) - space
}

func (c *candidate) GetHeight() int {
	space := c.currentStyle.BlockSpace
	return len(c.shape)*(space+c.currentStyle.BlockSize) - space
}

func (c *candidate) renderCache(ops *op.Ops) {
	c.cache.Add(ops)
}

func (c *candidate) updateCache() {
	ops := new(op.Ops)
	macro := op.Record(ops)
	defer func() {
		c.cache = macro.Stop()
	}()

	adjust := c.currentStyle
	space := adjust.BlockSpace
	round := adjust.BlockRound
	blockSize := adjust.BlockSize

	area := clip.Rect(image.Rect(0, 0, c.GetWidth()+blockSize*2, c.GetHeight()+blockSize*2))

	defer area.Push(ops).Pop()
	event.Op(ops, c)

	defer op.Offset(image.Pt(blockSize, blockSize)).Push(ops).Pop() // inner offset

	blockColor := c.theme.Color
	bounds := image.Rect(0, 0, blockSize, blockSize)

	for r, row := range c.shape {
		rowOffset := op.Offset(image.Pt(0, r*(space+blockSize))).Push(ops)
		for c, dot := range row {
			colOffset := op.Offset(image.Pt(c*(space+blockSize), 0)).Push(ops)
			if dot {
				b := clip.RRect{Rect: bounds, SE: round, SW: round, NW: round, NE: round}.Push(ops)
				paint.ColorOp{Color: blockColor}.Add(ops)
				paint.PaintOp{}.Add(ops)
				b.Pop()
			}
			colOffset.Pop()
		}
		rowOffset.Pop()
	}
}

func (c *candidate) Render(ops *op.Ops) {
	adjust := c.currentStyle

	blockSize := adjust.BlockSize

	pos := c.presentPos
	left, top := int(pos.X), int(pos.Y)
	defer op.Offset(image.Pt(left-blockSize, top-blockSize)).Push(ops).Pop() // outer offset

	c.renderCache(ops)
}

type candidateGroup struct {
	style *styles.Style
	calm  styles.PileStyle
	drag  styles.PileStyle

	offsetTop int

	level int
	width int
	gap   int
	count int

	shapeGroups *blocks.ShapeGroups
}

func NewCandidateGroup(maxWidth, offsetTop, level, count int, style *styles.Style, rnd *rand.Rand) *candidateGroup {
	minGap := style.BlockSize / 2

	candidateWidth := (maxWidth - 2*minGap) / 3
	candidateSize := (candidateWidth - 4*style.BlockSpace) / 5

	return &candidateGroup{
		// style: style,
		width: candidateWidth,
		gap:   minGap, level: level,
		offsetTop: offsetTop, count: count,
		style: style,

		calm: styles.PileStyle{
			BlockSize:  candidateSize,
			BlockRound: style.BlockRound,
			BlockSpace: int(float32(style.BlockSpace) * 0.7),
		},
		drag: styles.PileStyle{
			BlockSize:  style.BlockSize - 3,
			BlockRound: style.BlockRound,
			BlockSpace: style.BlockSpace,
		},
		shapeGroups: blocks.NewShapeGroups(rnd),
	}
}

func (cg *candidateGroup) GenerateCandidates() []*candidate {
	result := make([]*candidate, cg.count)

	candidateOffset := cg.width + cg.gap
	offsetTop := cg.offsetTop + cg.style.BlockSize/3
	theme := cg.style.CurrentTheme

	for i := 0; i < cg.count; i++ {
		shapeIdx, shape := cg.shapeGroups.ChooseOneShape(0)
		result[i] = NewCandidate(shape,
			image.Point{X: i * candidateOffset, Y: offsetTop},
			cg.calm, cg.drag, theme.Shapes[shapeIdx])
	}
	return result
}
