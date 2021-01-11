package game

import (
	"github.com/arovesto/sdl/pkg/camera"
	"github.com/arovesto/sdl/pkg/game/global"
	"github.com/arovesto/sdl/pkg/input"
	"github.com/arovesto/sdl/pkg/math"
	"github.com/arovesto/sdl/pkg/sound"
	"github.com/arovesto/sdl/pkg/state"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"

	_ "github.com/arovesto/sdl/pkg/state/gameover"
	_ "github.com/arovesto/sdl/pkg/state/menu"
	_ "github.com/arovesto/sdl/pkg/state/pause"
	_ "github.com/arovesto/sdl/pkg/state/play"
)

type Game struct {
	renderer *sdl.Renderer
	window   *sdl.Window
	running  bool
	states   *state.Machine

	playerLives int

	levelComplete bool
	level         int

	fullscreen bool
	w          int32
	h          int32

	canFullScreen int32
}

type Opts struct {
	Height     int32
	Width      int32
	Fullscreen bool
	Title      string
	CamSpeedX  float64
	CamSpeedY  float64
}

func InitGame(opts Opts) (*Game, error) {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return nil, err
	}
	if err := ttf.Init(); err != nil {
		return nil, err
	}

	sdlOpts := sdl.WINDOW_SHOWN
	if opts.Fullscreen {
		sdlOpts |= sdl.WINDOW_FULLSCREEN_DESKTOP
	}
	window, err := sdl.CreateWindow(opts.Title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, opts.Width, opts.Height, uint32(sdlOpts))
	if err != nil {
		return nil, err
	}

	camera.Camera = camera.NewCamera(opts.Width, opts.Height, math.NewVec(opts.CamSpeedX, opts.CamSpeedY))

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return nil, err
	}

	global.Renderer = renderer

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "2")

	if err = sound.Load("assets/heroes.flac", sound.MUSIC, "heroes"); err != nil {
		return nil, err
	}
	if err = sound.Load("assets/shot.wav", sound.SFX, "shot"); err != nil {
		return nil, err
	}

	w, h := window.GetSize()
	machine := state.NewMachine()
	g := &Game{renderer: renderer, window: window, running: true, states: machine, fullscreen: opts.Fullscreen, w: w, h: h}
	global.SetGame(g)

	return g, machine.PushState(state.Menu)
}

func (g *Game) Render() (err error) {
	if err = g.renderer.Clear(); err != nil {
		return
	}

	if err = g.states.Render(); err != nil {
		return
	}

	g.renderer.Present()
	return
}

func (g *Game) Update() (err error) {
	// TODO create some kind of storage for such events and run them all here
	if input.IsKeyDown(sdl.SCANCODE_F11) && g.canFullScreen >= 20 {
		g.canFullScreen = 0
		if err = g.ToggleFullscreen(); err != nil {
			return err
		}
	}
	if g.canFullScreen < 20 {
		g.canFullScreen++
	}
	return g.states.Update()
}

func (g *Game) Running() bool {
	return g.running
}

func (g *Game) HandleEvents() error {
	return input.Update()
}

func (g *Game) Destroy() {
	_ = g.window.Destroy()
	_ = g.renderer.Clear()
	_ = g.renderer.Destroy()
	_ = input.Destroy()
}

func (g *Game) Quit() {
	g.running = false
}

func (g *Game) GetMachine() *state.Machine {
	return g.states
}

func (g *Game) GetSize() (int32, int32) {
	return g.w, g.h
}

func (g *Game) ToggleFullscreen() error {
	g.fullscreen = !g.fullscreen

	var o uint32
	if g.fullscreen {
		o = sdl.WINDOW_FULLSCREEN_DESKTOP
	}
	err := g.window.SetFullscreen(o)
	g.w, g.h = g.window.GetSize()
	return err
}
