package level

import (
	"github.com/arovesto/sdl/pkg/camera"
	"github.com/arovesto/sdl/pkg/game/global"
	"github.com/arovesto/sdl/pkg/input"
	"github.com/arovesto/sdl/pkg/math"
	"github.com/arovesto/sdl/pkg/object"
	"github.com/arovesto/sdl/pkg/sound"
	"github.com/arovesto/sdl/pkg/state"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	gunSpeed = 1
)

type player struct {
	shooterObject

	shootAt uint32
}

var MainPlayer *player

func NewPlayer(st object.Properties) object.GameObject {
	p := &player{shooterObject: newShooterObj(st)}
	MainPlayer = p
	return p
}

func (p *player) Update() (err error) {
	if err = p.handleInput(); err != nil {
		return err
	}
	return p.shooterObject.Update()
}

func (p *player) handleInput() error {
	now := sdl.GetTicks()
	if input.IsKeyDown(sdl.SCANCODE_SPACE) && now-p.shootAt > 1000 {
		if err := sound.PlaySound("shot", 0); err != nil {
			return err
		}
		p.shootAt = now
	}
	if input.IsKeyDown(sdl.SCANCODE_W) {
		p.vel.Y = -4
	}
	requiredAngle := math.AngleOn(p.model.Parts[0].GetPivotPoint(p.pos.IntVector(), &p.model).FloatV(), input.GetMousePositionInCamera().FloatV())

	p.model.Parts[0].Angle += math.ClampAngle(requiredAngle-p.model.Parts[0].Angle, gunSpeed)

	camera.Camera.MainTarget = p.pos.Add(math.NewVec(200, 1000*p.acc.Y))
	var player math.Vector2D
	switch {
	case input.IsKeyDown(sdl.SCANCODE_D):
		player.X = 0.1
		p.changeSprite()
		p.model.Flip = sdl.FLIP_NONE
	case input.IsKeyDown(sdl.SCANCODE_A):
		player.X = -0.1
		p.changeSprite()
		p.model.Flip = sdl.FLIP_HORIZONTAL
	default:
	}

	friction := p.vel.Mul(-0.1)
	if player.X == 0 || (friction.X > 0) == (player.X > 0) {
		p.acc = friction.Add(player)
	} else {
		p.acc = player
	}

	return nil
}

func (p *player) GetType() object.Type {
	return PlayerType
}

func (p *player) Collide(other object.GameObject) error {
	return global.GetMachine().ChangeState(state.GameOver)
}
