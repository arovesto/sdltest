package camera

import (
	"github.com/arovesto/sdl/pkg/game/global"
	"github.com/arovesto/sdl/pkg/math"
)

type Target struct {
	Pos    math.Vector2D
	Weight float64
}

var Camera *camera

type camera struct {
	Pos        math.Vector2D
	Targets    map[global.ID]Target
	MainTarget math.Vector2D

	vel math.Vector2D

	center math.Vector2D

	approachSpeed math.Vector2D
}

func NewCamera(w, h int32, approachSpeed math.Vector2D) *camera {
	return &camera{center: math.NewIntVector(w, h).FloatV().Div(2), approachSpeed: approachSpeed, Targets: map[global.ID]Target{}}
}

func (c *camera) GetRect() math.Rect {
	w, h := global.GetSize()
	return math.Rect{X: int32(c.Pos.X), Y: int32(c.Pos.Y), W: w, H: h}
}

func (c *camera) Reset() {
	c.Targets = map[global.ID]Target{}
	c.Pos = math.ZeroVec()
}

func (c *camera) Update() {
	var target math.Vector2D
	var haveTarget bool
	var div float64
	for _, t := range c.Targets {
		if math.PointInRect(c.GetRect(), t.Pos) {
			target = target.Add(t.Pos).Mul(t.Weight)
			div += t.Weight
			haveTarget = true
		}
	}
	target = c.MainTarget.Add(target.Div(div).DivComponents(c.approachSpeed))

	if !haveTarget || !math.PointInRect(c.GetRect(), target) {
		target = c.MainTarget
	}

	// TODO if target not in cam rec - some margin do just move camera so it is in that margin

	c.vel = target.Sub(c.center).Sub(c.Pos).DivComponents(c.approachSpeed)
	c.Pos = c.Pos.Add(c.vel)
	c.Targets = map[global.ID]Target{}
}
