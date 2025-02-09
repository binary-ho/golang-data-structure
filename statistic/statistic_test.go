package statistic_test

import (
	"fmt"
	"go-data-structure/statistic"
	"testing"

	"github.com/stretchr/testify/assert"
)

type stringValue struct {
	content string
	removed bool
}

func newStatisticNode(content string) *statistic.Node[*stringValue] {
	return statistic.NewNode(&stringValue{
		content: content,
	})
}

func (v *stringValue) Len() int {
	if v.removed {
		return 0
	}
	return len(v.content)
}

func (v *stringValue) String() string {
	return v.content
}

func TestLLRBTree(t *testing.T) {
	t.Run("insert and find test", func(t *testing.T) {
		// create empty LLRB
		tree := statistic.NewTree[*stringValue](nil)

		// initially empty, Find(0) => no node
		node, idx, err := tree.Find(0)
		assert.Nil(t, node)
		assert.NoError(t, err)
		assert.Equal(t, 0, idx)

		// Insert "A2"
		nodeA := tree.Insert(newStatisticNode("A2"))
		fmt.Println(tree.Print())
		// Insert "B23"
		nodeB := tree.Insert(newStatisticNode("B23"))
		fmt.Println(tree.Print())
		// Insert "C234"
		nodeC := tree.Insert(newStatisticNode("C234"))
		fmt.Println(tree.Print())
		// Insert "D2345"
		nodeD := tree.Insert(newStatisticNode("D2345"))
		fmt.Println(tree.Print())

		// Check their index
		assert.Equal(t, 0, tree.IndexOf(nodeA))
		assert.Equal(t, 2, tree.IndexOf(nodeB))
		assert.Equal(t, 5, tree.IndexOf(nodeC))
		assert.Equal(t, 9, tree.IndexOf(nodeD))

		// Try Find
		node, offset, err := tree.Find(1)
		assert.Equal(t, nodeA, node)
		assert.Equal(t, 1, offset)
		assert.NoError(t, err)

		node, offset, err = tree.Find(7)
		assert.Equal(t, nodeC, node)
		assert.Equal(t, 2, offset)
		assert.NoError(t, err)

		node, offset, err = tree.Find(11)
		assert.Equal(t, nodeD, node)
		assert.Equal(t, 2, offset)
		assert.NoError(t, err)
	})

	t.Run("deletion test", func(t *testing.T) {
		tree := statistic.NewTree[*stringValue](nil)

		nodeH := tree.Insert(newStatisticNode("H"))
		fmt.Println(tree.ToTestString())
		assert.Equal(t, 1, tree.Len())

		nodeE := tree.Insert(newStatisticNode("E"))
		fmt.Println(tree.ToTestString())
		assert.Equal(t, 2, tree.Len())

		nodeL := tree.Insert(newStatisticNode("LL"))
		fmt.Println(tree.ToTestString())
		assert.Equal(t, 4, tree.Len()) // "LL" has length 2 => total=1+1+2=4

		nodeO := tree.Insert(newStatisticNode("O"))
		fmt.Println(tree.ToTestString())
		assert.Equal(t, 5, tree.Len()) // +1 => total 5

		// delete E
		fmt.Println("Deleting E " + tree.ToTestString())
		fmt.Println(tree.Print())
		tree.Delete(nodeE)
		fmt.Println("After Delete E" + tree.ToTestString())
		fmt.Println(tree.Print())
		// check len after delete
		// we had 5 total, removing E => 4 total
		assert.Equal(t, 4, tree.Len())

		// check indexOf
		assert.Equal(t, 0, tree.IndexOf(nodeH))
		assert.Equal(t, -1, tree.IndexOf(nodeE)) // deleted
		// nodeL => "LL" => length=2 => presumably at index 1..2 or so
		idxL := tree.IndexOf(nodeL)
		assert.True(t, idxL >= 1, "L's index should be 1 or 2")
		nodeLPos := tree.IndexOf(nodeL)
		t.Logf("nodeL index is %d", nodeLPos)

		assert.Equal(t, 3, tree.IndexOf(nodeO)) // maybe index 3 if len so far is 4
	})

	t.Run("single node index test", func(t *testing.T) {
		tree := statistic.NewTree[*stringValue](nil)
		node := tree.Insert(newStatisticNode("A"))
		assert.Equal(t, 0, tree.IndexOf(node))
		tree.Delete(node)
		assert.Equal(t, -1, tree.IndexOf(node))
		assert.Equal(t, 0, tree.Len())
	})
}

// Additional helper sample code if needed
func makeSampleTree() (*statistic.Tree[*stringValue], []*statistic.Node[*stringValue]) {
	tree := statistic.NewTree[*stringValue](nil)
	var nodes []*statistic.Node[*stringValue]

	nodes = append(nodes, tree.Insert(newStatisticNode("A")))
	nodes = append(nodes, tree.Insert(newStatisticNode("BB")))
	nodes = append(nodes, tree.Insert(newStatisticNode("CCC")))
	nodes = append(nodes, tree.Insert(newStatisticNode("DDDD")))
	nodes = append(nodes, tree.Insert(newStatisticNode("EEEEE")))
	nodes = append(nodes, tree.Insert(newStatisticNode("FFFF")))
	nodes = append(nodes, tree.Insert(newStatisticNode("GGG")))
	nodes = append(nodes, tree.Insert(newStatisticNode("HH")))
	nodes = append(nodes, tree.Insert(newStatisticNode("I")))

	return tree, nodes
}
