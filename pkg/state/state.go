package state

type ID int

type State interface {
	Update() error
	Render() error
	OnEnter() error
	OnExit() error
	GetID() ID
}

var (
	Pause    State
	Menu     State
	Play     State
	GameOver State
	Between  State
)
