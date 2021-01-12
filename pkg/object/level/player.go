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
	gunSpeed = 0.5
)

type player struct {
	shooterObject

	shootAt    uint32
	rotatingAt uint32
	rotating   bool
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
	// TODO predicate of action
	if input.GetMousePressed(input.LEFT) && now-p.shootAt > 1000 && !p.rotating {
		// TODO setting action off
		p.shootAt = now
		// TODO action callback
		if err := sound.PlaySound("shot", 0); err != nil {
			return err
		}
	}
	// TODO predicate of action
	if now-p.rotatingAt > 400 && p.rotating {
		// TODO setting action off
		p.rotating = false
		// TODO action callback
		p.model.Parts[0].Hidden = false
		p.model.Parts[1].Frame = 0
	}
	if input.IsKeyDown(sdl.SCANCODE_W) {
		p.vel.Y = -4
	}
	mousePos := input.GetMousePositionInCamera().FloatV()
	pivotPos := p.model.Parts[0].GetPivotPoint(p.pos.IntVector()).FloatV()
	requiredAngle := math.AngleOn(pivotPos, mousePos)

	p.model.Parts[0].Angle += math.ClampAngle(requiredAngle-p.model.Parts[0].Angle, gunSpeed)
	oldFlip := p.model.Parts[0].Flip
	// TODO action callbackNotHappen
	if mousePos.X-pivotPos.X > 40 && !p.rotating {
		p.model.Parts[0].Flip = sdl.FLIP_NONE
	}
	if mousePos.X-pivotPos.X < -40 && !p.rotating {
		p.model.Parts[0].Flip = sdl.FLIP_HORIZONTAL
	}

	if oldFlip != p.model.Parts[0].Flip {
		// TODO setting action off
		p.rotating = true
		p.rotatingAt = now
	}
	// TODO action CallbackNotHappen
	if !p.rotating {
		p.model.Parts[1].Flip = p.model.Parts[0].Flip
	} else {
		// TODO action CallbackHappen
		p.model.Parts[0].Hidden = true
		p.model.Parts[0].Angle = 0
		p.model.Parts[1].Frame = 1
		p.model.Parts[1].Flip = sdl.FLIP_NONE
	}

	camera.Camera.MainTarget = p.pos.Add(math.NewVec(200, 1000*p.acc.Y))
	var player math.Vector2D
	switch {
	case input.IsKeyDown(sdl.SCANCODE_D):
		player.X = 0.1
		p.changeSprite(2)
		p.model.Parts[2].Flip = sdl.FLIP_NONE
		p.model.GlobalFlip = sdl.FLIP_NONE
	case input.IsKeyDown(sdl.SCANCODE_A):
		player.X = -0.1
		p.changeSprite(2)
		p.model.Parts[2].Flip = sdl.FLIP_HORIZONTAL
		p.model.GlobalFlip = sdl.FLIP_HORIZONTAL
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
