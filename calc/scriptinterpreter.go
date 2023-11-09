package calc


type ScriptInterpreter interface {
	Run(script Script) []CellResult
}

type CellResult struct {
	ok string
	err error
}

func (r CellResult) Ok() string {
	return r.ok
}

func (r CellResult) Err() error {
	return r.err
}

func (r CellResult) IsOK() bool {
	return r.err == nil
}

func EvalOK(ok string) CellResult {
	return CellResult{
		ok: ok,
	}
}

func EvalError(err error) CellResult {
	return CellResult{
		err: err,
	}
}


