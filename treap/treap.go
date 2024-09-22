package treap

import "go-data-structure/random"

type Treap struct {
	root   *Node[Value]
	random *random.Tree
}

func NewTreap() *Treap {
	return &Treap{
		root:   nil,
		random: random.NewUint32Tree(16),
	}
}

func (treap *Treap) Find(index int) *Node[Value] {
	return treap.find(treap.root, index)
}

func (treap *Treap) find(nodeNow *Node[Value], index int) *Node[Value] {
	leftSize := 1
	if nodeNow.left != nil {
		leftSize += nodeNow.left.size
	}

	if index == leftSize {
		return nodeNow
	}

	// 왼쪽으로 이동
	if index < leftSize {
		return treap.find(nodeNow.left, index)
	}
	// 오른쪽으로 이동
	return treap.find(nodeNow.right, index-leftSize)
}

func (treap *Treap) Insert(key int, value Value) {
	root := treap.root
	newNode := treap.NewNode(key, value)
	treap.root = treap.insert(root, newNode)
}

func (treap *Treap) insert(current, newNode *Node[Value]) *Node[Value] {
	if current == nil {
		return newNode
	}

	if current.priority < newNode.priority {
		left, right := treap.split(current, newNode)
		newNode.left = left
		newNode.right = right
		return newNode
	}

	if newNode.key < current.key {
		insert := treap.insert(current.left, newNode)
		current.setLeft(insert)
		return current
	}
	insert := treap.insert(current.right, newNode)
	current.setRight(insert)
	return current
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
		baseNode.setLeft(right)
		return left, baseNode
	}
}

func (treap *Treap) Remove(key int) {
	treap.remove(treap.root, key)
}

func (treap *Treap) remove(baseNode *Node[Value], key int) (returnNode *Node[Value]) {
	returnNode = baseNode
	if treap.root == nil {
		return
	}

	if baseNode.key == key {
		merge := treap.merge(baseNode.left, baseNode.right)
		baseNode = nil
		returnNode = merge
	} else if key < baseNode.key {
		baseNode.left = treap.remove(baseNode.left, key)
	} else {
		baseNode.right = treap.remove(baseNode.right, key)
	}
	return
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
