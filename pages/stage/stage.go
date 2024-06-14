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
	"github.com/hellflame/matrix-ranger/pages/stage/candidate"
	"github.com/hellflame/matrix-ranger/styles"
)

type Stage struct {
	arena *blocks.Arena

	candidates []*candidate.Candidate

	watcher *framework.Watcher

	rnd *rand.Rand

	style *styles.Style

	candidateGroup *candidate.CandidateGroup
}

func NewStage(s *styles.Style, seed int64) *Stage {
	rnd := rand.New(rand.NewSource(seed))

	watcher := new(framework.Watcher)
	stage := &Stage{
		style: s, rnd: rnd,
		watcher: watcher,
	}
	watcher.Add(stage)
	return stage
}

func (s *Stage) OnCreate(ctx *framework.Context) {
	println("px per dp", ctx.Event.Metric.PxPerDp)
	// s.style.BlockSize = int(float32(s.style.BlockSize) * (ctx.Event.Metric.PxPerDp / 4))
	println("block size", s.style.BlockSize)
	s.arena = blocks.NewArena(s.style)

	size := s.arena.GetAreaSize()
	s.candidateGroup = candidate.NewCandidateGroup(size, size, 0, 3, s.style, s.rnd)
	s.GenerateCandidates()
}

func (s *Stage) Interest() event.Filter {
	return pointer.Filter{
		Target: s,
		Kinds:  pointer.Drag | pointer.Release, // | pointer.Press,
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
			if c.IsChosen() {
				c.ToggleDrag(false)
				// left, top := c.GetCenterOffset()
				// centerPosition := f32.Point{X: c.presentPos.X + float32(left), Y: c.presentPos.Y + float32(top)}
				// fmt.Println("release pointer pos:", x.Position, "top left pos:", c.presentPos)
				p, _, _ := c.GetStatus()
				fmt.Println("release point", p)
				if s.arena.Place(c.GetStatus()) {
					c.Consume()
					if s.arena.CheckErasable() {
						s.arena.Erase()
					}
					if !s.checkAnyPlaceable() {
						println("game over")
						s.reset()
					}
				} else {
					c.BackToDefault()
				}
				c.ToggleChosen(false)
			}
		}
	case pointer.Drag:
		for _, c := range s.candidates {
			if c.IsChosen() {
				c.ToggleDrag(true)
				innerOffset := s.style.StageInnerOffset
				left, top := c.GetCenterOffset()
				c.UpdatePosition(f32.Point{X: x.Position.X - float32(left+innerOffset), Y: x.Position.Y - float32(top+innerOffset)})
			}
		}
	}
}

func (s *Stage) checkAnyPlaceable() bool {
	placeable := false
	hasCandidate := false
	for _, c := range s.candidates {
		if !c.IsConsumed() {
			hasCandidate = true
			if r, col, ok := s.arena.FindPlace(c.GetShape()); ok {
				placeable = true
				println("find place", r, col)
				println(c.GetShape().Desc())
			}
		}
	}
	if hasCandidate {
		return placeable
	}
	return true
}

func (s *Stage) reset() {
	println("reset")
	s.arena.Reset()
	for _, c := range s.candidates {
		c.Consume()
	}
}

func (s *Stage) Render(ctx *framework.Context) {
	s.watcher.Trigger(ctx.Event.Source)
	ops := ctx.Ops

	innerOffset := s.style.StageInnerOffset
	defer op.Offset(image.Pt(s.style.OffsetLeft-innerOffset, s.style.OffsetTop-innerOffset)).Push(ops).Pop()

	w, h := s.GetSize()
	area := clip.Rect(image.Rect(0, 0, w+innerOffset*2, h+innerOffset*2))
	defer area.Push(ops).Pop()
	event.Op(ops, s)

	defer op.Offset(image.Pt(innerOffset, innerOffset)).Push(ops).Pop()

	s.GenerateCandidatesIfNeed()
	s.arena.Render(ctx)

	for _, candidate := range s.candidates {
		if !candidate.IsConsumed() {
			candidate.Render(ctx)
		}
	}
	println("after stage render")
}

func (s *Stage) GenerateCandidates() {
	for _, c := range s.candidates {
		s.watcher.Remove(c)
	}
	s.candidates = s.candidateGroup.GenerateCandidates()
	for _, c := range s.candidates {
		s.watcher.Add(c)
	}

	if !s.checkAnyPlaceable() {
		println("game over after generation")
		s.reset()
	}
}

func (s *Stage) GenerateCandidatesIfNeed() {
	allConsumed := true
	for _, c := range s.candidates {
		if !c.IsConsumed() {
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
	return areaSize, areaSize + s.candidateGroup.GetOffset()
}
