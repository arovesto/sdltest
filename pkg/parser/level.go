package parser

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

	"github.com/arovesto/sdl/pkg/model"
	"gopkg.in/yaml.v2"

	"github.com/arovesto/sdl/pkg/level"

	"github.com/arovesto/sdl/pkg/math"
	"github.com/arovesto/sdl/pkg/object"

	"github.com/arovesto/sdl/pkg/game/global"
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
	TileSets []*level.TileSet `xml:"tileset"`
	Nodes    []struct {
		XMLName    xml.Name
		Objects    []LevelObject `xml:"object"`
		Data       string        `xml:"data"`
		Properties struct {
			Property []struct {
				Name  string `xml:"name,attr"`
				Type  string `xml:"type,attr"`
				Value string `xml:"value,attr"`
			} `xml:"property"`
		} `xml:"properties"`
	} `xml:",any"`
}

type LevelObject struct {
	XMLName    xml.Name `xml:"object"`
	Name       string   `xml:"name,attr"`
	Type       string   `xml:"type,attr"`
	X          float64  `xml:"x,attr"`
	Y          float64  `xml:"y,attr"`
	Width      float64  `xml:"width,attr"`
	Height     float64  `xml:"height,attr"`
	Properties struct {
		Property []struct {
			Name  string `xml:"name,attr"`
			Type  string `xml:"type,attr"`
			Value string `xml:"value,attr"`
		} `xml:"property"`
	} `xml:"properties"`
}

func parseModels(file string) (map[string]model.Model, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	res := map[string]model.Model{}
	return res, yaml.Unmarshal(data, &res)
}

func ParseLevel(levelFile, modelsFile string) (*level.Level, error) {
	models, err := parseModels(modelsFile)
	if err != nil {
		return nil, err
	}

	contents, err := ioutil.ReadFile(levelFile)
	if err != nil {
		return nil, err
	}
	var format Format
	if err = xml.Unmarshal(contents, &format); err != nil {
		return nil, err
	}

	for _, t := range format.TileSets {
		t.DrawWidth = format.TileWidth
		t.DrawHeight = format.TileHeight
	}

	var collisionObjects []object.GameObject
	var collisionTileLayers []*level.TileLayer

	var layers []level.Layer
	for _, n := range format.Nodes {
		c := false
		for _, p := range n.Properties.Property {
			if p.Name == "collision" {
				c = true
			}
		}
		switch n.XMLName.Local {
		case "objectgroup":
			var objects []object.GameObject
			for _, o := range n.Objects {
				state := object.Properties{Pos: math.NewVec(o.X, o.Y)}
				for _, p := range o.Properties.Property {
					switch p.Name {
					case "model":
						state.Model = models[p.Value].GetCopy()
					case "callbackID":
						state.Callback = global.ID(mustInt(p.Value))
					case "animSpeed":
						state.AnimSpeed = uint32(mustInt(p.Value))
					}
				}
				if obj, err := object.Create(o.Type, state); err == nil {
					objects = append(objects, obj)
					if c {
						collisionObjects = append(collisionObjects, obj)
					}
				} else {
					return nil, err
				}
			}
			layers = append(layers, level.NewObjectLayer(objects))
		case "layer":
			data, err := base64.StdEncoding.DecodeString(strings.TrimSpace(n.Data))
			if err != nil {
				return nil, err
			}
			m := make(level.Tiles, format.Height)
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

			layer := level.NewTileLayer(format.TileWidth, format.TileSets, m, c)
			layers = append(layers, layer)
			if c {
				collisionTileLayers = append(collisionTileLayers, layer)
			}
		}
	}
	return level.NewLevel(format.TileSets, append([]level.Layer{level.NewCollisionLayer(collisionObjects, collisionTileLayers)}, layers...)), nil
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
