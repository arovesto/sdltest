package menu

import (
	"github.com/arovesto/sdl/pkg/game/global"
	"github.com/arovesto/sdl/pkg/object"
	menu2 "github.com/arovesto/sdl/pkg/object/menu"
	"github.com/arovesto/sdl/pkg/parser"
	"github.com/arovesto/sdl/pkg/sound"
	"github.com/arovesto/sdl/pkg/state"
)

var callbacks = []menu2.Callback{
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
		// TODO now mouse object really should be separated thing
		if err = o.Update(0); err != nil {
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

	menu2.SetCallbacks(m.objects, callbacks, nil)
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
