package object

type animation struct {
	shooterObject
	animChanged uint32
	animSpeed   uint32
}

func NewAnimation(st Properties) GameObject {
	st.IgnoreCam = true
	return &animation{shooterObject: newShooterObj(st), animSpeed: st.AnimSpeed}
}

func (a *animation) Update() error {
	a.changeSprite()
	return nil
}
