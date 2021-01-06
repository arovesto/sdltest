package global

import (
	"path"

	"github.com/arovesto/sdl/pkg/state"
)

type ID int

const (
	MapPath    = "assets/map.tmx"
	AssetsPath = "assets/game.yaml"
	AssetsDir  = "assets"
)

func GetAssetsPath(name string) string {
	return path.Join(AssetsDir, name)
}

type Game interface {
	Quit()
	Render() error
	Update() error
	Running() bool
	HandleEvents() error
	Destroy()
	GetMachine() *state.Machine
	GetSize() (int32, int32)
	GetScrollSpeed() float64
	DecreasePlayerLives()
	GetPlayerLives() int
	IncreaseLevel() error
	LevelComplete() bool
}

var game Game

func SetGame(g Game) {
	game = g
}

func Render() error {
	return game.Render()
}

func Update() error {
	return game.Update()
}

func Running() bool {
	return game.Running()
}

func HandleEvents() error {
	return game.HandleEvents()
}

func Destroy() {
	game.Destroy()
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

func GetScrollSpeed() float64 {
	return game.GetScrollSpeed()
}

func DecreasePlayerLives() {
	game.DecreasePlayerLives()
}

func GetPlayerLives() int {
	return game.GetPlayerLives()
}

func IncreaseLevel() error {
	return game.IncreaseLevel()
}

func LevelComplete() bool {
	return game.LevelComplete()
}
