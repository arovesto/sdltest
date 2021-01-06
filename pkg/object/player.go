package object

import (
	"github.com/arovesto/sdl/pkg/game/global"
	"github.com/arovesto/sdl/pkg/input"
	"github.com/arovesto/sdl/pkg/math"
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
		return global.GetMachine().ChangeState(state.GameOver)
	}

	return nil
}

func (p *player) GetType() Type {
	return Player
}

func (p *player) Collide() error {
	if p.invulnerable || global.LevelComplete() {
		return nil
	}
	// TODO store this metadata in texture manager
	p.id = "large-explosion"
	p.size = math.NewVec(128, 128)
	p.frame = 0
	p.frames = 9

	return p.shooterObject.Collide()
}
