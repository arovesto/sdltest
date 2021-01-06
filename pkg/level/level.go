package level

type Level struct {
	sets   []TileSet
	layers []Layer
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
	firstGID   int
	tileWidth  int32
	tileHeight int32
	spacing    int32
	margin     int32
	w          int32
	h          int32
	cols       int32
	name       string
}
