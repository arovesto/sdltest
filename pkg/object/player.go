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

	invulnerable        bool
	invulnerableTime    int
	invulnerableCounter int

	alpha int
	angle float64

	shootAt uint32
}

func NewPlayer(st Properties) GameObject {
	return &player{shooterObject: newShooterObj(st)}
}

func (p *player) Update() (err error) {
	p.changeSprite()
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
	if input.IsKeyDown(sdl.SCANCODE_SPACE) && now-p.shootAt > 1000 {
		if err := sound.PlaySound("shot", 0); err != nil {
			return err
		}
		p.shootAt = now
	}
	if input.IsKeyDown(sdl.SCANCODE_W) {
		p.vel.Y = -4
	}

	camera.GoTo(math.Add(p.pos, math.NewVec(200, -100)))
	var player math.Vector2D
	switch {
	case input.IsKeyDown(sdl.SCANCODE_D):
		player.X = 0.05
		p.flip = sdl.FLIP_NONE
	case input.IsKeyDown(sdl.SCANCODE_A):
		player.X = -0.05
		p.flip = sdl.FLIP_HORIZONTAL
	default:
	}

	friction := math.Mul(p.vel, -0.1)
	if player.X == 0 || (friction.X > 0) == (player.X > 0) {
		p.acc = math.Add(friction, player)
	} else {
		p.acc = player
	}

	return nil
}

func (p *player) GetType() Type {
	return Player
}

func (p *player) Collide(other GameObject) {
	_ = global.GetMachine().ChangeState(state.GameOver)
	if p.invulnerable || global.LevelComplete() {
		return
	}
	// TODO store this metadata in texture manager
	//p.id = "large-explosion"
	//p.size = math.NewVec(128, 128)
	//p.frame = 0
	//p.frames = 9

	p.shooterObject.Collide(other)
}
