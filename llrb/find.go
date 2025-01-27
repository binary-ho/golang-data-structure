package llrb

import "fmt"

// Find returns the Node and offset of the given index. (Not Key)
func (tree *Tree[K, V]) Find(index int) (*Node[K, V], int, error) {
	if tree.root == nil {
		return nil, 0, nil
	}

	node := tree.root
	offset := index
	for {
		if node.left != nil && offset <= node.leftWeight() {
			node = node.left
		} else if node.right != nil && node.leftWeight()+node.value.Len() < offset {
			offset -= node.leftWeight() + node.value.Len()
			node = node.right
		} else {
			offset -= node.leftWeight()
			break
		}
	}

	if offset > node.value.Len() {
		return nil, 0, fmt.Errorf("node length %d, index %d: %w", node.value.Len(), offset, ErrOutOfIndex)
	}

	return node, offset, nil
}

// IndexOf Find the index of the given node.
func (tree *Tree[K, V]) IndexOf(node *Node[K, V]) int {
	rank := 1
	if node.left != nil {
		rank += node.left.weight
	}

	currentNode := node
	for currentNode.parent != nil {
		if currentNode.parent.right == currentNode {
			if currentNode.parent.left != nil {
				rank += currentNode.parent.left.weight
			}
			rank += 1
		}
		currentNode = currentNode.parent
	}
	return rank
}
