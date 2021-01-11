package play

import (
	"github.com/arovesto/sdl/pkg/camera"
	"github.com/arovesto/sdl/pkg/game/global"
	"github.com/arovesto/sdl/pkg/input"
	"github.com/arovesto/sdl/pkg/level"
	"github.com/arovesto/sdl/pkg/parser"
	"github.com/arovesto/sdl/pkg/state"
	"github.com/veandco/go-sdl2/sdl"
)

type play struct {
	level *level.Level

	canPause uint32
}

func init() {
	state.Play = &play{}
}

func (p *play) Update() (err error) {
	if input.IsKeyDown(sdl.SCANCODE_ESCAPE) && p.canPause >= 20 {
		p.canPause = 0
		if err = global.GetMachine().PushState(state.Pause); err != nil {
			return
		}
	}
	if p.canPause < 20 {
		p.canPause++
	}
	camera.Camera.Update()
	return p.level.Update()
}

func (p *play) Render() (err error) {
	return p.level.Render()
}

func (p *play) OnEnter() (err error) {
	camera.Camera.Reset()
	p.level, err = parser.ParseLevel(global.MapPath, global.ModelsPath)
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
	return p.level.Destroy()
}
