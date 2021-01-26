package gameover

import (
	"github.com/arovesto/sdl/pkg/game/global"
	"github.com/arovesto/sdl/pkg/input"
	"github.com/arovesto/sdl/pkg/object"
	"github.com/arovesto/sdl/pkg/object/menu"
	"github.com/arovesto/sdl/pkg/parser"
	"github.com/arovesto/sdl/pkg/state"
)

var callbacks = []menu.Callback{
	func() error {
		return global.GetMachine().ChangeState(state.Play)
	},
	func() error {
		return global.GetMachine().ChangeState(state.Menu)
	},
}

type gameover struct {
	objects []object.GameObject
}

func init() {
	state.GameOver = &gameover{}
}

func (g *gameover) Update() (err error) {
	for _, o := range g.objects {
		// TODO now mouse object really should be separated thing
		if err = o.Update(0); err != nil {
			return
		}
	}
	return
}

func (g *gameover) Render() (err error) {
	for _, o := range g.objects {
		if err = o.Draw(); err != nil {
			return
		}
	}
	return
}

func (g *gameover) OnEnter() (err error) {
	g.objects, err = parser.Parse(global.MenusPath, "game_over")
	if err != nil {
		return
	}

	menu.SetCallbacks(g.objects, callbacks, nil)

	return nil
}

func (g *gameover) OnSwitch() error {
	return nil
}

func (g *gameover) OnContinue() error {
	return nil
}

func (g *gameover) OnExit() (err error) {
	for _, o := range g.objects {
		if err = o.Destroy(); err != nil {
			return err
		}
	}
	input.ResetMouse()
	return
}
