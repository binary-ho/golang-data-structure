package llrb3

import (
	"fmt"
	"strings"
)

// Print returns a visually appealing string representation of the tree
func (t *Tree[K, V]) Print() string {
	if t.root == nil {
		return "Empty tree"
	}

	// 중위순회로 노드의 인덱스 매핑 생성
	nodeIndices := make(map[*Node[K, V]]int)
	currentIndex := 0
	t.inorderIndexing(t.root, &currentIndex, nodeIndices)

	var sb strings.Builder
	t.print(t.root, "", true, &sb, nodeIndices)
	return sb.String()
}

// inorderIndexing performs an inorder traversal and assigns indices to nodes
func (t *Tree[K, V]) inorderIndexing(node *Node[K, V], currentIndex *int, nodeIndices map[*Node[K, V]]int) {
	if node == nil {
		return
	}
	t.inorderIndexing(node.left, currentIndex, nodeIndices)
	nodeIndices[node] = *currentIndex
	*currentIndex++
	t.inorderIndexing(node.right, currentIndex, nodeIndices)
}

func (t *Tree[K, V]) print(node *Node[K, V], prefix string, isLeft bool, sb *strings.Builder, nodeIndices map[*Node[K, V]]int) {
	if node == nil {
		return
	}

	// 현재 노드의 색상과 값을 표시할 문자열 준비
	nodeColor := "⚫" // 검은색 노드
	if node.isRed {
		nodeColor = "🔴" // 빨간색 노드
	}

	// 오른쪽 서브트리 먼저 출력 (트리를 왼쪽으로 90도 회전한 형태로 출력)
	t.print(node.right, prefix+(map[bool]string{true: "│   ", false: "    "})[isLeft], false, sb, nodeIndices)

	// 현재 노드 출력
	sb.WriteString(prefix)
	sb.WriteString(map[bool]string{true: "└── ", false: "┌── "}[isLeft])
	sb.WriteString(fmt.Sprintf("%v %d: %v\n", nodeColor, nodeIndices[node], node.value))

	// 왼쪽 서브트리 출력
	t.print(node.left, prefix+(map[bool]string{true: "│   ", false: "    "})[isLeft], true, sb, nodeIndices)
}
