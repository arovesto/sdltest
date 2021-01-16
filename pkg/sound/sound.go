package sound

import (
	"errors"

	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
)

type Type int

const (
	MUSIC Type = iota
	SFX
)

var m *Manager

type Manager struct {
	sfxs  map[string]*mix.Chunk
	music map[string]*mix.Music

	volume int
}

func init() {
	if err := mix.OpenAudio(mix.DEFAULT_FREQUENCY, sdl.AUDIO_S16, 2, 4096); err != nil {
		panic(err)
	}
	if err := mix.Init(mix.INIT_FLAC | mix.INIT_MP3); err != nil {
		panic(err)
	}
	//if err := mix.OpenAudio(mix.DEFAULT_FREQUENCY, sdl.AUDIO_S16, 2, 4096); err != nil {
	//	panic(err)
	//}
	mix.Volume(-1, 16)
	m = &Manager{sfxs: map[string]*mix.Chunk{}, music: map[string]*mix.Music{}, volume: mix.Volume(-1, -1)}
}

func Load(file string, t Type, id string) (err error) {
	switch t {
	case MUSIC:
		m.music[id], err = mix.LoadMUS(file)
	case SFX:
		m.sfxs[id], err = mix.LoadWAV(file)
	default:
		return errors.New("unknown type")
	}
	return
}

func PlaySound(id string, loop int) error {
	_, err := m.sfxs[id].Play(-1, loop)
	return err
}

func PlayMusic(id string, loop int) error {
	return m.music[id].Play(loop)
}

func HaltMusic() {
	mix.HaltMusic()
}

func IncVolume() {
	m.volume += 4
	if m.volume > mix.MAX_VOLUME {
		m.volume = mix.MAX_VOLUME
	}
	mix.Volume(-1, m.volume)
}

func DecVolume() {
	m.volume -= 4
	if m.volume < 0 {
		m.volume = 0
	}
	mix.Volume(-1, m.volume)
}

func Destroy() {
	mix.CloseAudio()
}

func GetVolume() int {
	return m.volume
}
