package stage

import (
	"fmt"
	"image"
	"math/rand"

	"gioui.org/f32"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/op"
	"gioui.org/op/clip"
	"github.com/hellflame/matrix-ranger/blocks"
	"github.com/hellflame/matrix-ranger/framework"
	"github.com/hellflame/matrix-ranger/styles"
)

type Stage struct {
	arena *blocks.Arena

	candidates []*candidate

	watcher *framework.Watcher

	rnd *rand.Rand

	style *styles.Style

	candidateGroup *candidateGroup

	shapeGroups *blocks.ShapeGroups
}

func NewStage(s *styles.Style, seed int64) *Stage {
	rnd := rand.New(rand.NewSource(seed))
	arena := blocks.NewArena(s)
	size := arena.GetAreaSize()
	watcher := new(framework.Watcher)
	stage := &Stage{
		style: s, rnd: rnd, arena: arena,
		watcher: watcher,

		candidateGroup: NewCandidateGroup(size, size, 0, 3, s, rnd),
		shapeGroups:    blocks.NewShapeGroups(rnd),
	}
	watcher.Add(stage)
	stage.GenerateCandidates()
	return stage
}

func (s *Stage) OnCreate() {
}

func (s *Stage) Interest() event.Filter {
	return pointer.Filter{
		Target: s,
		Kinds:  pointer.Drag | pointer.Release,
	}
}

func (s *Stage) OnEvent(ev event.Event) {
	x, ok := ev.(pointer.Event)
	if !ok {
		return
	}
	switch x.Kind {
	case pointer.Release:
		for _, c := range s.candidates {
			if c.chosen {
				c.ToggleDrag(false)
				// left, top := c.GetCenterOffset()
				// centerPosition := f32.Point{X: c.presentPos.X + float32(left), Y: c.presentPos.Y + float32(top)}
				fmt.Println("release pointer pos:", x.Position, "top left pos:", c.presentPos)
				if s.arena.Place(c.presentPos, c.shape, c.theme) {
					c.consumed = true
					if s.arena.CheckErasable() {
						s.arena.Erase()
					}
					// fmt.Println("erasable:", s.arena.CheckErasable())
				} else {
					c.presentPos = c.defaultPos
				}
				c.chosen = false
			}
		}
	case pointer.Drag:
		for _, c := range s.candidates {
			if c.chosen {
				c.ToggleDrag(true)
				left, top := c.GetCenterOffset()
				c.presentPos = f32.Point{X: x.Position.X - float32(left), Y: x.Position.Y - float32(top)}
			}
		}
	}
}

func (s *Stage) Render(ctx *framework.Context) {
	s.watcher.Trigger(ctx.Event.Source)
	ops := ctx.Ops
	w, h := s.GetSize()
	defer op.Offset(image.Pt(s.style.OffsetLeft, s.style.OffsetTop)).Push(ops).Pop()
	area := clip.Rect(image.Rect(0, 0, w, h))
	defer area.Push(ops).Pop()
	event.Op(ops, s)

	s.arena.Render(ops)

	s.GenerateCandidatesIfNeed()
	for _, candidate := range s.candidates {
		if !candidate.consumed {
			candidate.Render(ops)
		}
	}
}

func (s *Stage) GenerateCandidates() {
	for _, c := range s.candidates {
		s.watcher.Remove(c)
	}
	s.candidates = s.candidateGroup.GenerateCandidates()
	for _, c := range s.candidates {
		s.watcher.Add(c)
	}
}

func (s *Stage) GenerateCandidatesIfNeed() {
	allConsumed := true
	for _, c := range s.candidates {
		if !c.consumed {
			allConsumed = false
			break
		}
	}
	if allConsumed {
		s.GenerateCandidates()
	}
}

func (s *Stage) GetSize() (w, h int) {
	areaSize := s.arena.GetAreaSize()
	return areaSize, areaSize + s.candidateGroup.gap + s.candidateGroup.width
}
