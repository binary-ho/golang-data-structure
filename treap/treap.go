package treap

type Treap struct {
	root *Node[Value]
}

func (treap *Treap) Insert(key int, value Value) {
	root := treap.root
	newNode := NewNode(key, value)
	if root == nil {
		root = newNode
		return
	}

	// set newNode to root
	if root.priority < newNode.priority {
		treap.setRoot(newNode)
		return
	}

	if newNode.key < root.key {
		root.setLeft(newNode)
		return
	}
	root.setRight(newNode)
}

func (treap *Treap) setRoot(node *Node[Value]) {
	left, right := treap.split(treap.root, node)
	node.left = left
	node.right = right
	treap.root = node
}

func (treap *Treap) split(baseNode, newNode *Node[Value]) (*Node[Value], *Node[Value]) {
	if baseNode == nil {
		return nil, nil
	}

	if baseNode.key < newNode.key {
		left, right := treap.split(baseNode.right, newNode)
		baseNode.setRight(left)
		return baseNode, right
	} else {
		left, right := treap.split(baseNode.left, newNode)
		// 여기 right 오는게 맞는거 같아
		baseNode.setLeft(right)
		return left, baseNode
	}
}

func (treap *Treap) Remove(key int) {
	treap.remove(treap.root, key)
}

func (treap *Treap) remove(baseNode *Node[Value], key int) *Node[Value] {
	if treap.root == nil {
		return baseNode
	}

	if baseNode.key == key {
		merge := treap.merge(baseNode.left, baseNode.right)
		baseNode = nil
		return merge
	}

	if key < baseNode.key {
		baseNode.left = treap.remove(baseNode.left, key)
	} else {
		baseNode.right = treap.remove(baseNode.right, key)
	}
	return baseNode
}

func (treap *Treap) merge(left, right *Node[Value]) *Node[Value] {
	if left == nil {
		return right
	}

	if right == nil {
		return left
	}

	if left.priority < right.priority {
		right.left = treap.merge(left, right.left)
		return right
	}

	left.right = treap.merge(left.right, right)
	return left
}
