package object

import (
	"github.com/arovesto/sdl/pkg/game/global"
	"github.com/arovesto/sdl/pkg/math"
	"github.com/arovesto/sdl/pkg/model"
)

type GameObject interface {
	Draw() error
	Update() error
	Destroy() error
	Collide(other GameObject) error
	GetPosition() math.Vector2D
	GetCollider() math.Rect
	GetObjectCollider() math.Rect
	GetType() Type
	BackOff(info BackOffInfo)
}

type BackOffInfo struct {
	LeftGrounded  bool
	RightGrounded bool
	UpGrounded    bool
	DownGrounded  bool
	Delta         math.Vector2D
}

type Type int

type Properties struct {
	Pos       math.Vector2D
	Model     *model.Model
	AnimSpeed uint32
	Callback  global.ID
	ID        global.ID
}

const (
	NOType Type = iota
	EnemyType
	PlayerType
	BulletType
	NoTypeMenuObject
)
