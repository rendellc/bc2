package langs


type ScriptInterpreter interface {
	Run(lines []string) []lineEvalResult
}

type lineEvalResult struct {
	ok string
	err error
}

func (r lineEvalResult) Ok() string {
	return r.ok
}

func (r lineEvalResult) Err() error {
	return r.err
}

func (r lineEvalResult) IsOK() bool {
	return r.err == nil
}

func EvalOK(ok string) lineEvalResult {
	return lineEvalResult{
		ok: ok,
	}
}

func EvalError(err error) lineEvalResult {
	return lineEvalResult{
		err: err,
	}
}


