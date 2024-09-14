package splay

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

type Tree struct {
	root *Node
}

func (tree *Tree) Rotate(node *Node) {
	parent := node.parent
	if isNil(parent) {
		return
	}

	parent.pushLazyValue()
	node.pushLazyValue()

	node.setParentToChild()
	node.setGrandParentToParent()

	if isNil(node.parent) {
		tree.root = node
	}

	// update Child Count
	parent.updateTreeSizeAndSum()
	node.updateTreeSizeAndSum()
}

func (tree *Tree) Splay(node *Node) {
	if isNil(tree.root) || isNil(node) {
		//} node.isNil() || node == tree.root {
		return
	}

	if node.parent == tree.root {
		tree.Rotate(node)
		return
	}

	for isPresent(node.parent) {
		parent := node.parent
		grandParent := parent.parent

		if isPresent(grandParent) {
			if checkSameDirectionChildWithParent(node) {
				// Zig-Zig
				tree.Rotate(parent)
			} else {
				// Zig-Zag
				tree.Rotate(node)
			}
		}
		// 공통 rotate 작업
		tree.Rotate(node)
	}
}

func (tree *Tree) GetRangeSubtreeRootWithGather(start, end int64) *Node {
	tree.gather(start, end)

	subtreeRoot := tree.root.right.left
	subtreeRoot.pushLazyValue()
	return subtreeRoot
}

func (tree *Tree) gather(start, end int64) {
	tree.GetKthNodeAndPush(end + 1)
	endNode := tree.root

	tree.GetKthNodeAndPush(start - 1)
	startNode := tree.root

	//tree.PrintDFS()
	tree.splayAndSetChild(startNode, endNode)
}

func (tree *Tree) splayAndSetChild(rootNode *Node, child *Node) {
	// TODO : 여기 의심
	if isNil(tree.root) || isNil(child) {
		return
	}

	for child.parent != rootNode && isPresent(child.parent) {
		parent := child.parent
		grandParent := parent.parent

		if grandParent == rootNode {
			tree.Rotate(child)
			break
		}

		if checkSameDirectionChildWithParent(child) {
			// Zig-Zig
			tree.Rotate(parent)
			tree.Rotate(child)
		} else {
			// Zig-Zag
			tree.Rotate(child)
			tree.Rotate(child)
		}
	}

	tree.root.updateTreeSizeAndSum()
	tree.root.pushLazyValue()
	rootNode.updateTreeSizeAndSum()
	rootNode.pushLazyValue()
	child.updateTreeSizeAndSum()
	child.pushLazyValue()
	if isNil(rootNode) {
		tree.root = child
	}
}

func (tree *Tree) Find(key int64) *Node {
	//fmt.Println("Find key: ", key)
	if isNil(tree.root) {
		return nil
	}

	node, parent := tree.findNodeAndParent(key)
	if isPresent(node) {
		tree.Splay(node)
		return node
	} else {
		tree.Splay(parent)
		return nil
	}
}

func (tree *Tree) findNodeAndParent(key int64) (node, parent *Node) {
	node = tree.root
	for isPresent(node) && key != node.key {
		parent = node
		if key < node.key {
			node = node.left
		} else {
			node = node.right
		}
	}

	return node, parent
}

func (tree *Tree) Insert(key int64) {
	//fmt.Println("Insert key: ", key)
	if isNil(tree.root) {

		tree.root = &Node{key: key, size: 1, sum: key}
		return
	}

	_, parent := tree.findNodeAndParent(key)

	newNode := &Node{key: key, parent: parent, size: 1, sum: key}
	if key < parent.key {
		parent.left = newNode
	} else {
		parent.right = newNode
	}
	tree.splayAndSetChild(nil, newNode)
}

func (tree *Tree) Delete(key int64) {
	//fmt.Println("Delete key: ", key)
	if isNil(tree.Find(key)) {
		return
	}

	switch root := tree.root; true {
	case isPresent(root.left) && isPresent(root.right):
		tree.root = root.left
		tree.root.parent = nil

		node := tree.root
		for isPresent(node.right) {
			node = node.right
		}
		node.right = root.right
		root.right.parent = node

	case isPresent(root.left):
		tree.root = root.left
		tree.root.parent = nil

	case isPresent(root.right):
		tree.root = root.right
		tree.root.parent = nil

	default:
		tree.root = nil
	}
}

func (tree *Tree) SumRange(start, end, value int64) {
	subtreeRoot := tree.GetRangeSubtreeRootWithGather(start, end)
	tree.root.updateTreeSizeAndSum()
	if isNil(subtreeRoot) {
		return
	}
	subtreeRoot.sum += subtreeRoot.size * value
	subtreeRoot.lazy += value
}

func (tree *Tree) GetKthNode(k int64) int64 {
	//fmt.Printf("Get %d th Node", k)
	k -= 1
	node := tree.root
	for isPresent(node) {
		for isPresent(node.left) && node.left.size > k {
			node = node.left
		}

		if isPresent(node.left) {
			k -= node.left.size
		}

		if k == 0 {
			break
		}

		k--
		node = node.right
	}

	tree.Splay(node)
	//fmt.Printf(" -> %d\n", node.key)
	return node.key
}

func (tree *Tree) GetKthNodeAndPush(k int64) {
	//fmt.Printf("Get %d th Node And Push\n", k)
	//k -= 1	// 더미 떄문에 삭제
	node := tree.root
	node.pushLazyValue()

	for isPresent(node) {
		for isPresent(node.left) && node.left.size > k {
			node = node.left
			node.pushLazyValue()
		}

		if isPresent(node.left) {
			k -= node.left.size
		}

		if k == 0 {
			break
		}

		k--
		node = node.right
		if isPresent(node) {
			node.pushLazyValue()
		}
	}

	tree.splayAndSetChild(nil, node)
	//fmt.Printf(" -> %d\n", node.key)
	return
}

func checkSameDirectionChildWithParent(node *Node) bool {
	parent := node.parent
	grandParent := parent.parent

	isNodeLeft := node == parent.left
	isParentLeft := parent == grandParent.left
	return isNodeLeft == isParentLeft
}

func (tree *Tree) PrintDFS() {
	printDFS(tree.root, 0, "root")
	fmt.Println() // 보기 편하려고 넣음
}

func printDFS(node *Node, level int, direction string) {
	if isNil(node) {
		return
	}

	fmt.Printf("%s[%s] : Node %d\n", getIndent(level), direction, node.key)

	printDFS(node.left, level+1, "left")
	printDFS(node.right, level+1, "right")
}

// 출력 시 들여쓰기를 위한 함수
func getIndent(level int) string {
	indent := ""
	for i := 0; i < level; i++ {
		indent += "  "
	}
	return indent
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)

	var N, M, K int
	//fmt.Fscanf(reader, "%d %d %d", &maxValue, &M, &K)
	fmt.Fscanln(reader, &N, &M, &K)

	tree := &Tree{root: nil}
	tree.Insert(math.MinInt32)
	tree.root.updateTreeSizeAndSum()
	for i := 0; i < N; i++ {
		var value int64
		//fmt.Fscanf(reader, "%d", &value)
		fmt.Fscanln(reader, &value)
		//fmt.Println(value)
		tree.Insert(value)
		tree.root.updateTreeSizeAndSum()
	}
	tree.Insert(math.MaxInt32)
	tree.root.updateTreeSizeAndSum()

	const (
		SUM     = 1
		GET_SUM = 2
	)

	for i := 0; i < M+K; i++ {
		var command, start, end, value int64
		//fmt.Fscanf(reader, "%d %d %d", &command, &start, &end)
		fmt.Fscanln(reader, &command, &start, &end, &value)
		//fmt.Printf("%d %d %d %d\n", command, start, end, value)

		switch command {
		case SUM:
			//fmt.Fscanf(reader, "%d", &value)
			tree.SumRange(start, end, value)
			//fmt.Fprintf(writer, "%d %d %d %d\n", command, start, end, value)

		case GET_SUM:
			subtreeRoot := tree.GetRangeSubtreeRootWithGather(start, end)
			fmt.Fprintf(writer, "%d\n", subtreeRoot.sum)
			//fmt.Println(subtreeRoot.sum)
			//fmt.Fprintf(writer, "%d %d %d %d\n", command, start, end, value)
		}
	}

	writer.Flush()
}
