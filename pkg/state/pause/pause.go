package pause

import (
	"fmt"

	"github.com/arovesto/sdl/pkg/object/menu"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/arovesto/sdl/pkg/game/global"
	"github.com/arovesto/sdl/pkg/input"
	"github.com/arovesto/sdl/pkg/object"
	"github.com/arovesto/sdl/pkg/parser"
	"github.com/arovesto/sdl/pkg/sound"
	"github.com/arovesto/sdl/pkg/state"
	"github.com/veandco/go-sdl2/mix"
)

const (
	stateID = 2
)

var callbacks = []menu.Callback{
	func() error {
		return global.GetMachine().PopState()
	},
	func() error {
		return global.GetMachine().ChangeState(state.Menu)
	},
	func() error {
		sound.DecVolume()
		return nil
	},
	func() error {
		sound.IncVolume()
		return nil
	},
}

var textCallbacks = []menu.TextCallback{
	func() (string, error) {
		return fmt.Sprintf("sound %.2f%%", float64(sound.GetVolume())/float64(mix.MAX_VOLUME)), nil
	},
}

type pause struct {
	objects []object.GameObject

	canPause uint32
}

func init() {
	state.Pause = &pause{}
}

func (p *pause) Update() (err error) {
	if input.IsKeyDown(sdl.SCANCODE_ESCAPE) && p.canPause >= 20 {
		p.canPause = 0
		if err = global.GetMachine().PopState(); err != nil {
			return
		}
	}

	if p.canPause < 20 {
		p.canPause++
	}

	for _, o := range p.objects {
		// TODO now mouse object really should be separated thing
		if err = o.Update(0); err != nil {
			return
		}
	}
	return
}

func (p *pause) Render() (err error) {
	for _, o := range p.objects {
		if err = o.Draw(); err != nil {
			return
		}
	}
	return
}

func (p *pause) OnEnter() (err error) {
	p.objects, err = parser.Parse(global.MenusPath, "pause")
	if err != nil {
		return
	}

	menu.SetCallbacks(p.objects, callbacks, textCallbacks)

	return nil
}

func (p *pause) OnSwitch() error {
	return nil
}

func (p *pause) OnContinue() error {
	return nil
}

func (p *pause) OnExit() (err error) {
	p.canPause = 0
	for _, o := range p.objects {
		if err = o.Destroy(); err != nil {
			return err
		}
	}
	input.ResetMouse()
	return
}

func (p *pause) GetID() state.ID {
	return stateID
}
