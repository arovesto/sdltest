package object

import "errors"

type Creator func(p Properties) GameObject

type Factory struct {
	creators map[string]Creator
}

var f = NewFactory()

func init() {
	Register("button", NewButton)
	Register("player", NewPlayer)
	Register("enemy", NewEnemy)
	Register("animation", NewAnimation)
	Register("background", NewBackground)
	Register("text", NewText)
}

func NewFactory() *Factory {
	return &Factory{creators: map[string]Creator{}}
}

func Create(id string, s Properties) (GameObject, error) {
	if c, ok := f.creators[id]; ok {
		return c(s), nil
	} else {
		return nil, errors.New("id not found")
	}
}

func Register(id string, c Creator) {
	f.creators[id] = c
}
