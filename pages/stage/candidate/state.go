package candidate

type State int

const (
	stable State = iota
	dragging
	consumed
)
