package level

import (
	"github.com/arovesto/sdl/pkg/camera"
	"github.com/arovesto/sdl/pkg/math"
	"github.com/arovesto/sdl/pkg/object"
)

const (
	margin = 5
)

type Layer interface {
	Render() error
	Update() error
}

type Tiles [][]int

type tileLayer struct {
	tileSize int32

	sets  []*TileSet
	tiles Tiles

	collision bool
}

func NewTileLayer(tileSize int32, set []*TileSet, tiles Tiles, collision bool) *tileLayer {
	return &tileLayer{sets: set, tiles: tiles, tileSize: tileSize, collision: collision}
}

func (l *tileLayer) Render() (err error) {
	xCam, yCam, w, h := camera.Camera.GetRect().Values()
	cols, rows := w/l.tileSize+1, h/l.tileSize+1
	x, y := xCam/l.tileSize, yCam/l.tileSize
	x2, y2 := xCam%l.tileSize, yCam%l.tileSize
	for i := int32(0); i < rows; i++ {
		for j := int32(0); j < cols; j++ {
			if l.outside(j+x, i+y) {
				continue
			}
			id := l.tiles[i+y][j+x]
			ts := l.getTileSetByID(id)
			if err = ts.Draw(math.NewIntVector((j*l.tileSize)-x2, (i*l.tileSize)-y2),
				math.NewIntVector(int32(id-ts.FirstGID)%ts.Cols, int32(id-ts.FirstGID)/ts.Cols)); err != nil {
				return
			}
		}
	}
	return
}

func (l *tileLayer) Update() error {
	return nil
}

func (l *tileLayer) outside(b, a int32) bool {
	return a >= int32(len(l.tiles)) || b >= int32(len(l.tiles[0])) || a < 0 || b < 0
}

func (l *tileLayer) BackOffVector(o object.GameObject) (isGroundedPositive, isGroundedNegative, delta math.Vector2D) {
	xP, yP, wP, hP := o.GetCollider().Values()

	// TODO refactor me
	x, y := xP/l.tileSize, yP/l.tileSize
	w, h := math.DivRoundUp(wP, l.tileSize)+2, math.DivRoundUp(hP, l.tileSize)+2
	for i := x - 2; i <= w+x; i++ {
		for j := y - 2; j <= h+y; j++ {
			if !l.outside(i, j) && l.tiles[j][i] != 0 {
				if math.Abs(j*l.tileSize-yP-hP) <= margin && xP < (i+1)*l.tileSize && xP+wP > i*l.tileSize {
					isGroundedPositive.Y = 1
					delta.Y = float64(j*l.tileSize - yP - hP)
				}
				if math.Abs(i*l.tileSize-xP-wP) <= margin && yP < (j+1)*l.tileSize && yP+hP > j*l.tileSize {
					isGroundedPositive.X = 1
					delta.X = float64(i*l.tileSize - xP - wP)
				}
				if math.Abs(xP-(i+1)*l.tileSize) <= margin && yP < (j+1)*l.tileSize && yP+hP > j*l.tileSize {
					isGroundedNegative.X = -1
					delta.X = float64((i+1)*l.tileSize - xP)
				}
				if math.Abs(yP-(j+1)*l.tileSize) <= margin && xP < (i+1)*l.tileSize && xP+wP > i*l.tileSize {
					isGroundedNegative.Y = -1
					delta.Y = float64((j+1)*l.tileSize - yP)
				}
			}
		}
	}
	return
}

func (l *tileLayer) getTileSetByID(id int) *TileSet {
	for i := 0; i < len(l.sets)-1; i++ {
		if id >= l.sets[i].FirstGID && id < l.sets[i+1].FirstGID {
			return l.sets[i]
		}
	}
	return l.sets[len(l.sets)-1]
}

type objectLayer struct {
	objects []object.GameObject

	collision bool
}

func NewObjectLayer(obj []object.GameObject, collision bool) *objectLayer {
	return &objectLayer{objects: obj, collision: collision}
}

func (o *objectLayer) Render() (err error) {
	for _, o := range o.objects {
		if err = o.Draw(); err != nil {
			return
		}
	}
	return
}

func (o *objectLayer) Update() (err error) {
	for _, o := range o.objects {
		if err = o.Update(); err != nil {
			return
		}
	}
	return
}

func (o *objectLayer) Collision(tileLayers []*tileLayer) (err error) {
	for _, o := range o.objects {
		for _, tl := range tileLayers {
			if !tl.collision {
				continue
			}
			o.BackOff(tl.BackOffVector(o))
		}
	}
	for _, o1 := range o.objects {
		for _, o2 := range o.objects {
			if o1 != o2 && math.Collide(o1.GetObjectCollider(), o2.GetObjectCollider()) {
				if err = o1.Collide(o2); err != nil {
					return
				}
				if err = o2.Collide(o1); err != nil {
					return
				}
			}
		}
	}
	return
}
