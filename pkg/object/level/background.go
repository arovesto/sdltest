package level

import (
	"github.com/arovesto/sdl/pkg/camera"
	"github.com/arovesto/sdl/pkg/math"
	"github.com/arovesto/sdl/pkg/object"
)

const (
	backgroundSpeed = 11
)

type background struct {
	shooterObject

	scrollLine int32

	prevCamPos int32
	speed      int32
}

func NewBackground(st object.Properties) object.GameObject {
	if st.AnimSpeed == 0 {
		st.AnimSpeed = backgroundSpeed
	}
	res := &background{
		shooterObject: newShooterObj(st),
		speed:         int32(st.AnimSpeed),
	}
	res.init()
	return res
}

func (b *background) init() {
	_, _, width, height := camera.Camera.GetRect().Values()
	x := int32(b.pos.X)
	y := int32(b.pos.Y)

	b.model.Parts[0].OnTexture = math.Rect{W: b.model.Collider.W, H: b.model.Collider.H}
	b.model.Parts[0].OnModel = math.Rect{X: x, Y: y, W: width, H: height}
	b.model.Parts[1].OnTexture = math.Rect{H: b.model.Collider.H}
	b.model.Parts[1].OnModel = math.Rect{X: x + width, Y: y, H: height}

	b.scrollLine = width
}

func (b *background) Update() error {
	pos, _, width, height := camera.Camera.GetRect().Values()
	if math.Abs(pos-b.prevCamPos) >= b.speed {
		b.scrollLine -= (pos - b.prevCamPos) / b.speed
		b.prevCamPos = pos
	}

	if b.scrollLine > width {
		b.scrollLine -= width
	}
	if b.scrollLine < 0 {
		b.scrollLine = width - b.scrollLine
	}

	b.model.Parts[0].OnTexture.X = (width - b.scrollLine) * b.model.Collider.W / width
	b.model.Parts[0].OnTexture.W = b.scrollLine * b.model.Collider.W / width
	b.model.Parts[0].OnModel.W = b.scrollLine
	b.model.Parts[0].OnModel.H = height

	b.model.Parts[1].OnTexture.W = (width - b.scrollLine) * b.model.Collider.W / width
	b.model.Parts[1].OnModel.W = width - b.scrollLine
	b.model.Parts[1].OnModel.H = height
	b.model.Parts[1].OnModel.X = b.scrollLine

	return nil
}
