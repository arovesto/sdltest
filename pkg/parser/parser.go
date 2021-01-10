package parser

import (
	"errors"
	"io/ioutil"

	"github.com/arovesto/sdl/pkg/model"

	"github.com/arovesto/sdl/pkg/game/global"
	"github.com/arovesto/sdl/pkg/math"

	"github.com/arovesto/sdl/pkg/object"
	"gopkg.in/yaml.v2"
)

type State struct {
	Objects []Object               `yaml:"objects"`
	Models  map[string]model.Model `yaml:"models"`
}

type Object struct {
	Type      string    `yaml:"type"`
	X         int       `yaml:"x"`
	Y         int       `yaml:"y"`
	Model     string    `yaml:"model"`
	Callback  global.ID `yaml:"callback"`
	AnimSpeed int       `yaml:"animspeed"`
	Flip      string    `yaml:"flip"`
}

func Parse(file, state string) (res []object.GameObject, err error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	s := map[string]State{}
	if err = yaml.Unmarshal(content, &s); err != nil {
		return
	}

	if st, ok := s[state]; ok {
		for _, o := range st.Objects {
			var obj object.GameObject
			obj, err = object.Create(o.Type, object.Properties{
				Pos:       math.NewVec(float64(o.X), float64(o.Y)),
				Model:     st.Models[o.Model].GetCopy(),
				AnimSpeed: uint32(o.AnimSpeed),
				Callback:  o.Callback,
			})
			if err != nil {
				return
			}
			res = append(res, obj)
		}
		return
	} else {
		return nil, errors.New("state not found")

	}
}
