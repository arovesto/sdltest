package math

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

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

func AbsAngle(a AngleDeg) AngleDeg {
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

func SDLRect(r Rect) *sdl.Rect {
	return &sdl.Rect{X: r.X, Y: r.Y, W: r.W, H: r.H}
}

func SDLPoint(v IntVector) *sdl.Point {
	return &sdl.Point{X: v.X, Y: v.Y}
}

func ClampAngle(a, max AngleDeg) AngleDeg {
	if a > max {
		return max
	}
	if a < -max {
		return -max
	}
	return a
}

type AngleDeg float64 // clockwise from X axis

type AngleRad float64 // clockwise from X axis

func (a AngleRad) Deg() AngleDeg {
	return AngleDeg(180 * a / math.Pi)
}

func (a AngleDeg) ToVec() Vector2D {
	raw := float64(a.ToRad())
	return NewVec(math.Cos(raw), math.Sin(raw))
}

func (a AngleDeg) ToRad() AngleRad {
	return AngleRad(a * math.Pi / 180)
}

func AngleOn(from, to Vector2D) AngleDeg {
	return AngleRad(math.Atan((to.Y - from.Y) / math.Abs(to.X-from.X))).Deg()
}
