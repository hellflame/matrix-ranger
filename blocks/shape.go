package blocks

import (
	"math"
	"math/rand"
	"strings"
)

type Shape [][]bool

type ShapeGroup struct {
	Idx    int
	Level  int
	Shapes []Shape
}

func NewShapeGroup(shape Shape, rotates, level int) *ShapeGroup {
	return &ShapeGroup{
		Level:  level,
		Shapes: shape.GenerateRotateGroup(rotates),
	}
}

func (s Shape) Desc() string {
	result := []string{}
	for _, row := range s {
		tmp := []string{}
		for _, col := range row {
			if col {
				tmp = append(tmp, "*")
			} else {
				tmp = append(tmp, " ")
			}
		}
		result = append(result, strings.Join(tmp, " "))
	}
	return strings.Join(result, "\n")
}

func (s Shape) Rotate90() Shape {
	rows := len(s)
	cols := len(s[0])

	result := [][]bool{}

	for c := 0; c < cols; c++ {
		result = append(result, make([]bool, rows))
	}

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			result[c][r] = s[rows-r-1][c]
		}
	}

	return result
}

func (s Shape) GenerateRotateGroup(rotates int) []Shape {
	result := []Shape{s}
	for i, last := 1, s; i < rotates; i++ {
		last = last.Rotate90()
		result = append(result, last)
	}
	return result
}

type ShapeGroups struct {
	rnd *rand.Rand

	maxLevel int
	groupMap map[int][]*ShapeGroup
}

func NewShapeGroups(rnd *rand.Rand) *ShapeGroups {
	// level should fit in all along side digits
	groups := []*ShapeGroup{
		// *
		NewShapeGroup(Shape{{true}}, 1, 1),
		// * *
		NewShapeGroup(Shape{{true, true}}, 2, 2),
		// * * *
		NewShapeGroup(Shape{{true, true, true}}, 2, 3),
		// * * * *
		NewShapeGroup(Shape{{true, true, true, true}}, 2, 4),
		// * * * * *
		NewShapeGroup(Shape{{true, true, true, true, true}}, 2, 5),
		// *
		// * *
		NewShapeGroup(Shape{{true, false}, {true, true}}, 4, 3),
		// * *
		// * *
		NewShapeGroup(Shape{{true, true}, {true, true}}, 1, 2),
		// *
		// *
		// * * *
		NewShapeGroup(Shape{{true, false, false}, {true, false, false}, {true, true, true}}, 4, 4),
		// * * *
		// * * *
		// * * *
		NewShapeGroup(Shape{{true, true, true}, {true, true, true}, {true, true, true}}, 1, 5),
		// *
		// * *
		// * * *
		NewShapeGroup(Shape{{true, false, false}, {true, true, false}, {true, true, true}}, 4, 6),
		//   *
		// * * *
		//   *
		NewShapeGroup(Shape{{false, true, false}, {true, true, true}, {false, true, false}}, 1, 7),
	}
	maxLevel := 1
	groupMap := make(map[int][]*ShapeGroup)
	for idx, g := range groups {
		l := g.Level
		g.Idx = idx
		groupMap[l] = append(groupMap[l], g)
		if l > maxLevel {
			maxLevel = g.Level
		}
	}
	return &ShapeGroups{
		rnd: rnd, maxLevel: maxLevel, groupMap: groupMap,
	}
}

// hardness: 0 => Normal
// hardness: -1 => Easy
// hardness: 1 => Hard
func (sg *ShapeGroups) ChooseOneShape(hardness int) (int, Shape) {
	rnd := sg.rnd.Float64()

	switch hardness {
	case 1:
		// rnd ^ 0.5
		rnd = math.Sqrt(rnd)
	case -1:
		// rnd ^ 2
		rnd = rnd * rnd
	}
	// [1, maxLevel)
	level := int(math.Round(rnd*float64(sg.maxLevel-1) + 1))
	target := sg.groupMap[level]
	group := target[sg.rnd.Intn(len(target))]
	shapes := group.Shapes

	return group.Idx, shapes[sg.rnd.Intn(len(shapes))]
}
