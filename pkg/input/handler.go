package input

import (
	"github.com/arovesto/sdl/pkg/camera"
	"github.com/arovesto/sdl/pkg/game/global"
	"github.com/arovesto/sdl/pkg/math"
	"github.com/veandco/go-sdl2/sdl"
)

var h = &Handler{}

const (
	LEFT = iota
	MIDDLE
	RIGHT
)

type Handler struct {
	mousePressed  [3]bool
	mousePos      math.Vector2D
	keyboardState []uint8
}

func Update() error {
	h.keyboardState = sdl.GetKeyboardState()
	for {
		e := sdl.PollEvent()
		if e == nil {
			return nil
		}
		switch event := e.(type) {
		case *sdl.QuitEvent:
			global.Quit()
		case *sdl.MouseButtonEvent:
			switch event.Button {
			case sdl.BUTTON_LEFT:
				h.mousePressed[LEFT] = event.Type == sdl.MOUSEBUTTONDOWN
			case sdl.BUTTON_MIDDLE:
				h.mousePressed[MIDDLE] = event.Type == sdl.MOUSEBUTTONDOWN
			case sdl.BUTTON_RIGHT:
				h.mousePressed[RIGHT] = event.Type == sdl.MOUSEBUTTONDOWN
			}
		case *sdl.MouseMotionEvent:
			h.mousePos.X = float64(event.X)
			h.mousePos.Y = float64(event.Y)
		}
	}
}

func GetMousePressed(id int) bool {
	return h.mousePressed[id]
}

func GetMousePosition() math.Vector2D {
	return h.mousePos
}

func GetMousePositionInCamera() math.Vector2D {
	return math.Sub(h.mousePos, camera.GetCamPos())
}

func IsKeyDown(scancode sdl.Scancode) bool {
	if h.keyboardState == nil {
		return false
	}
	return h.keyboardState[scancode] == 1
}

func UpdateKeyboard() {
	h.keyboardState = sdl.GetKeyboardState()
}

func Destroy() error {
	return nil
}

func ResetMouse() {
	h.mousePressed = [3]bool{}
}
