package ui

import (
	"fmt"
	"rendellc/bc2/langs"

	// "github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type scriptEditor struct {
	cellEditors      []textinput.Model
	focusedCellIndex int
}

func (s *scriptEditor) getFocusedCell() textinput.Model {
	return s.cellEditors[s.focusedCellIndex]
}

func (s *scriptEditor) removeTrailingEmptyCells() {
	for i := len(s.cellEditors) - 1; i > s.focusedCellIndex; i-- {
		isEmpty := len(s.cellEditors[i].Value()) == 0

		if isEmpty {
			s.removeCellAt(i)
		} else {
			return
		}
	}
}

func CreateScriptEditor() scriptEditor {
	initialCell := createCellEditor()
	initialCell.Focus()

	return scriptEditor{
		cellEditors: []textinput.Model{
			initialCell,
		},
	}
}

func createCellEditor() textinput.Model {
	t := textinput.New()
	t.Placeholder = ""
	t.Prompt = ""

	t.Blur()
	return t
}

func (s *scriptEditor) setNumberOfCells(count int) {
	cellsToAdd := count - len(s.cellEditors)

	if cellsToAdd < 0 {
		s.cellEditors = s.cellEditors[:count]
	}
	if cellsToAdd > 0 {
		for i := 0; i < cellsToAdd; i++ {
			s.cellEditors = append(s.cellEditors, createCellEditor())
		}
	}
}

func (s *scriptEditor) insertCellAfter(index int) {
	s.cellEditors = append(s.cellEditors, createCellEditor())

	for i := len(s.cellEditors) - 2; i > index; i-- {
		s.cellEditors[i+1] = s.cellEditors[i]
	}

	s.cellEditors[index+1] = createCellEditor()
}

func (s *scriptEditor) removeCellAt(index int) {
	s.cellEditors = append(s.cellEditors[:index], s.cellEditors[index+1:]...)
}

func (s scriptEditor) Init() tea.Cmd {
	return nil
}

func (s scriptEditor) Update(msg tea.Msg) (scriptEditor, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			if s.focusedCellIndex > 0 {
				s.cellEditors[s.focusedCellIndex].Blur()
				s.focusedCellIndex = s.focusedCellIndex - 1
			}

			s.cellEditors[s.focusedCellIndex].Focus()
			s.removeTrailingEmptyCells()
			return s, nil
		case tea.KeyDown:
			if s.focusedCellIndex < len(s.cellEditors)-1 {
				s.cellEditors[s.focusedCellIndex].Blur()
				s.focusedCellIndex = s.focusedCellIndex + 1
			}

			s.cellEditors[s.focusedCellIndex].Focus()
			return s, nil
		case tea.KeyEnter:
			s.insertCellAfter(s.focusedCellIndex)
			s.cellEditors[s.focusedCellIndex].Blur()
			s.focusedCellIndex += 1
			s.cellEditors[s.focusedCellIndex].Focus()
		}
	}

	ta, cmd := s.cellEditors[s.focusedCellIndex].Update(msg)
	s.cellEditors[s.focusedCellIndex] = ta

	s.removeTrailingEmptyCells()

	return s, cmd
}

func (s scriptEditor) View() string {

	allCellView := ""
	for _, cell := range s.cellEditors {
		// isFinal := i < len(s.cellEditors)
		allCellView += cell.View() + "\n"
	}

	numberOfCells := len(s.cellEditors)
	allCellView += fmt.Sprintf("\n\n\nNumber of cells: %d", numberOfCells)

	return allCellView
}

func (s *scriptEditor) Reset(script langs.Script) {
	cells := script.Cells()
	s.setNumberOfCells(len(cells))

	for i, cell := range cells {
		s.cellEditors[i].SetValue(string(cell))
	}
}
