package level

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/binary"
	"encoding/xml"
	"io"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/arovesto/sdl/pkg/math"
	"github.com/arovesto/sdl/pkg/object"

	"github.com/arovesto/sdl/pkg/game/global"
	"github.com/arovesto/sdl/pkg/texturemanager"
)

type Format struct {
	XMLName xml.Name `xml:"map"`

	Width      int32 `xml:"width,attr"`
	Height     int32 `xml:"height,attr"`
	TileWidth  int32 `xml:"tilewidth,attr"`
	TileHeight int32 `xml:"tileheight,attr"`
	Props      struct {
		Properties []struct {
			Name  string `xml:"name,attr"`
			Value string `xml:"value,attr"`
		} `xml:"property"`
	} `xml:"properties"`
	TileSets []struct {
		FirstGID   int    `xml:"firstgid,attr"`
		TileWidth  int32  `xml:"tilewidth,attr"`
		TileHeight int32  `xml:"tileheight,attr"`
		Spacing    int32  `xml:"spacing,attr"`
		Margin     int32  `xml:"margin,attr"`
		Columns    int32  `xml:"columns,attr"`
		Name       string `xml:"name,attr"`
		Image      struct {
			Width  int32  `xml:"width,attr"`
			Height int32  `xml:"height,attr"`
			Source string `xml:"source,attr"`
		} `xml:"image"`
	} `xml:"tileset"`
	Nodes []struct {
		XMLName xml.Name
		Objects []Object `xml:"object"`
		Data    string   `xml:"data"`
	} `xml:",any"`
}

type Object struct {
	XMLName    xml.Name `xml:"object"`
	Name       string   `xml:"name,attr"`
	Type       string   `xml:"type,attr"`
	X          int32    `xml:"x,attr"`
	Y          int32    `xml:"y,attr"`
	Width      int32    `xml:"width,attr"`
	Height     int32    `xml:"height,attr"`
	Properties struct { // TODO store additional pngs in properties like "texture-..."
		Property []struct {
			Name  string `xml:"name,attr"`
			Type  string `xml:"type,attr"`
			Value string `xml:"value,attr"`
		} `xml:"property"`
	} `xml:"properties"`
}

func Parse(name string) (*Level, error) {
	contents, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	var format Format
	if err = xml.Unmarshal(contents, &format); err != nil {
		return nil, err
	}
	var tileSets []TileSet

	for _, t := range format.Props.Properties {
		if err = texturemanager.Load(texturemanager.LoadOpts{Path: global.GetAssetsPath(t.Value), ID: t.Name}); err != nil {
			return nil, err
		}
	}

	for _, s := range format.TileSets {
		if err = texturemanager.Load(texturemanager.LoadOpts{
			Path: global.GetAssetsPath(s.Image.Source),
			ID:   s.Name,
		}); err != nil {
			return nil, err
		}
		tileSets = append(tileSets, TileSet{
			firstGID:   s.FirstGID,
			tileWidth:  s.TileWidth,
			tileHeight: s.TileHeight,
			spacing:    s.Spacing,
			margin:     s.Margin,
			w:          s.Image.Width,
			h:          s.Image.Height,
			cols:       s.Columns,
			name:       s.Name,
		})
	}

	var layers []Layer
	for _, n := range format.Nodes {
		switch n.XMLName.Local {
		case "objectgroup":
			var objects []object.GameObject
			for _, o := range n.Objects {
				state := object.Properties{Pos: math.NewVecInt(o.X, o.Y), Size: math.NewVecInt(o.Width, o.Height), ID: o.Name}
				for _, p := range o.Properties.Property {
					switch p.Name {
					case "frames":
						state.Cols = int32(mustInt(p.Value))
					case "texture":
						state.ID = p.Value
					case "callbackID":
						state.Callback = global.ID(mustInt(p.Value))
					case "animSpeed":
						state.AnimSpeed = uint32(mustInt(p.Value))
					case "maxSpeed":
						state.MaxSpeed = mustFloat(p.Value)
					}
				}
				if obj, err := object.Create(o.Type, state); err == nil {
					objects = append(objects, obj)
				} else {
					return nil, err
				}
			}
			layers = append(layers, NewObjectLayer(objects))
		case "layer":
			data, err := base64.StdEncoding.DecodeString(strings.TrimSpace(n.Data))
			if err != nil {
				return nil, err
			}
			m := make(Tiles, format.Height)
			for i := range m {
				m[i] = make([]int, format.Width)
			}

			reader, err := zlib.NewReader(bytes.NewReader(data))
			if err != nil {
				return nil, err
			}
			ids, err := readInts(reader)
			if err != nil {
				return nil, err
			}

			for row := int32(0); row < format.Height; row++ {
				for col := int32(0); col < format.Width; col++ {
					m[row][col] = ids[row*format.Width+col]
				}
			}
			layers = append(layers, NewTileLayer(format.TileWidth, tileSets, m))
		}
	}
	return &Level{sets: tileSets, layers: layers}, nil
}

func readInts(r io.Reader) (ids []int, err error) {
	for {
		var n int32
		if err = binary.Read(r, binary.LittleEndian, &n); err != nil {
			if err == io.EOF {
				err = nil
			}
			return
		} else {
			ids = append(ids, int(n))
		}
	}
}

func mustInt(c string) int {
	res, err := strconv.Atoi(c)
	if err != nil {
		panic(err)
	}
	return res
}

func mustFloat(c string) float64 {
	res, err := strconv.ParseFloat(c, 64)
	if err != nil {
		panic(err)
	}
	return res
}
