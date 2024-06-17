package candidate

type State int

const (
	stable State = iota
	movingToCursor
	dragging
	consumed
)
