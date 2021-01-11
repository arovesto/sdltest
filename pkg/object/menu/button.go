package menu

import (
	"github.com/arovesto/sdl/pkg/game/global"
	"github.com/arovesto/sdl/pkg/input"
	"github.com/arovesto/sdl/pkg/math"
	"github.com/arovesto/sdl/pkg/object"
)

const (
	OUT = iota
	HOVER
	CLICKED
)

type button struct {
	menuObject
	c          Callback
	released   bool
	wasPressed bool
	cID        global.ID
}

func NewButton(st object.Properties) object.GameObject {
	return &button{menuObject: newMenuObject(st), cID: st.Callback}
}

func (b *button) Update() error {
	// TODO: delay can be introduced to allow button be visible UP before callback
	if b.wasPressed && b.released {
		b.wasPressed = false
		return b.c()
	}

	pos := input.GetMousePosition()
	if math.InsideRect(b.model.Collider.Add(b.GetPosition().IntVector()), pos) {
		pressed := input.GetMousePressed(input.LEFT)
		if pressed {
			b.wasPressed = true
			b.model.ChangeSprites(CLICKED)
			b.released = false
		} else {
			b.model.ChangeSprites(HOVER)
			b.released = true
		}
	} else {
		b.wasPressed = false
		b.model.ChangeSprites(OUT)
	}

	return b.menuObject.Update()
}
