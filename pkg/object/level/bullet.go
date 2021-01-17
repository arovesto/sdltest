package level

import (
	"github.com/arovesto/sdl/pkg/level"
	"github.com/arovesto/sdl/pkg/math"
	"github.com/arovesto/sdl/pkg/model"
	"github.com/arovesto/sdl/pkg/object"
	"github.com/veandco/go-sdl2/sdl"
)

type bullet struct {
	shooterObject
}

func NewBullet(from, vel math.Vector2D) object.GameObject {
	m := model.AvailableModels["bullet"].GetCopy()
	delta := m.Center().FloatV()
	if vel.X < 0 {
		delta.X *= -1
	}
	return &bullet{shooterObject: shooterObject{pos: from.Sub(delta), vel: vel, model: m, maxSpeed: 100}}
}

func (b *bullet) GetType() object.Type {
	return object.BulletType
}

func (b *bullet) Update() error {
	b.model.Parts[0].Angle = math.AngleOn(math.ZeroVec(), b.vel)
	if b.vel.X < 0 {
		b.model.Parts[0].Flip = sdl.FLIP_HORIZONTAL
	}
	return b.shooterObject.Update()
}

func (b *bullet) BackOff(isGroundedP, isGroundedN, delta math.Vector2D) {
	if isGroundedP.X != 0 || isGroundedP.Y != 0 || isGroundedN.X != 0 || isGroundedN.Y != 0 {
		level.CurrentLevel.DelObject(b)
	}
}
