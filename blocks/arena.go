package blocks

import (
	"image"

	"gioui.org/f32"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"github.com/hellflame/matrix-ranger/styles"
)

type Arena struct {
	style *styles.Style

	active bool

	bricks [][]brick
}

// Create an arena to place piles of bricks
func NewArena(s *styles.Style) *Arena {
	size := s.BlockCount
	defaultBlockStyle := s.CurrentTheme.Block
	bricks := make([][]brick, size)
	for i := 0; i < size; i++ {
		tmp := make([]brick, size)
		for j := 0; j < size; j++ {
			tmp[j] = brick{
				solid: false,
				style: *defaultBlockStyle,
			}
		}
		bricks[i] = tmp
	}
	return &Arena{bricks: bricks, style: s}
}

func (a *Arena) Place(start f32.Point, shape Shape, theme *styles.Block) bool {
	blockSize := float32(a.style.BlockSize)
	blockCount := a.style.BlockCount
	colIdx := int(start.X / blockSize)
	rowIdx := int(start.Y / blockSize)
	shapeRows := len(shape)
	shapeCols := len(shape[0])
	if colIdx > blockCount || colIdx+shapeCols > blockCount ||
		rowIdx > blockCount || rowIdx+shapeRows > blockCount {
		return false
	}
	for r, shapeR := rowIdx, 0; r < blockCount && shapeR < shapeRows; r, shapeR = r+1, shapeR+1 {
		for c, shapeC := colIdx, 0; c < blockCount && shapeC < shapeCols; c, shapeC = c+1, shapeC+1 {
			if a.bricks[r][c].solid && shape[shapeR][shapeC] {
				return false
			}
		}
	}

	for r, shapeR := rowIdx, 0; r < blockCount && shapeR < shapeRows; r, shapeR = r+1, shapeR+1 {
		for c, shapeC := colIdx, 0; c < blockCount && shapeC < shapeCols; c, shapeC = c+1, shapeC+1 {
			if shape[shapeR][shapeC] {
				a.bricks[r][c].solid = true
				a.bricks[r][c].style = *theme
			}
		}
	}

	return true
}

func (a *Arena) OnEvent(e event.Event) {
	x, ok := e.(pointer.Event)
	if !ok {
		return
	}
	switch x.Kind {
	case pointer.Enter:
		a.active = true
		// println("arena Enter")
	case pointer.Leave:
		a.active = false
		// println("arena Leave")
	case pointer.Move:
		// println("arena pointer:", x.Position.X, x.Position.Y)
	}
}

func (a *Arena) Render(ops *op.Ops) {
	space := a.style.BlockSpace
	blockCount := a.style.BlockCount
	blockSize := a.style.BlockSize
	size := blockCount*(space+blockSize) - space
	area := clip.Rect(image.Rect(0, 0, size, size))
	defer area.Push(ops).Pop()
	event.Op(ops, a)
	round := a.style.BlockRound
	bricks := a.bricks
	for r := 0; r < blockCount; r++ {
		rowOffset := op.Offset(image.Pt(0, r*(space+blockSize))).Push(ops)
		for c := 0; c < blockCount; c++ {
			brick := bricks[r][c]
			colOffset := op.Offset(image.Pt(c*(space+blockSize), 0)).Push(ops)

			b := clip.RRect{Rect: image.Rect(0, 0, blockSize, blockSize),
				SE: round, SW: round, NW: round, NE: round}.Push(ops)

			paint.ColorOp{Color: brick.style.Color}.Add(ops)
			paint.PaintOp{}.Add(ops)

			b.Pop()
			colOffset.Pop()
		}
		rowOffset.Pop()
	}
}

func (a *Arena) GetAreaSize() int {
	return (a.style.BlockSize+a.style.BlockSpace)*a.style.BlockCount - a.style.BlockSpace
}
