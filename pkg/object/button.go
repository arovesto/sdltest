package object

import (
	"github.com/arovesto/sdl/pkg/game/global"
	"github.com/arovesto/sdl/pkg/input"
	"github.com/arovesto/sdl/pkg/math"
)

const (
	OUT = iota
	HOVER
	CLICKED
)

type button struct {
	shooterObject
	c          Callback
	released   bool
	wasPressed bool
	cID        global.ID
}

type Callback func() error

func NewButton(st Properties) GameObject {
	return &button{shooterObject: newShooterObj(st), cID: st.Callback}
}

func SetCallbacks(obj []GameObject, clb []Callback) {
	for _, o := range obj {
		if btn, ok := o.(*button); ok {
			btn.c = clb[btn.cID]
		}
	}
}

func (b *button) Update() error {
	// TODO: delay can be introduced to allow button be visible UP before callback
	if b.wasPressed && b.released {
		b.wasPressed = false
		return b.c()
	}

	pos := input.GetMousePosition()
	if math.Inside(b.pos, pos, math.Add(b.pos, b.size)) {
		pressed := input.GetMousePressed(input.LEFT)
		if pressed {
			b.wasPressed = true
			b.frame = CLICKED
			b.released = false
		} else {
			b.frame = HOVER
			b.released = true
		}
	} else {
		b.wasPressed = false
		b.frame = OUT
	}

	return nil
}
