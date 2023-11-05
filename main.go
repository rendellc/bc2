package main

import (
	"fmt"
	"log"
	"os"

	"rendellc/bc2/storage"
	"rendellc/bc2/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	storage := storage.LoadStorage()
	logpath := storage.GetLogPath()
	f, err := tea.LogToFile(logpath, "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	// script, err := storage.LoadScript("test2")
	// // script, err := store.LoadScript(scriptNames[0])
	// if err != nil {
	// 	log.Fatalf("Failed to load script: %s", script)
	// }
	// scriptLines := langs.SplitLines(script)

	// luaInterpreter := langs.CreateLuaScriptInterpreter()
	// defer luaInterpreter.Close()

	// var interpreter langs.ScriptInterpreter = luaInterpreter

	// results := interpreter.Run(scriptLines)
	// for i := range results {
	// 	if results[i].IsOK() {
	// 		fmt.Printf("%v\t%v\n", scriptLines[i], results[i].Ok())
	// 	} else {
	// 		fmt.Printf("%v\t<%v>\n", scriptLines[i], results[i].Err())
	// 	}
	// }

	app := ui.CreateApp(&storage)
	if _, err := tea.NewProgram(app).Run(); err != nil {
		log.Fatal(err)
	}

}
