package texturemanager

import (
	"errors"

	"github.com/veandco/go-sdl2/ttf"

	"github.com/arovesto/sdl/pkg/camera"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"

	"golang.org/x/xerrors"
)

var t *manager

type manager struct {
	textureMap map[string]*sdl.Texture
	fonts      map[string]*ttf.Font
	renderer   *sdl.Renderer
}

func InitManager(r *sdl.Renderer) error {
	if err := ttf.Init(); err != nil {
		return err
	}
	t = &manager{textureMap: map[string]*sdl.Texture{}, fonts: map[string]*ttf.Font{}, renderer: r}
	return nil
}

type LoadOpts struct {
	Path string
	ID   string
}

func Load(opts LoadOpts) error {
	if t == nil {
		return errors.New("manager not initialized")
	}
	s, err := img.Load(opts.Path)
	if err != nil {
		return err
	}

	sprites, err := t.renderer.CreateTextureFromSurface(s)
	if err != nil {
		return err
	}
	t.textureMap[opts.ID] = sprites
	s.Free()
	return nil
}

// TODO add "ignore cam" option, for overlay menus
type DrawOpts struct {
	ID        string
	X         int32
	Y         int32
	W         int32
	H         int32
	Row       int32
	Col       int32
	Flip      sdl.RendererFlip
	Angle     float64
	Alpha     uint8
	IgnoreCam bool
}

func Draw(opts DrawOpts) error {
	if t == nil {
		return errors.New("manager not initialized")
	}
	if texture, ok := t.textureMap[opts.ID]; ok {
		camX, camY := camera.GetCamPos().IntPos()
		if opts.IgnoreCam {
			camX, camY = 0, 0
		}
		src := sdl.Rect{X: opts.W * opts.Col, Y: opts.H * opts.Row, W: opts.W, H: opts.H}
		dst := sdl.Rect{X: opts.X - camX, Y: opts.Y - camY, W: opts.W, H: opts.H}
		if err := texture.SetAlphaMod(opts.Alpha); err != nil {
			return err
		}
		return t.renderer.CopyEx(texture, &src, &dst, opts.Angle, nil, opts.Flip)
	} else {
		return xerrors.Errorf("texture not found: %v", opts.ID)
	}
}

type DrawTileOpts struct {
	ID      string
	Spacing int32
	Margin  int32
	X       int32
	Y       int32
	W       int32
	H       int32
	Row     int32
	Col     int32
}

func DrawTile(opts DrawTileOpts) error {
	if t == nil {
		return errors.New("manager not initialized")
	}
	if texture, ok := t.textureMap[opts.ID]; ok {
		src := sdl.Rect{
			X: opts.Margin + (opts.Spacing+opts.W)*opts.Col,
			Y: opts.Margin + (opts.Spacing+opts.H)*opts.Row,
			W: opts.W,
			H: opts.H,
		}
		dst := sdl.Rect{X: opts.X, Y: opts.Y, W: opts.W, H: opts.H}
		return t.renderer.CopyEx(texture, &src, &dst, 0, nil, sdl.FLIP_NONE)
	} else {
		return xerrors.Errorf("texture not found: %v", opts.ID)
	}
}

func DrawRect(id string, src, dst sdl.Rect) error {
	if t == nil {
		return errors.New("manager not initialized")
	}
	if texture, ok := t.textureMap[id]; ok {
		return t.renderer.Copy(texture, &src, &dst)
	} else {
		return xerrors.Errorf("texture not found: %v", id)
	}
}

func Destroy() error {
	if t == nil {
		return errors.New("manager not initialized")
	}
	for _, p := range t.textureMap {
		if err := p.Destroy(); err != nil {
			return err
		}
	}
	return nil
}

func Delete(id string) {
	if t != nil {
		delete(t.textureMap, id)
	}
}

type LoadFontOpts struct {
	ID   string
	Path string
	Size int
}

func LoadFont(o LoadFontOpts) (err error) {
	if t == nil {
		return errors.New("manager not initialized")
	}
	t.fonts[o.ID], err = ttf.OpenFont(o.Path, o.Size)
	return
}

type DrawMessageOpts struct {
	ID      string
	Message string
	X       int32
	Y       int32
	W       int32
	H       int32
	Color   sdl.Color
}

func DrawMessage(o DrawMessageOpts) error {
	if t == nil {
		return errors.New("manager not initialized")
	}
	if fnt, ok := t.fonts[o.ID]; ok {
		srf, err := fnt.RenderUTF8Solid(o.Message, o.Color)
		if err != nil {
			return err
		}
		txt, err := t.renderer.CreateTextureFromSurface(srf)
		if err != nil {
			return err
		}
		srf.Free()
		dst := sdl.Rect{X: o.X, Y: o.Y, W: o.W, H: o.H}
		return t.renderer.CopyEx(txt, nil, &dst, 0, nil, sdl.FLIP_NONE)
	} else {
		return xerrors.Errorf("font not found: %v", o.ID)
	}
}
