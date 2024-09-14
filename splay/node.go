package splay

type Node struct {
	key                 int64
	left, right, parent *Node
	size, sum, lazy     int64
}

func (node *Node) pushLazyValue() {
	lazyValue := node.lazy
	node.lazy = 0

	//if node.key != math.MinInt64 && node.key != math.MaxInt64 {
	node.key += lazyValue
	//}

	if left := node.left; isPresent(left) {
		left.lazy += lazyValue
		left.sum += left.size * lazyValue
	}

	if right := node.right; isPresent(right) {
		right.lazy += lazyValue
		right.sum += right.size * lazyValue
	}
}

func (node *Node) setGrandParentToParent() {
	parent := node.parent
	grandParent := parent.parent

	// change Parent
	node.parent = grandParent
	parent.parent = node

	if isNil(node.parent) {
		return
	}

	// TODO : grandParent로 바꾸기
	if parent == node.parent.left {
		node.parent.left = node
	} else {
		node.parent.right = node
	}

	// TODO : 중복 제거
	node.updateTreeSizeAndSum()
	parent.updateTreeSizeAndSum()
}

func (node *Node) setParentToChild() {
	parent := node.parent
	var newChild *Node

	if node == parent.left {
		newChild = node.right
		parent.left = newChild
		node.right = parent
	} else {
		newChild = node.left
		parent.right = newChild
		node.left = parent
	}

	if newChild != nil {
		newChild.parent = parent
	}
}

func (node *Node) updateTreeSizeAndSum() {
	node.size = 1
	node.sum = node.key

	if isPresent(node.left) {
		node.size += node.left.size
		node.sum += node.left.sum
	}

	if isPresent(node.right) {
		node.size += node.right.size
		node.sum += node.right.sum
	}
}

func isNil(node *Node) bool {
	return node == nil
}

func isPresent(node *Node) bool {
	return node != nil
}
