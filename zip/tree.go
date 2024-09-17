package zip

type Tree struct {
	root  *Node[Value]
	nodes map[Key]*Node[Value]
}

func (tree *Tree) Insert(key int, value Value) {
	newNode := newNode(key, value)
	tree.nodes[newNode.key] = newNode
	tree.root = tree.insert(tree.root, newNode)
}

// insert 이후의 Root를 반환
func (tree *Tree) insert(baseNode, newNode *Node[Value]) *Node[Value] {
	if baseNode == nil {
		return newNode
	}

	if newNode.Key() < baseNode.Key() {
		// insert 이후 SubTree의 Root 노드가 새로 삽입된 노드가 아니다.
		newRoot := tree.insert(baseNode.left, newNode)
		if newRoot != newNode {
			return baseNode
		}

		if newNode.rank < baseNode.rank {
			baseNode.left = newNode
			return baseNode
		}

		// NewNode가 BaseNode보다 위에 있어야 할 때
		baseNode.left = newNode.right
		newNode.right = baseNode
		return newNode
	}

	newRoot := tree.insert(baseNode.right, newNode)
	if newRoot != newNode {
		return baseNode
	}

	// TODO: 왜 여기는 다른거
	if newNode.rank <= baseNode.rank {
		baseNode.right = newNode
		return baseNode
	}

	// NewNode가 BaseNode보다 위에 있어야 할 때
	baseNode.right = newNode.left
	newNode.left = baseNode
	return newNode
}

func (tree *Tree) Find(key int) *Node[Value] {
	return tree.nodes[Key(key)]
}

func (tree *Tree) Delete(key int) {
	tree.root = tree.delete(tree.root, Key(key))
}

func (tree *Tree) delete(root *Node[Value], key Key) *Node[Value] {
	if key == root.key {
		delete(tree.nodes, key)
		return tree.zip(root.left, root.right)
	}

	if key < root.key {
		left := root.left
		if key == left.key {
			delete(tree.nodes, key)
			root.left = tree.zip(left.left, left.right)
		} else {
			tree.delete(root.left, key)
		}
	} else {
		right := root.right
		if key == right.key {
			delete(tree.nodes, key)
			root.right = tree.zip(right.left, right.right)
		} else {
			tree.delete(root.right, key)
		}
	}
	return root
}

func (tree *Tree) zip(left, right *Node[Value]) *Node[Value] {
	if left == nil {
		return right
	}

	if right == nil {
		return left
	}

	if left.rank < right.rank {
		right.left = tree.zip(left, right.left)
		return right
	} else {
		left.right = tree.zip(left.right, right)
		return left
	}
}
