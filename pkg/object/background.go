package object

import (
	"github.com/arovesto/sdl/pkg/camera"
	"github.com/arovesto/sdl/pkg/math"
	"github.com/arovesto/sdl/pkg/texturemanager"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	backgroundSpeed = 5
)

type background struct {
	shooterObject

	from1 sdl.Rect
	from2 sdl.Rect
	to1   sdl.Rect
	to2   sdl.Rect

	scrollLine int32
	width      int32

	prevCamPos int32
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
	b.scrollLine = width
	b.width = width
}

func (b *background) Draw() (err error) {
	if err = texturemanager.DrawRect(b.id, b.from1, b.to1); err != nil {
		return
	}
	return texturemanager.DrawRect(b.id, b.from2, b.to2)
}

func (b *background) Update() error {
	pos := int32(camera.GetCamPos().X)
	if math.Abs(pos-b.prevCamPos) < backgroundSpeed {
		return nil
	}

	b.scrollLine -= (pos - b.prevCamPos) / backgroundSpeed
	b.prevCamPos = pos
	if b.scrollLine > b.width {
		b.scrollLine -= b.width
	}
	if b.scrollLine < 0 {
		b.scrollLine = b.width - b.scrollLine
	}

	b.from1.X = b.width - b.scrollLine
	b.from1.W = b.scrollLine
	b.to1.W = b.scrollLine

	b.from2.W = b.width - b.scrollLine
	b.to2.W = b.width - b.scrollLine
	b.to2.X = b.scrollLine

	return nil
}
