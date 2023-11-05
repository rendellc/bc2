package ui

import (
	"fmt"
	"log"
	"rendellc/bc2/langs"
	"rendellc/bc2/storage"

	tea "github.com/charmbracelet/bubbletea"
)

type Editor struct {
	storage      *storage.Storage
	scriptInfo   *storage.ScriptInfo
	scriptEditor scriptEditor
	hasFocus     bool
}

type scriptLoadedMsg langs.Script

func LoadScriptCmd(storage *storage.Storage, scriptInfo storage.ScriptInfo) tea.Cmd {
	return func() tea.Msg {
		content, err := storage.LoadScript(scriptInfo)
		if err != nil {
			log.Printf("Unable to load script: %v", err)
		}

		return scriptLoadedMsg(langs.LuaScript(content))
	}
}

func CreateEditor(storage *storage.Storage) Editor {
	return Editor{
		storage:      storage,
		scriptEditor: CreateScriptEditor(),
	}
}

func (e Editor) Init() tea.Cmd {
	log.Printf("Initialize editor with %v", e.storage)
	return nil
}

func (e Editor) Update(msg tea.Msg) (Editor, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case focusChangeMsg:
		e.hasFocus = (focusElement(msg) == focusEditor)
		return e, nil
	case historySelectMsg:
		log.Printf("Editor history selected: %T %+v", msg, msg)
		e.scriptInfo = new(storage.ScriptInfo)
		*e.scriptInfo = storage.ScriptInfo(msg)

		return e, LoadScriptCmd(e.storage, *e.scriptInfo)
	case scriptLoadedMsg:

		log.Printf("Script loaded: %T %+v", msg, msg)
		e.scriptEditor.Reset(langs.Script(msg))

		return e, focusChangeCmd(focusEditor)
	case tea.KeyMsg:
		if !e.hasFocus {
			return e, nil
		}
		switch msg.String() {
		case "ctrl+s":
			if e.scriptInfo == nil {
				log.Printf("Saving new script")
			} else {
				log.Printf("Saving %s", (*e.scriptInfo).Name())
			}

		}
	}

	e.scriptEditor, cmd = e.scriptEditor.Update(msg)
	cmds = append(cmds, cmd)

	return e, tea.Batch(cmds...)
}

func (e Editor) View() string {
	return fmt.Sprintf("Editor\n\n%s", e.scriptEditor.View())
}
