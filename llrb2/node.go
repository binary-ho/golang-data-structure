package llrb2

import (
	"fmt"
	"strings"
)

// Key represents key of Tree.
type Key interface {
	// Compare gives the result of a 3-way comparison
	// a.Compare(b) = 1 => a > b
	// a.Compare(b) = 0 => a == b
	// a.Compare(b) = -1 => a < b
	Compare(k Key) int
}

// Value represents the data stored in the nodes of Tree.
type Value interface {
	String() string
}

// Node is a node of Tree.
type Node[K Key, V Value] struct {
	key   K
	value V

	// parent pointer is present but not strictly maintained
	// in all rotation steps in this demo code.
	parent *Node[K, V]

	left  *Node[K, V]
	right *Node[K, V]
	isRed bool

	// size is the total number of nodes in this subtree (including self)
	size int
}

// NewNode creates a new instance of Node.
func NewNode[K Key, V Value](key K, value V, isRed bool) *Node[K, V] {
	return &Node[K, V]{
		key:   key,
		value: value,
		isRed: isRed,
		size:  1, // new node => size=1
	}
}

// Tree is an implementation of Left-learning Red-Black Tree.
type Tree[K Key, V Value] struct {
	root *Node[K, V]
	size int // total number of nodes
}

// NewTree creates a new instance of Tree.
func NewTree[K Key, V Value]() *Tree[K, V] {
	return &Tree[K, V]{}
}

func nodeSize[K Key, V Value](n *Node[K, V]) int {
	if n == nil {
		return 0
	}
	return n.size
}

func updateSize[K Key, V Value](n *Node[K, V]) {
	if n == nil {
		return
	}
	n.size = nodeSize(n.left) + nodeSize(n.right) + 1
}

func rotateLeft[K Key, V Value](node *Node[K, V]) *Node[K, V] {
	right := node.right
	node.right = right.left
	right.left = node

	// 색상 계승
	right.isRed = node.isRed
	node.isRed = true

	// size 갱신
	updateSize(node)
	updateSize(right)

	return right
}

func rotateRight[K Key, V Value](node *Node[K, V]) *Node[K, V] {
	left := node.left
	node.left = left.right
	left.right = node

	left.isRed = node.isRed
	node.isRed = true

	updateSize(node)
	updateSize(left)

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
	updateSize(node) // 마지막에 size 갱신
	return node
}

func isRed[K Key, V Value](node *Node[K, V]) bool {
	return node != nil && node.isRed
}

// Put puts the value of the given key.
func (t *Tree[K, V]) Put(k K, v V) *Node[K, V] {
	newNode := NewNode(k, v, true)
	t.root = t.put(t.root, newNode)
	t.root.isRed = false
	return newNode
}

func (t *Tree[K, V]) put(node *Node[K, V], newNode *Node[K, V]) *Node[K, V] {
	if node == nil {
		t.size++
		return newNode
	}

	compare := newNode.key.Compare(node.key)
	if compare < 0 {
		node.left = t.put(node.left, newNode)
	} else if compare > 0 {
		node.right = t.put(node.right, newNode)
	} else {
		// same key -> update value
		node.value = newNode.value
	}

	node = fixUp(node)
	return node
}

// Remove removes the value of the given key.
func (t *Tree[K, V]) Remove(key K) {
	if t.root == nil {
		return
	}

	if !isRed(t.root.left) && !isRed(t.root.right) {
		t.root.isRed = true
	}

	t.root = t.remove(t.root, key)
	if t.root != nil {
		t.root.isRed = false
	}
}

func (t *Tree[K, V]) remove(node *Node[K, V], key K) *Node[K, V] {
	if node == nil {
		return nil
	}

	if key.Compare(node.key) < 0 {
		if !isRed(node.left) && !isRed(node.left.left) {
			node = moveRedLeft(node)
		}
		node.left = t.remove(node.left, key)
	} else {
		if isRed(node.left) {
			node = rotateRight(node)
		}
		if key.Compare(node.key) == 0 && node.right == nil {
			t.size--
			return nil
		}
		if !isRed(node.right) && !isRed(node.right.left) {
			node = moveRedRight(node)
		}
		if key.Compare(node.key) == 0 {
			t.size--
			smallest := min(node.right)
			node.key = smallest.key
			node.value = smallest.value
			node.right = removeMin(node.right)
		} else {
			node.right = t.remove(node.right, key)
		}
	}

	return fixUp(node)
}

func removeMin[K Key, V Value](node *Node[K, V]) *Node[K, V] {
	if node == nil {
		return nil
	}
	if node.left == nil {
		return nil
	}
	if !isRed(node.left) && !isRed(node.left.left) {
		node = moveRedLeft(node)
	}
	node.left = removeMin(node.left)
	return fixUp(node)
}

func min[K Key, V Value](node *Node[K, V]) *Node[K, V] {
	if node == nil {
		return nil
	}
	for node.left != nil {
		node = node.left
	}
	return node
}

// Size returns the total number of nodes in the tree
func (t *Tree[K, V]) Size() int {
	return t.size
}

// Root returns the root node of the tree
func (t *Tree[K, V]) Root() *NodeInfo[K, V] {
	if t.root == nil {
		return nil
	}
	return &NodeInfo[K, V]{
		Key:   t.root.key,
		Value: t.root.value,
		IsRed: t.root.isRed,
	}
}

// NodeInfo represents the public information of a node
type NodeInfo[K Key, V Value] struct {
	Key   K
	Value V
	IsRed bool
}

// String returns an inorder list of node values.
func (t *Tree[K, V]) String() string {
	var str []string
	traverseInOrder(t.root, func(node *Node[K, V]) {
		str = append(str, node.value.String())
	})
	return strings.Join(str, ",")
}

// StringTree prints a simple multiline view of the tree.
func (t *Tree[K, V]) StringTree() string {
	var sb strings.Builder
	printTree(t.root, 0, &sb)
	return sb.String()
}

// PrintPrettyTree returns a visually appealing string representation of the tree
func (t *Tree[K, V]) PrintPrettyTree() string {
	if t.root == nil {
		return "Empty tree"
	}

	var sb strings.Builder
	t.prettyPrintTree(t.root, "", true, &sb)
	return sb.String()
}

func (t *Tree[K, V]) prettyPrintTree(node *Node[K, V], prefix string, isLeft bool, sb *strings.Builder) {
	if node == nil {
		return
	}

	// 현재 노드의 색상과 값을 표시할 문자열 준비
	nodeColor := "⚫" // 검은색 노드
	if node.isRed {
		nodeColor = "🔴" // 빨간색 노드
	}

	// 오른쪽 서브트리 먼저 출력 (트리를 왼쪽으로 90도 회전한 형태로 출력)
	t.prettyPrintTree(node.right, prefix+(map[bool]string{true: "│   ", false: "    "})[isLeft], false, sb)

	// 현재 노드 출력
	sb.WriteString(prefix)
	sb.WriteString(map[bool]string{true: "└── ", false: "┌── "}[isLeft])
	sb.WriteString(fmt.Sprintf("%v %v: %v\n", nodeColor, node.key, node.value))

	// 왼쪽 서브트리 출력
	t.prettyPrintTree(node.left, prefix+(map[bool]string{true: "│   ", false: "    "})[isLeft], true, sb)
}

// printTree: pre-order traversal with indentation to show hierarchy.
func printTree[K Key, V Value](node *Node[K, V], level int, sb *strings.Builder) {
	if node == nil {
		return
	}
	indent := strings.Repeat("  ", level)
	color := "B"
	if node.isRed {
		color = "R"
	}
	fmt.Fprintf(sb, "%s- [%v, size=%d, %s]\n", indent, node.key, node.size, color)
	printTree(node.left, level+1, sb)
	printTree(node.right, level+1, sb)
}

func traverseInOrder[K Key, V Value](node *Node[K, V], callback func(node *Node[K, V])) {
	if node == nil {
		return
	}
	traverseInOrder(node.left, callback)
	callback(node)
	traverseInOrder(node.right, callback)
}
