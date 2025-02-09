package statistic

import "fmt"

func (tree *Tree[V]) Insert(newNode *Node[V]) *Node[V] {
	// 1) 트리가 비어 있으면 => root = newNode
	if tree.root == nil {
		tree.root = newNode
		// 루트이므로 레드-블랙 트리 규칙상 블랙으로
		tree.root.isRed = false
		tree.root.parent = nil
		return newNode
	}

	// 2) 트리가 있으면 => "맨 끝 인덱스" = root.weight
	idx := tree.root.weight // root subtree 전체 길이
	// 인덱스 기반 삽입
	tree.root = tree.insertByIndex(tree.root, idx, newNode)

	// 삽입 끝나고 루트를 블랙 처리
	tree.root.isRed = false
	tree.root.parent = nil

	return newNode
}

func (tree *Tree[V]) InsertAfter(prev, newNode *Node[V]) *Node[V] {
	r := tree.IndexOf(prev)
	if r < 0 {
		// prev not found => 끝에 삽입
		r = 0
		if tree.root != nil {
			r = tree.root.weight
		}
	} else {
		// prev 뒤 => r + prev.value.Len()
		r = r + prev.value.Len()
	}

	tree.root = tree.insertByIndex(tree.root, r, newNode)
	tree.root.isRed = false
	tree.root.parent = nil
	fmt.Printf("======= InsertAfter: newNode=%s (index=%d) =======\n", newNode.value.String(), r)
	return newNode
}

// insertByIndex: BST 키 비교 없이 weight 기반으로 idx 위치에 삽입
func (tree *Tree[V]) insertByIndex(node *Node[V], idx int, newNode *Node[V]) *Node[V] {
	if node == nil {
		newNode.isRed = true
		newNode.parent = nil
		return newNode
	}

	lw := node.leftWeight()
	valLen := node.value.Len()
	if idx <= lw {
		leftInserted := tree.insertByIndex(node.left, idx, newNode)
		node.left = leftInserted
		leftInserted.parent = node
	} else if idx > lw+valLen-1 {
		newIdx := idx - (lw + valLen)
		rightInserted := tree.insertByIndex(node.right, newIdx, newNode)
		node.right = rightInserted
		rightInserted.parent = node
	} else {
		// idx가 lw~lw+valLen-1 사이면 => node 바로 뒤(= lw+nodeLen)
		// 사실상 lw+nodeLen == lw+valLen
		newIdx := idx - (lw + valLen)
		rightInserted := tree.insertByIndex(node.right, newIdx, newNode)
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
