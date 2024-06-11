package stage

import (
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
	arena := blocks.NewArena(s)
	size := arena.GetAreaSize()
	watcher := new(framework.Watcher)
	stage := &Stage{
		style: s, rnd: rnd, arena: arena,
		watcher: watcher,

		candidateGroup: candidate.NewCandidateGroup(size, size, 0, 3, s, rnd),
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
				if s.arena.Place(c.GetStatus()) {
					c.Consume()
					if s.arena.CheckErasable() {
						s.arena.Erase()
					}
					// fmt.Println("erasable:", s.arena.CheckErasable())
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

	s.arena.Render(ops)

	s.GenerateCandidatesIfNeed()
	for _, candidate := range s.candidates {
		if !candidate.IsConsumed() {
			candidate.Render(ctx)
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
