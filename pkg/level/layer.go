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

func (l *tileLayer) BackOffVector(o object.GameObject) (r object.BackOffInfo) {
	xP, yP, wP, hP := o.GetCollider().Values()

	// TODO refactor me
	x, y := xP/l.tileSize, yP/l.tileSize
	w, h := math.DivRoundUp(wP, l.tileSize)+2, math.DivRoundUp(hP, l.tileSize)+2
	for i := x - 2; i <= w+x; i++ {
		for j := y - 2; j <= h+y; j++ {
			if !l.outside(i, j) && l.tiles[j][i] != 0 {
				if !math.RectIntersect(math.Rect{X: xP, Y: yP, W: wP, H: hP}, math.Rect{X: i * l.tileSize, Y: j * l.tileSize, W: l.tileSize, H: l.tileSize}) {
					continue
				}
				tileX1, tileY1, tileX2, tileY2 := i*l.tileSize, j*l.tileSize, (i+1)*l.tileSize, (j+1)*l.tileSize
				actions := []object.BackOffInfo{
					{DownGrounded: true, Delta: math.NewIntVector(0, tileY1-(yP+hP)).FloatV()},
					{UpGrounded: true, Delta: math.NewIntVector(0, tileY2-yP).FloatV()},
					{LeftGrounded: true, Delta: math.NewIntVector(tileX2-xP, 0).FloatV()},
					{RightGrounded: true, Delta: math.NewIntVector(tileX1-(xP+wP), 0).FloatV()},
				}
				res := actions[0]
				for _, a := range actions {
					if a.Delta.Abs() < res.Delta.Abs() {
						res = a
					}
				}

				r.Delta = r.Delta.Add(res.Delta)
				xP += int32(res.Delta.X)
				yP += int32(res.Delta.Y)
				r.DownGrounded = res.DownGrounded || r.DownGrounded
				r.UpGrounded = res.UpGrounded || r.UpGrounded
				r.LeftGrounded = res.LeftGrounded || r.LeftGrounded
				r.RightGrounded = res.RightGrounded || r.RightGrounded
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

	deletedIDs map[object.GameObject]struct{}

	collision bool
}

func NewObjectLayer(obj []object.GameObject, collision bool) *objectLayer {
	return &objectLayer{objects: obj, collision: collision, deletedIDs: map[object.GameObject]struct{}{}}
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
	if len(o.deletedIDs) != 0 {
		newOBJ := make([]object.GameObject, 0, len(o.objects))
		for _, obj := range o.objects {
			if _, ok := o.deletedIDs[obj]; !ok {
				newOBJ = append(newOBJ, obj)
			}
		}
		o.objects = newOBJ
		o.deletedIDs = map[object.GameObject]struct{}{}
	}
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
