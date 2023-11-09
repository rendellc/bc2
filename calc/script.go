package calc


type Script interface {
	Cells() []ScriptCell
	ToString() string
}

type ScriptCell string
