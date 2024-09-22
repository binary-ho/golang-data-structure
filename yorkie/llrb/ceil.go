package yorkie_llrb

// Ceil returns the smallest key greater than or equal to the given key.
func (tree *Tree[K, V]) Ceil(key K) (K, V) {
	node := tree.root

	for node != nil {
		compare := key.Compare(node.key)
		if compare < 0 {
			if node.left != nil {
				node.left.parent = node
				node = node.left
			} else {
				return node.key, node.value
			}
		} else if compare > 0 {
			if node.right != nil {
				node.right.parent = node
				node = node.right
			} else {
				parent := node.parent
				child := node
				for parent != nil && child == parent.right {
					child = parent
					parent = parent.parent
				}

				if parent != nil {
					return parent.key, parent.value
				} else {
					var zeroK K
					var zeroV V
					return zeroK, zeroV
				}
			}
		} else {
			return node.key, node.value
		}
	}

	var zeroK K
	var zeroV V
	return zeroK, zeroV
}
