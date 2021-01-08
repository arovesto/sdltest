package collision

import (
	"github.com/arovesto/sdl/pkg/level"
	"github.com/arovesto/sdl/pkg/math"
	"github.com/arovesto/sdl/pkg/object"
)

// TODO made into it's own layer
var m = &manager{}

var margin = math.NewVec(3, 3)

type manager struct {
	tileLayers []*level.TileLayer
	objects    []object.GameObject
}

func RegisterTileLayer(l *level.TileLayer) {
	m.tileLayers = append(m.tileLayers, l)
}

func RegisterObject(o object.GameObject) {
	m.objects = append(m.objects, o)
}

func RunCollision() {
	for _, o := range m.objects {
		for _, tl := range m.tileLayers {
			o.BackOff(tl.BackOffVector(o))
		}
	}

	for _, o1 := range m.objects {
		for _, o2 := range m.objects {
			if o1 != o2 && math.CollideMargin(o1.GetPosition(), o1.GetSize(), o2.GetPosition(), o2.GetSize(), margin) {
				o1.Collide(o2)
				o2.Collide(o1)
			}
		}
	}
}

func Clear() {
	m.objects = nil
	m.tileLayers = nil
}

// TODO func CanGoThere(from, size, to math.Vector2D) bool
