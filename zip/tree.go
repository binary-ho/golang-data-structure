package zip

import "go-data-structure/random"

type Tree struct {
	root   *Node[Value]
	nodes  map[Key]*Node[Value]
	random *random.Tree
}

func NewTree(maxRank int) *Tree {
	return &Tree{
		nodes:  make(map[Key]*Node[Value], maxRank),
		random: random.NewUint32Tree(maxRank),
	}
}

func (tree *Tree) Insert(key int, value Value) {
	newNode := tree.newNode(key, value)
	tree.root = tree.insert(tree.root, newNode)
	tree.nodes[newNode.key] = newNode
}

// insert 이후의 Root를 반환
func (tree *Tree) insert(current, newNode *Node[Value]) *Node[Value] {
	if current == nil {
		return newNode
	}

	if newNode.Key() < current.Key() {
		// insert 이후 SubTree의 Root 노드가 새로 삽입된 노드가 아니다.
		newRoot := tree.insert(current.left, newNode)
		if newRoot != newNode {
			return current
		}

		if newNode.rank < current.rank {
			current.left = newNode
			return current
		}

		// NewNode가 BaseNode보다 위에 있어야 할 때
		current.left = newNode.right
		newNode.right = current
		return newNode
	}

	newRoot := tree.insert(current.right, newNode)
	if newRoot != newNode {
		return current
	}

	// TODO: 왜 여기는 다른거
	if newNode.rank <= current.rank {
		current.right = newNode
		return current
	}

	// NewNode가 BaseNode보다 위에 있어야 할 때
	current.right = newNode.left
	newNode.left = current
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
