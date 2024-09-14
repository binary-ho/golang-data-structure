package skip

import (
	"fmt"
)

const (
	ArrowBody   = "="
	ArrowHead   = ">"
	Arrow       = ArrowBody + ArrowHead
	LowestLevel = 0
	Padding     = 2
)

type visualizer[V Value] struct {
	list *List[V]
}

func (v *visualizer[V]) visualizer() (output string) {
	lowestLevel := v.getAllValues()

	for level := v.list.height - 1; level >= 0; level-- {
		output += fmt.Sprintf("LEVEL %02d ", level)
		for i, next := 0, v.list.head.tower[level]; next != nil; i, next = i+1, next.tower[level] {
			key := next.String()
			for key = next.String(); lowestLevel[i] != key; i++ {
				output += v.getArrowBody(len(lowestLevel[i]))
			}
			output += fmt.Sprintf("%v %v ", Arrow, key)
		}
		output += "\n"
	}

	return output
}

func (v *visualizer[V]) getAllValues() (lowestLevel []string) {
	node := v.list.head.tower[LowestLevel]
	for node != nil {
		lowestLevel = append(lowestLevel, node.String())
		node = node.tower[LowestLevel]
	}
	return lowestLevel
}

func (v *visualizer[V]) getArrowBody(contentLength int) (arrowBody string) {
	for count := contentLength + len(Arrow) + Padding; count > 0; count-- {
		arrowBody += ArrowBody
	}
	return arrowBody
}
