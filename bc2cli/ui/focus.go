package ui

import tea "github.com/charmbracelet/bubbletea"

type focusElement int

const (
	focusHistory focusElement = iota
	focusEditor
	focusFilename
)

type focusChangeMsg focusElement

func focusChangeCmd(newFocusElement focusElement) tea.Cmd {
	return func() tea.Msg {
		return focusChangeMsg(newFocusElement)
	}
}
