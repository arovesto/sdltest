package play

import (
	"github.com/arovesto/sdl/pkg/camera"
	"github.com/arovesto/sdl/pkg/game/global"
	"github.com/arovesto/sdl/pkg/input"
	"github.com/arovesto/sdl/pkg/level"
	"github.com/arovesto/sdl/pkg/object"
	"github.com/arovesto/sdl/pkg/state"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	stateID = 1
)

type play struct {
	objects []object.GameObject
	level   *level.Level
}

func init() {
	state.Play = &play{}
}

func (p *play) Update() (err error) {
	if input.IsKeyDown(sdl.SCANCODE_ESCAPE) {
		if err = global.GetMachine().PushState(state.Pause); err != nil {
			return
		}
	}

	if err = camera.Update(); err != nil {
		return
	}

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
	if err = camera.SwitchCam(1); err != nil {
		return err
	}
	p.level, err = level.Parse(global.MapPath)
	return
}

func (p *play) OnSwitch() error {
	return camera.SwitchCam(0)
}

func (p *play) OnContinue() error {
	return camera.SwitchCam(1)
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
