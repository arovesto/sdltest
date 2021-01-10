package global

import (
	"path"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/arovesto/sdl/pkg/state"
)

type ID int

var curID ID

func NewID() ID {
	r := curID
	curID++
	return r
}

const (
	MapPath    = "assets/map.tmx"
	ModelsPath = "assets/models.yaml"
	MenusPath  = "assets/menus.yaml"
	AssetsDir  = "assets"
)

var Renderer *sdl.Renderer

func GetAssetsPath(name string) string {
	return path.Join(AssetsDir, name)
}

type Game interface {
	Quit()
	GetMachine() *state.Machine
	GetSize() (int32, int32)
	ToggleFullscreen() error
}

var game Game

func SetGame(g Game) {
	game = g
}

func Quit() {
	game.Quit()
}

func GetMachine() *state.Machine {
	return game.GetMachine()
}

func GetSize() (int32, int32) {
	return game.GetSize()
}

func ToggleFullscreen() error {
	return game.ToggleFullscreen()
}
