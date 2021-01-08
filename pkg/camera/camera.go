package camera

import (
	"github.com/arovesto/sdl/pkg/math"
	"golang.org/x/xerrors"
)

type CamType int

const (
	STATIONARY CamType = iota
	MOVING
)

type camera struct {
	basePos math.Vector2D
	pos     math.Vector2D
	vel     math.Vector2D

	target math.Vector2D

	size math.Vector2D

	approachSpeed float64

	tp CamType
}

func (c *camera) Update() error {
	switch c.tp {
	case MOVING:
		c.vel = math.Div(math.Sub(c.target, c.pos), c.approachSpeed)
		c.pos = math.Add(c.pos, c.vel)
		return nil
	case STATIONARY:
		return nil
	default:
		return xerrors.New("unknown camera type")
	}
}

func (c *camera) Center() math.Vector2D {
	return math.Div(math.Add(c.basePos, c.size), 2)
}
