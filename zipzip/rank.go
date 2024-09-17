package zipzip

import "github.com/bytedance/gopkg/lang/fastrand"

func getRank() (rank int) {
	for fastrand.Int()&1 == 0 {
		rank++
	}
	return rank
}
