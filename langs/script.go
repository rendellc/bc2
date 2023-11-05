package langs


type Script interface {
	Cells() []ScriptCell
	ToString() string
}

type ScriptCell string
