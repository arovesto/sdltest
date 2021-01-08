package parser

import (
	"errors"
	"io/ioutil"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/arovesto/sdl/pkg/game/global"
	"github.com/arovesto/sdl/pkg/math"
	"github.com/arovesto/sdl/pkg/texturemanager"

	"github.com/arovesto/sdl/pkg/object"
	"gopkg.in/yaml.v2"
)

type FileStructure struct {
	States []State `yaml:"states"`
}

type State struct {
	ID       global.ID `yaml:"id"`
	Textures []Texture `yaml:"textures"`
	Fonts    []Font    `yaml:"fonts"`
	Objects  []Object  `yaml:"objects"`
}

type Object struct {
	Type      string    `yaml:"type"`
	X         int       `yaml:"x"`
	Y         int       `yaml:"y"`
	W         int       `yaml:"w"`
	H         int       `yaml:"h"`
	Frames    int       `yaml:"frames"`
	Texture   string    `yaml:"texture"`
	Callback  global.ID `yaml:"callback"`
	AnimSpeed int       `yaml:"animspeed"`
	Flip      string    `yaml:"flip"`
}

type Texture struct {
	ID     string `yaml:"id"`
	Path   string `yaml:"path"`
	Width  int    `yaml:"w"`
	Height int    `yaml:"h"`
}

type Font struct {
	ID   string `yaml:"id"`
	Path string `yaml:"path"`
	Size int    `yaml:"size"`
}

func Parse(file string, state global.ID) (res []object.GameObject, err error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	var s FileStructure
	if err = yaml.Unmarshal(content, &s); err != nil {
		return
	}

	for _, s := range s.States {
		if s.ID != state {
			continue
		}
		for _, t := range s.Textures {
			if err = texturemanager.Load(texturemanager.LoadOpts{Path: t.Path, ID: t.ID}); err != nil {
				return
			}
		}
		for _, f := range s.Fonts {
			if err = texturemanager.LoadFont(texturemanager.LoadFontOpts{ID: f.ID, Path: f.Path, Size: f.Size}); err != nil {
				return
			}
		}
		for _, o := range s.Objects {
			var obj object.GameObject
			obj, err = object.Create(o.Type, object.Properties{
				Pos:       math.NewVec(float64(o.X), float64(o.Y)),
				Size:      math.NewVec(float64(o.W), float64(o.H)),
				Cols:      int32(o.Frames),
				ID:        o.Texture,
				AnimSpeed: uint32(o.AnimSpeed),
				Callback:  o.Callback,
				Flip:      getFlip(o.Flip),
			})
			if err != nil {
				return
			}
			res = append(res, obj)
		}
		return
	}
	return nil, errors.New("state not found")
}

func getFlip(raw string) sdl.RendererFlip {
	switch raw {
	case "horizontal":
		return sdl.FLIP_HORIZONTAL
	case "vertical":
		return sdl.FLIP_VERTICAL
	case "both":
		return sdl.FLIP_HORIZONTAL | sdl.FLIP_HORIZONTAL
	default:
		return sdl.FLIP_NONE
	}
}
