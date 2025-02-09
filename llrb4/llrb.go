package llrb4

import (
	"fmt"
	"strings"
	"time"
)

// =================== 인터페이스 ===================
// Value: 노드 값. Len()은 weight 계산에 사용
type Value interface {
	Len() int
	String() string
}

// ================== Node, Tree ==================

// Node: LLRB 노드
type Node[V Value] struct {
	key    time.Time // 키: Insert() 시 BST 비교에 사용
	value  V
	weight int // node+자식들의 총 길이
	parent *Node[V]
	left   *Node[V]
	right  *Node[V]
	isRed  bool
}

// NewNode: 새 노드를 만드는 시점에 time.Now()로 key를 생성 (또는 외부 주입도 가능)
func NewNode[V Value](value V) *Node[V] {
	now := time.Now()
	return &Node[V]{
		key:    now,
		value:  value,
		isRed:  true,
		weight: value.Len(),
	}
}

// Tree: LLRB 전체
type Tree[V Value] struct {
	root *Node[V]
	size int // 전체 노드 수
}

// NewTree: 생성
func NewTree[V Value](root *Node[V]) *Tree[V] {
	return &Tree[V]{
		root: root,
	}
}

// Len: SplayTree와 동일 인터페이스, 전체 노드 수
func (t *Tree[V]) Len() int {
	return t.size
}

// ============ 1. 주요 헬퍼 ============

func isRed[V Value](n *Node[V]) bool {
	return n != nil && n.isRed
}
func isNotRed[V Value](n *Node[V]) bool {
	return !isRed(n)
}

func (n *Node[V]) leftWeight() int {
	if n.left == nil {
		return 0
	}
	return n.left.weight
}
func (n *Node[V]) rightWeight() int {
	if n.right == nil {
		return 0
	}
	return n.right.weight
}

// updateWeight: leftWeight + rightWeight + node.value.Len()
func (n *Node[V]) updateWeight() {
	if n == nil {
		return
	}
	w := n.value.Len()
	if n.left != nil {
		w += n.left.weight
	}
	if n.right != nil {
		w += n.right.weight
	}
	n.weight = w
}

// UpdateTreeWeight: 이 노드에서 루트까지 올라가면서 weight 재계산
func (n *Node[V]) UpdateTreeWeight() {
	for cur := n; cur != nil; cur = cur.parent {
		cur.updateWeight()
	}
}

// rotateLeft, rotateRight
func rotateLeft[V Value](node *Node[V]) *Node[V] {
	right := node.right
	node.right = right.left
	if right.left != nil {
		right.left.parent = node
	}
	right.left = node
	right.parent = node.parent
	node.parent = right

	right.isRed = node.isRed
	node.isRed = true

	node.updateWeight()
	right.updateWeight()
	return right
}
func rotateRight[V Value](node *Node[V]) *Node[V] {
	left := node.left
	node.left = left.right
	if left.right != nil {
		left.right.parent = node
	}
	left.right = node
	left.parent = node.parent
	node.parent = left

	left.isRed = node.isRed
	node.isRed = true

	node.updateWeight()
	left.updateWeight()
	return left
}

// flipColors
func flipColors[V Value](node *Node[V]) {
	node.isRed = !node.isRed
	node.left.isRed = !node.left.isRed
	node.right.isRed = !node.right.isRed
}

// fixUp
func fixUp[V Value](node *Node[V]) *Node[V] {
	if isRed(node.right) && isNotRed(node.left) {
		node = rotateLeft(node)
	}
	if isRed(node.left) && isRed(node.left.left) {
		node = rotateRight(node)
	}
	if isRed(node.left) && isRed(node.right) {
		flipColors(node)
	}
	node.updateWeight()
	return node
}

// moveRedLeft, moveRedRight
func moveRedLeft[V Value](node *Node[V]) *Node[V] {
	flipColors(node)
	if isRed(node.right.left) {
		node.right = rotateRight(node.right)
		node = rotateLeft(node)
		flipColors(node)
	}
	return node
}
func moveRedRight[V Value](node *Node[V]) *Node[V] {
	flipColors(node)
	if isRed(node.left.left) {
		node = rotateRight(node)
		flipColors(node)
	}
	return node
}

// min, removeMin
func min[V Value](node *Node[V]) *Node[V] {
	for node.left != nil {
		node = node.left
	}
	return node
}
func removeMin[V Value](node *Node[V]) *Node[V] {
	if node.left == nil {
		return nil
	}
	if isNotRed(node.left) && isNotRed(node.left.left) {
		node = moveRedLeft(node)
	}
	node.left = removeMin(node.left)
	return fixUp(node)
}

// =========== 2.  Insert(Key= time.Now), Delete(node) ===========

// Insert(node): splay와 동일하게 "노드를 추가"
func (t *Tree[V]) Insert(newNode *Node[V]) {
	// BST: compare "newNode.key" with node.key
	t.root = t.insertInternal(t.root, newNode)
	t.root.isRed = false
	t.root.parent = nil
}

// insertInternal: BST 로직(time.Time 비교) + LLRB
func (t *Tree[V]) insertInternal(node *Node[V], newNode *Node[V]) *Node[V] {
	if node == nil {
		t.size++
		newNode.isRed = true
		newNode.parent = nil
		return newNode
	}
	if newNode.key.Before(node.key) {
		// left
		leftInserted := t.insertInternal(node.left, newNode)
		node.left = leftInserted
		leftInserted.parent = node
	} else if newNode.key.After(node.key) {
		// right
		rightInserted := t.insertInternal(node.right, newNode)
		node.right = rightInserted
		rightInserted.parent = node
	} else {
		// same key => update value
		node.value = newNode.value
	}

	// LLRB balancing
	if isRed(node.right) && isNotRed(node.left) {
		node = rotateLeft(node)
	}
	if isRed(node.left) && isRed(node.left.left) {
		node = rotateRight(node)
	}
	if isRed(node.left) && isRed(node.right) {
		flipColors(node)
	}

	node.updateWeight()
	return node
}

// Delete(node): splay 와 동일. node.key 로 remove.
func (t *Tree[V]) Delete(n *Node[V]) {
	if n == nil {
		return
	}
	t.remove(n.key)
}

// remove(key): 내부 BST 로직으로 해당 key 제거
func (t *Tree[V]) remove(key time.Time) {
	if t.root == nil {
		return
	}
	if isNotRed(t.root.left) && isNotRed(t.root.right) {
		t.root.isRed = true
	}
	t.root = t.removeInternal(t.root, key)
	if t.root != nil {
		t.root.isRed = false
		t.root.parent = nil
	}
}

func (t *Tree[V]) removeInternal(node *Node[V], key time.Time) *Node[V] {
	if node == nil {
		return nil
	}
	if key.Before(node.key) {
		if isNotRed(node.left) && isNotRed(node.left.left) {
			node = moveRedLeft(node)
		}
		node.left = t.removeInternal(node.left, key)
	} else {
		if isRed(node.left) {
			node = rotateRight(node)
		}
		if key.Equal(node.key) && node.right == nil {
			t.size--
			return nil
		}
		if isNotRed(node.right) && isNotRed(node.right.left) {
			node = moveRedRight(node)
		}
		if key.Equal(node.key) {
			t.size--
			sm := min(node.right)
			node.key, node.value = sm.key, sm.value
			node.right = removeMin(node.right)
		} else {
			node.right = t.removeInternal(node.right, key)
		}
	}
	return fixUp(node)
}

// ============ 3. Find(index), InsertAfter(prev, node) ============

// Find(index): (node, offset, error)
func (t *Tree[V]) Find(index int) (*Node[V], int, error) {
	if t.root == nil {
		return nil, 0, nil
	}
	return findByIndex(t.root, index)
}

// findByIndex => (node, offset, error)
func findByIndex[V Value](root *Node[V], idx int) (*Node[V], int, error) {
	if root == nil {
		return nil, 0, fmt.Errorf("out of index: idx=%d", idx)
	}
	lw := root.leftWeight()
	valLen := root.value.Len() // 이 노드 자체가 차지하는 길이 (보통 1)
	if idx < lw {
		return findByIndex(root.left, idx)
	} else if idx >= lw+valLen {
		return findByIndex(root.right, idx-(lw+valLen))
	} else {
		// idx 범위가 [lw, lw+valLen)
		offset := idx - lw
		return root, offset, nil
	}
}

// IndexOf(node): pointer 기반 bottom-up
func (t *Tree[V]) IndexOf(n *Node[V]) int {
	if n == nil {
		return -1
	}
	rank := n.leftWeight()
	cur := n
	for cur.parent != nil {
		p := cur.parent
		if p.right == cur {
			rank += p.leftWeight() + p.value.Len() // 이 노드도 count
		}
		cur = p
	}
	return rank
}

// InsertAfter(prev, newNode):
func (t *Tree[V]) InsertAfter(prev, newNode *Node[V]) *Node[V] {
	if prev == nil {
		// nil => 끝에 삽입
		endIdx := 0
		if t.root != nil {
			endIdx = t.root.weight
		}
		t.root = t.insertByIndex(t.root, endIdx, newNode)
		t.root.isRed = false
		t.root.parent = nil
		fmt.Println("InsertAfter(nil) => append at the end")
		return newNode
	}

	r := t.IndexOf(prev)
	if r < 0 {
		// prev not found => 끝에 삽입
		r = 0
		if t.root != nil {
			r = t.root.weight
		}
	} else {
		// prev 뒤 => r + prev.value.Len()
		r = r + prev.value.Len()
	}

	t.root = t.insertByIndex(t.root, r, newNode)
	t.root.isRed = false
	t.root.parent = nil
	fmt.Printf("======= InsertAfter: newNode=%s (index=%d) =======\n", newNode.value.String(), r)
	return newNode
}

// insertByIndex: BST 키 비교 없이 weight 기반으로 idx 위치에 삽입
func (t *Tree[V]) insertByIndex(node *Node[V], idx int, newNode *Node[V]) *Node[V] {
	if node == nil {
		t.size++
		newNode.isRed = true
		newNode.parent = nil
		return newNode
	}
	lw := node.leftWeight()
	valLen := node.value.Len()
	if idx <= lw {
		leftInserted := t.insertByIndex(node.left, idx, newNode)
		node.left = leftInserted
		leftInserted.parent = node
	} else if idx > lw+valLen-1 {
		newIdx := idx - (lw + valLen)
		rightInserted := t.insertByIndex(node.right, newIdx, newNode)
		node.right = rightInserted
		rightInserted.parent = node
	} else {
		// idx가 lw~lw+valLen-1 사이면 => node 바로 뒤(= lw+nodeLen)
		// 사실상 lw+nodeLen == lw+valLen
		newIdx := idx - (lw + valLen)
		rightInserted := t.insertByIndex(node.right, newIdx, newNode)
		node.right = rightInserted
		rightInserted.parent = node
	}

	// LLRB rotate/flip
	if isRed(node.right) && isNotRed(node.left) {
		node = rotateLeft(node)
	}
	if isRed(node.left) && isRed(node.left.left) {
		node = rotateRight(node)
	}
	if isRed(node.left) && isRed(node.right) {
		flipColors(node)
	}
	node.updateWeight()
	return node
}

// =========== 4. ToTestString & String ===========

// ToTestString: splay처럼 노드의 weight, len, value를 inorder로
func (t *Tree[V]) ToTestString() string {
	var sb strings.Builder
	traverseInOrder(t.root, func(n *Node[V]) {
		sb.WriteString(fmt.Sprintf("[%d,%d]%s", n.weight, n.value.Len(), n.value.String()))
	})
	return sb.String()
}

// String: 간단히 값들만
func (t *Tree[V]) String() string {
	var arr []string
	traverseInOrder(t.root, func(n *Node[V]) {
		arr = append(arr, n.value.String())
	})
	return strings.Join(arr, ",")
}

// traverseInOrder
func traverseInOrder[V Value](node *Node[V], visit func(*Node[V])) {
	if node == nil {
		return
	}
	traverseInOrder(node.left, visit)
	visit(node)
	traverseInOrder(node.right, visit)
}
