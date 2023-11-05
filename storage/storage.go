package storage

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"time"
)

type Storage struct {
	storeDir string
}

func LoadStorage() Storage {
	storeDir := getConfigDirectory()

	return Storage{
		storeDir: storeDir,
	}
}

func (s Storage) GetScriptInfos() []ScriptInfo {
	scriptDir := s.getScriptDir()
	dirEntries, err := os.ReadDir(scriptDir)
	if err != nil {
		log.Fatalf("Unable to ReadDir in script path: %v", err)
	}

	infos := []ScriptInfo{}
	for _, dirEntry := range dirEntries {
		info, err := dirEntry.Info()
		if err != nil {
			log.Printf("Error getting file info for %v", dirEntry)
		}
		infos = append(infos, ScriptInfo{
			name: dirEntry.Name(),
			fileInfo: info,
		})
	}

	return infos
}

func (s Storage) getScriptDir() string {
	return path.Join(s.storeDir, "scripts")
}

func (s Storage) getScriptPath(name string) string {
	return path.Join(s.storeDir, "scripts", name)
}

func (s Storage) LoadScript(info ScriptInfo) (string, error) {
	path := s.getScriptPath(info.Name())
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (s Storage) SaveScript(name string, content string) error {
	scriptpath := s.getScriptPath(name)

	err := os.WriteFile(scriptpath, []byte(content), 0666)
	if err != nil {
		log.Printf("Failed to write %s", scriptpath)
	}

	return nil
}

func (s Storage) SaveNewScript(content string) error {
	now := time.Now()
	scriptname := now.Format("2006-01-02-15-04-05")
	log.Printf("Saving new script with name: %s", scriptname)
	return s.SaveScript(scriptname, content)
}
 
func (s Storage) GetLogPath() string {
	return path.Join(s.storeDir, "debug.log")
}

func getConfigDirectory() string {
	switch runtime.GOOS {
	case "linux":
		return fmt.Sprintf("%s/.config/bc2", os.Getenv("HOME"))
	}

	log.Fatalf("Config directory not implemented for os '%s'", runtime.GOOS)
	return "" // note: unreachable
}
