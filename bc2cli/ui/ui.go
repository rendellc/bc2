package ui

import (
	"log"

	"rendellc/bc2/storage"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type app struct {
	storage      storage.Store
	history      historyBrowser
	editor       Editor
	focusElement focusElement
}

func CreateApp(storage storage.Store) app {
	return app{
		storage:      storage,
		history:      CreateHistoryBrowser(storage),
		editor:       CreateEditor(storage),
		focusElement: focusEditor,
	}
}

func (a app) Init() tea.Cmd {
	log.Printf("Initialize app")
	cmds := []tea.Cmd{
		tea.EnterAltScreen,
		focusChangeCmd(focusEditor),
		a.history.Init(),
		a.editor.Init(),
	}
	return tea.Batch(cmds...)
}

func (a app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return a, tea.Quit
		case tea.KeyTab:
			switch a.focusElement {
			case focusEditor:
				cmds = append(cmds, focusChangeCmd(focusFilename))
			case focusFilename:
				cmds = append(cmds, focusChangeCmd(focusHistory))
			case focusHistory:
				cmds = append(cmds, focusChangeCmd(focusEditor))
			}
		}
	case focusChangeMsg:
		a.focusElement = focusElement(msg)
	}

	a.history, cmd = a.history.Update(msg)
	cmds = append(cmds, cmd)
	a.editor, cmd = a.editor.Update(msg)
	cmds = append(cmds, cmd)

	return a, tea.Batch(cmds...)
}

func (a app) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Top, a.history.View(), a.editor.View())
}
