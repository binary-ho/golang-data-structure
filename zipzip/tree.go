package zipzip

import "go-data-structure/random"

type Tree struct {
	root   *Node[Value]
	size   int
	random *random.Tree
}

func NewTree(maxRank int) *Tree {
	return &Tree{
		root:   nil,
		random: random.NewUint32Tree(maxRank),
	}
}

func (tree *Tree) Insert(key int, value Value) {
	newNode := tree.newNode(key, value)
	tree.insert(newNode)
	tree.size++
}

func (tree *Tree) insert(newNode *Node[Value]) {
	currentNode := tree.root
	prev := currentNode

	for currentNode != nil &&
		(newNode.rank < currentNode.rank ||
			(newNode.rank == currentNode.rank && newNode.key > currentNode.key)) {
		prev = currentNode
		if newNode.key < currentNode.key {
			currentNode = currentNode.left
			continue
		}
		currentNode = currentNode.right
	}

	if currentNode == tree.root {
		tree.root = newNode
	} else if newNode.key < prev.key {
		prev.left = newNode
	} else {
		prev.right = newNode
	}

	if currentNode == nil {
		newNode.left = nil
		newNode.right = nil
		return
	}

	if newNode.key < currentNode.key {
		newNode.right = currentNode
	} else {
		newNode.left = currentNode
	}

	prev = currentNode
	for currentNode != nil {
		if currentNode.key < newNode.key {
			//	until: currentNode == nil || currentNode.key > newNode.key
			for currentNode != nil && currentNode.key <= newNode.key {
				prev = currentNode
				currentNode = currentNode.right
			}
		} else {
			// unitl: currentNode == nil || currentNode.key < newNode.key
			for currentNode != nil && currentNode.key >= newNode.key {
				prev = currentNode
				currentNode = currentNode.left
			}
		}

		// 왜 Fix이고, 어떤 역할인가
		if fix := prev; fix.key > newNode.key || (fix == newNode && prev.key > newNode.key) {
			fix.left = currentNode
		} else {
			fix.right = currentNode
		}
	}
}

func (tree *Tree) Remove(key int, value Value) {
	newNode := tree.newNode(key, value)
	tree.remove(newNode)
	tree.size--
}

func (tree *Tree) remove(newNode *Node[Value]) {
	currentNode := tree.root
	prev := currentNode

	for currentNode != nil && newNode.key != currentNode.key {
		prev = currentNode
		if newNode.key < currentNode.key {
			currentNode = currentNode.left
		} else {
			currentNode = currentNode.right
		}
	}

	if currentNode == nil {
		return
	}

	left := currentNode.left
	right := currentNode.right
	if left == nil {
		currentNode = right
	} else if right == nil {
		currentNode = left
	} else if left.rank >= right.rank {
		currentNode = left
	} else {
		currentNode = right
	}

	if tree.root == newNode {
		tree.root = currentNode
	} else if newNode.key < prev.key {
		prev.left = currentNode
	} else {
		prev.right = currentNode
	}

	for left != nil && right != nil {
		if left.rank >= right.rank {
			// until: left == nil || left.rank < right.rank
			for left != nil && left.rank >= right.rank {
				prev = left
				left = left.right
			}
			prev.right = right
		} else {
			// until: right == nil || left.rank >= right.rank
			for right != nil && left.rank < right.rank {
				prev = right
				right = right.left
			}
			prev.left = left
		}
	}
}

func (tree *Tree) Find(key int) *Node[Value] {
	return tree.find(tree.root, Key(key))
}

func (tree *Tree) find(currentNode *Node[Value], key Key) *Node[Value] {
	if currentNode == nil {
		return nil
	}

	if currentNode.key == key {
		return currentNode
	} else if key < currentNode.key {
		return tree.find(currentNode.left, key) // 왼쪽 서브트리로 재귀적 탐색
	} else {
		return tree.find(currentNode.right, key) // 오른쪽 서브트리로 재귀적 탐색
	}
}
