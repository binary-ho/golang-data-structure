package llrb4

import (
	"fmt"
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

func (v StringValue) Len() int {
	return 1
}

func (k StringKey) Compare(other Key) int {
	otherKey := other.(StringKey)
	if k < otherKey {
		return -1
	}
	if k > otherKey {
		return 1
	}
	return 0
}

//
//func TestPut(t *testing.T) {
//	// 트리 생성
//	tree := NewTree[StringKey, StringValue]()
//
//	fmt.Println("=== LLRB 트리 시각화 테스트 ===")
//	fmt.Println("각 노드는 다음과 같이 표시됩니다:")
//	fmt.Println("🔴: Red 노드")
//	fmt.Println("⚫: Black 노드")
//	fmt.Println()
//
//	// 10번의 삽입 수행
//	for i := 0; i < 10; i++ {
//		value := fmt.Sprintf("값-%d", i)
//		newNode := NewNode(StringValue{value}, true)
//		tree.InsertAfter(nil, newNode)
//
//		fmt.Printf("\n=== %d번째 삽입 후 트리 상태 ===\n", i+1)
//		fmt.Println(tree.Print())
//	}
//
//	// 트리 구조체의 필드 출력
//	fmt.Println("\n=== 최종 트리 정보 ===")
//	fmt.Printf("트리 전체 크기: %d\n", tree.size)
//	fmt.Printf("루트 노드 키: %v\n", tree.root.key)
//	fmt.Printf("루트 노드 값: %v\n", tree.root.value)
//	fmt.Printf("루트 노드 색상: %v\n", map[bool]string{true: "Red", false: "Black"}[tree.root.isRed])
//}

func TestInsertAfter(t *testing.T) {
	// 트리 생성
	tree := NewTree[StringValue]()

	fmt.Println("=== LLRB 트리 시각화 테스트 ===")
	fmt.Println("각 노드는 다음과 같이 표시됩니다:")
	fmt.Println("🔴: Red 노드")
	fmt.Println("⚫: Black 노드")
	fmt.Println()

	// 10번의 삽입 수행
	var nodes []*Node[StringValue]
	for i := 0; i < 10; i++ {
		//key := StringKey(fmt.Sprintf("%d", i))
		value := fmt.Sprintf("값-%d", i)
		newNode := NewNode(StringValue{value}, true)
		insertedNode := tree.InsertAfter(nil, newNode)
		nodes = append(nodes, insertedNode)

		fmt.Printf("\n=== %d번째 삽입 후 트리 상태 ===\n", i+1)
		fmt.Println(tree.Print())
	}

	// 트리 구조체의 필드 출력
	fmt.Println("\n=== 최종 트리 정보 ===")
	fmt.Printf("트리 전체 크기: %d\n", tree.size)
	fmt.Printf("루트 노드 키: %v\n", tree.root.key)
	fmt.Printf("루트 노드 값: %v\n", tree.root.value)
	fmt.Printf("루트 노드 색상: %v\n", map[bool]string{true: "Red", false: "Black"}[tree.root.isRed])

	for i, node := range nodes {
		fmt.Printf("노드 %d: %v\n", i, node)
	}

	//node0 := nodes[0]
	node1 := nodes[1]
	//node2 := nodes[2]
	node3 := nodes[3]
	//node4 := nodes[4]
	//node5 := nodes[5]
	node6 := nodes[6]
	node7 := nodes[7]
	//node8 := nodes[8]
	//node9 := nodes[9]

	newNode10 := tree.InsertAfter(node3, NewNode(StringValue{"값-10"}, true))
	nodes = append(nodes, newNode10)
	fmt.Println("\n=== 3번 노드 뒤에 10번 노드 삽입 ===")
	fmt.Println(tree.Print())

	newNode11 := tree.InsertAfter(node7, NewNode(StringValue{"값-11"}, true))
	nodes = append(nodes, newNode11)
	fmt.Println("\n=== 7번 노드 뒤에 11번 노드 삽입 ===")
	fmt.Println(tree.Print())

	newNode12 := tree.InsertAfter(nodes[10], NewNode(StringValue{"값-12"}, true))
	nodes = append(nodes, newNode12)
	fmt.Println("\n=== 10번 노드 뒤에 12번 노드 삽입 ===")
	fmt.Println(tree.Print())

	newNode13 := tree.InsertAfter(node1, NewNode(StringValue{"값-13"}, true))
	nodes = append(nodes, newNode13)
	fmt.Println("\n=== 1번 노드 뒤에 13번 노드 삽입 ===")
	fmt.Println(tree.Print())

	newNode14 := tree.InsertAfter(node6, NewNode(StringValue{"값-14"}, true))
	nodes = append(nodes, newNode14)
	fmt.Println("\n=== 6번 노드 뒤에 14번 노드 삽입 ===")
	fmt.Println(tree.Print())

	newNode15 := tree.InsertAfter(node3, NewNode(StringValue{"값-15"}, true))
	nodes = append(nodes, newNode15)
	fmt.Println("\n=== 3번 노드 뒤에 15번 노드 삽입 ===")
	fmt.Println(tree.Print())

	newNode16 := tree.InsertAfter(nodes[10], NewNode(StringValue{"값-16"}, true))
	nodes = append(nodes, newNode16)
	fmt.Println("\n=== 10번 노드 뒤에 16번 노드 삽입 ===")
	fmt.Println(tree.Print())
}
