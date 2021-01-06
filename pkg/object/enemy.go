package object

import (
	"github.com/arovesto/sdl/pkg/math"
	"github.com/veandco/go-sdl2/sdl"
)

type enemy struct {
	shooterObject
}

func NewEnemy(st Properties) GameObject {
	obj := newShooterObj(st)
	obj.vel = math.NewVec(0, 2)
	return &enemy{shooterObject: obj}
}

func (e *enemy) Update() error {
	now := sdl.GetTicks()
	if now-e.spriteChanged > animationTime {
		e.spriteChanged = now
		e.frame += 1
		if e.frame >= e.frames {
			e.frame = 0
		}
	}
	if e.pos.Y < 0 || e.pos.Y > 1000 {
		e.vel.Y *= -1
	}
	return e.shooterObject.Update()
}

func (e *enemy) GetType() Type {
	return Enemy
}

func (e *enemy) Collide() error {
	// TODO store this metadata in texture manager
	e.id = "largeexplosion"
	e.size = math.NewVec(128, 128)
	e.frame = 0
	e.frames = 9
	return e.shooterObject.Collide()
}
