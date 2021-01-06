package math

import (
	"math"
)

type Vector2D struct {
	X float64
	Y float64
}

func NewVec(X, Y float64) Vector2D {
	return Vector2D{X: X, Y: Y}
}

func NewVecInt(X, Y int32) Vector2D {
	return NewVec(float64(X), float64(Y))
}

func ZeroVec() Vector2D {
	return NewVec(0, 0)
}

func (v Vector2D) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v Vector2D) Normalized() Vector2D {
	return Div(v, v.Length())
}

func Add(v1, v2 Vector2D) Vector2D {
	return Vector2D{X: v1.X + v2.X, Y: v1.Y + v2.Y}
}

func Prod(v Vector2D, c float64) Vector2D {
	return Vector2D{X: v.X * c, Y: v.Y * c}
}

func Sub(from, which Vector2D) Vector2D {
	return Vector2D{X: from.X - which.X, Y: from.Y - which.Y}
}

func Div(v Vector2D, c float64) Vector2D {
	return Vector2D{X: v.X / c, Y: v.Y / c}
}

/*
   A*----------------+
    |      *B        |
    |                |
    +----------------*C
*/
func Inside(A, B, C Vector2D) bool {
	return A.X < B.X && B.X < C.X && A.Y < B.Y && B.Y < C.Y
}

func Collide(A, ASize, B, BSize Vector2D) bool {
	aTop, aBottom, aLeft, aRight := A.Y, A.Y+ASize.Y, A.X, A.X+ASize.X
	bTop, bBottom, bLeft, bRight := B.Y, B.Y+BSize.Y, B.X, B.X+BSize.X
	return aBottom >= bTop && aTop <= bBottom && aRight >= bLeft && aLeft <= bRight
}
