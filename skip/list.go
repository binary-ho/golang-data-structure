package skip

import (
	"errors"
)

/*
	필요한 것
	0. 노드의 인덱스를 어떻게 유일하게 식별할 것인가??????????? -> 타워 Weight의 합???
		-> 일단은 고려하지 않고 스킵 리스트를 만들어보자.
	1. Insertion, Insertion After
	2. Find
	3. Deletion
	4. 초기 구성
*/

type List[V Value] struct {
	head   *Node[V]
	height int
}

func NewList[V Value]() *List[V] {
	return &List[V]{
		head:   &Node[V]{},
		height: 1,
	}
}

func (list *List[V]) Find(key int) (*Node[V], error) {
	find, _ := list.get(key)
	if find == nil {
		return nil, errors.New("node not found")
	}
	return find, nil
}

func (list *List[V]) get(key int) (next *Node[V], route [MaxHeight]*Node[V]) {
	now := list.head
	for level := list.height; level >= 0; level-- {
		for next = now.tower[level]; next != nil; next = now.tower[level] {
			if key <= next.Weight() {
				break
			}
			now = next
		}
		route[level] = now
	}

	if next != nil && key == next.Weight() {
		return next, route
	}
	return nil, route
}

func (list *List[V]) Insert(key int, value V) {
	node, route := list.get(key)
	if node != nil {
		node.value = value
		return
	}

	height := getNewHeight()
	newNode := &Node[V]{value: value, key: key}

	for level := 0; level < height; level++ {
		prevNode := route[level]

		if prevNode == nil {
			prevNode = list.head
		}

		// 다음 노드를 이어준다.
		newNode.tower[level] = prevNode.tower[level]

		// 현재 노드의 다음이 새로운 노드
		prevNode.tower[level] = newNode
	}

	// Update height
	if height > list.height {
		list.height = height
	}
}

func (list *List[V]) Delete(key int) {
	node, route := list.get(key)
	if node == nil {
		return
	}

	// route 상의 노드가 삭제할 노드와 같은 동안 삭제
	for level := 0; level < list.height && route[level].tower[level] == node; level++ {
		route[level].tower[level] = node.tower[level]
		node.tower[level] = nil
	}

	node = nil
	list.shrink()
}

func (list *List[V]) shrink() {
	for level := list.height - 1; level >= 0; level-- {
		// head의 tower가 가리키는 Node가 없는 경우 shrink
		if list.head.tower[level] == nil {
			list.height--
		}
	}
}
