package factory

import (
	"errors"

	"github.com/arovesto/sdl/pkg/object"

	"github.com/arovesto/sdl/pkg/object/level"
	"github.com/arovesto/sdl/pkg/object/menu"
)

type Creator func(p object.Properties) object.GameObject

type Factory struct {
	creators map[string]Creator
}

var f = NewFactory()

func init() {
	Register("button", menu.NewButton)
	Register("animation", menu.NewAnimation)
	Register("text", menu.NewText)
	Register("player", level.NewPlayer)
	Register("enemy", level.NewEnemy)
	Register("background", level.NewBackground)
}

func NewFactory() *Factory {
	return &Factory{creators: map[string]Creator{}}
}

func Create(id string, s object.Properties) (object.GameObject, error) {
	if c, ok := f.creators[id]; ok {
		return c(s), nil
	} else {
		return nil, errors.New("id not found")
	}
}

func Register(id string, c Creator) {
	f.creators[id] = c
}
