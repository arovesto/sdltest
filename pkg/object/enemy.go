package object

import (
	"github.com/veandco/go-sdl2/sdl"
)

type enemy struct {
	shooterObject
}

func NewEnemy(st Properties) GameObject {
	obj := newShooterObj(st)
	//obj.vel = math.NewVec(2, 0)
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
	//if e.pos.X < 500 || e.pos.X > 2000 {
	//	e.vel.X *= -1
	//	if e.vel.X > 0 {
	//		e.flip = sdl.FLIP_NONE
	//	} else {
	//		e.flip = sdl.FLIP_HORIZONTAL
	//	}
	//}
	return e.shooterObject.Update()
}

func (e *enemy) GetType() Type {
	return Enemy
}

func (e *enemy) Collide(o GameObject) {
	// TODO store this metadata in texture manager
	//e.id = "largeexplosion"
	//e.size = math.NewVec(128, 128)
	//e.frame = 0
	//e.frames = 9
	e.shooterObject.Collide(o)
}
