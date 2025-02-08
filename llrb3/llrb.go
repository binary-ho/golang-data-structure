package llrb3

import (
	"errors"
	"fmt"
	"strings"
)

// Key, Value, Node, Tree are from the original code, but we'll add 'size' to Node.
type Key interface {
	Compare(k Key) int
}
type Value interface {
	String() string
}

// Node is extended with `size`.
type Node[K Key, V Value] struct {
	key    K
	value  V
	parent *Node[K, V]
	left   *Node[K, V]
	right  *Node[K, V]
	isRed  bool

	// size: number of nodes in this subtree (for Order Statistic).
	// We'll maintain this in put/remove/InsertAfter.
	size int
}

func NewNode[K Key, V Value](key K, value V, isRed bool) *Node[K, V] {
	return &Node[K, V]{
		key:   key,
		value: value,
		isRed: isRed,
		size:  1, // new node => size=1
	}
}

type Tree[K Key, V Value] struct {
	root *Node[K, V]
	size int // total node count
}

// NewTree creates a new LLRB Tree.
func NewTree[K Key, V Value]() *Tree[K, V] {
	return &Tree[K, V]{}
}

// Len() implements the same usage as in splay's "Len() int"
func (t *Tree[K, V]) Len() int {
	// We also maintain t.size (for quick access),
	// but let's trust root.size if we want the entire tree's node count.
	if t.root == nil {
		return 0
	}
	return t.root.size
}

// ========== Common helper ==========

func nodeSize[K Key, V Value](n *Node[K, V]) int {
	if n == nil {
		return 0
	}
	return n.size
}

// updateSize: recalc node.size = left.size + right.size + 1
func updateSize[K Key, V Value](n *Node[K, V]) {
	if n == nil {
		return
	}
	n.size = nodeSize(n.left) + nodeSize(n.right) + 1
}

// isRed: same as original
func isRed[K Key, V Value](node *Node[K, V]) bool {
	return node != nil && node.isRed
}

// rotateLeft with size updates
func rotateLeft[K Key, V Value](node *Node[K, V]) *Node[K, V] {
	right := node.right
	node.right = right.left
	right.left = node

	right.isRed = node.isRed
	node.isRed = true

	updateSize(node)
	updateSize(right)

	return right
}

// rotateRight with size updates
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

// fixUp: final re-balance after remove
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
	updateSize(node)
	return node
}

// traverseInOrder used in String, etc.
func traverseInOrder[K Key, V Value](node *Node[K, V], cb func(*Node[K, V])) {
	if node == nil {
		return
	}
	traverseInOrder(node.left, cb)
	cb(node)
	traverseInOrder(node.right, cb)
}

// ================== PART 1: normal BST Put/Remove by Key (createdAt) ==================

// Put inserts (k,v) by Key-based BST logic (like the original).
func (t *Tree[K, V]) Put(k K, v V) {
	t.root = t.putInternal(t.root, k, v)
	t.root.isRed = false
}

func (t *Tree[K, V]) putInternal(node *Node[K, V], key K, value V) *Node[K, V] {
	if node == nil {
		t.size++
		return NewNode(key, value, true)
	}

	cmp := key.Compare(node.key)
	if cmp < 0 {
		node.left = t.putInternal(node.left, key, value)
	} else if cmp > 0 {
		node.right = t.putInternal(node.right, key, value)
	} else {
		// update value
		node.value = value
	}

	// LLRB balancing steps inline (instead of fixUp)
	if isRed(node.right) && !isRed(node.left) {
		node = rotateLeft(node)
	}
	if isRed(node.left) && isRed(node.left.left) {
		node = rotateRight(node)
	}
	if isRed(node.left) && isRed(node.right) {
		flipColors(node)
	}

	updateSize(node)
	return node
}

// Remove by Key
func (t *Tree[K, V]) Remove(key K) {
	if t.root == nil {
		return
	}
	if !isRed(t.root.left) && !isRed(t.root.right) {
		t.root.isRed = true
	}
	t.root = t.removeInternal(t.root, key)
	if t.root != nil {
		t.root.isRed = false
	}
}

func (t *Tree[K, V]) removeInternal(node *Node[K, V], key K) *Node[K, V] {
	if node == nil {
		return nil
	}
	cmp := key.Compare(node.key)
	if cmp < 0 {
		if !isRed(node.left) && !isRed(node.left.left) {
			node = moveRedLeft(node)
		}
		node.left = t.removeInternal(node.left, key)
	} else {
		if isRed(node.left) {
			node = rotateRight(node)
		}
		if cmp == 0 && node.right == nil {
			t.size--
			return nil
		}
		if !isRed(node.right) && !isRed(node.right.left) {
			node = moveRedRight(node)
		}
		if cmp == 0 {
			t.size--
			smallest := min(node.right)
			node.key, node.value = smallest.key, smallest.value
			node.right = removeMin(node.right)
		} else {
			node.right = t.removeInternal(node.right, key)
		}
	}
	return fixUp(node)
}

func removeMin[K Key, V Value](node *Node[K, V]) *Node[K, V] {
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
	for node.left != nil {
		node = node.left
	}
	return node
}

// ========== PART 2: "Splay-like" interface ==========

// Find(index): returns (node, offset, error) by in-order rank.
func (t *Tree[K, V]) Find(index int) (*Node[K, V], int, error) {
	if index < 0 || index >= nodeSize(t.root) {
		return nil, 0, errors.New("out of index")
	}
	// We'll do a downward search by rank
	return findByRank(t.root, index)
}

// findByRank returns (foundNode, offset, nil) or (nil, 0, err).
// offset is typically "in-node offset" concept if a node had multiple items.
// But here each node is just 1 item, so offset=0 or 1. We'll keep offset=0 for success.
func findByRank[K Key, V Value](node *Node[K, V], rank int) (*Node[K, V], int, error) {
	if node == nil {
		return nil, 0, errors.New("unexpected nil node")
	}

	leftSz := nodeSize(node.left)
	if rank < leftSz {
		return findByRank(node.left, rank)
	} else if rank > leftSz {
		return findByRank(node.right, rank-leftSz-1)
	} else {
		// rank == leftSz => node is the target
		return node, 0, nil
	}
}

// InsertAfter: inserts newNode so that it is right after 'prev' in in-order.
func (t *Tree[K, V]) InsertAfter(prev *Node[K, V], newNode *Node[K, V]) *Node[K, V] {
	if prev == nil {
		// if we treat 'nil' as "insert at the end", or "insert at root"?
		// let's do "append at the end".
		fmt.Println("InsertAfter(nil) => append at the end")
		idx := nodeSize(t.root)
		t.root = t.insertAtIndex(t.root, idx, newNode)
		t.root.isRed = false
		return newNode
	}

	// find rank of prev
	r := rankOfNode(t.root, prev)
	// 노드와 랭크 정보를 출력한다.
	fmt.Printf("prev=%v, rank=%d\n", prev.value, r)

	if r < 0 {
		// not found => treat as end
		r = nodeSize(t.root)
	} else {
		r = r + 1
	}

	t.root = t.insertAtIndex(t.root, r, newNode)
	t.root.isRed = false
	return newNode
}

// rankOfNode: returns in-order rank of the target node, or -1 if not found
func rankOfNode[K Key, V Value](root *Node[K, V], target *Node[K, V]) int {
	if root == nil || target == nil {
		return -1
	}
	cmp := target.key.Compare(root.key)
	if cmp < 0 {
		return rankOfNode(root.left, target)
	} else if cmp > 0 {
		rR := rankOfNode(root.right, target)
		if rR < 0 {
			return -1
		}
		return nodeSize(root.left) + 1 + rR
	} else {
		// same key => check pointer
		if root == target {
			return nodeSize(root.left)
		}
		// or search both sides if duplicates possible
		rL := rankOfNode(root.left, target)
		if rL >= 0 {
			return rL
		}
		rR := rankOfNode(root.right, target)
		if rR < 0 {
			return -1
		}
		return nodeSize(root.left) + 1 + rR
	}
}

// insertAtIndex: core function for "in-order rank" insertion.
func (t *Tree[K, V]) insertAtIndex(node *Node[K, V], idx int, newNode *Node[K, V]) *Node[K, V] {
	if node == nil {
		return newNode
	}
	leftSz := nodeSize(node.left)

	if idx <= leftSz {
		if node.left == nil {
			fmt.Printf("=========== INSERT %v ============\n", node.value)
			fmt.Println(t.Print())
			fmt.Printf("=========== INSERT %v ============\n", node.value)
		}

		node.left = t.insertAtIndex(node.left, idx, newNode)
	} else {
		if node.right == nil {
			fmt.Printf("=========== INSERT %v ============\n", node.value)
			fmt.Println(t.Print())
			fmt.Printf("=========== INSERT %v ============\n", node.value)
		}

		node.right = t.insertAtIndex(node.right, idx-leftSz-1, newNode)
	}

	// LLRB balancing
	if isRed(node.right) && !isRed(node.left) {
		node = rotateLeft(node)
	}
	if isRed(node.left) && isRed(node.left.left) {
		node = rotateRight(node)
	}
	if isRed(node.left) && isRed(node.right) {
		flipColors(node)
	}

	updateSize(node)
	return node
}

// Delete: deletes the given node from the tree (like splay's Delete(node)).
func (t *Tree[K, V]) Delete(node *Node[K, V]) {
	if node == nil {
		return
	}
	// we find node.key, remove by the standard remove logic
	t.Remove(node.key)
}

// ============ PART 3: "ToTestString" & debugging ============

// ToTestString: similar to splay's "ToTestString()" => print [size, ???key]value
// We'll do an inorder to match the splay style.
func (t *Tree[K, V]) ToTestString() string {
	var sb strings.Builder
	traverseInOrder(t.root, func(n *Node[K, V]) {
		sb.WriteString(fmt.Sprintf("[%d]%s", n.size, n.value.String()))
	})
	return sb.String()
}

// String: just for printing node values in order
func (t *Tree[K, V]) String() string {
	var arr []string
	traverseInOrder(t.root, func(n *Node[K, V]) {
		arr = append(arr, n.value.String())
	})
	return strings.Join(arr, ",")
}
