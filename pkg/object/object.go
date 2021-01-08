package object

import (
	"github.com/arovesto/sdl/pkg/game/global"
	"github.com/arovesto/sdl/pkg/math"
	"github.com/arovesto/sdl/pkg/texturemanager"
	"github.com/veandco/go-sdl2/sdl"
)

type Type int

const (
	NOType Type = iota
	Player
	Enemy
)

type GameObject interface {
	Draw() error
	Update() error
	Destroy() error
	Collide(other GameObject)
	GetPosition() math.Vector2D
	GetVelocity() math.Vector2D
	GetSize() math.Vector2D
	GetType() Type
	BackOff(isGroundedP, isGroundedN, delta math.Vector2D)
}

type Properties struct {
	Pos       math.Vector2D
	Size      math.Vector2D
	Rows      int32
	Cols      int32
	ID        string
	AnimSpeed uint32
	XMaxSpeed float64
	YMaxSpeed float64
	Callback  global.ID
	Flip      sdl.RendererFlip
	IgnoreCam bool
}

type shooterObject struct {
	pos  math.Vector2D
	vel  math.Vector2D
	acc  math.Vector2D
	size math.Vector2D
	grP  math.Vector2D
	grN  math.Vector2D

	ignoreCam bool

	xMaxSpeed float64
	yMaxSpeed float64

	alpha uint8
	angle float64

	updating bool
	dead     bool
	dying    bool

	flip sdl.RendererFlip

	frame  int32
	frames int32

	fireSpeed     int
	bulletCounter int

	speed int

	dyingTime    int
	dyingCounter int

	playedDeathSound bool

	spriteChanged uint32

	id string
}

func newShooterObj(st Properties) shooterObject {
	if st.XMaxSpeed > 5 || st.XMaxSpeed == 0 {
		st.XMaxSpeed = 5
	}
	if st.YMaxSpeed > 5 || st.YMaxSpeed == 0 {
		st.YMaxSpeed = 5
	}
	return shooterObject{
		pos:       st.Pos,
		size:      st.Size,
		frames:    st.Cols,
		id:        st.ID,
		alpha:     255,
		flip:      st.Flip,
		xMaxSpeed: st.XMaxSpeed,
		yMaxSpeed: st.YMaxSpeed,
		ignoreCam: st.IgnoreCam,
	}
}

func (s *shooterObject) Update() error {
	s.acc.Y = 0.1

	s.acc = math.CampDirection(s.acc, s.grP, s.grN)
	s.vel = math.CampDirection(s.vel, s.grP, s.grN)

	s.vel = math.Add(s.vel, s.acc)
	if s.xMaxSpeed != 0 && math.AbsF(s.vel.X) > math.AbsF(s.xMaxSpeed) {
		s.vel.X = s.xMaxSpeed * math.SignF(s.vel.X)
	}
	if s.yMaxSpeed != 0 && math.AbsF(s.vel.Y) > math.AbsF(s.yMaxSpeed) {
		s.vel.Y = s.yMaxSpeed * math.SignF(s.vel.Y)
	}
	s.pos = math.Add(s.pos, s.vel)
	return nil
}

func (s *shooterObject) Draw() error {
	return texturemanager.Draw(texturemanager.DrawOpts{
		ID:        s.id,
		X:         int32(s.pos.X),
		Y:         int32(s.pos.Y),
		W:         int32(s.size.X),
		H:         int32(s.size.Y),
		Col:       s.frame,
		Flip:      s.flip,
		Alpha:     s.alpha,
		Angle:     s.angle,
		IgnoreCam: s.ignoreCam,
	})
}

func (s *shooterObject) changeSprite() {
	now := sdl.GetTicks()
	if now-s.spriteChanged > animationTime {
		s.spriteChanged = now
		s.frame += 1
		if s.frame >= s.frames {
			s.frame = 0
		}
	}
}

func (s *shooterObject) dyingAnim() error {
	if s.dyingCounter == s.dyingTime {
		s.dead = true
	}
	s.dyingCounter++
	return nil
}

func (s *shooterObject) Destroy() error {
	texturemanager.Delete(s.id)
	return nil
}

func (s *shooterObject) GetPosition() math.Vector2D {
	return s.pos
}

func (s *shooterObject) GetSize() math.Vector2D {
	return s.size
}

func (s *shooterObject) Collide(other GameObject) {
	s.dying = true
}

func (s *shooterObject) Updating() bool {
	return s.updating
}

func (s *shooterObject) Dead() bool {
	return s.dead
}

func (s *shooterObject) Dying() bool {
	return s.dying
}

func (s *shooterObject) StartUpdate(b bool) {
	s.updating = true
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

func (s *shooterObject) GetVelocity() math.Vector2D {
	return s.vel
}
