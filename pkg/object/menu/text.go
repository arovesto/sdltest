package menu

import (
	"github.com/arovesto/sdl/pkg/game/global"
	"github.com/arovesto/sdl/pkg/math"
	"github.com/arovesto/sdl/pkg/object"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type text struct {
	menuObject
	c   TextCallback
	cID global.ID
	txt string

	font *ttf.Font
}

func NewText(st object.Properties) object.GameObject {
	// TODO switch to render text based on picture, then remove this shit and use proper model
	f, _ := ttf.OpenFont("assets/minecraft.ttf", 30)
	return &text{menuObject: newMenuObject(st), cID: st.Callback, font: f}
}

func (t *text) Update(tD float64) (err error) {
	t.txt, err = t.c()
	if err != nil {
		return err
	}
	return t.menuObject.Update(tD)
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
	return global.Renderer.CopyEx(txt, nil, math.SDLRect(t.model.Collider.Add(t.GetPosition().IntVector())), 0, nil, sdl.FLIP_NONE)
}
