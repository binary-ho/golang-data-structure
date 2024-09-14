package segment

type Tree struct {
	children []*Node
}

func InitSegmentTree(N int, values []int64) (tree *Tree) {
	size := 4 * N
	tree = &Tree{children: make([]*Node, size+1)}
	for i := 0; i <= size; i++ {
		tree.children[i] = &Node{}
	}
	tree.initSegmentTree(1, 1, N, values)
	return tree
}

func (tree *Tree) initSegmentTree(nodeNumber, start, end int, values []int64) int64 {
	if start == end {
		node := tree.children[nodeNumber]
		node.setValues(values[start], start, end)
		return node.sum
	}

	mid := (start + end) / 2
	leftSum := tree.initSegmentTree(nodeNumber*2, start, mid, values)    // left child
	rightSum := tree.initSegmentTree(nodeNumber*2+1, mid+1, end, values) // right child

	node := tree.children[nodeNumber]
	node.setValues(leftSum+rightSum, start, end)
	return node.sum
}

func (tree *Tree) Query(start, end int) int64 {
	return tree.query(1, start, end)
}

func (tree *Tree) query(nodeNumber, start, end int) int64 {
	node := tree.children[nodeNumber]
	tree.pushLazy(nodeNumber)

	if start > node.end || node.start > end {
		return 0
	}

	if start <= node.start && node.end <= end {
		return node.sum
	}

	leftSum := tree.query(nodeNumber*2, start, end)
	rightSum := tree.query(nodeNumber*2+1, start, end)
	return leftSum + rightSum
}

func (tree *Tree) Update(targetStart, targetEnd int, newValue int64) {
	tree.update(1, targetStart, targetEnd, newValue)
}

func (tree *Tree) update(nodeNumber, targetStart, targetEnd int, newValue int64) int64 {
	node := tree.children[nodeNumber]
	tree.pushLazy(nodeNumber)

	if node.start > targetEnd || targetStart > node.end {
		return node.sum
	}

	if targetStart <= node.start && node.end <= targetEnd {
		node.sum += newValue * int64(node.end-node.start+1)
		if tree.IsNotLeaf(nodeNumber) {
			tree.children[nodeNumber*2].lazy += newValue
			tree.children[nodeNumber*2+1].lazy += newValue
		}
		return node.sum
	}

	leftSum := tree.update(nodeNumber*2, targetStart, targetEnd, newValue)
	rightSum := tree.update(nodeNumber*2+1, targetStart, targetEnd, newValue)
	node.sum = leftSum + rightSum
	return node.sum
}

func (tree *Tree) pushLazy(nodeNumber int) {
	node := tree.children[nodeNumber]
	if node.lazy == 0 {
		return
	}

	node.sum += node.lazy * int64(node.end-node.start+1)
	if tree.IsNotLeaf(nodeNumber) {
		tree.children[nodeNumber*2].lazy += node.lazy
		tree.children[nodeNumber*2+1].lazy += node.lazy
	}
	node.lazy = 0
}

func (tree *Tree) GetSize() int {
	return len(tree.children) - 1
}

func (tree *Tree) IsNotLeaf(nodeNumber int) bool {
	end := tree.GetSize()
	start := (end + 1) / 2
	return !(start <= nodeNumber && nodeNumber <= end)
}
