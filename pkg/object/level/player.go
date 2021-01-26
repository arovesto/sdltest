package level

import (
	"github.com/arovesto/sdl/pkg/camera"
	"github.com/arovesto/sdl/pkg/game/global"
	"github.com/arovesto/sdl/pkg/input"
	"github.com/arovesto/sdl/pkg/level"
	"github.com/arovesto/sdl/pkg/math"
	"github.com/arovesto/sdl/pkg/model"
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

	SpawnedAt uint32
}

func NewPlayer(st object.Properties) object.GameObject {
	return &player{shooterObject: newShooterObj(st)}
}

func (p *player) Update(tDelta float64) (err error) {
	if err = p.handleInput(tDelta); err != nil {
		return err
	}
	return p.shooterObject.Update(tDelta)
}

func (p *player) handleInput(tDelta float64) error {
	now := sdl.GetTicks()
	mousePos := input.GetMousePositionInCamera().FloatV()
	pivotPos := p.model.Parts[0].GetPivot(p.pos.IntVector()).FloatV()

	// TODO predicate of action
	if input.GetMousePressed(input.LEFT) && now-p.shootAt > 1000 && !p.rotating {
		// TODO setting action off
		p.shootAt = now
		// TODO action callback
		if err := sound.PlaySound("shot", 0); err != nil {
			return err
		}
		level.CurrentLevel.NewObj(NewBullet(pivotPos, p.model.Parts[0].GetAngle().ToVec().Mul(1500)))
	}

	if input.GetMousePressed(input.RIGHT) && now-p.SpawnedAt > 1000 {
		p.SpawnedAt = now
		level.CurrentLevel.NewObj(NewEnemy(object.Properties{
			Pos:   mousePos,
			Model: model.AvailableModels["evil_tank"].GetCopy(),
			ID:    global.NewID(),
		}))
	}

	if input.IsKeyDown(sdl.SCANCODE_E) {
		p.model.Angle -= 0.05
	}
	if input.IsKeyDown(sdl.SCANCODE_Q) {
		p.model.Angle += 0.05
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
		p.vel.Y = -400
	}

	oldFlip := p.model.Parts[0].Flip
	// TODO action callbackNotHappen
	if mousePos.X-pivotPos.X > 40 && !p.rotating {
		p.model.Parts[0].Flip = sdl.FLIP_NONE
	}
	if mousePos.X-pivotPos.X < -40 && !p.rotating {
		p.model.Parts[0].Flip = sdl.FLIP_HORIZONTAL
	}

	requiredAngle := math.AngleOn(pivotPos, mousePos)
	p.model.Parts[0].Angle += math.ClampAngle((requiredAngle-p.model.Parts[0].Angle)*math.AngleDeg(5*tDelta), gunSpeed)
	p.model.Parts[0].ClampAngle()

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

	var player math.Vector2D
	switch {
	case input.IsKeyDown(sdl.SCANCODE_D):
		player.X = 1000
		p.changeSprite(2)
		p.model.Parts[2].Flip = sdl.FLIP_NONE
		p.model.GlobalFlip = sdl.FLIP_NONE
	case input.IsKeyDown(sdl.SCANCODE_A):
		player.X = -1000
		p.changeSprite(2)
		p.model.Parts[2].Flip = sdl.FLIP_HORIZONTAL
		p.model.GlobalFlip = sdl.FLIP_HORIZONTAL
	}

	if p.model.Parts[0].Flip&sdl.FLIP_HORIZONTAL == 0 {
		camera.Camera.MainTarget = p.pos.Add(math.NewVec(600+p.vel.X/1000, -200))
	} else {
		camera.Camera.MainTarget = p.pos.Add(math.NewVec(-500+p.vel.X/1000, -200))
	}

	friction := p.vel.Mul(-500000 * tDelta)
	if player.X == 0 || (friction.X > 0) == (player.X > 0) {
		p.acc = friction.Add(player)
	} else {
		p.acc = player
	}

	return nil
}

func (p *player) GetType() object.Type {
	return object.PlayerType
}

func (p *player) Collide(other object.GameObject) error {
	if other.GetType() == object.EnemyType {
		return global.GetMachine().ChangeState(state.GameOver)
	}
	return nil
}
