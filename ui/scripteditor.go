package ui

import (
	"fmt"
	"rendellc/bc2/langs"

	tea "github.com/charmbracelet/bubbletea"
)

type scriptEditor struct {
	cells        []cell
	focusedCellIndex   int
	interpreterBuilder func() langs.LuaScriptInterpreter
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
		cells:        []cell{initialCell},
		interpreterBuilder: langs.CreateLuaScriptInterpreter,
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

func (s scriptEditor) Init() tea.Cmd {
	return nil
}

func (s scriptEditor) Update(msg tea.Msg) (scriptEditor, tea.Cmd) {
	switch msg := msg.(type) {
	case evaluateScriptMsg:
		cellResults := []langs.CellResult(msg)
		for i := range cellResults {
			s.cells[i].SetResult(cellResults[i].Ok())
		}

		return s, nil
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			if s.focusedCellIndex > 0 {
				s.cells[s.focusedCellIndex].Blur()
				s.focusedCellIndex = s.focusedCellIndex - 1
			}

			s.cells[s.focusedCellIndex].Focus()
			s.removeTrailingEmptyCells()
			return s, nil
		case tea.KeyDown:
			if s.focusedCellIndex < len(s.cells)-1 {
				s.cells[s.focusedCellIndex].Blur()
				s.focusedCellIndex = s.focusedCellIndex + 1
			}

			s.cells[s.focusedCellIndex].Focus()
			return s, nil
		case tea.KeyEnter:
			s.insertCellAfter(s.focusedCellIndex)
			s.cells[s.focusedCellIndex].Blur()
			s.focusedCellIndex += 1
			s.cells[s.focusedCellIndex].Focus()
			return s, nil
		}
	}

	ta, cmd := s.cells[s.focusedCellIndex].Update(msg)
	s.cells[s.focusedCellIndex] = ta

	s.removeTrailingEmptyCells()

	return s, cmd
}

func (s scriptEditor) View() string {

	allCellView := ""
	for _, cell := range s.cells {
		// isFinal := i < len(s.cellEditors)
		allCellView += cell.View() + "\n"
	}

	numberOfCells := len(s.cells)
	allCellView += fmt.Sprintf("\n\n\nNumber of cells: %d", numberOfCells)

	return allCellView
}

func (s *scriptEditor) Reset(script langs.Script) {
	cells := script.Cells()
	s.setNumberOfCells(len(cells))

	for i, cell := range cells {
		s.cells[i].SetValue(string(cell))
	}

	s.relimFocusedCellIndex()
}
