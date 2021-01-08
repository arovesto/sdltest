package input

import (
	"github.com/arovesto/sdl/pkg/camera"
	"github.com/arovesto/sdl/pkg/game/global"
	"github.com/arovesto/sdl/pkg/math"
	"github.com/veandco/go-sdl2/sdl"
)

var handler = &Handler{}

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

func (h *Handler) Update() error {
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

func (h *Handler) GetMousePressed(id int) bool {
	return h.mousePressed[id]
}

func (h *Handler) GetMousePosition() math.Vector2D {
	return math.Sub(h.mousePos, camera.GetCamPos())
}

func (h *Handler) IsKeyDown(scancode sdl.Scancode) bool {
	if h.keyboardState == nil {
		return false
	}
	return h.keyboardState[scancode] == 1
}

func (h *Handler) Destroy() error {
	return nil
}

func (h *Handler) ResetMouse() {
	h.mousePressed = [3]bool{}
}

func GetHandler() *Handler {
	return handler
}

func Update() error {
	return handler.Update()
}

func Destroy() error {
	return handler.Destroy()
}

func GetMousePressed(id int) bool {
	return handler.GetMousePressed(id)
}

func GetMousePosition() math.Vector2D {
	return handler.GetMousePosition()
}

func IsKeyDown(scancode sdl.Scancode) bool {
	return handler.IsKeyDown(scancode)
}

func ResetMouse() {
	handler.ResetMouse()
}
