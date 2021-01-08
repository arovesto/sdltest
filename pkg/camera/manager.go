package camera

import (
	"github.com/arovesto/sdl/pkg/math"
	"golang.org/x/xerrors"
)

var m = &manager{cams: []*camera{}, mainCam: 0}

type manager struct {
	cams    []*camera
	mainCam int
}

type Opts struct {
	X int32
	Y int32
	W int32
	H int32

	ApproachSpeed float64
}

func RegisterCam(o Opts) {
	if o.ApproachSpeed == 0 {
		o.ApproachSpeed = 50
	}
	m.cams = append(m.cams, &camera{pos: math.NewVecInt(o.X, o.Y), basePos: math.NewVecInt(o.X, o.Y), size: math.NewVecInt(o.W, o.H), approachSpeed: o.ApproachSpeed})
}

func GetCamPos() math.Vector2D {
	return m.cams[m.mainCam].pos
}

func GetCamVel() math.Vector2D {
	return m.cams[m.mainCam].vel
}

func SwitchCam(to int) error {
	if len(m.cams) <= to {
		return xerrors.Errorf("can't switch to cam %v", to)
	}
	m.mainCam = to
	return nil
}

func Update() error {
	for _, c := range m.cams {
		if err := c.Update(); err != nil {
			return err
		}
	}
	return nil
}

func Reset() {
	c := m.cams[m.mainCam]
	c.pos = c.basePos
}

func GoTo(p math.Vector2D) {
	c := m.cams[m.mainCam]
	c.target = math.Sub(p, c.Center())
}
