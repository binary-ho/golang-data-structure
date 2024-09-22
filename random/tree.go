package random

import (
	yorkiellrb "go-data-structure/yorkie/llrb"
	"math"
	"math/rand"
	"strconv"
)

type Tree struct {
	tree *yorkiellrb.Tree[*Key, *Value]
}

func NewUint32Tree(maxRank int) *Tree {
	maxRank = min(maxRank, MaxRank)
	tree := &Tree{tree: yorkiellrb.NewTree[*Key, *Value]()}
	key := math.MaxUint32
	for rank := 0; rank <= maxRank; rank++ {
		tree.put(key, rank)
		key >>= 1
	}
	tree.put(0, maxRank)
	return tree
}

func (tree *Tree) put(key, value int) {
	treeKey := Key(key)
	treeValue := Value(value)
	tree.tree.Put(&treeKey, &treeValue)
}

func (tree *Tree) Rank() int {
	uint32Number := rand.Uint32()
	key := Key(uint32Number)
	_, value := tree.tree.Floor(&key)
	return int(*value)
}

type Key uint32

func (key Key) Compare(other yorkiellrb.Key) int {
	otherKey := uint32(*(other.(*Key)))
	uint32Key := uint32(key)
	if uint32Key > otherKey {
		return 1
	} else if uint32Key < otherKey {
		return -1
	}
	return 0
}

type Value uint8

func (value *Value) String() string {
	return strconv.Itoa(int(*value))
}
