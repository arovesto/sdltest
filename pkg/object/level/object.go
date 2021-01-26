package level

import (
	"github.com/arovesto/sdl/pkg/game/global"
	"github.com/arovesto/sdl/pkg/math"
	"github.com/arovesto/sdl/pkg/model"
	"github.com/arovesto/sdl/pkg/object"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	maxSpeed      = 10000
	animationTime = 200
)

type shooterObject struct {
	pos math.Vector2D
	vel math.Vector2D
	acc math.Vector2D

	gr object.BackOffInfo

	updating bool
	dead     bool
	dying    bool

	model         *model.Model
	spriteChanged uint32

	maxSpeed float64

	id global.ID
}

func newShooterObj(st object.Properties) shooterObject {
	return shooterObject{
		pos:      st.Pos,
		model:    st.Model,
		id:       st.ID,
		maxSpeed: maxSpeed,
	}
}

func (s *shooterObject) Update(delta float64) error {
	s.acc.Y = 1000
	s.acc = math.ClampDirection(s.acc, s.gr.UpGrounded, s.gr.DownGrounded, s.gr.LeftGrounded, s.gr.RightGrounded)
	s.vel = math.ClampDirection(s.vel, s.gr.UpGrounded, s.gr.DownGrounded, s.gr.LeftGrounded, s.gr.RightGrounded)

	s.vel = s.vel.Add(s.acc.Mul(delta))
	if math.AbsF(s.vel.X) > math.AbsF(s.maxSpeed) {
		s.vel.X = s.maxSpeed * math.SignF(s.vel.X)
	}
	if math.AbsF(s.vel.Y) > math.AbsF(s.maxSpeed) {
		s.vel.Y = s.maxSpeed * math.SignF(s.vel.Y)
	}
	s.pos = s.pos.Add(s.vel.Mul(delta))
	return nil
}

func (s *shooterObject) Draw() error {
	return s.model.Draw(s.pos.IntVector())
}

func (s *shooterObject) changeSprite(part int) {
	now := sdl.GetTicks()
	if now-s.spriteChanged > animationTime {
		s.spriteChanged = now
		if part < 0 {
			s.model.ChangeSprites(-1)
		} else {
			s.model.Parts[part].ChangeSprites(-1)
		}

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

func (s *shooterObject) GetType() object.Type {
	return object.NOType
}

func (s *shooterObject) BackOff(info object.BackOffInfo) {
	s.gr = info
	s.pos = s.pos.Add(s.gr.Delta)
}

func (s *shooterObject) Collide(o object.GameObject) error {
	return nil
}
