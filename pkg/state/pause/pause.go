package pause

import (
	"github.com/arovesto/sdl/pkg/game/global"
	"github.com/arovesto/sdl/pkg/input"
	"github.com/arovesto/sdl/pkg/object"
	"github.com/arovesto/sdl/pkg/parser"
	"github.com/arovesto/sdl/pkg/state"
)

const (
	stateID = 2
)

var callbacks = []object.Callback{
	func() error {
		return global.GetMachine().PopState()
	},
	func() error {
		return global.GetMachine().ChangeState(state.Menu)
	},
}

type pause struct {
	objects []object.GameObject
}

func init() {
	state.Pause = &pause{}
}

func (m *pause) Update() (err error) {
	for _, o := range m.objects {
		if err = o.Update(); err != nil {
			return
		}
	}
	return
}

func (m *pause) Render() (err error) {
	for _, o := range m.objects {
		if err = o.Draw(); err != nil {
			return
		}
	}
	return
}

func (m *pause) OnEnter() (err error) {
	m.objects, err = parser.Parse(global.AssetsPath, stateID)
	if err != nil {
		return
	}

	object.SetCallbacks(m.objects, callbacks)

	return nil
}

func (m *pause) OnExit() (err error) {
	for _, o := range m.objects {
		if err = o.Destroy(); err != nil {
			return err
		}
	}
	input.ResetMouse()
	return
}

func (m *pause) GetID() state.ID {
	return stateID
}
