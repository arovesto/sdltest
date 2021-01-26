package math

import (
	"fmt"
	"math"
)

type Vector2D struct {
	X float64
	Y float64
}

func NewVec(x, y float64) Vector2D {
	return Vector2D{X: x, Y: y}
}

func NewVecInt(x, y int32) Vector2D {
	return NewVec(float64(x), float64(y))
}

func ZeroVec() Vector2D {
	return NewVec(0, 0)
}

func (v Vector2D) IntPos() (int32, int32) {
	return int32(v.X), int32(v.Y)
}

func (v Vector2D) IntVector() IntVector {
	return NewIntVector(v.IntPos())
}

func (v Vector2D) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v Vector2D) Normalized() Vector2D {
	return v.Div(v.Length())
}

func (v Vector2D) Add(v2 Vector2D) Vector2D {
	return Vector2D{X: v.X + v2.X, Y: v.Y + v2.Y}
}

func (v Vector2D) Mul(c float64) Vector2D {
	return Vector2D{X: v.X * c, Y: v.Y * c}
}

func (v Vector2D) Sub(which Vector2D) Vector2D {
	return Vector2D{X: v.X - which.X, Y: v.Y - which.Y}
}

func (v Vector2D) Div(c float64) Vector2D {
	return Vector2D{X: v.X / c, Y: v.Y / c}
}

func (v Vector2D) DivComponents(other Vector2D) Vector2D {
	return Vector2D{X: v.X / other.X, Y: v.Y / other.Y}
}

func (v Vector2D) String() string {
	return fmt.Sprintf("(%.2f, %.2f)", v.X, v.Y)
}

func Near(a, b, distance Vector2D) bool {
	return AbsF(a.X-b.X) <= distance.X && AbsF(a.Y-b.Y) <= distance.Y
}

/*
   a*----------------+
    |      *b        |
    |                |
    +----------------*c
*/
func Inside(a, b, c Vector2D) bool {
	return a.X < b.X && b.X < c.X && a.Y < b.Y && b.Y < c.Y
}

func InsideRect(r Rect, p IntVector) bool {
	return r.X < p.X && p.X < r.X+r.W && r.Y < p.Y && p.Y < r.Y+r.H
}

func Collide(a, b Rect) bool {
	aTop, aBottom, aLeft, aRight := a.Y, a.Y+a.H, a.X, a.X+a.W
	bTop, bBottom, bLeft, bRight := b.Y, b.Y+b.H, b.X, b.X+b.W
	return aBottom >= bTop && aTop <= bBottom && aRight >= bLeft && aLeft <= bRight
}

func ClampDirection(a Vector2D, blockUp, blockDown, blockLeft, blockRight bool) Vector2D {
	if blockRight && a.X > 0 {
		a.X = 0
	}
	if blockLeft && a.X < 0 {
		a.X = 0
	}

	if blockDown && a.Y > 0 {
		a.Y = 0
	}
	if blockUp && a.Y < 0 {
		a.Y = 0
	}
	return a
}

type IntVector struct {
	X int32 `yaml:"x"`
	Y int32 `yaml:"y"`
}

func (v IntVector) Mul(axis IntVector) IntVector {
	return IntVector{X: v.X * axis.X, Y: v.Y * axis.Y}
}

func (v IntVector) Div(axis IntVector) IntVector {
	return IntVector{X: v.X / axis.X, Y: v.Y / axis.Y}
}

func (v IntVector) Add(other IntVector) IntVector {
	return IntVector{X: v.X + other.X, Y: v.Y + other.Y}
}

func (v IntVector) Sub(other IntVector) IntVector {
	return IntVector{X: v.X - other.X, Y: v.Y - other.Y}
}

func (v Vector2D) Abs() float64 {
	return math.Abs(v.X) + math.Abs(v.Y)
}

func (v IntVector) Values() (int32, int32) {
	return v.X, v.Y
}

func (v IntVector) FloatV() Vector2D {
	return Vector2D{X: float64(v.X), Y: float64(v.Y)}
}

func NewIntVector(x, y int32) IntVector {
	return IntVector{X: x, Y: y}
}

type Rect struct {
	X int32 `yaml:"x"`
	Y int32 `yaml:"y"`
	W int32 `yaml:"w"`
	H int32 `yaml:"h"`
}

func (r Rect) Mul(axis IntVector) Rect {
	return Rect{X: r.X * axis.X, Y: r.Y * axis.Y, W: r.W * axis.X, H: r.H * axis.Y}
}

func (r Rect) Div(axis IntVector) Rect {
	return Rect{X: r.X / axis.X, Y: r.Y / axis.Y, W: r.W / axis.X, H: r.H / axis.Y}
}

func (r Rect) GetPos() IntVector {
	return IntVector{X: r.X, Y: r.Y}
}

func (r Rect) GetSize() IntVector {
	return IntVector{X: r.W, Y: r.H}
}

func (r Rect) Values() (int32, int32, int32, int32) {
	return r.X, r.Y, r.W, r.H
}

func NewRect(a, size IntVector) Rect {
	return Rect{X: a.X, Y: a.Y, W: size.X, H: size.Y}
}

func (r Rect) Add(v IntVector) Rect {
	return NewRect(r.GetPos().Add(v), r.GetSize())
}

func (r Rect) Empty() bool {
	return r.X == 0 && r.Y == 0 && r.W == 0 && r.H == 0
}
