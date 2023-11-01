package langs

import (
	"fmt"
	"regexp"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

type luaScriptInterpreter struct {
	state *lua.LState
}

const (
	opAssignment = iota		// line of form `x = 12*y*x + 2`
	opExpression			// line of form `12*x - 2` (not valid lua)
	opNone					// whitespace
)

func CreateLuaScriptInterpreter() luaScriptInterpreter {
	return luaScriptInterpreter{
		state: lua.NewState(),
	}
}
func (h luaScriptInterpreter) Close() {
	defer h.state.Close()
}
func (h luaScriptInterpreter) Run(lines []string) []lineEvalResult {
	results := []lineEvalResult{}

	for i, line := range lines {
		operationType := detectOperationType(line)

		if operationType == opExpression {
			// this is not valid lua, need to prepend it 
			// to make it an expression
			line = fmt.Sprintf("unnamedVariable%d = %s", i, line) 
		}

		err := h.state.DoString(line)
		if err != nil {
			results = append(results, EvalError(err))
			continue
		} 


		hasValue := operationType == opAssignment || operationType == opExpression
		if hasValue {
			line = strings.TrimSpace(line)
			if varName, err := extractVariableName(line); err != nil {
				results = append(results, EvalError(fmt.Errorf("value parse error")))
			} else {
				value := h.state.GetGlobal(varName).String()
				results = append(results, EvalOK(value))
			}
		} else {
			results = append(results, EvalOK("<no value>"))
		}
	}

	return results
}

func detectOperationType(luaLine string) uint {
	if len(strings.TrimSpace(luaLine)) == 0 {
		return opNone
	}
	if strings.ContainsRune(luaLine, '=') {
		return opAssignment
	}

	return opExpression
}


func extractVariableName(line string) (string, error) {
	re := regexp.MustCompile(`^([a-zA-Z]\w*)\s*=`)

	matches := re.FindStringSubmatch(line)
	if len(matches) != 2 {
		return "", fmt.Errorf("Expected exactly 2 matches in %s", line)
	}

	return matches[1], nil
}

