package camera

import (
	"github.com/arovesto/sdl/pkg/math"
)

type camera struct {
	basePos math.Vector2D
	pos     math.Vector2D
	vel     math.Vector2D

	target math.Vector2D

	size math.Vector2D

	approachSpeed float64
}

func (c *camera) Update() error {
	c.vel = math.Div(math.Sub(c.target, c.pos), c.approachSpeed)
	c.pos = math.Add(c.pos, c.vel)
	return nil
}

func (c *camera) Center() math.Vector2D {
	return math.Div(math.Add(c.basePos, c.size), 2)
}
