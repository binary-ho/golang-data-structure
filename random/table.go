package random

import (
	"math"
	"math/rand"
)

type Table struct {
	maxRank     int8
	randomTable *[]uint32
}

func NewTable(maxRank int8) *Table {
	maxRank = min(maxRank, MaxRank)
	table := make([]uint32, maxRank+1)
	value := uint32(math.MaxUint32)
	for rank := int8(0); rank <= maxRank; rank++ {
		table[rank] = value
		value >>= 1
	}
	return &Table{maxRank, &table}
}

func (t *Table) Rank() int {
	uint32Number := rand.Uint32()
	table := *t.randomTable
	rank := int8(0)
	for rank <= t.maxRank && uint32Number <= table[rank] {
		rank++
	}
	return int(rank)
}
