package zip

type Node[V Value] struct {
	key         Key
	value       V
	rank        int
	left, right *Node[V]
}

type Key int

func newNode(key int, value Value) *Node[Value] {
	return &Node[Value]{
		key:   Key(key),
		value: value,
		rank:  getRank(),
	}
}

func (node *Node[Value]) Key() int {
	return int(node.key)
}

type Value interface {
	String() string
}
