package model

import (
	"github.com/arovesto/sdl/pkg/camera"
	"github.com/arovesto/sdl/pkg/game/global"
	"github.com/arovesto/sdl/pkg/math"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	alpha = 255
)

// TODO some kind ot relative model for scaling menus
type Model struct {
	// Model contains it's parts. Parts are rendered in order them appear in array
	Parts []*Part `yaml:"parts"`
	// relative to position collision rectangle. Should be used on collisions with entities
	EntityCollider math.Rect `yaml:"entity_collider"`
	// relative to position collision rectangle. Should be used on collisions with terrain
	Collider math.Rect `yaml:"collider"`
	// FIXME (not working) model rotation angle, around the center of Collider (collision is not affected by this)
	Angle float64 `yaml:"angle"`
	// model brightness setting
	Alpha uint8 `yaml:"alpha"`
	// should we ignore camera completely
	IgnoreCam bool `yaml:"ignore_cam"`
	// texture path
	Path string `yaml:"texture_path"`
	// flip value, 0 = none; 1 = horizontal; 2 = vertical; 3 = both
	Flip sdl.RendererFlip `yaml:"render_flip"`

	t *sdl.Texture
}

func (m Model) GetCopy() Model {
	n := m
	n.Parts = make([]*Part, 0, len(m.Parts))
	n.Alpha = alpha
	if n.EntityCollider.Empty() {
		n.EntityCollider = n.Collider
	}
	for _, p := range m.Parts {
		prt := *p
		if prt.OnModel.Empty() {
			prt.OnModel = n.Collider
		}
		if prt.OnTexture.Empty() {
			prt.OnTexture = n.Collider
		}
		n.Parts = append(n.Parts, &prt)
	}

	return n
}

// if to < 0 || to >= frames - go to next sprite
func (m *Model) ChangeSprites(to int32) {
	for _, s := range m.Parts {
		s.Frame++
		if s.Frame >= s.Frames {
			s.Frame = 0
		}
		if to >= 0 && to < s.Frames {
			s.Frame = to
		}
	}
}

func (m *Model) load() error {
	s, err := img.Load(m.Path)
	if err != nil {
		return err
	}
	sprite, err := global.Renderer.CreateTextureFromSurface(s)
	m.t = sprite
	s.Free()
	m.Path = ""
	return err
}

func (m *Model) Draw(where math.IntVector) error {
	if m.t == nil {
		if err := m.load(); err != nil {
			return err
		}
	}
	if !m.IgnoreCam {
		where = where.Sub(camera.Camera.Pos.IntVector())
	}
	for _, p := range m.Parts {
		if err := p.draw(m, where); err != nil {
			return err
		}
	}
	return nil
}

func (m *Model) Destroy() error {
	if m.t != nil {
		return m.t.Destroy()
	}
	return nil
}

type Part struct {
	// Rectangle used on real world
	OnModel math.Rect `yaml:"on_model"`
	// Rectangle used on texture place
	OnTexture math.Rect `yaml:"on_texture"`
	// point to rotate this part
	Pivot math.IntVector `yaml:"pivot"`
	// current rotation of this part
	Angle float64 `yaml:"angle"`
	// lightning of texture
	Alpha uint8 `yaml:"alpha"`
	// amount of frames in part's texture
	Frames int32 `yaml:"frames"`
	// current frame of animation
	Frame int32 `yaml:"frame"`
}

func (p Part) draw(m *Model, where math.IntVector) error {
	mod, err := m.t.GetAlphaMod()
	if err != nil {
		return err
	}
	if err = m.t.SetAlphaMod(p.Alpha + m.Alpha); err != nil {
		return err
	}
	// TODO not precise, need better method of combining rotations
	// TODO add flip heredddda
	p.OnTexture.X += p.OnTexture.W * p.Frame
	if err = global.Renderer.CopyEx(m.t, math.SDLRect(p.OnTexture), math.SDLRect(p.OnModel.Add(where)), p.Angle+m.Angle, math.SDLPoint(p.Pivot.Add(where)), m.Flip); err != nil {
		return err
	}
	return m.t.SetAlphaMod(mod)
}
