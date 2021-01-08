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
	Collide() error
	GetPosition() math.Vector2D
	GetSize() math.Vector2D
	GetType() Type
}

type Properties struct {
	Pos       math.Vector2D
	Size      math.Vector2D
	Rows      int32
	Cols      int32
	ID        string
	AnimSpeed uint32
	MaxSpeed  float64
	Callback  global.ID
	Flip      sdl.RendererFlip
}

type shooterObject struct {
	pos  math.Vector2D
	vel  math.Vector2D
	acc  math.Vector2D
	size math.Vector2D

	maxSpeed float64

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
	return shooterObject{
		pos:      st.Pos,
		size:     st.Size,
		frames:   st.Cols,
		id:       st.ID,
		alpha:    255,
		flip:     st.Flip,
		maxSpeed: st.MaxSpeed,
	}
}

func (s *shooterObject) Update() error {
	s.vel = math.Add(s.vel, s.acc)
	if s.maxSpeed != 0 {
		l := s.vel.Length()
		if l > s.maxSpeed {
			s.vel = math.Div(s.vel, l/s.maxSpeed)
		}
	}
	s.pos = math.Add(s.pos, s.vel)
	return nil
}

func (s *shooterObject) Draw() error {
	return texturemanager.Draw(texturemanager.DrawOpts{
		ID:    s.id,
		X:     int32(s.pos.X),
		Y:     int32(s.pos.Y),
		W:     int32(s.size.X),
		H:     int32(s.size.Y),
		Col:   s.frame,
		Flip:  s.flip,
		Alpha: s.alpha,
		Angle: s.angle,
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
	s.Scroll(global.GetScrollSpeed())
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

func (s *shooterObject) Scroll(speed float64) {
	s.pos.X -= speed
}

func (s *shooterObject) Collide() error {
	s.dying = true
	return nil
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
