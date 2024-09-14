package skip

import "math/rand"

func getNewHeight() int {
	return rand.Int() % MaxHeight
}
