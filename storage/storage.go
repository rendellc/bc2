package storage

import (
	"net/url"
)


type Store interface {
	GetConfigParams() []string
	GetConfigParam(paramName string) string
	GetScriptInfos() ([]StoreMeta, error)
	GetScript(scriptMeta StoreMeta) (string, error)
	SaveScript(scriptMeta StoreMeta, content string) error
}


type StoreMeta struct {
	name string
	url  url.URL
}

func (s StoreMeta) Name() string { return s.name }
func (s StoreMeta) URL() url.URL { return s.url }

func CreateNewMetaByName(name string) StoreMeta {
	return StoreMeta{
		name: name,
	}
}


