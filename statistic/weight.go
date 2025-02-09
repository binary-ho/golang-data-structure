package statistic

func (node *Node[V]) InitWeight() {
	node.weight = node.value.Len()
}

func (node *Node[V]) UpdateTreeWeight() {
	currentNode := node
	for currentNode != nil {
		currentNode.updateWeight()
		currentNode = currentNode.parent
	}
}

// updateWeight recalculates the weight of this node with the value and children.
func (node *Node[V]) updateWeight() {
	node.InitWeight()

	if node.left != nil {
		node.increaseWeight(node.leftWeight())
	}

	if node.right != nil {
		node.increaseWeight(node.rightWeight())
	}
}

func (node *Node[V]) increaseWeight(weight int) {
	node.weight += weight
}

func (node *Node[V]) leftWeight() int {
	if node.left == nil {
		return 0
	}
	return node.left.weight
}

func (node *Node[V]) rightWeight() int {
	if node.right == nil {
		return 0
	}
	return node.right.weight
}
