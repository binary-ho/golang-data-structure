package llrb2

func (t *Tree[K, V]) insertAtIndex(node *Node[K, V], idx int, key K, value V) *Node[K, V] {
	if node == nil {
		t.size++
		return NewNode(key, value, true)
	}

	leftSize := nodeSize(node.left)

	if idx <= leftSize {
		// 삽입 위치가 왼쪽 서브트리에 있음
		node.left = t.insertAtIndex(node.left, idx, key, value)
	} else {
		// 삽입 위치가 오른쪽 서브트리에 있음
		// idx를 조정: 왼쪽 서브트리 + 현재 노드 1개를 건너뛰었다
		node.right = t.insertAtIndex(node.right, idx-leftSize-1, key, value)
	}

	node = fixUp(node)
	return node
}

// RankOf returns the 0-based index (in-order) of `targetNode`.
func (t *Tree[K, V]) RankOf(targetNode *Node[K, V]) int {
	rank, found := rankOfRec(t.root, targetNode)
	if !found {
		return -1
	}
	return rank
}

// rankOfRec returns (rank, found).
// This is a BST search by key, plus an extra "are we the exact pointer?" check.
func rankOfRec[K Key, V Value](root *Node[K, V], target *Node[K, V]) (int, bool) {
	if root == nil || target == nil {
		return 0, false
	}

	cmp := target.key.Compare(root.key)
	if cmp < 0 {
		return rankOfRec(root.left, target)
	} else if cmp > 0 {
		rankRight, foundRight := rankOfRec(root.right, target)
		if !foundRight {
			return 0, false
		}
		// skip all left subtree + root itself
		return nodeSize(root.left) + 1 + rankRight, true
	} else {
		// cmp == 0 => possible match
		// check pointer identity
		if root == target {
			return nodeSize(root.left), true
		}
		// if there are duplicates or same key, we might check both sides:
		rL, fL := rankOfRec(root.left, target)
		if fL {
			return rL, true
		}
		rR, fR := rankOfRec(root.right, target)
		if fR {
			return nodeSize(root.left) + 1 + rR, true
		}
		return 0, false
	}
}

// InsertAfter inserts (key, value) so that it comes immediately
// after `prev` in the in-order ordering.
func (t *Tree[K, V]) InsertAfter(prev *Node[K, V], key K, value V) *Node[K, V] {
	if prev == nil {
		// 그냥 끝에 넣는 케이스로 정의할 수도 있음
		t.root = t.insertAtIndex(t.root, t.size, key, value)
		t.root.isRed = false
		return nil
	}

	r := t.RankOf(prev)
	if r < 0 {
		// prev not found in tree => handle error, or insert at end
		r = t.size - 1
	}

	t.root = t.insertAtIndex(t.root, r+1, key, value)
	if t.root != nil {
		t.root.isRed = false
	}

	// 여기서는 방금 삽입된 노드(pointer)를 찾으려면 직접 rank+1로 select해 오는 로직이 필요.
	// 간단히는 nil을 반환하거나, "newly inserted node" 포인터를 구하는 별도 selectByIndex 함수를 쓰면 됨.
	return nil
}
