package object

import "errors"

type Creator func(p Properties) GameObject

type Factory struct {
	creators map[string]Creator
}

var factory = NewFactory()

func NewFactory() *Factory {
	return &Factory{creators: map[string]Creator{}}
}

func (f *Factory) Create(id string, s Properties) (GameObject, error) {
	if c, ok := f.creators[id]; ok {
		return c(s), nil
	} else {
		return nil, errors.New("id not found")
	}
}

func (f *Factory) Register(id string, c Creator) {
	f.creators[id] = c
}

func Create(id string, s Properties) (GameObject, error) {
	return factory.Create(id, s)
}

func Register(id string, c Creator) {
	factory.Register(id, c)
}
