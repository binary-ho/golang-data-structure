package llrb

// Node is a node of Tree.
type Node[K Key, V Value] struct {
	key    K
	value  V
	weight int
	parent *Node[K, V]
	left   *Node[K, V]
	right  *Node[K, V]
	isRed  bool
}

// NewNode creates a new instance of Node.
func NewNode[K Key, V Value](key K, value V, isRed bool) *Node[K, V] {
	return &Node[K, V]{
		key:   key,
		value: value,
		isRed: isRed,
	}
}

func (node *Node[K, V]) hasLinks() bool {
	return node.parent != nil || node.left != nil || node.right != nil
}

// Key represents key of Tree.
type Key interface {
	// Compare gives the result of a 3-way comparison
	// a.Compare(b) = 1 => a > b
	// a.Compare(b) = 0 => a == b
	// a.Compare(b) = -1 => a < b
	Compare(k Key) int
}

// Value represents the data stored in the nodes of Tree.
type Value interface {
	Len() int
	String() string
}
