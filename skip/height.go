package skip

import "github.com/bytedance/gopkg/lang/fastrand"

func getNewHeight() int {
	return fastrand.Int() % MaxHeight
}
