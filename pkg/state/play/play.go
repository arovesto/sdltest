package play

import (
	"github.com/arovesto/sdl/pkg/camera"
	"github.com/arovesto/sdl/pkg/collision"
	"github.com/arovesto/sdl/pkg/game/global"
	"github.com/arovesto/sdl/pkg/input"
	"github.com/arovesto/sdl/pkg/level"
	"github.com/arovesto/sdl/pkg/object"
	"github.com/arovesto/sdl/pkg/parser"
	"github.com/arovesto/sdl/pkg/state"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	stateID = 1
)

type play struct {
	objects []object.GameObject
	level   *level.Level

	canPause uint32
}

func init() {
	state.Play = &play{}
}

func (p *play) Update() (err error) {
	if input.IsKeyDown(sdl.SCANCODE_ESCAPE) && p.canPause >= 10 {
		if err = global.GetMachine().PushState(state.Pause); err != nil {
			return
		}
	}
	if p.canPause < 10 {
		p.canPause++
	}

	if err = camera.Update(); err != nil {
		return
	}
	collision.RunCollision()

	for _, o := range p.objects {
		if err = o.Update(); err != nil {
			return
		}
	}

	return p.level.Update()
}

func (p *play) Render() (err error) {
	return p.level.Render()
}

func (p *play) OnEnter() (err error) {
	p.level, err = parser.ParseLevel(global.MapPath)
	return
}

func (p *play) OnSwitch() error {
	p.canPause = 0
	return nil
}

func (p *play) OnContinue() error {
	return nil
}

func (p *play) OnExit() (err error) {
	for _, o := range p.objects {
		if err = o.Destroy(); err != nil {
			return err
		}
	}
	return
}

func (p *play) GetID() state.ID {
	return stateID
}
