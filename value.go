package intership_task_go

type Value struct {
	// 255 * 10 = 2550 << 2**31
	Sum         *int32
	BrokenCount *int32
	Count       uint32
}

func NewValue(count uint32) Value {
	sum := int32(0)
	brokenCount := int32(0)
	return Value{&sum, &brokenCount, count}
}
