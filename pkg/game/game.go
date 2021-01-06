package game

import (
	"github.com/arovesto/sdl/pkg/game/global"
	"github.com/arovesto/sdl/pkg/input"
	"github.com/arovesto/sdl/pkg/object"
	"github.com/arovesto/sdl/pkg/state"
	"github.com/arovesto/sdl/pkg/texturemanager"
	"github.com/veandco/go-sdl2/sdl"

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

	scrollSpeed float64
	playerLives int

	levelComplete bool
	level         int

	opts Opts
}

type Opts struct {
	Height     int32
	Width      int32
	Fullscreen bool
	Title      string
}

func InitGame(opts Opts) error {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return err
	}

	sdlOpts := sdl.WINDOW_SHOWN
	if opts.Fullscreen {
		sdlOpts |= sdl.WINDOW_FULLSCREEN
	}
	window, err := sdl.CreateWindow(opts.Title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, opts.Width, opts.Height, uint32(sdlOpts))
	if err != nil {
		return err
	}

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return err
	}

	if err = renderer.SetDrawColor(150, 0, 0, 255); err != nil {
		return err
	}
	texturemanager.InitManager(renderer)

	object.Register("button", object.NewButton)
	object.Register("player", object.NewPlayer)
	object.Register("enemy", object.NewEnemy)
	object.Register("animation", object.NewAnimation)
	object.Register("background", object.NewBackground)

	machine := state.NewMachine()
	if err := machine.PushState(state.Menu); err != nil {
		return err
	}

	global.SetGame(&Game{renderer: renderer, window: window, running: true, states: machine, opts: opts, scrollSpeed: 5})

	return nil
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
	_ = texturemanager.Destroy()
	_ = input.Destroy()
}

func (g *Game) Quit() {
	g.running = false
}

func (g *Game) GetMachine() *state.Machine {
	return g.states
}

func (g *Game) GetSize() (int32, int32) {
	return g.opts.Width, g.opts.Height
}

func (g *Game) GetScrollSpeed() float64 {
	return g.scrollSpeed
}

func (g *Game) DecreasePlayerLives() {
	g.playerLives--
}

func (g *Game) GetPlayerLives() int {
	return g.playerLives
}

func (g *Game) IncreaseLevel() error {
	g.level++
	g.levelComplete = false
	return g.states.ChangeState(state.Between)
}

func (g *Game) LevelComplete() bool {
	return g.levelComplete
}
