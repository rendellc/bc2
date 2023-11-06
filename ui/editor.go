package ui

import (
	"fmt"
	"log"
	"rendellc/bc2/langs"
	"rendellc/bc2/storage"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var scriptEditorWidth int = 30

type Editor struct {
	filename	textinput.Model
	storage      *storage.Storage
	scriptInfo   *storage.ScriptInfo
	scriptEditor scriptEditor
	focusElement focusElement
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


type fileSavedMsg string
func fileSavedCmd(filename string) tea.Cmd {
	return func() tea.Msg {
		return fileSavedMsg(filename)
	}
}

func (e *Editor) trySave() tea.Cmd {
	content := e.scriptEditor.GetScriptString()

	if len(content) == 0 {
		log.Printf("Script is empty. Ignoring save")
		return nil
	}

	name := ""
	if e.scriptInfo == nil {
		log.Printf("Saving new script")
		if len(e.filename.Value()) == 0 {
			// no filename provided
			log.Printf("No filename provided. Please specify")
			return focusChangeCmd(focusFilename)
		}

		name = e.filename.Value()
	} else {
		name = (*e.scriptInfo).Name()
	}

	// TODO: Handle e.scriptInfo != nil and 
	// e.filename.Value() != e.scriptInfo.Name()
	// The user has then changed the filename so the old file should be removed
	log.Printf("Saving %s", name)
	e.storage.SaveScript(name, content)

	return fileSavedCmd(name)
}



func CreateEditor(storage *storage.Storage) Editor {
	now := time.Now()
	defaultscriptname := now.Format("2006-01-02-15-04-05")
	filename := textinput.New()
	filename.Placeholder = defaultscriptname
	filename.Prompt = "Filename: "
	filename.Blur()

	return Editor{
		filename: filename,
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
		e.focusElement = focusElement(msg)
		if e.focusElement == focusFilename {
			 e.filename.Focus()
		} else {
			 e.filename.Blur()
		}
		return e, nil
	case historySelectMsg:
		log.Printf("Editor history selected: %T %+v", msg, msg)
		e.scriptInfo = new(storage.ScriptInfo)
		*e.scriptInfo = storage.ScriptInfo(msg)

		return e, LoadScriptCmd(e.storage, *e.scriptInfo)
	case scriptLoadedMsg:

		log.Printf("Script loaded: %T %+v", msg, msg)
		e.scriptEditor.Reset(langs.Script(msg))

		return e, tea.Batch(focusChangeCmd(focusEditor), e.evaluateScriptCmd())
	case tea.KeyMsg:
		if !(e.focusElement == focusEditor || e.focusElement == focusFilename) {
			return e, nil
		}
		switch msg.String() {
		case "ctrl+s":
			cmd := e.trySave()
			return e, cmd
		}
	}

	switch e.focusElement {
	case focusEditor:
		e.scriptEditor, cmd = e.scriptEditor.Update(msg)
		cmds = append(cmds, cmd, e.evaluateScriptCmd())
	case focusFilename:
		e.filename, cmd = e.filename.Update(msg)
		cmds = append(cmds, cmd)
	}

	return e, tea.Batch(cmds...)
}

func (e Editor) View() string {
	return fmt.Sprintf("Editor\n%s\n%s", e.filename.View(), e.scriptEditor.View(scriptEditorWidth))
}

type evaluateScriptMsg []langs.CellResult
func (e Editor) evaluateScriptCmd() tea.Cmd {
	interpreter := langs.CreateLuaScriptInterpreter()
	script := langs.LuaScript(e.scriptEditor.GetScriptString())

	return func() tea.Msg {
		results, err := interpreter.Run(script)
		if err != nil {
			log.Printf("error evaluating script: %s", err.Error())
		}

		return evaluateScriptMsg(results)
	}
}
