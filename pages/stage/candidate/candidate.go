package candidate

import (
	"image"
	"math/rand"
	"time"

	"gioui.org/f32"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"github.com/hellflame/matrix-ranger/blocks"
	"github.com/hellflame/matrix-ranger/framework"
	"github.com/hellflame/matrix-ranger/styles"
)

// type candidateState int

// const (
// 	calm candidateState = iota
// 	dragging
// )

type Candidate struct {
	shape blocks.Shape
	theme *styles.Block

	consumed bool
	chosen   bool

	cache op.CallOp

	state  State
	styles map[State]styles.PileStyle

	currentStyle styles.PileStyle

	defaultPos  f32.Point
	presentPos  f32.Point
	transitions []styles.Transition
}

func NewCandidate(shape blocks.Shape, pos image.Point, calm, drag styles.PileStyle, theme *styles.Block) *Candidate {
	c := &Candidate{
		shape: shape, theme: theme,
		state: stable,
		styles: map[State]styles.PileStyle{
			stable:   calm,
			dragging: drag,
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

func (c *Candidate) ToPos(p f32.Point) {
	c.state = movingToCursor

	startPos := c.presentPos
	start := time.Now()
	duration := float32(300.)
	vx := (p.X - startPos.X) / duration
	vy := (p.Y - startPos.Y) / duration
	c.transitions = append(c.transitions, func() bool {
		elapsed := float32(time.Since(start).Milliseconds())
		if elapsed > duration {
			return false
		}
		c.presentPos.X = startPos.X + vx*elapsed
		c.presentPos.Y = startPos.Y + vy*elapsed
		return true
	})

	// targetStyle := c.styles[dragging]
	// startStyle := c.styles[stable]
	// d := float32(100.)
	// vr := float32(targetStyle.BlockRound-startStyle.BlockRound) / d
	// vs := float32(targetStyle.BlockSpace-startStyle.BlockSpace) / d
	// vsize := float32(targetStyle.BlockSize-startStyle.BlockSize) / d
	// c.transitions = append(c.transitions, func() bool {
	// 	elapsed := float32(time.Since(start).Milliseconds())
	// 	if elapsed > d {
	// 		return false
	// 	}
	// 	c.currentStyle.BlockRound = int(float32(startStyle.BlockRound) + vr*elapsed)
	// 	c.currentStyle.BlockSpace = int(float32(startStyle.BlockSpace) + vs*elapsed)
	// 	c.currentStyle.BlockSize = int(float32(startStyle.BlockSize) + vsize*elapsed)
	// 	c.updateCache()
	// 	return true
	// })
}

func (c *Candidate) GetStatus() (f32.Point, blocks.Shape, *styles.Block) {
	return c.presentPos, c.shape, c.theme
}

func (c *Candidate) GetShape() blocks.Shape {
	return c.shape
}

func (c *Candidate) IsConsumed() bool {
	return c.consumed
}

func (c *Candidate) Consume() {
	c.consumed = true
}

func (c *Candidate) IsChosen() bool {
	return c.chosen
}

func (c *Candidate) ToggleDrag(drag bool) {
	if drag {
		c.state = dragging
	} else {
		c.state = stable
	}
	c.currentStyle = c.styles[c.state]
	c.updateCache()
}

func (c *Candidate) Interest() event.Filter {
	return pointer.Filter{
		Target: c,
		Kinds:  pointer.Press,
	}
}

func (c *Candidate) OnEvent(ev event.Event) {
	x, ok := ev.(pointer.Event)
	if !ok {
		return
	}
	switch x.Kind {
	case pointer.Press:
		c.ToggleChosen(true)
		// ox, oy := c.GetInnerOffset()
		adjust := c.currentStyle

		blockSize := adjust.BlockSize
		cx, cy := c.GetCenterOffset()
		c.ToPos(f32.Point{X: c.presentPos.X + x.Position.X - float32(cx) - float32(blockSize),
			Y: c.presentPos.Y + x.Position.Y - float32(cy) - float32(blockSize)})
	}
	println(c.shape.Desc())
}

func (c *Candidate) UpdatePosition(p f32.Point) {
	c.presentPos = p
}

func (c *Candidate) BackToDefault() {
	// c.presentPos = c.defaultPos
	c.ToPos(c.defaultPos)
}

func (c *Candidate) ToggleChosen(chosen bool) {
	c.chosen = chosen
	c.updateCache()
}

func (c *Candidate) GetMaxWidth() int {
	return c.currentStyle.BlockSize*5 + 4*c.currentStyle.BlockSpace
}

func (c *Candidate) GetInnerOffset() (int, int) {
	maxWidth := c.GetMaxWidth()
	leftOffset := (maxWidth - c.GetWidth()) / 2
	topOffset := (maxWidth - c.GetHeight()) / 2
	return leftOffset, topOffset
}

// offset to top-left corner point
func (c *Candidate) GetCenterOffset() (int, int) {
	leftOffset := c.GetWidth() / 2
	topOffset := c.GetHeight() / 2
	return leftOffset, topOffset
}

func (c *Candidate) GetWidth() int {
	space := c.currentStyle.BlockSpace
	return len(c.shape[0])*(space+c.currentStyle.BlockSize) - space
}

func (c *Candidate) GetHeight() int {
	space := c.currentStyle.BlockSpace
	return len(c.shape)*(space+c.currentStyle.BlockSize) - space
}

func (c *Candidate) renderCache(ops *op.Ops) {
	c.cache.Add(ops)
}

func (c *Candidate) updateCache() {
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

func (c *Candidate) Render(ctx *framework.Context) {
	ops := ctx.Ops
	adjust := c.currentStyle

	blockSize := adjust.BlockSize

	pos := c.presentPos
	left, top := int(pos.X), int(pos.Y)
	defer op.Offset(image.Pt(left-blockSize, top-blockSize)).Push(ops).Pop() // outer offset

	// switch c.state {
	// case movingToCursor:
	// 	// ctx.DeltaT

	// }
	transformed := false
	finished := []int{}
	for idx, trans := range c.transitions {
		if trans() {
			transformed = true
		} else {
			finished = append(finished, idx)
		}
	}
	for i := len(finished) - 1; i >= 0; i-- {
		c.transitions = append(c.transitions[:finished[i]], c.transitions[finished[i]+1:]...)
	}
	if transformed {
		println("refreshing")
		ctx.Refresh()
	}

	c.renderCache(ops)
}

type CandidateGroup struct {
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

func NewCandidateGroup(maxWidth, offsetTop, level, count int, style *styles.Style, rnd *rand.Rand) *CandidateGroup {
	minGap := style.BlockSize / 2

	candidateWidth := (maxWidth - 2*minGap) / 3
	candidateSize := (candidateWidth - 4*style.BlockSpace) / 5

	return &CandidateGroup{
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

func (cg *CandidateGroup) GetOffset() int {
	return cg.width + cg.gap
}

func (cg *CandidateGroup) GenerateCandidates() []*Candidate {
	result := make([]*Candidate, cg.count)

	candidateOffset := cg.GetOffset()
	offsetTop := cg.offsetTop + cg.style.BlockSize/3
	theme := cg.style.CurrentTheme

	for i := 0; i < cg.count; i++ {
		shapeIdx, shape := cg.shapeGroups.ChooseOneShape(0.3)
		result[i] = NewCandidate(shape,
			image.Point{X: i * candidateOffset, Y: offsetTop},
			cg.calm, cg.drag, theme.Shapes[shapeIdx])
	}
	return result
}
