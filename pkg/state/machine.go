package state

type Machine struct {
	states     []State
	onTopState []State
}

// TODO make this into a manager
func NewMachine() *Machine {
	return &Machine{}
}

func (m *Machine) PushState(s State) error {
	if len(m.states) != 0 {
		if err := m.states[len(m.states)-1].OnSwitch(); err != nil {
			return err
		}
	}
	m.states = append(m.states, s)
	return s.OnEnter()
}

// TODO implement "PushOnTop" mechanics here so can have menu above game, previous staff is drawn but not updated
func (m *Machine) PushOnTop() error {
	return nil
}

func (m *Machine) ChangeState(s State) error {
	if err := m.PopState(); err != nil {
		return err
	}
	return m.PushState(s)
}

func (m *Machine) PopState() error {
	if len(m.states) == 0 {
		return nil
	}
	if err := m.states[len(m.states)-1].OnExit(); err != nil {
		return err
	}
	m.states = m.states[:len(m.states)-1]
	if len(m.states) == 0 {
		return nil
	}
	return m.states[len(m.states)-1].OnContinue()
}

func (m *Machine) Update() error {
	if len(m.states) != 0 {
		return m.states[len(m.states)-1].Update()
	}
	return nil
}

func (m *Machine) Render() error {
	if len(m.states) != 0 {
		return m.states[len(m.states)-1].Render()
	}
	return nil
}
