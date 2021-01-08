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

func AbsF(a float64) float64 {
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

func SignF(a float64) float64 {
	if a > 0 {
		return 1
	} else {
		return -1
	}
}

func DivRoundUp(a, b int32) int32 {
	return (a + b - 1) / b
}

func IsInside(x, y, w, h, xP, yP, margin int32) bool {
	return xP-x >= margin && x+w-xP >= margin && yP-y >= margin && y+h-yP >= margin
}

func LineBetween(y, yP, leftMargin, rightMargin int32) bool {
	if yP > y {
		return yP-y <= rightMargin
	} else {
		return y-yP <= leftMargin
	}
}
