package random

import (
	"math"
	"math/rand"
)

type Table struct {
	maxHeight   int8
	randomTable *[]uint32
}

func NewTable(maxHeight int8) *Table {
	maxHeight = min(maxHeight, MaxHeight)
	table := make([]uint32, maxHeight+1)
	value := uint32(math.MaxUint32)
	for height := int8(0); height <= maxHeight; height++ {
		table[height] = value
		value >>= 1
	}
	return &Table{maxHeight, &table}
}

func (t *Table) Height() int {
	uint32Number := rand.Uint32()
	table := *t.randomTable
	height := int8(0)
	for height <= t.maxHeight && uint32Number <= table[height] {
		height++
	}
	return int(height)
}
