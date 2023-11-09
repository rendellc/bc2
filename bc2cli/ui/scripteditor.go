package ui

import (
	"fmt"
	"rendellc/bc2/calc"
	"rendellc/bc2/calc/lua"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var scriptEditorStyle = lipgloss.NewStyle()


type scriptChangedMsg struct{}
func scriptChangedCmd() tea.Msg {
	return scriptChangedMsg{}
}

type scriptEditor struct {
	cells              []cell
	focusedCellIndex   int
	numberOfEvaluations int
	interpreterBuilder func() lua.LuaScriptInterpreter
}

func (s *scriptEditor) getFocusedCell() cell {
	return s.cells[s.focusedCellIndex]
}

func (s *scriptEditor) relimFocusedCellIndex() {
	if s.focusedCellIndex < 0 {
		s.focusedCellIndex = 0
	}

	if s.focusedCellIndex >= len(s.cells) {
		s.focusedCellIndex = len(s.cells) - 1
	}
}

func (s *scriptEditor) removeTrailingEmptyCells() {
	for i := len(s.cells) - 1; i > s.focusedCellIndex; i-- {
		isEmpty := len(s.cells[i].Value()) == 0

		if isEmpty {
			s.removeCellAt(i)
		} else {
			return
		}
	}
}

func (s scriptEditor) GetScriptString() string {
	content := ""
	for i, cell := range s.cells {
		isFinal := i == len(s.cells)-1
		content += cell.Value()

		if !isFinal {
			content += "\n"
		}
	}

	return content
}

func CreateScriptEditor() scriptEditor {
	initialCell := createCell()
	initialCell.Focus()

	return scriptEditor{
		cells:              []cell{initialCell},
		numberOfEvaluations: 0,
		interpreterBuilder: lua.CreateLuaScriptInterpreter,
	}
}

func (s *scriptEditor) setNumberOfCells(count int) {
	cellsToAdd := count - len(s.cells)

	if cellsToAdd < 0 {
		s.cells = s.cells[:count]
	}
	if cellsToAdd > 0 {
		for i := 0; i < cellsToAdd; i++ {
			s.cells = append(s.cells, createCell())
		}
	}
}

func (s *scriptEditor) insertCellAfter(index int) {
	s.cells = append(s.cells, createCell())

	for i := len(s.cells) - 2; i > index; i-- {
		s.cells[i+1] = s.cells[i]
	}

	s.cells[index+1] = createCell()
}

func (s *scriptEditor) removeCellAt(index int) {
	s.cells = append(s.cells[:index], s.cells[index+1:]...)
}

func (s scriptEditor) getMaxCellValueWidth() int {
	maxLength := 0
	for _, cell := range s.cells {
		cellWidth := len(cell.Value())
		if cellWidth > maxLength {
			maxLength = cellWidth
		}

	}

	return maxLength
}

func (s scriptEditor) Init() tea.Cmd {
	return nil
}

func (s scriptEditor) Update(msg tea.Msg) (scriptEditor, tea.Cmd) {
	scriptMayHaveChanged := false
	switch msg := msg.(type) {
	case evaluateScriptMsg:
		s.numberOfEvaluations += 1
		for _, lineResult := range []calc.InterpreterLineResult(msg) {
			s.SetResult(lineResult)
		}

		return s, nil
	case tea.KeyMsg:
		scriptMayHaveChanged = true
		switch msg.Type {
		case tea.KeyBackspace:
			if s.focusedCellIndex == 0 {
				break
			}
			if s.cells[s.focusedCellIndex].Position() != 0 {
				break
			}

			// pressed delete at beginning of cell and have cells above
			s.cells[s.focusedCellIndex].Blur()
			previousCellValue := s.cells[s.focusedCellIndex-1].Value()
			cellValue := s.cells[s.focusedCellIndex].Value()
			s.removeCellAt(s.focusedCellIndex)
			s.focusedCellIndex -= 1
			s.cells[s.focusedCellIndex].SetValue(previousCellValue + cellValue)
			s.cells[s.focusedCellIndex].SetCursor(len(previousCellValue))
			s.cells[s.focusedCellIndex].Focus()

			return s, scriptChangedCmd
		case tea.KeyUp:
			cellPos := s.getFocusedCell().Position()
			if s.focusedCellIndex > 0 {
				s.cells[s.focusedCellIndex].Blur()
				s.focusedCellIndex = s.focusedCellIndex - 1
			}

			s.cells[s.focusedCellIndex].SetCursor(cellPos)
			s.cells[s.focusedCellIndex].Focus()
			s.removeTrailingEmptyCells()
			return s, nil
		case tea.KeyDown:
			cellPos := s.getFocusedCell().Position()
			if s.focusedCellIndex < len(s.cells)-1 {
				s.cells[s.focusedCellIndex].Blur()
				s.focusedCellIndex = s.focusedCellIndex + 1
			}

			s.cells[s.focusedCellIndex].SetCursor(cellPos)
			s.cells[s.focusedCellIndex].Focus()
			return s, nil
		case tea.KeyEnter:
			s.insertCellAfter(s.focusedCellIndex)
			s.cells[s.focusedCellIndex].Blur()
			s.focusedCellIndex += 1
			s.cells[s.focusedCellIndex].Focus()
			return s, scriptChangedCmd
		}
	}

	preUpdateFocusCellLength := s.cells[s.focusedCellIndex].Length()
	ta, cmd := s.cells[s.focusedCellIndex].Update(msg)
	postUpdateFocusCellLength := s.cells[s.focusedCellIndex].Length()

	s.cells[s.focusedCellIndex] = ta

	if preUpdateFocusCellLength != postUpdateFocusCellLength {
		scriptMayHaveChanged = true
	}

	s.removeTrailingEmptyCells()

	cmds := []tea.Cmd{
		cmd,
	}
	if scriptMayHaveChanged {
		cmds = append(cmds, scriptChangedCmd)
	}
	return s, tea.Batch(cmds...)
}

func (s scriptEditor) View(width int) string {
	allCellView := ""
	for _, cell := range s.cells {
		allCellView += cell.View(width) + "\n"
	}

	numberOfCells := len(s.cells)
	allCellView = scriptEditorStyle.Render(allCellView)
	debugInformation := strings.Builder{}
	debugInformation.WriteString(fmt.Sprintf("Number of cells: %d", numberOfCells))
	debugInformation.WriteString("\n")
	debugInformation.WriteString(fmt.Sprintf("Number of evals: %d", s.numberOfEvaluations))

	return allCellView + "\n\n\n" + debugInformation.String()
}

func (s *scriptEditor) Reset(script string) {
	lines := calc.SplitLines(script)
	s.setNumberOfCells(len(lines))

	for i, line := range lines {
		s.cells[i].SetValue(line)
	}

	s.relimFocusedCellIndex()
}

func (s *scriptEditor) SetResult(result calc.InterpreterLineResult) {
	index := result.Line() - 1
	if index < 0 || index >= len(s.cells) {
		return
	}
	s.cells[index].SetResult(result.Message())
}
