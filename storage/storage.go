package storage

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
)

type storage struct {
	storeDir string
}

func LoadStorage() storage {
	storeDir := getConfigDirectory()

	return storage{
		storeDir: storeDir,
	}
}

func (s storage) GetScriptNames() []string {
	scriptDir := s.getScriptDir()
	dirEntries, err := os.ReadDir(scriptDir)
	if err != nil {
		log.Fatalf("Unable to ReadDir in script path: %v", err)
	}

	files := []string{}
	for _, dirEntry := range dirEntries {
		files = append(files, dirEntry.Name())
	}
	
	return files
}

func (s storage) getScriptDir() string {
	return path.Join(s.storeDir, "scripts")
}

func (s storage) getScriptPath(name string) string {
	return path.Join(s.storeDir, "scripts", name)
}

func (s storage) LoadScript(name string) (string, error) {
	path := s.getScriptPath(name)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}


func getConfigDirectory() string {
	switch runtime.GOOS {
	case "linux":
		return fmt.Sprintf("%s/.config/bc2", os.Getenv("HOME"))
	}

	log.Fatalf("Config directory not implemented for os '%s'", runtime.GOOS)
	return "" // note: unreachable
}
