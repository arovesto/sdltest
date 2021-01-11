package level

import (
	"github.com/arovesto/sdl/pkg/camera"
	"github.com/arovesto/sdl/pkg/math"
	"github.com/arovesto/sdl/pkg/object"
)

type enemy struct {
	shooterObject
}

func NewEnemy(st object.Properties) object.GameObject {
	obj := newShooterObj(st)
	return &enemy{shooterObject: obj}
}

func (e *enemy) GetType() object.Type {
	return EnemyType
}

func (e *enemy) Update() error {
	dist := math.AbsF(MainPlayer.pos.X - e.pos.X)
	var w float64
	switch {
	case dist <= 1000:
		w = 0.2
		e.model.Alpha = 255
	case dist <= 2000:
		w = 0.1
		e.model.Alpha = 100
	default:
		w = 0
		e.model.Alpha = 0
	}
	if w != 0 {
		camera.Camera.Targets[e.id] = camera.Target{Pos: e.pos, Weight: w}
	} else {
		delete(camera.Camera.Targets, e.id)
	}
	return e.shooterObject.Update()
}
