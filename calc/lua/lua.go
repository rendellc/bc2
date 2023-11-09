package lua

import (
	"fmt"
	"regexp"
	"rendellc/bc2/calc"

	glua "github.com/yuin/gopher-lua"
)

type LuaScriptInterpreter struct {
	state *glua.LState
}

const (
	opAssignment = iota		// line of form `x = 12*y*x + 2`
	opExpression			// line of form `12*x - 2` (not valid lua)
	opNone					// whitespace
)

func CreateLuaScriptInterpreter() LuaScriptInterpreter {
	return LuaScriptInterpreter{
		state: glua.NewState(),
	}
}
func (h LuaScriptInterpreter) Close() {
	defer h.state.Close()
}

	func (h LuaScriptInterpreter) Run(script string) calc.InterpreterResult {
	lines := calc.SplitLines(script)

	results := []calc.InterpreterLineResult{}

	for lineIndex, line := range lines {
		lineNumber := lineIndex + 1

		err := h.state.DoString(line)
		if err != nil {
			results = append(results, calc.NewLineResultError(lineNumber, err.Error()))
			continue
		} 

		varName, err := extractVariableName(line)
		if err != nil {
			results = append(results, calc.NewLineResultError(lineNumber, err.Error()))
		} else {
			value := h.state.GetGlobal(varName).String()
			results = append(results, calc.NewLineResultNormal(lineNumber, value))
		}
	}

	return calc.InterpreterResult(results)
}

func extractVariableName(line string) (string, error) {
	re := regexp.MustCompile(`^([a-zA-Z]\w*)\s*=`)

	matches := re.FindStringSubmatch(line)
	if len(matches) != 2 {
		return "", fmt.Errorf("Expected exactly 2 matches in %s", line)
	}

	return matches[1], nil
}

