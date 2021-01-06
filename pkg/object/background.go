package object

import (
	"github.com/arovesto/sdl/pkg/texturemanager"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	animDelay   = 10
	scrollSpeed = 1
)

type background struct {
	shooterObject

	from1 sdl.Rect
	from2 sdl.Rect
	to1   sdl.Rect
	to2   sdl.Rect
}

func NewBackground(st Properties) GameObject {
	res := &background{
		shooterObject: newShooterObj(st),
	}
	res.init()
	return res
}

func (b *background) init() {
	height := int32(b.size.Y)
	width := int32(b.size.X)
	x := int32(b.pos.X)
	y := int32(b.pos.Y)
	b.from1 = sdl.Rect{W: width, H: height}
	b.to1 = sdl.Rect{X: x, Y: y, W: width, H: height}
	b.from2 = sdl.Rect{H: height}
	b.to2 = sdl.Rect{X: x + width, Y: y, H: height}
}

func (b *background) Draw() (err error) {
	if err = texturemanager.DrawRect(b.id, b.from1, b.to1); err != nil {
		return
	}
	return texturemanager.DrawRect(b.id, b.from2, b.to2)
}

func (b *background) Scroll(speed float64) {
	s := int32(speed / 5)
	if s == 0 {
		s = 1
	}
	b.from1.X += s
	b.from1.W -= s
	b.to1.W -= s

	b.from2.W += s
	b.to2.W += s
	b.to2.X -= s

	if b.to2.W > int32(b.size.X) {
		b.init()
	}
}
