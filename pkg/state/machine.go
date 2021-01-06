package state

type Machine struct {
	states []State
}

func NewMachine() *Machine {
	return &Machine{}
}

func (m *Machine) PushState(s State) error {
	m.states = append(m.states, s)
	return s.OnEnter()
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
	err := m.states[len(m.states)-1].OnExit()
	m.states = m.states[:len(m.states)-1]
	return err
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
