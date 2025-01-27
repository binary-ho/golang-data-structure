package llrb

func (node *Node[K, V]) InitWeight() {
	node.weight = 1
}

func (node *Node[K, V]) UpdateTreeWeight() {
	currentNode := node
	for currentNode != nil {
		currentNode.updateWeight()
		currentNode = currentNode.parent
	}
}

// updateWeight recalculates the weight of this node with the value and children.
func (node *Node[K, V]) updateWeight() {
	node.InitWeight()

	if node.left != nil {
		node.increaseWeight(node.leftWeight())
	}

	if node.right != nil {
		node.increaseWeight(node.rightWeight())
	}
}

func (node *Node[K, V]) increaseWeight(weight int) {
	node.weight += weight
}

func (node *Node[K, V]) leftWeight() int {
	if node.left == nil {
		return 0
	}
	return node.left.weight
}

func (node *Node[K, V]) rightWeight() int {
	if node.right == nil {
		return 0
	}
	return node.right.weight
}
