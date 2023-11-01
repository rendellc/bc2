package main

import (
	"bufio"
	"fmt"
	"log"
	"strings"

	"rendellc/bc2/langs"
	"rendellc/bc2/storage"
)

func splitLines(script string) []string {
	lines := []string{}
	scanner := bufio.NewScanner(strings.NewReader(script))

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}

func main() {
	store := storage.LoadStorage()
	// scriptNames := store.GetScriptNames()

	script, err := store.LoadScript("test2")
	// script, err := store.LoadScript(scriptNames[0])
	if err != nil {
		log.Fatalf("Failed to load script: %s", script)
	}
	scriptLines := splitLines(script)

	luaInterpreter := langs.CreateLuaScriptInterpreter()
	defer luaInterpreter.Close()
	var interpreter langs.ScriptInterpreter = luaInterpreter

	results := interpreter.Run(scriptLines)
	for i := range results {
		fmt.Printf("%v\t%v\n", scriptLines[i], results[i].Ok(), results[i].Err())
	}

}
