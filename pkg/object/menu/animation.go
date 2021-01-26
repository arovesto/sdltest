package menu

import (
	"github.com/arovesto/sdl/pkg/object"
	"github.com/veandco/go-sdl2/sdl"
)

type animation struct {
	menuObject
	animChanged uint32
	animSpeed   uint32
}

func NewAnimation(st object.Properties) object.GameObject {
	return &animation{menuObject: newMenuObject(st), animSpeed: st.AnimSpeed}
}

func (a *animation) Update(t float64) error {
	now := sdl.GetTicks()
	if now-a.animChanged > a.animSpeed {
		a.animChanged = now
		a.model.ChangeSprites(-1)
	}
	return a.menuObject.Update(t)
}
