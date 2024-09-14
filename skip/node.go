package skip

import "fmt"

const MaxHeight = 16

type Node[V Value] struct {
	value V
	key   int
	tower Tower[V]
}

func (node *Node[V]) Value() V {
	return node.value
}

func (node *Node[V]) Weight() int {
	return node.key
}

func (node *Node[V]) String() string {
	return fmt.Sprintf("{ %d: %s }", node.key, node.value.String())
}

type Tower[V Value] [MaxHeight]*Node[V]

type Value interface {
	String() string
}
