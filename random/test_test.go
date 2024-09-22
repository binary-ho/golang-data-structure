package random

import "testing"

func Test(t *testing.T) {
	table := NewTable(16)
	println(table)

	tree := NewUint32Tree(16)
	println(tree)
}
