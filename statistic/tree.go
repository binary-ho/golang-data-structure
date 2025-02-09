package statistic

import (
	"fmt"
	"strings"
)

// ErrOutOfIndex is returned when the given index is out of index.
var ErrOutOfIndex = fmt.Errorf("out of index")

type Value interface {
	Len() int
	String() string
}

type Node[V Value] struct {
	value  V
	weight int
	parent *Node[V]
	left   *Node[V]
	right  *Node[V]
	isRed  bool
}

func (node *Node[V]) Value() V {
	return node.value
}

func NewNode[V Value](value V) *Node[V] {
	node := &Node[V]{
		value: value,
		isRed: true,
	}
	node.InitWeight()
	return node
}

type Tree[V Value] struct {
	root *Node[V]
}

func NewTree[V Value](root *Node[V]) *Tree[V] {
	return &Tree[V]{
		root: root,
	}
}

func (tree *Tree[V]) Len() int {
	if tree.root == nil {
		return 0
	}
	return tree.root.weight
}

func isRed[V Value](n *Node[V]) bool {
	return n != nil && n.isRed
}
func isNotRed[V Value](n *Node[V]) bool {
	return !isRed(n)
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
	if node == nil {
		return nil
	}

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

// Delete(node): "node -> indexOf -> removeAtIndex"
func (tree *Tree[V]) Delete(node *Node[V]) {
	if node == nil {
		return
	}

	nodeIndex := tree.IndexOf(node)
	if nodeIndex < 0 {
		return
	}

	tree.root = tree.removeAtIndex(tree.root, nodeIndex)
	if tree.root != nil {
		tree.root.isRed = false
		tree.root.parent = nil
	}

	//node.unlink()
	if tree.root != nil {
		tree.root.updateWeight()
	}
}

func (tree *Tree[V]) removeAtIndex(node *Node[V], idx int) *Node[V] {
	if node == nil {
		return nil
	}

	leftSz := 0
	if node.left != nil {
		leftSz = node.left.weight
	}
	nodeLen := node.value.Len() // 보통 1

	// 1) idx가 leftSz보다 작으면 => 왼쪽 서브트리로
	if idx < leftSz {
		if isNotRed(node.left) && isNotRed(node.left.left) {
			node = moveRedLeft(node)
		}
		fmt.Println("remove 1")
		node.left = tree.removeAtIndex(node.left, idx)
	} else {
		// 2) idx >= leftSz
		//   => "왼쪽 서브트리 + 내 노드 범위"를 넘는다면 => 오른쪽
		if isRed(node.left) {
			node = rotateRight(node)
		}

		// 만약 "이 노드를 지우는 경우"이고 "오른쪽이 비었으면" => 그냥 free
		// "삭제 대상 노드"이고 "오른쪽 서브트리가 없다" => 그냥 노드 제거
		if idx >= leftSz && idx < leftSz+nodeLen && node.right == nil {
			//tree.size--
			fmt.Println("remove 2")
			return nil
		}

		if isNotRed(node.right) && isNotRed(node.right.left) {
			node = moveRedRight(node)
		}

		if idx >= leftSz && idx < leftSz+nodeLen {
			// => 이 노드를 삭제하는 상황
			//tree.size--

			minNode, minRight := tree.extractMinNode(node.right)
			if minNode == nil {
				fmt.Println("remove 3")
				return nil
			}

			minNode.isRed = node.isRed
			minNode.left = node.left
			if minNode.left != nil {
				minNode.left.parent = minNode
			}

			minNode.right = minRight
			if minNode.right != nil {
				minNode.right.parent = minNode
			}

			node.UpdateTreeWeight()
			node.unlink()
			node = nil

			fmt.Println(tree.Print())
			fmt.Println("remove 4")
			return minNode
			//node = minNode

			//만약 node.right == nil => 위에서 return nil 됨
			//if node.right == nil {
			//	fmt.Println("remove 3")
			//	return nil
			//}
			//아니면 "오른쪽 서브트리에서 min"을 가져와 교체
			//왼쪽 자식이 null인 위치
			//sm := min(node.right)
			//fmt.Printf("뭐 찾았어 시키야!!! %v, %v, %d\n", node.value, sm.value, idx)
			//node.value = sm.value
			//node.key = sm.key (키가 없으므로 스킵)
			//node.right = removeMin(node.right)
			//
			//sm.UpdateTreeWeight()
			//sm.unlink()
		} else {
			// idx가 왼+내 길이보다 크면 => 오른쪽으로
			newIdx := idx - (leftSz + nodeLen)
			fmt.Println("remove 5")
			node.right = tree.removeAtIndex(node.right, newIdx)
		}
	}
	return fixUp(node)
}

func (tree *Tree[V]) extractMinNode(node *Node[V]) (minNode *Node[V], minNodeRight *Node[V]) {
	for node == nil {
		return nil, nil
	}

	if node.left == nil {
		minNodeRight = node.right
		if minNodeRight != nil {
			minNodeRight.parent = node.parent
		}
		node.unlink()
		return node, node.right
	}

	if isNotRed(node.left) && isNotRed(node.left.left) {
		node = moveRedLeft(node)
	}

	extractedMinNode, extractedMinNodeRight := tree.extractMinNode(node.left)
	node.left = extractedMinNodeRight

	fixUp(node)
	return extractedMinNode, node
}

// Find returns the Node and offset of the given index. (Not Key)
func (tree *Tree[V]) Find(index int) (*Node[V], int, error) {
	if tree.root == nil {
		return nil, 0, nil
	}

	node := tree.root
	offset := index
	for {
		if node.left != nil && offset <= node.leftWeight() {
			node = node.left
		} else if node.right != nil && node.leftWeight()+node.value.Len() < offset {
			offset -= node.leftWeight() + node.value.Len()
			node = node.right
		} else {
			offset -= node.leftWeight()
			break
		}
	}

	if offset > node.value.Len() {
		return nil, 0, fmt.Errorf("node length %d, index %d: %w", node.value.Len(), offset, ErrOutOfIndex)
	}

	return node, offset, nil
}

// IndexOf(node): pointer 기반 bottom-up
func (tree *Tree[V]) IndexOf(node *Node[V]) int {
	if node == nil || tree.isNotTreeNode(node) {
		return -1
	}

	rank := node.leftWeight()
	cur := node
	for cur.parent != nil {
		p := cur.parent
		if p.right == cur {
			rank += p.leftWeight() + p.value.Len()
		}
		cur = p
	}
	return rank
}

func (tree *Tree[V]) ToTestString() string {
	//return tree.Print()
	var sb strings.Builder
	traverseInOrder(tree.root, func(n *Node[V]) {
		//sb.WriteString(fmt.Sprintf("[weight: %d, len: %d] value: %s", n.weight, n.value.Len(), n.value.String()))
		sb.WriteString(fmt.Sprintf("[%d,%d]%s", n.weight, n.value.Len(), n.value.String()))
	})
	return sb.String()
}

func (tree *Tree[V]) String() string {
	var arr []string
	traverseInOrder(tree.root, func(n *Node[V]) {
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

func (node *Node[V]) unlink() {
	node.parent = nil
	node.right = nil
	node.left = nil
}

func (tree *Tree[V]) isNotTreeNode(node *Node[V]) bool {
	return node != tree.root && node.hasNotLink()
}

func (node *Node[V]) hasNotLink() bool {
	return node.parent == nil && node.right == nil && node.left == nil
}
