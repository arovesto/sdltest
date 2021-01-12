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
	Angle math.AngleDeg `yaml:"angle"`
	// model brightness setting
	Alpha uint8 `yaml:"alpha"`
	// should we ignore camera completely
	IgnoreCam bool `yaml:"ignore_cam"`
	// texture path
	Path string `yaml:"texture_path"`
	// window size this model was created for
	BaseSize math.IntVector `yaml:"base_size"`
	// !TBA! shell all model be flipped around center of collider
	GlobalFlip sdl.RendererFlip `yaml:"global_flip"`

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
		switch {
		case prt.Frames > 1 && len(prt.OnTexture) == 1:
			if len(prt.OnModel) == 0 {
				prt.OnModel = make([]math.Rect, 0, prt.Frames)
				prt.OnModel = append(prt.OnModel, n.Collider)
			}
			for i := 1; i < prt.Frames; i++ {
				prt.OnTexture = append(prt.OnTexture, prt.OnTexture[0].Add(math.NewIntVector(n.Collider.W*int32(i), 0)))
				prt.OnModel = append(prt.OnModel, prt.OnModel[0])
			}
		case prt.Frames != 0:
			prt.OnModel = make([]math.Rect, p.Frames)
			prt.OnTexture = make([]math.Rect, p.Frames)
			for i := range prt.OnTexture {
				prt.OnTexture[i] = n.Collider.Add(math.NewIntVector(n.Collider.W*int32(i), 0))
				prt.OnModel[i] = n.Collider
			}
		}

		if prt.CenterPoint.X == 0 && prt.CenterPoint.Y == 0 {
			prt.CenterPoint = n.Center()
		}
		if prt.Pivot.X == 0 && prt.Pivot.Y == 0 {
			prt.Pivot = n.Center()
		}
		n.Parts = append(n.Parts, &prt)
	}
	return n
}

func (m *Model) Center() math.IntVector {
	return m.Collider.GetPos().Add(m.Collider.GetSize()).Div(math.NewIntVector(2, 2))
}

// if to < 0 || to >= frames - go to next sprite
func (m *Model) ChangeSprites(to int) {
	for _, s := range m.Parts {
		s.ChangeSprites(to)
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
	OnModel []math.Rect `yaml:"on_model"`
	// Rectangle used on texture place
	OnTexture []math.Rect `yaml:"on_texture"`
	// point to rotate this part around
	Pivot math.IntVector `yaml:"pivot"`
	// current rotation of this part
	Angle math.AngleDeg `yaml:"angle"`
	// maximum allowed angle
	MaxAngle math.AngleDeg `yaml:"max_angle"`
	// minimum allowed angle
	MinAngle math.AngleDeg `yaml:"min_angle"`
	// lightning of texture
	Alpha uint8 `yaml:"alpha"`
	// current frame of animation
	Frame int `yaml:"frame"`
	// used if all frames are copy of each other but moved right
	Frames int `yaml:"frames"`
	// should part be in flipped state, 0 = none; 1 = horizontal; 2 = vertical; 3 = both
	Flip sdl.RendererFlip `yaml:"flip"`
	// point to be flipped around
	CenterPoint math.IntVector `yaml:"center"`
	// should model be drawn
	Hidden bool `yaml:"hidden"`
}

func (p Part) draw(m *Model, where math.IntVector) error {
	if p.Hidden {
		return nil
	}
	mod, err := m.t.GetAlphaMod()
	if err != nil {
		return err
	}
	if err = m.t.SetAlphaMod(p.Alpha + m.Alpha); err != nil {
		return err
	}
	// TODO not precise, need better method of combining rotations
	// TODO add flip here
	rect := p.OnModel[p.Frame]
	if p.Flip&sdl.FLIP_HORIZONTAL != 0 {
		cX := p.CenterPoint.X
		rect.X = cX - (rect.X - cX) - rect.W
		p.Pivot.X = p.OnModel[p.Frame].W - p.Pivot.X
		tmp := p.MaxAngle
		p.MaxAngle = -p.MinAngle
		p.MinAngle = -tmp
		p.Angle = -p.Angle
	}

	if p.Angle > p.MaxAngle && p.MaxAngle-p.MinAngle != 0 {
		p.Angle = p.MaxAngle
	}
	if p.Angle < p.MinAngle && p.MaxAngle-p.MinAngle != 0 {
		p.Angle = p.MinAngle
	}

	if err = global.Renderer.CopyEx(m.t, math.SDLRect(p.OnTexture[p.Frame]), math.SDLRect(rect.Add(where)), float64(p.Angle+m.Angle), math.SDLPoint(p.Pivot), p.Flip); err != nil {
		return err
	}
	return m.t.SetAlphaMod(mod)
}

func (p Part) GetPivotPoint(where math.IntVector) math.IntVector {
	rect := p.OnModel[p.Frame]
	if p.Flip&sdl.FLIP_HORIZONTAL != 0 {
		cX := p.CenterPoint.X
		rect.X = cX - (rect.X - cX) - rect.W
		p.Pivot.X = p.OnModel[p.Frame].W - p.Pivot.X
	}
	return rect.GetPos().Add(where).Add(p.Pivot)
}

// if to < 0 || to >= frames - go to next sprite
func (p *Part) ChangeSprites(to int) {
	p.Frame++
	if p.Frame >= len(p.OnModel) {
		p.Frame = 0
	}
	if to >= 0 && to < len(p.OnModel) {
		p.Frame = to
	}
}
