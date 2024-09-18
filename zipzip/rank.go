package zipzip

import (
	"math"
	"math/rand"
)

type rank struct {
	rank1, rank2 int
}

func (tree *Tree) getRank() rank {
	return rank{tree.randRank1(), tree.randRank2()}
}

func (tree *Tree) randRank1() int {
	return int(math.Floor(math.Log(1-rand.Float64()) / math.Log(0.5)))
}

func (tree *Tree) randRank2() int {
	if tree.size == 0 {
		return 1
	}
	maxRank2 := int(math.Ceil(math.Pow(math.Log(float64(tree.size)), 3)))
	return rand.Intn(maxRank2) + 1
}
