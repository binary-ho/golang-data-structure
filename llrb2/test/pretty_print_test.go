package test

import (
	"fmt"
	"go-data-structure/llrb2"
	"testing"
)

// StringKey는 LLRB 트리의 Key 인터페이스를 구현합니다
type StringKey string

type StringValue struct {
	Value string
}

func (v StringValue) String() string {
	return v.Value
}

func (k StringKey) Compare(other llrb2.Key) int {
	otherKey := other.(StringKey)
	if k < otherKey {
		return -1
	}
	if k > otherKey {
		return 1
	}
	return 0
}

func TestPrettyPrint(t *testing.T) {
	// 트리 생성
	tree := llrb2.NewTree[StringKey, StringValue]()

	fmt.Println("=== LLRB 트리 시각화 테스트 ===")
	fmt.Println("각 노드는 다음과 같이 표시됩니다:")
	fmt.Println("🔴: Red 노드")
	fmt.Println("⚫: Black 노드")
	fmt.Println()

	// 10번의 삽입 수행
	var nodes []*llrb2.Node[StringKey, StringValue]
	for i := 0; i < 10; i++ {
		key := StringKey(fmt.Sprintf("%d", i))
		value := fmt.Sprintf("값-%d", i)

		node := tree.Put(key, StringValue{value})
		nodes = append(nodes, node)

		fmt.Printf("\n=== %d번째 삽입 후 트리 상태 ===\n", i+1)
		fmt.Println(tree.PrintPrettyTree())
	}

	// 트리 구조체의 필드 출력
	fmt.Println("\n=== 최종 트리 정보 ===")
	fmt.Printf("트리 전체 크기: %d\n", tree.Size())
	fmt.Printf("루트 노드 키: %v\n", tree.Root().Key)
	fmt.Printf("루트 노드 값: %v\n", tree.Root().Value)
	fmt.Printf("루트 노드 색상: %v\n", map[bool]string{true: "Red", false: "Black"}[tree.Root().IsRed])

	for i, node := range nodes {
		fmt.Printf("노드 %d: %v\n", i, node)
	}

	node3 := nodes[3]
	tree.InsertAfter(node3, StringKey("10"), StringValue{"값-10"})
	fmt.Println("\n=== 3번 노드 뒤에 10번 노드 삽입 ===")
	fmt.Println(tree.PrintPrettyTree())

	node7 := nodes[7]
	tree.InsertAfter(node7, StringKey("11"), StringValue{"값-11"})
	fmt.Println("\n=== 7번 노드 뒤에 11번 노드 삽입 ===")
	fmt.Println(tree.PrintPrettyTree())
}
