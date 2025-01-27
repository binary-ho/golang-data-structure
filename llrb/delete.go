package llrb

// Delete removes the value of the given key.
func (tree *Tree[K, V]) Delete(key K) {
	if isNotRed(tree.root.left) && isNotRed(tree.root.right) {
		tree.root.isRed = true
	}

	tree.root = tree.delete(tree.root, key)

	if tree.root != nil {
		tree.root.parent = nil
		tree.root.isRed = false
	}
}

func (tree *Tree[K, V]) delete(node *Node[K, V], key K) *Node[K, V] {
	if key.Compare(node.key) < 0 {
		if isNotRed(node.left) && isNotRed(node.left.left) {
			node = moveRedLeft(node)
		}
		node.left = tree.delete(node.left, key)
		return fixUp(node)
	}

	if isRed(node.left) {
		node = rotateRight(node)
	}

	if key.Compare(node.key) == 0 && node.right == nil {
		tree.size--
		return nil
	}

	if isNotRed(node.right) && isNotRed(node.right.left) {
		node = moveRedRight(node)
	}

	if key.Compare(node.key) == 0 {
		tree.size--
		smallest := min(node.right)
		node.value = smallest.value
		node.key = smallest.key
		node.right = removeMin(node.right)
	} else {
		node.right = tree.delete(node.right, key)
	}
	return fixUp(node)
}
