package object

import (
	"github.com/veandco/go-sdl2/sdl"
)

type animation struct {
	shooterObject
	animChanged uint32
	animSpeed   uint32
}

func NewAnimation(st Properties) GameObject {
	return &animation{shooterObject: newShooterObj(st), animSpeed: st.AnimSpeed}
}

func (a *animation) Update() error {
	now := sdl.GetTicks()
	if now-a.animChanged > a.animSpeed {
		a.frame += 1
		if a.frame >= a.frames {
			a.frame = 0
		}
		a.animChanged = now
	}
	return a.shooterObject.Update()
}
