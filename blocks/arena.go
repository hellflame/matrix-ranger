package blocks

import (
	"image"

	"gioui.org/f32"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"github.com/hellflame/matrix-ranger/framework"
	"github.com/hellflame/matrix-ranger/styles"
)

type Arena struct {
	style *styles.Style

	active bool
	cache  op.CallOp

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
	arena := &Arena{bricks: bricks, style: s}
	arena.updateCache()
	return arena
}

func (a *Arena) Reset() {
	blockCount := a.style.BlockCount
	for r := 0; r < blockCount; r++ {
		for c := 0; c < blockCount; c++ {
			if a.bricks[r][c].solid {
				a.bricks[r][c].solid = false
				a.bricks[r][c].style = *a.style.CurrentTheme.Block
			}
		}
	}
}

func (a *Arena) renderCache(ops *op.Ops) {
	a.cache.Add(ops)
}

func (a *Arena) updateCache() {
	ops := new(op.Ops)
	macro := op.Record(ops)
	defer func() {
		a.cache = macro.Stop()
	}()

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

func (a *Arena) Erase() {
	defer a.updateCache()
	rowIndices := []int{}
	colIndices := []int{}

	blockCount := a.style.BlockCount
	for r := 0; r < blockCount; r++ {
		rowIsFull := true
		for c := 0; c < blockCount; c++ {
			if !a.bricks[r][c].solid {
				rowIsFull = false
				break
			}
		}
		if rowIsFull {
			rowIndices = append(rowIndices, r)
		}
	}

	for c := 0; c < blockCount; c++ {
		colIsFull := true
		for r := 0; r < blockCount; r++ {
			if !a.bricks[r][c].solid {
				colIsFull = false
				break
			}
		}
		if colIsFull {
			colIndices = append(colIndices, c)
		}
	}

	for _, idx := range rowIndices {
		for c := 0; c < blockCount; c++ {
			a.bricks[idx][c].solid = false
			a.bricks[idx][c].style = *a.style.CurrentTheme.Block
		}
	}

	for _, idx := range colIndices {
		for r := 0; r < blockCount; r++ {
			a.bricks[r][idx].solid = false
			a.bricks[r][idx].style = *a.style.CurrentTheme.Block
		}
	}

}

func (a *Arena) CheckErasable() bool {
	blockCount := a.style.BlockCount
	for r := 0; r < blockCount; r++ {
		rowIsFull := true
		for c := 0; c < blockCount; c++ {
			if !a.bricks[r][c].solid {
				rowIsFull = false
				break
			}
		}
		if rowIsFull {
			return true
		}
	}

	for c := 0; c < blockCount; c++ {
		colIsFull := true
		for r := 0; r < blockCount; r++ {
			if !a.bricks[r][c].solid {
				colIsFull = false
				break
			}
		}
		if colIsFull {
			return true
		}
	}

	return false
}

func (a *Arena) HasPlace(rowIdx, colIdx int, shape Shape) bool {
	blockCount := a.style.BlockCount
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
	return true
}

func (a *Arena) FindPlace(shape Shape) (int, int, bool) {
	blockCount := a.style.BlockCount
	for r := 0; r < blockCount; r++ {
		for c := 0; c < blockCount; c++ {
			if !a.bricks[r][c].solid && a.HasPlace(r, c, shape) {
				return r, c, true
			}
		}
	}
	return 0, 0, false
}

func (a *Arena) Place(start f32.Point, shape Shape, theme *styles.Block) bool {
	defer a.updateCache()
	blockSize := float32(a.style.BlockSize)
	blockCount := a.style.BlockCount
	colIdx := int(start.X / blockSize)
	rowIdx := int(start.Y / blockSize)
	shapeRows := len(shape)
	shapeCols := len(shape[0])

	if !a.HasPlace(rowIdx, colIdx, shape) {
		return false
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

func (a *Arena) Render(ctx *framework.Context) {
	a.renderCache(ctx.Ops)
}

func (a *Arena) GetAreaSize() int {
	println("arena block size", a.style.BlockSize)
	return (a.style.BlockSize+a.style.BlockSpace)*a.style.BlockCount - a.style.BlockSpace
}
