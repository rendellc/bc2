package storage

import "io/fs"

type ScriptInfo struct {
	name     string
	fileInfo fs.FileInfo
}

func (s ScriptInfo) Name() string { return s.name }

func (s ScriptInfo) FileInfo() fs.FileInfo { return s.fileInfo }
