package calc


type Interpreter interface {
	Run(script string) InterpreterResult
}

type InterpreterResult []InterpreterLineResult


type LineResultType int
const (
	LineResultNormal LineResultType = iota
	LineResultError
)

type InterpreterLineResult struct {
	lineNumber int
	resultType LineResultType
	message string
}

func (r InterpreterLineResult) Line() int { return r.lineNumber }
func (r InterpreterLineResult) ResultType() LineResultType { return r.resultType }
func (r InterpreterLineResult) Message() string { return r.message }

func NewLineResultNormal(lineNumber int, message string) InterpreterLineResult {
	return InterpreterLineResult{
		lineNumber: lineNumber,
		resultType: LineResultNormal,
		message: message,
	}
}

func NewLineResultError(lineNumber int, message string) InterpreterLineResult {
	return InterpreterLineResult{
		lineNumber: lineNumber,
		resultType: LineResultError,
		message: message,
	}
}
