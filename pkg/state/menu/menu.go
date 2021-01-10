package menu

import (
	"github.com/arovesto/sdl/pkg/game/global"
	"github.com/arovesto/sdl/pkg/object"
	"github.com/arovesto/sdl/pkg/parser"
	"github.com/arovesto/sdl/pkg/sound"
	"github.com/arovesto/sdl/pkg/state"
)

var callbacks = []object.Callback{
	func() error {
		return global.GetMachine().ChangeState(state.Play)
	},
	func() error {
		global.Quit()
		return nil
	},
}

type menu struct {
	objects []object.GameObject
}

func init() {
	state.Menu = &menu{}
}

func (m *menu) Update() (err error) {
	for _, o := range m.objects {
		if err = o.Update(); err != nil {
			return
		}
	}
	return
}

func (m *menu) Render() (err error) {
	for _, o := range m.objects {
		if err = o.Draw(); err != nil {
			return
		}
	}
	return
}

func (m *menu) OnSwitch() error {
	return nil
}

func (m *menu) OnContinue() error {
	return nil
}

func (m *menu) OnEnter() (err error) {
	m.objects, err = parser.Parse(global.MenusPath, "menu")
	if err != nil {
		return
	}

	object.SetCallbacks(m.objects, callbacks, nil)
	return sound.PlayMusic("heroes", -1)
}

func (m *menu) OnExit() (err error) {
	sound.HaltMusic()
	for _, o := range m.objects {
		if err = o.Destroy(); err != nil {
			return err
		}
	}
	return
}
