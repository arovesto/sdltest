package main

import (
	"log"

	"github.com/arovesto/sdl/pkg/game"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	// set this high to remove fixed cap
	FPS = 60

	PerFrameTime = uint32(1000 / FPS)

	DEBUG = true
)

// Now on "scrolling background" of 8 chapter
func main() {
	g, err := game.InitGame(game.Opts{
		Height:    1080, // 34
		Width:     1920, // 60
		CamSpeedX: 50,
		CamSpeedY: 50,
		Title:     "sdlgame",
	})
	if err != nil {
		panic(err)
	}
	defer g.Destroy()
	for g.Running() {
		frameStart := sdl.GetTicks()
		if err := g.HandleEvents(); err != nil {
			panic(err)
		}
		if err := g.Update(); err != nil {
			panic(err)
		}
		if err := g.Render(); err != nil {
			panic(err)
		}
		delayRequired := PerFrameTime - (sdl.GetTicks() - frameStart)
		if int32(delayRequired) > 0 {
			sdl.Delay(delayRequired)
		} else {
			if DEBUG {
				log.Println("can't keep up")
			}
		}
	}
}
