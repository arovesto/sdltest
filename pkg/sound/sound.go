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

var manager *Manager

type Manager struct {
	sfxs  map[string]*mix.Chunk
	music map[string]*mix.Music
}

func init() {
	if err := mix.OpenAudio(22050, sdl.AUDIO_S16, 2, 4096); err != nil {
		panic(err)
	}
	manager = &Manager{sfxs: map[string]*mix.Chunk{}, music: map[string]*mix.Music{}}
}

func (m *Manager) Load(file string, t Type, id string) (err error) {
	switch t {
	case MUSIC:
		if music, err := mix.LoadMUS(file); err == nil {
			m.music[id] = music
			return
		} else {
			return err
		}
	case SFX:
		if chunk, err := mix.LoadWAV(file); err == nil {
			m.sfxs[id] = chunk
			return
		} else {
			return err
		}
	default:
		return errors.New("unknown type")
	}
}

func (m *Manager) PlaySound(id string, loop int) error {
	_, err := m.sfxs[id].Play(-1, loop)
	return err
}

func (m *Manager) PlayMusic(id string, loop int) error {
	return m.music[id].Play(loop)
}

func (m *Manager) Destroy() {
	mix.CloseAudio()
}

func Load(file string, t Type, id string) error {
	return manager.Load(file, t, id)
}

func PlaySound(id string, loop int) error {
	return manager.PlaySound(id, loop)
}

func PlayMusic(id string, loop int) error {
	return manager.PlayMusic(id, loop)
}
