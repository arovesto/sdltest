package level

type Level struct {
	sets   []TileSet
	layers []Layer
}

func NewLevel(s []TileSet, l []Layer) *Level {
	return &Level{sets: s, layers: l}
}

func (l *Level) Update() (err error) {
	for _, l := range l.layers {
		if err = l.Update(); err != nil {
			return
		}
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

type TileSet struct {
	FirstGID int
	TWidth   int32
	THeight  int32
	Spacing  int32
	Margin   int32
	W        int32
	H        int32
	Cols     int32
	Name     string
}
