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

	FULLSCREEN = !DEBUG
)

// Now on "scrolling background" of 8 chapter
func main() {
	g, err := game.InitGame(game.Opts{
		Height:     1080, // 34
		Width:      1920, // 60
		CamSpeedX:  50,
		CamSpeedY:  20,
		Title:      "sdlgame",
		Fullscreen: FULLSCREEN,
	})
	if err != nil {
		panic(err)
	}
	var avg, scans uint32
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
			avg += delayRequired
			scans++
		} else {
			if DEBUG {
				log.Println("can't keep up")
			}
		}
		if scans == 60 && DEBUG {
			log.Printf("delay is %.2f", float64(avg)/float64(scans))
			avg = 0
			scans = 0
		}
	}
}
