package math

func DivIfCan(a, b int32) int32 {
	absA, absB := Abs(a), Abs(b)
	if absA > absB {
		return a / b
	}
	if absA == 0 {
		return 0
	}
	return Sign(a)
}

func Abs(a int32) int32 {
	if a < 0 {
		return -a
	}
	return a
}

func Sign(a int32) int32 {
	if a > 0 {
		return 1
	} else {
		return -1
	}
}
