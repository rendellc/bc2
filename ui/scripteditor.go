package ui

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)


type scriptEditor struct {
	textarea textarea.Model
}

func (s scriptEditor) Init() tea.Cmd {
	return nil
}

func (s scriptEditor) Update(msg tea.Msg) (scriptEditor, tea.Cmd) {
	ta, cmd := s.textarea.Update(msg)
	s.textarea = ta

	return s, cmd
}

func (s scriptEditor) View() string {
	return s.textarea.View()
}


func (s *scriptEditor) Reset(script string) {
	s.textarea.SetValue(script)
}
