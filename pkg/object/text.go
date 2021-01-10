package object

import (
	"github.com/arovesto/sdl/pkg/game/global"
	"github.com/arovesto/sdl/pkg/math"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type text struct {
	shooterObject
	c   TextCallback
	cID global.ID
	txt string

	font *ttf.Font
}

type TextCallback func() (string, error)

func NewText(st Properties) GameObject {
	// TODO switch to render text based on picture, then remove this shit and use proper model
	f, _ := ttf.OpenFont("assets/minecraft.ttf", 30)
	return &text{shooterObject: newShooterObj(st), cID: st.Callback, font: f}
}

func (t *text) Update() (err error) {
	t.txt, err = t.c()
	return
}

func (t *text) Draw() error {
	if t.txt == "" {
		return nil
	}
	srf, err := t.font.RenderUTF8Solid(t.txt, sdl.Color{R: 255, G: 255, B: 255, A: 255})
	if err != nil {
		return err
	}
	txt, err := global.Renderer.CreateTextureFromSurface(srf)
	if err != nil {
		return err
	}
	srf.Free()
	return global.Renderer.CopyEx(txt, nil, math.SDLRect(t.model.Collider), 0, nil, sdl.FLIP_NONE)
}
