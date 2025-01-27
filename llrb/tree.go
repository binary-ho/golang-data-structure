package llrb

import (
	"fmt"
	"strings"
)

// ErrOutOfIndex is returned when the given index is out of index.
var ErrOutOfIndex = fmt.Errorf("out of index")

// Tree is an implementation of Left-learning Red-Black Tree.
// Original paper on Left-leaning Red-Black Trees:
// http://www.cs.princeton.edu/~rs/talks/LLRB/LLRB.pdf
//
// Invariant 1: No red node has a red child
// Invariant 2: Every leaf path has the same number of black nodes
// Invariant 3: Only the left child can be red (left leaning)
type Tree[K Key, V Value] struct {
	root *Node[K, V]
	size int
}

// NewTree creates a new instance of Tree.
func NewTree[K Key, V Value]() *Tree[K, V] {
	return &Tree[K, V]{}
}

// Len returns the length of the tree.
func (tree *Tree[K, V]) Len() int {
	return tree.size
}

func rotateLeft[K Key, V Value](node *Node[K, V]) *Node[K, V] {
	right := node.right
	node.right = right.left
	right.left = node
	right.isRed = right.left.isRed
	right.left.isRed = true

	// 논문 rotateLeft 참고, 하위 노드 갯수가 바뀌는 두 노드 Update
	node.UpdateTreeWeight()
	right.UpdateTreeWeight()
	return right
}

func rotateRight[K Key, V Value](node *Node[K, V]) *Node[K, V] {
	left := node.left
	node.left = left.right
	left.right = node
	left.isRed = left.right.isRed
	left.right.isRed = true

	// 논문 rotateRight 참고, 하위 노드 갯수가 바뀌는 두 노드 Update
	node.UpdateTreeWeight()
	left.UpdateTreeWeight()
	return left
}

func flipColors[K Key, V Value](node *Node[K, V]) {
	node.isRed = !node.isRed
	node.left.isRed = !node.left.isRed
	node.right.isRed = !node.right.isRed
}

func moveRedLeft[K Key, V Value](node *Node[K, V]) *Node[K, V] {
	flipColors(node)
	if isRed(node.right.left) {
		node.right = rotateRight(node.right)
		node = rotateLeft(node)
		flipColors(node)
	}
	return node
}

func moveRedRight[K Key, V Value](node *Node[K, V]) *Node[K, V] {
	flipColors(node)
	if isRed(node.left.left) {
		node = rotateRight(node)
		flipColors(node)
	}
	return node
}

func removeMin[K Key, V Value](node *Node[K, V]) *Node[K, V] {
	if node.left == nil {
		return nil
	}

	if isNotRed(node.left) && isNotRed(node.left.left) {
		node = moveRedLeft(node)
	}

	node.left = removeMin(node.left)
	return fixUp(node)

}

func min[K Key, V Value](node *Node[K, V]) *Node[K, V] {
	if node.left == nil {
		return node
	}

	return min(node.left)
}

func fixUp[K Key, V Value](node *Node[K, V]) *Node[K, V] {
	if isRed(node.right) {
		node = rotateLeft(node)
	}

	if isRed(node.left) && isRed(node.left.left) {
		node = rotateRight(node)
	}

	if isRed(node.left) && isRed(node.right) {
		flipColors(node)
	}

	return node
}

func isRed[K Key, V Value](node *Node[K, V]) bool {
	return node != nil && node.isRed
}

func isNotRed[K Key, V Value](node *Node[K, V]) bool {
	return !isRed(node)
}

func (tree *Tree[K, V]) String() string {
	var str []string
	traverseInOrder(tree.root, func(node *Node[K, V]) {
		str = append(str, node.value.String())
	})
	return strings.Join(str, ",")
}

func traverseInOrder[K Key, V Value](node *Node[K, V], callback func(node *Node[K, V])) {
	if node == nil {
		return
	}

	traverseInOrder(node.left, callback)
	callback(node)
	traverseInOrder(node.right, callback)
}
