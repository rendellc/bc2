package main

import (
	"fmt"
	"log"
	"os"

	"rendellc/bc2/storage"
	"rendellc/bc2/bc2cli/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	storage := storage.LoadFilesystemStorage()
	logpath := storage.GetLogPath()
	f, err := tea.LogToFile(logpath, "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	app := ui.CreateApp(storage)
	if _, err := tea.NewProgram(app).Run(); err != nil {
		log.Fatal(err)
	}

}
