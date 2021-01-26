package level

import (
	"github.com/arovesto/sdl/pkg/game/global"
	"github.com/arovesto/sdl/pkg/math"
	"github.com/arovesto/sdl/pkg/object"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type Level struct {
	sets   []*TileSet
	layers []Layer

	mainLayer           *objectLayer
	collisionTileLayers []*tileLayer

	player object.GameObject
}

var CurrentLevel *Level

func NewLevel(s []*TileSet, l []Layer) *Level {
	lvl := &Level{sets: s, layers: l}

	for _, layer := range l {
		switch l := layer.(type) {
		case *objectLayer:
			if l.collision {
				lvl.mainLayer = l
			}
			for _, o := range l.objects {
				if o.GetType() == object.PlayerType {
					lvl.player = o
				}
			}
		case *tileLayer:
			if l.collision {
				lvl.collisionTileLayers = append(lvl.collisionTileLayers, l)
			}
		}
	}
	return lvl
}

func (l *Level) Update() (err error) {
	for _, l := range l.layers {
		if err = l.Update(); err != nil {
			return
		}
	}
	if err = l.mainLayer.Collision(l.collisionTileLayers); err != nil {
		return err
	}
	return
}

func (l *Level) Render() (err error) {
	for _, l := range l.layers {
		if err = l.Render(); err != nil {
			return
		}
	}
	return
}

func (l *Level) GetPlayer() object.GameObject {
	return l.player
}

func (l *Level) NewObj(gameObject object.GameObject) {
	l.mainLayer.objects = append(l.mainLayer.objects, gameObject)
}

// TODO implement some better datastructure to allow order + fast inserts and deletes
func (l *Level) DelObject(gameObject object.GameObject) {
	l.mainLayer.deletedIDs[gameObject] = struct{}{}
}

func (l *Level) Destroy() error {
	for _, t := range l.sets {
		if err := t.Destroy(); err != nil {
			return err
		}
	}
	return nil
}

type TileSet struct {
	FirstGID int   `xml:"firstgid,attr"`
	TWidth   int32 `xml:"tilewidth,attr"`
	THeight  int32 `xml:"tileheight,attr"`
	Spacing  int32 `xml:"spacing,attr"`
	Margin   int32 `xml:"margin,attr"`
	Image    struct {
		W    int32  `xml:"width,attr"`
		H    int32  `xml:"height,attr"`
		Path string `xml:"source,attr"`
	} `xml:"image"`
	Cols int32  `xml:"columns,attr"`
	Name string `xml:"name,attr"`

	DrawWidth  int32
	DrawHeight int32

	t *sdl.Texture
}

func (t *TileSet) Draw(where, which math.IntVector) error {
	if t.t == nil {
		s, err := img.Load(global.GetAssetsPath(t.Image.Path))
		if err != nil {
			return err
		}
		t.t, err = global.Renderer.CreateTextureFromSurface(s)
		if err != nil {
			return err
		}
		s.Free()
	}
	src := sdl.Rect{
		X: t.Margin + (t.Spacing+t.TWidth)*which.X,
		Y: t.Margin + (t.Spacing+t.THeight)*which.Y,
		W: t.TWidth,
		H: t.THeight,
	}
	dst := sdl.Rect{X: where.X, Y: where.Y, W: t.DrawWidth, H: t.DrawHeight}
	return global.Renderer.CopyEx(t.t, &src, &dst, 0, nil, sdl.FLIP_NONE)
}

func (t *TileSet) Destroy() error {
	if t.t != nil {
		return t.t.Destroy()
	}
	return nil
}
