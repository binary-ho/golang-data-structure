package segment

type Node struct {
	start, end int
	sum, lazy  int64
}

func (node *Node) setValues(value int64, start, end int) {
	node.sum = value
	node.start = start
	node.end = end
}
