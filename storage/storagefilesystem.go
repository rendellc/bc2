package storage

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
)

type FilesystemStorage struct {
	storeDir string
}

func (s FilesystemStorage) GetConfigParams() []string {
	return []string{}
}

func (s FilesystemStorage) GetConfigParam(paramName string) string {
	log.Fatalf("No config parameters for FilesystemStorage")
	return ""
}

func (s FilesystemStorage) GetScript(scriptMeta StoreMeta) (string, error) {
	path := s.getScriptPath(scriptMeta.Name())
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (s FilesystemStorage) SaveScript(scriptMeta StoreMeta, content string) error {
	scriptpath := s.getScriptPath(scriptMeta.Name())
	err := os.WriteFile(scriptpath, []byte(content), 0666)
	if err != nil {
		return fmt.Errorf("Failed to write %s", scriptpath)
	}

	return nil
}

func LoadFilesystemStorage() FilesystemStorage {
	fs := FilesystemStorage{
		storeDir: getConfigDirectory(),
	}

	if err := ensureExists(fs.storeDir); err != nil {
		log.Fatalf(err.Error())
	}
	if err := ensureExists(fs.getScriptDir()); err != nil {
		log.Fatalf(err.Error())
	}
	if err := ensureExists(fs.getLogDir()); err != nil {
		log.Fatalf(err.Error())
	}

	return fs
}

func (s FilesystemStorage) GetScriptInfos() ([]StoreMeta, error) {
	scriptDir := s.getScriptDir()
	dirEntries, err := os.ReadDir(scriptDir)
	if err != nil {
		return []StoreMeta{}, fmt.Errorf("unable to ReadDir in script path: %s", err.Error())
	}

	infos := []StoreMeta{}
	for _, dirEntry := range dirEntries {
		info, err := dirEntry.Info()
		if err != nil {
			return []StoreMeta{}, fmt.Errorf("Error getting file info for %s", dirEntry)
		}

		path := path.Join(scriptDir, info.Name())
		url, err := url.Parse("file://" + path)
		if err != nil {
			return []StoreMeta{}, fmt.Errorf("Error parsing %s to url: ", err.Error())
		}

		infos = append(infos, StoreMeta{
			name: dirEntry.Name(),
			url: *url,
		})
	}

	return infos, nil
}

func (s FilesystemStorage) getScriptDir() string {
	return path.Join(s.storeDir, "scripts")
}

func (s FilesystemStorage) getLogDir() string {
	return path.Join(s.storeDir, "log")
}

func (s FilesystemStorage) getScriptPath(name string) string {
	return path.Join(s.storeDir, "scripts", name)
}

func (s FilesystemStorage) GetLogPath() string {
	return path.Join(s.getLogDir(), "debug.log")
}

func getConfigDirectory() string {
	confdir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("Unable to determine config dir: %s", err.Error())
	}

	return path.Join(confdir, "bc2")
}

func ensureExists(dirPath string) error {
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		return fmt.Errorf("Unable to mkdirall for %s: %v", dirPath, err.Error())
	}

	return nil
}
