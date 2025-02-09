package statistic

import (
	"fmt"
	"strings"
)

// Print returns a visually appealing string representation of the tree
func (tree *Tree[V]) Print() string {
	if tree.root == nil {
		return "Empty tree"
	}

	nodeIndices := make(map[*Node[V]]int)
	currentIndex := 0
	tree.inorderIndexing(tree.root, &currentIndex, nodeIndices)

	var sb strings.Builder
	tree.print(tree.root, "", true, &sb, nodeIndices)
	return sb.String()
}

// inorderIndexing performs an inorder traversal and assigns indices to nodes
func (tree *Tree[V]) inorderIndexing(node *Node[V], currentIndex *int, nodeIndices map[*Node[V]]int) {
	if node == nil {
		return
	}
	tree.inorderIndexing(node.left, currentIndex, nodeIndices)
	nodeIndices[node] = *currentIndex
	*currentIndex++
	tree.inorderIndexing(node.right, currentIndex, nodeIndices)
}

func (tree *Tree[V]) print(node *Node[V], prefix string, isLeft bool, sb *strings.Builder, nodeIndices map[*Node[V]]int) {
	if node == nil {
		return
	}

	// 현재 노드의 색상과 값을 표시할 문자열 준비
	nodeColor := "⚫" // 검은색 노드
	if node.isRed {
		nodeColor = "🔴" // 빨간색 노드
	}

	// 오른쪽 서브트리 먼저 출력 (트리를 왼쪽으로 90도 회전한 형태로 출력)
	tree.print(node.right, prefix+(map[bool]string{true: "│   ", false: "    "})[isLeft], false, sb, nodeIndices)

	// 현재 노드 출력
	sb.WriteString(prefix)
	sb.WriteString(map[bool]string{true: "└── ", false: "┌── "}[isLeft])
	sb.WriteString(fmt.Sprintf("%v %d: %v, indexOf: %v, weight: %d\n", nodeColor, nodeIndices[node], node.value, tree.IndexOf(node), node.weight))

	// 왼쪽 서브트리 출력
	tree.print(node.left, prefix+(map[bool]string{true: "│   ", false: "    "})[isLeft], true, sb, nodeIndices)
}
