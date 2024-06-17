package styles

type Animation struct {
	Duration float64
}

type Transition func() bool
