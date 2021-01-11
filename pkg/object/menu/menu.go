package menu

import (
	"github.com/arovesto/sdl/pkg/game/global"
	"github.com/arovesto/sdl/pkg/math"
	"github.com/arovesto/sdl/pkg/model"
	"github.com/arovesto/sdl/pkg/object"
)

const (
	NoTypeMenuObject object.Type = iota
)

type Callback func() error
type TextCallback func() (string, error)

func SetCallbacks(obj []object.GameObject, clb []Callback, tlb []TextCallback) {
	for _, o := range obj {
		if btn, ok := o.(*button); ok {
			btn.c = clb[btn.cID]
		}
		if txt, ok := o.(*text); ok {
			txt.c = tlb[txt.cID]
		}
	}
}

type menuObject struct {
	relPos math.Vector2D
	model  model.Model

	screen math.IntVector
}

func newMenuObject(st object.Properties) menuObject {
	st.Model.IgnoreCam = true
	m := menuObject{relPos: st.Pos, model: st.Model, screen: math.NewIntVector(global.GetSize())}
	if st.Model.BaseSize.X != 0 && st.Model.BaseSize.Y != 0 {
		m.updateModelAfterScaleChange(m.model.BaseSize, m.screen)
	}
	return m
}

func (m *menuObject) Draw() error {
	w, h := int32(float64(m.screen.X)*m.relPos.X), int32(float64(m.screen.Y)*m.relPos.Y)
	return m.model.Draw(math.NewIntVector(w, h))
}

func (m *menuObject) Destroy() error {
	return m.model.Destroy()
}

func (m *menuObject) Collide(other object.GameObject) error {
	return nil
}

func (m *menuObject) GetPosition() math.Vector2D {
	return math.NewVec(float64(m.screen.X)*m.relPos.X, float64(m.screen.Y)*m.relPos.Y)
}

func (m *menuObject) GetCollider() math.Rect {
	return math.Rect{}
}

func (m *menuObject) GetObjectCollider() math.Rect {
	return math.Rect{}
}

func (m *menuObject) GetType() object.Type {
	return NoTypeMenuObject
}

func (m *menuObject) BackOff(isGroundedP, isGroundedN, delta math.Vector2D) {}

func (m *menuObject) updateModelAfterScaleChange(oldScale, newScale math.IntVector) {
	m.model.Collider = m.model.Collider.Mul(newScale).Div(oldScale)
	for _, p := range m.model.Parts {
		p.Pivot = p.Pivot.Mul(newScale).Div(oldScale)
		p.OnModel = p.OnModel.Mul(newScale).Div(oldScale)
	}
}

func (m *menuObject) Update() error {
	w, h := global.GetSize()
	if w != m.screen.X || h != m.screen.Y {
		old := m.screen
		m.screen = math.NewIntVector(w, h)
		m.updateModelAfterScaleChange(old, m.screen)
	}
	return nil
}
