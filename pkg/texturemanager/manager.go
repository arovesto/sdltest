package texturemanager

import (
	"errors"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"

	"golang.org/x/xerrors"
)

var manager *TextureManager

type TextureManager struct {
	textureMap map[string]*sdl.Texture
	renderer   *sdl.Renderer
}

func InitManager(r *sdl.Renderer) {
	manager = &TextureManager{textureMap: map[string]*sdl.Texture{}, renderer: r}
}

func Load(opts LoadOpts) error {
	if manager == nil {
		return errors.New("manager not initialized")
	}
	return manager.Load(opts)
}

func Draw(opts DrawOpts) error {
	if manager == nil {
		return errors.New("manager not initialized")
	}
	return manager.Draw(opts)
}

func DrawTile(opts DrawTileOpts) error {
	if manager == nil {
		return errors.New("manager not initialized")
	}
	return manager.DrawTile(opts)
}

func DrawRect(id string, src, dst sdl.Rect) error {
	if manager == nil {
		return errors.New("manager not initialized")
	}
	return manager.DrawRect(id, src, dst)
}

func Destroy() error {
	if manager == nil {
		return errors.New("manager not initialized")
	}
	return manager.Destroy()
}

func Delete(id string) {
	if manager != nil {
		manager.Delete(id)
	}
}

type LoadOpts struct {
	Path string
	ID   string
}

func (t *TextureManager) Load(opts LoadOpts) error {
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

type DrawOpts struct {
	ID    string
	X     int32
	Y     int32
	W     int32
	H     int32
	Row   int32
	Col   int32
	Flip  sdl.RendererFlip
	Angle float64
	Alpha uint8
}

func (t *TextureManager) Draw(opts DrawOpts) error {
	if texture, ok := t.textureMap[opts.ID]; ok {
		src := sdl.Rect{X: opts.W * opts.Col, Y: opts.H * opts.Row, W: opts.W, H: opts.H}
		dst := sdl.Rect{X: opts.X, Y: opts.Y, W: opts.W, H: opts.H}
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

func (t *TextureManager) DrawRect(id string, src, dst sdl.Rect) error {
	if texture, ok := t.textureMap[id]; ok {
		return t.renderer.Copy(texture, &src, &dst)
	} else {
		return xerrors.Errorf("texture not found: %v", id)
	}

}

func (t *TextureManager) DrawTile(opts DrawTileOpts) error {
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

func (t *TextureManager) Delete(id string) {
	delete(t.textureMap, id)
}

func (t *TextureManager) Destroy() error {
	for _, p := range t.textureMap {
		if err := p.Destroy(); err != nil {
			return err
		}
	}
	return nil
}
