package level

import (
	"github.com/arovesto/sdl/pkg/game/global"
	"github.com/arovesto/sdl/pkg/math"
	"github.com/arovesto/sdl/pkg/object"
	"github.com/arovesto/sdl/pkg/texturemanager"
)

type Layer interface {
	Render() error
	Update() error
}

type Tiles [][]int

type tileLayer struct {
	cols     int32
	rows     int32
	tileSize int32

	pos math.Vector2D
	vel math.Vector2D

	sets  []TileSet
	tiles Tiles
}

func NewTileLayer(tileSize int32, set []TileSet, tiles Tiles) *tileLayer {
	w, h := global.GetSize()
	return &tileLayer{sets: set, tiles: tiles, tileSize: tileSize, cols: w/tileSize + 1, rows: h/tileSize + 1}
}

func (l *tileLayer) Render() (err error) {
	x, y := int32(l.pos.X)/l.tileSize, int32(l.pos.Y)/l.tileSize
	x2, y2 := int32(l.pos.X)%l.tileSize, int32(l.pos.Y)%l.tileSize
	for i := int32(0); i < l.rows; i++ {
		for j := int32(0); j < l.cols; j++ {
			if i+y >= int32(len(l.tiles)) || j+x >= int32(len(l.tiles[0])) || i+y < 0 || j+x < 0 {
				continue
			}
			id := l.tiles[i+y][j+x]
			ts := l.getTileSetByID(id)
			if err = texturemanager.DrawTile(texturemanager.DrawTileOpts{
				ID:      ts.name,
				Spacing: ts.spacing,
				Margin:  ts.margin,
				X:       (j * l.tileSize) - x2,
				Y:       (i * l.tileSize) - y2,
				W:       l.tileSize,
				H:       l.tileSize,
				Row:     int32(id-ts.firstGID) / ts.cols,
				Col:     int32(id-ts.firstGID) % ts.cols,
			}); err != nil {
				return
			}
		}
	}
	return
}

func (l *tileLayer) Update() error {
	return nil
}

func (l *tileLayer) getTileSetByID(id int) TileSet {
	for i := 0; i < len(l.sets)-1; i++ {
		if id >= l.sets[i].firstGID && id < l.sets[i+1].firstGID {
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
