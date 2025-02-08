package llrb4

import (
	"fmt"
	"strings"
	"time"
)

// ========== 0. 인터페이스 및 구조체 ==========

// Key: 키 비교(Insert에서만 사용)
type Key interface {
	Compare(k Key) int
}

// Value: 노드 값. Len()은 weight 계산에 사용 가능
type Value interface {
	Len() int
	String() string
}

// Node: LLRB 노드 구조
type Node[V Value] struct {
	key    time.Time
	value  V
	weight int // 이 노드(및 하위 자식)가 몇 개의 원소(길이)를 담는지
	parent *Node[V]
	left   *Node[V]
	right  *Node[V]
	isRed  bool
}

// NewNode: 새 노드 생성 시, 일단 weight = value.Len()
func NewNode[V Value](value V, isRed bool) *Node[V] {
	now := time.Now()
	return &Node[V]{
		key:    now,
		value:  value,
		isRed:  isRed,
		weight: value.Len(), // 초기 weight
	}
}

// Tree: LLRB 트리
type Tree[V Value] struct {
	root *Node[V]
	size int // 전체 노드 수 (노드 개수)
}

// NewTree creates a new LLRB Tree.
func NewTree[V Value]() *Tree[V] {
	return &Tree[V]{}
}

// Len(): 전체 노드 수 (weight 합이 아님)
func (t *Tree[V]) Len() int {
	return t.size
}

// ================== 1. LLRB 기본 헬퍼 ==================
func isRed[V Value](n *Node[V]) bool {
	return n != nil && n.isRed
}

func isNotRed[V Value](n *Node[V]) bool {
	return !isRed(n)
}

// leftWeight/rightWeight: 자식의 weight
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

// updateWeight: 이 노드의 weight = leftWeight + rightWeight + value.Len()
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

// UpdateTreeWeight: 현재 노드 ~ 루트까지 올라가며 weight 갱신
func (n *Node[V]) UpdateTreeWeight() {
	cur := n
	for cur != nil {
		cur.updateWeight()
		cur = cur.parent
	}
}

// rotateLeft / rotateRight
func rotateLeft[V Value](node *Node[V]) *Node[V] {
	right := node.right
	node.right = right.left
	if right.left != nil {
		right.left.parent = node
	}
	right.left = node
	right.parent = node.parent
	node.parent = right

	// 색상 이동
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

// fixUp: removeMin 등에서 마지막 균형
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

// moveRedLeft/moveRedRight
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

// min/removeMin
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

// ============== 2. 일반 키 기반 삽입/삭제 =============

// Insert: (key,value) LLRB 삽입
func (t *Tree[V]) Insert(newNode *Node[V]) *Node[V] {
	t.root = t.insertInternal(t.root, newNode)
	t.root.isRed = false
	t.root.parent = nil
	return t.root
}
func (t *Tree[V]) insertInternal(node *Node[V], newNode *Node[V]) *Node[V] {
	if node == nil {
		t.size++
		newNode.isRed = true
		newNode.parent = nil
		return newNode
	}
	cmp := newNode.key.Compare(node.key)
	if cmp < 0 {
		put := t.insertInternal(node.left, newNode)
		node.left = put
		if put != nil {
			put.parent = node
		}
	} else if cmp > 0 {
		put := t.insertInternal(node.right, newNode)
		node.right = put
		if put != nil {
			put.parent = node
		}
	} else {
		// 동일 키 -> 값만 갱신
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

// Remove by key
func (t *Tree[V]) Remove(createdAt time.Time) {
	if t.root == nil {
		return
	}

	if isNotRed(t.root.left) && isNotRed(t.root.right) {
		t.root.isRed = true
	}

	t.root = t.removeInternal(t.root, createdAt)
	if t.root != nil {
		t.root.isRed = false
		t.root.parent = nil
	}
}

func (t *Tree[V]) removeInternal(node *Node[V], createdAt time.Time) *Node[V] {
	if node == nil {
		return nil
	}

	cmp := createdAt.Compare(node.key)
	if cmp < 0 {
		if isNotRed(node.left) && isNotRed(node.left.left) {
			node = moveRedLeft(node)
		}
		node.left = t.removeInternal(node.left, createdAt)
	} else {
		if isRed(node.left) {
			node = rotateRight(node)
		}
		if cmp == 0 && node.right == nil {
			t.size--
			return nil
		}
		if isNotRed(node.right) && isNotRed(node.right.left) {
			node = moveRedRight(node)
		}
		if cmp == 0 {
			t.size--
			sm := min(node.right)
			node.key, node.value = sm.key, sm.value
			node.right = removeMin(node.right)
		} else {
			node.right = t.removeInternal(node.right, createdAt)
		}
	}
	return fixUp(node)
}

// ============ 3. IndexOf(node) & InsertAfter(prev,newNode) ============

// IndexOf(node): pointer 기반, bottom-up 계산.
// node가 "인오더 상 몇 번째(0-based)인지"를 구함.
// (node.value.Len()이 여러 개 아이템이라면 1이 아닌 node.value.Len()만큼 자리를 차지하나,
// 여기서는 "offset"까지 더 세밀하게 처리할 수도 있음. 일단 여기선 "노드 한 덩이=1칸" 으로 가정하면 node.value.Len()=1 식)
func (t *Tree[V]) IndexOf(n *Node[V]) int {
	if n == nil {
		return -1
	}
	// 1) 왼쪽 weight => "내 왼쪽 subtree에 몇 개나 있나"
	rank := n.leftWeight()
	// 2) 위로 올라가면서, "내가 부모의 오른쪽 자식이면" => 부모의 왼쪽 subtree + 부모 1칸도 건너뜀
	cur := n
	for cur.parent != nil {
		parent := cur.parent
		if parent.right == cur {
			rank += parent.leftWeight() + 1 // parent 그 자체 1개
		}
		cur = parent
	}
	return rank
}

// InsertAfter(prev, newNode): pointer 기반 rank계산 + index삽입
func (t *Tree[V]) InsertAfter(prev, newNode *Node[V]) *Node[V] {
	if prev == nil {
		// nil => "맨 끝 append"
		idx := 0
		if t.root != nil {
			idx = t.root.weight // root의 전체 weight (맨 뒤)
		}
		t.root = t.insertByIndex(t.root, idx, newNode)
		t.root.isRed = false
		t.root.parent = nil
		fmt.Println("InsertAfter(nil) => append at the end")
		return newNode
	}

	r := t.IndexOf(prev)
	if r < 0 {
		// prev가 트리에 없다고 판단 => 그냥 끝에 삽입
		r = 0
		if t.root != nil {
			r = t.root.weight
		}
	} else {
		// prev 뒤 => r + 1
		// (만약 node.value.Len()>1이면 r+ node.value.Len() 해야 할 수도 있음)
		r = r + 1
	}

	// 인덱스 r 위치에 newNode 삽입
	t.root = t.insertByIndex(t.root, r, newNode)
	t.root.isRed = false
	t.root.parent = nil

	fmt.Printf("=========== INSERT %s ============\n", newNode.value.String())
	return newNode
}

// insertByIndex: "weight 기반"으로 인덱스=idx 위치에 newNode 삽입
// => BST key 비교를 전혀 안 함!
func (t *Tree[V]) insertByIndex(node *Node[V], idx int, newNode *Node[V]) *Node[V] {
	if node == nil {
		t.size++
		newNode.isRed = true
		newNode.parent = nil
		// newNode.weight는 이미 newNode.value.Len()으로 설정
		return newNode
	}

	leftSz := node.leftWeight()
	if idx <= leftSz {
		// 왼쪽 서브트리에 삽입
		leftInserted := t.insertByIndex(node.left, idx, newNode)
		node.left = leftInserted
		leftInserted.parent = node
	} else {
		// 오른쪽 서브트리에 (idx - leftSz -1) 위치로 삽입
		newIdx := idx - leftSz - 1
		rightInserted := t.insertByIndex(node.right, newIdx, newNode)
		node.right = rightInserted
		rightInserted.parent = node
	}

	// LLRB 회전 / flip
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

// =========== 4. 디버깅용 String() ===========

func (t *Tree[V]) String() string {
	// inorder로 value만
	var arr []string
	traverseInOrder(t.root, func(n *Node[V]) {
		arr = append(arr, n.value.String())
	})
	return strings.Join(arr, ",")
}

func traverseInOrder[V Value](node *Node[V], visit func(*Node[V])) {
	if node == nil {
		return
	}
	traverseInOrder(node.left, visit)
	visit(node)
	traverseInOrder(node.right, visit)
}
