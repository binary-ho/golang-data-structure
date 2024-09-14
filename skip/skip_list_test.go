package skip

import (
	"fmt"
	"testing"
)

func TestSkipList(t *testing.T) {
	list := NewList[TestValue]()
	visualizer := visualizer[TestValue]{list: list}
	fmt.Println(visualizer.visualizer())

	printMessage("Insert Test")
	list.Insert(5, "E")
	fmt.Println(visualizer.visualizer())

	list.Insert(4, "D")
	fmt.Println(visualizer.visualizer())

	list.Insert(1, "A")
	fmt.Println(visualizer.visualizer())

	list.Insert(3, "C")
	fmt.Println(visualizer.visualizer())

	list.Insert(6, "F")
	fmt.Println(visualizer.visualizer())

	list.Insert(2, "B")
	fmt.Println(visualizer.visualizer())

	printMessage("Delete Test")
	list.Delete(4)
	fmt.Println(visualizer.visualizer())

	list.Delete(2)
	fmt.Println(visualizer.visualizer())

	printMessage("Insert Test")

	list.Insert(8, "H")
	fmt.Println(visualizer.visualizer())

	list.Insert(4, "DDDDD")
	fmt.Println(visualizer.visualizer())

	list.Insert(10, "J")
	fmt.Println(visualizer.visualizer())

	list.Insert(7, "G")
	fmt.Println(visualizer.visualizer())

	list.Insert(9, "I")
	fmt.Println(visualizer.visualizer())

	list.Insert(2, "BBBBB")
	fmt.Println(visualizer.visualizer())

	printMessage("Find Test")
	fmt.Println(list.Find(1))
	fmt.Println(list.Find(2))
	fmt.Println(list.Find(3))
	fmt.Println(list.Find(4))
	fmt.Println(list.Find(5))
	fmt.Println(list.Find(6))
	visualizer.visualizer()
}

type TestValue string

func (testValue TestValue) Len() int {
	return len(testValue)
}

func (testValue TestValue) String() string {
	return string(testValue)
}

func printMessage(message string) {
	fmt.Println()
	fmt.Printf("======================================== %s ========================================\n", message)
	fmt.Println()
}
