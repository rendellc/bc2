package main

import (
	"fmt"
	"log"

	"rendellc/bc2/langs"
	"rendellc/bc2/storage"
)

func main() {
	store := storage.LoadStorage()

	script, err := store.LoadScript("test2")
	// script, err := store.LoadScript(scriptNames[0])
	if err != nil {
		log.Fatalf("Failed to load script: %s", script)
	}
	scriptLines := langs.SplitLines(script)

	luaInterpreter := langs.CreateLuaScriptInterpreter()
	defer luaInterpreter.Close()
	var interpreter langs.ScriptInterpreter = luaInterpreter

	results := interpreter.Run(scriptLines)
	for i := range results {
		if results[i].IsOK() {
			fmt.Printf("%v\t%v\n", scriptLines[i], results[i].Ok())
		} else {
			fmt.Printf("%v\t<%v>\n", scriptLines[i], results[i].Err())
		}
	}

}
