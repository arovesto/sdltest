package object

import (
	"github.com/arovesto/sdl/pkg/game/global"
	"github.com/arovesto/sdl/pkg/texturemanager"
	"github.com/veandco/go-sdl2/sdl"
)

type text struct {
	shooterObject
	c   TextCallback
	cID global.ID
	txt string
}

type TextCallback func() (string, error)

func NewText(st Properties) GameObject {
	return &text{shooterObject: newShooterObj(st), cID: st.Callback}
}

func (t *text) Update() (err error) {
	t.txt, err = t.c()
	return
}

func (t *text) Draw() error {
	if t.txt == "" {
		return nil
	}
	return texturemanager.DrawMessage(texturemanager.DrawMessageOpts{
		ID:      t.id,
		Message: t.txt,
		Color:   sdl.Color{R: 255, G: 255, B: 255, A: 255},
		X:       int32(t.pos.X),
		Y:       int32(t.pos.Y),
		W:       int32(t.size.X),
		H:       int32(t.size.Y),
	})
}
