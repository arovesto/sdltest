package object

import (
	"github.com/arovesto/sdl/pkg/model"

	"github.com/arovesto/sdl/pkg/game/global"
	"github.com/arovesto/sdl/pkg/math"
	"github.com/veandco/go-sdl2/sdl"
)

type Type int

const (
	NOType Type = iota
	EnemyType
	PlayerType

	maxSpeed = 10
)

var MainPlayer *player

type GameObject interface {
	Draw() error
	Update() error
	Destroy() error
	Collide(other GameObject) error
	GetPosition() math.Vector2D
	GetCollider() math.Rect
	GetObjectCollider() math.Rect
	GetType() Type
	BackOff(isGroundedP, isGroundedN, delta math.Vector2D)
}

type Properties struct {
	Pos       math.Vector2D
	Model     model.Model
	AnimSpeed uint32
	Callback  global.ID
}

type shooterObject struct {
	pos math.Vector2D
	vel math.Vector2D
	acc math.Vector2D

	grP math.Vector2D
	grN math.Vector2D

	updating bool
	dead     bool
	dying    bool

	model         model.Model
	spriteChanged uint32

	id global.ID
}

func newShooterObj(st Properties) shooterObject {
	return shooterObject{
		pos:   st.Pos,
		model: st.Model,
		id:    global.NewID(),
	}
}

func (s *shooterObject) Update() error {
	s.acc.Y = 0.1
	s.acc = math.ClampDirection(s.acc, s.grP, s.grN)
	s.vel = math.ClampDirection(s.vel, s.grP, s.grN)

	s.vel = s.vel.Add(s.acc)
	if math.AbsF(s.vel.X) > math.AbsF(maxSpeed) {
		s.vel.X = maxSpeed * math.SignF(s.vel.X)
	}
	if math.AbsF(s.vel.Y) > math.AbsF(maxSpeed) {
		s.vel.Y = maxSpeed * math.SignF(s.vel.Y)
	}
	s.pos = s.pos.Add(s.vel)
	return nil
}

func (s *shooterObject) Draw() error {
	return s.model.Draw(s.pos.IntVector())
}

func (s *shooterObject) changeSprite() {
	now := sdl.GetTicks()
	if now-s.spriteChanged > animationTime {
		s.spriteChanged = now
		s.model.ChangeSprites(-1)
	}
}

func (s *shooterObject) Destroy() error {
	return s.model.Destroy()
}

func (s *shooterObject) GetPosition() math.Vector2D {
	return s.pos
}

func (s *shooterObject) GetCollider() math.Rect {
	return s.model.Collider.Add(s.pos.IntVector())
}

func (s *shooterObject) GetObjectCollider() math.Rect {
	return s.model.EntityCollider.Add(s.pos.IntVector())
}

func (s *shooterObject) GetType() Type {
	return NOType
}

func (s *shooterObject) BackOff(isGroundedP, isGroundedN, delta math.Vector2D) {
	s.grP = isGroundedP
	s.grN = isGroundedN
	if delta.X != 0 && (delta.X > 0) == (s.vel.X > 0) {
		s.pos.X += delta.X
	}
	if delta.Y != 0 && (delta.Y > 0) == (s.vel.Y > 0) {
		s.pos.Y += delta.Y
	}
}

func (s *shooterObject) Collide(o GameObject) error {
	return nil
}
