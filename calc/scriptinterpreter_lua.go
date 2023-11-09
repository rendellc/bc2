package calc

import (
	"fmt"
	"regexp"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

type LuaScriptInterpreter struct {
	state *lua.LState
}

type LuaScript string

func (s LuaScript) Cells() []ScriptCell {
	lines := SplitLines(s.ToString())
	cells := []ScriptCell{}
	for _, line := range lines {
		cells = append(cells, ScriptCell(line))
	}

	return cells
}

func (s LuaScript) ToString() string {
	return string(s)
}

const (
	opAssignment = iota		// line of form `x = 12*y*x + 2`
	opExpression			// line of form `12*x - 2` (not valid lua)
	opNone					// whitespace
)

func CreateLuaScriptInterpreter() LuaScriptInterpreter {
	return LuaScriptInterpreter{
		state: lua.NewState(),
	}
}
func (h LuaScriptInterpreter) Close() {
	defer h.state.Close()
}

func (h LuaScriptInterpreter) Run(script Script) ([]CellResult, error) {
	switch script := script.(type) {
	case LuaScript:
		results := []CellResult{}
		cells := script.Cells()

		for i, cell := range cells {
			luaLine := string(cell)
			operationType := detectOperationType(luaLine)

			if operationType == opExpression {
				// this is not valid lua, need to prepend it 
				// to make it an expression
				luaLine = fmt.Sprintf("unnamedVariable%d = %s", i, luaLine) 
			}

			err := h.state.DoString(luaLine)
			if err != nil {
				results = append(results, EvalError(err))
				continue
			} 


			hasValue := operationType == opAssignment || operationType == opExpression
			if hasValue {
				luaLine = strings.TrimSpace(luaLine)
				if varName, err := extractVariableName(luaLine); err != nil {
					results = append(results, EvalError(fmt.Errorf("value parse error")))
				} else {
					value := h.state.GetGlobal(varName).String()
					results = append(results, EvalOK(value))
				}
			} else {
				results = append(results, EvalOK("<no value>"))
			}
		}

		return results, nil
	default:
		return nil, fmt.Errorf("Script is not a luaScript type")
	}
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

