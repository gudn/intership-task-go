package value

import (
	"sync/atomic"
)

type Value struct {
	// 255 * 10 = 2550 << 2**31
	Sum         *int32
	BrokenCount *int32
	Count       float64
}

func (v *Value) Average() (avg float64, broken bool) {
	brokenCount := atomic.LoadInt32(v.BrokenCount)
	sum := atomic.LoadInt32(v.Sum)
	if brokenCount != 0 {
		broken = true
	}
	avg = float64(sum) / v.Count
	return
}

func New(count int) *Value {
	sum := int32(0)
	brokenCount := int32(0)
	return &Value{&sum, &brokenCount, float64(count)}
}
