package treap

type Node[V Value] struct {
	key            int
	value          V
	priority, size int
	left, right    *Node[V]
}

func (treap *Treap) NewNode(key int, value Value) *Node[Value] {
	return &Node[Value]{
		key:      key,
		value:    value,
		priority: treap.random.Rank(),
		size:     1,
		left:     nil,
		right:    nil,
	}
}

func (node *Node[V]) setLeft(newLeft *Node[V]) {
	node.left = newLeft
	node.resize()
}

func (node *Node[V]) setRight(newRight *Node[V]) {
	node.right = newRight
	node.resize()
}

func (node *Node[V]) resize() {
	node.size = 1
	if node.left != nil {
		node.size += node.left.size
	}
	if node.right != nil {
		node.size += node.right.size
	}
}

type Value interface {
	String() string
}
