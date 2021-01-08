package level

import (
	"github.com/arovesto/sdl/pkg/camera"
	"github.com/arovesto/sdl/pkg/game/global"
	"github.com/arovesto/sdl/pkg/math"
	"github.com/arovesto/sdl/pkg/object"
	"github.com/arovesto/sdl/pkg/texturemanager"
)

const (
	margin = 3
)

type Layer interface {
	Render() error
	Update() error
}

type Tiles [][]int

type TileLayer struct {
	cols     int32
	rows     int32
	tileSize int32

	sets  []TileSet
	tiles Tiles

	collision bool
}

func NewTileLayer(tileSize int32, set []TileSet, tiles Tiles, collision bool) *TileLayer {
	w, h := global.GetSize()
	return &TileLayer{sets: set, tiles: tiles, tileSize: tileSize, cols: w/tileSize + 1, rows: h/tileSize + 1, collision: collision}
}

func (l *TileLayer) Render() (err error) {
	xCam, yCam := camera.GetCamPos().IntPos()
	x, y := xCam/l.tileSize, yCam/l.tileSize
	x2, y2 := xCam%l.tileSize, yCam%l.tileSize
	for i := int32(0); i < l.rows; i++ {
		for j := int32(0); j < l.cols; j++ {
			if l.outside(j+x, i+y) {
				continue
			}
			id := l.tiles[i+y][j+x]
			ts := l.getTileSetByID(id)
			if err = texturemanager.DrawTile(texturemanager.DrawTileOpts{
				ID:      ts.Name,
				Spacing: ts.Spacing,
				Margin:  ts.Margin,
				X:       (j * l.tileSize) - x2,
				Y:       (i * l.tileSize) - y2,
				W:       l.tileSize,
				H:       l.tileSize,
				Row:     int32(id-ts.FirstGID) / ts.Cols,
				Col:     int32(id-ts.FirstGID) % ts.Cols,
			}); err != nil {
				return
			}
		}
	}
	return
}

func (l *TileLayer) Update() error {
	return nil
}

func (l *TileLayer) outside(b, a int32) bool {
	return a >= int32(len(l.tiles)) || b >= int32(len(l.tiles[0])) || a < 0 || b < 0
}

func (l *TileLayer) BackOffVector(o object.GameObject) (isGroundedPositive, isGroundedNegative, delta math.Vector2D) {
	// TODO refactor me
	xP, yP := o.GetPosition().IntPos()
	x, y := xP/l.tileSize, yP/l.tileSize
	wP, hP := o.GetSize().IntPos()
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

func (l *TileLayer) getTileSetByID(id int) TileSet {
	for i := 0; i < len(l.sets)-1; i++ {
		if id >= l.sets[i].FirstGID && id < l.sets[i+1].FirstGID {
			return l.sets[i]
		}
	}
	return l.sets[len(l.sets)-1]
}

type objectLayer struct {
	objects []object.GameObject
}

func NewObjectLayer(obj []object.GameObject) *objectLayer {
	return &objectLayer{objects: obj}
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

/*
	for i := x - 2; i <= w+x; i++ {
		for j := y - 2; j <= h+y; j++ {
			if i == 5 && j == 24 {
				log.Println(xP, xP+wP, yP, yP+hP, x, x+w, y, y+h)
			}
			if !l.outside(i, j) && l.tiles[j][i] != 0 {
				if math.Abs(j*l.tileSize-yP-hP) <= margin && xP <= (i+1)*l.tileSize && xP+wP >= i*l.tileSize {
					isGrounded.Y = 1
					log.Println(i, j)
					delta.Y = float64(j*l.tileSize - yP - hP)
				}
				if math.Abs(i*l.tileSize-xP-wP) <= margin && yP <= (j+1)*l.tileSize && yP+hP >= j*l.tileSize {
					isGrounded.X = 1
					log.Println(i, j)
					delta.X = float64(i*l.tileSize - xP - wP)
				}
				if math.Abs(xP-(i+1)*l.tileSize) <= margin && yP <= (j+1)*l.tileSize && yP+hP >= j*l.tileSize {
					isGrounded.X = -1
					log.Println(i, j)
					delta.X = float64((i+1)*l.tileSize - xP)
				}
				if math.Abs(yP-(j+1)*l.tileSize) <= margin && xP <= (i+1)*l.tileSize && xP+wP >= i*l.tileSize {
					isGrounded.Y = -1
					log.Println(i, j)
					delta.Y = float64((j+1)*l.tileSize - yP)
				}
			}
		}
	}
*/
