package blocks

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestShape(t *testing.T) {
	for _, shape := range []Shape{
		[][]bool{{true, true}, {true, false}},
		[][]bool{{true, true, true, true}},
		[][]bool{
			{true, false, false},
			{true, true, false},
			{true, true, true},
		},
	} {
		println("origin:")
		println(shape.Desc())
		println("rotate:")
		println(shape.Rotate90().Desc())
		println("======")
	}
}

func TestRotateGroups(t *testing.T) {
	base := Shape([][]bool{
		{true, false},
		{true, true},
	})

	for _, shape := range base.GenerateRotateGroup(4) {
		println(shape.Desc())
	}
}

func TestNil(t *testing.T) {
	type x struct {
	}

	println(new(x))
}

func TestNorm(t *testing.T) {
	fmt.Println("rnd: ", rand.Float32())
}

func TestRemove(t *testing.T) {
	target := []int{1, 2, 3, 4, 5}
	for idx, v := range target {
		if v == 3 {
			target = append(target[:idx], target[idx+1:]...)
			break
		}
	}
	fmt.Println(target)
}
