package storage

import (
	"log"
	"os"
	"path"
)

type Storage struct {
	storeDir string
}

func LoadStorage() Storage {
	storeDir := getConfigDirectory()

	log.Printf("Checking if %v exists", storeDir)
	err := os.MkdirAll(storeDir, os.ModeDir)
	if err != nil {
		log.Fatalf("Unable to make config directory: %v", err.Error())
	}
	err = os.MkdirAll(path.Join(storeDir,"scripts"), os.ModeDir)
	if err != nil {
		log.Fatalf("Unable to make script directory: %v", err.Error())
	}

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

func (s Storage) GetLogPath() string {
	return path.Join(s.storeDir, "debug.log")
}

func getConfigDirectory() string {
	confdir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("Unable to determine config dir: %s", err.Error())
	}

	return path.Join(confdir, "bc2")
}
