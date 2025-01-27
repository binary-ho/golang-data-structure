package llrb

// Insert inserts the value of the given key.
func (tree *Tree[K, V]) Insert(newNode *Node[K, V]) *Node[K, V] {
	tree.root = tree.insert(tree.root, newNode)
	tree.root.isRed = false
	return tree.root
}

func (tree *Tree[K, V]) insert(node *Node[K, V], newNode *Node[K, V]) *Node[K, V] {
	if node == nil {
		tree.size++
		// TODO : 변경 필요
		return NewNode(newNode.key, newNode.value, true)
	}

	key := newNode.key
	compare := key.Compare(node.key)
	if compare < 0 {
		put := tree.insert(node.left, newNode)
		node.left = put
		put.parent = node
	} else if compare > 0 {
		put := tree.insert(node.right, newNode)
		node.right = put
		put.parent = node
	} else {
		node.value = newNode.value
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

	node.UpdateTreeWeight()
	return node
}
