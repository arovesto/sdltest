package object

import (
	"github.com/arovesto/sdl/pkg/camera"
	"github.com/arovesto/sdl/pkg/game/global"
	"github.com/arovesto/sdl/pkg/input"
	"github.com/arovesto/sdl/pkg/math"
	"github.com/arovesto/sdl/pkg/sound"
	"github.com/arovesto/sdl/pkg/state"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	animationTime = 200
)

type player struct {
	shooterObject

	shootAt uint32
}

func NewPlayer(st Properties) GameObject {
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
	if input.IsKeyDown(sdl.SCANCODE_E) {
		global.Quit()
	}
	now := sdl.GetTicks()
	if input.IsKeyDown(sdl.SCANCODE_SPACE) && now-p.shootAt > 500 {
		if err := sound.PlaySound("shot", 0); err != nil {
			return err
		}
		p.shootAt = now
	}
	if input.IsKeyDown(sdl.SCANCODE_W) {
		p.vel.Y = -4
	}

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

func (p *player) GetType() Type {
	return PlayerType
}

func (p *player) Collide(other GameObject) error {
	return global.GetMachine().ChangeState(state.GameOver)
}
