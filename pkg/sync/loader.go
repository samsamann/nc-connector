package sync

import (
	"encoding/json"
	"errors"
	"os"
	"syscall"
)

type Loader interface {
	reader
	writer
}

type reader interface {
	Load() (SearchableStorage, error)
}

type writer interface {
	Unload(SearchableStorage) error
}

type jsonFileLoader struct {
	localPath string
}

func NewJsonFileLoader(localPath string) Loader {
	return &jsonFileLoader{localPath: localPath}
}

func (j jsonFileLoader) Load() (SearchableStorage, error) {
	cachedDir := newDirEntry(nil, "/")
	data, err := os.ReadFile(j.localPath)
	if err != nil {
		if errors.Unwrap(err) == syscall.ERROR_FILE_NOT_FOUND {
			return cachedDir, nil
		}
		return nil, err
	}
	flattenedCache := make(map[string]*fileEntry)
	if err := json.Unmarshal(data, &flattenedCache); err != nil {
		return nil, err
	}
	for path, file := range flattenedCache {
		file.removable = true
		extendStorageTree(cachedDir, splitPath(path), file)
	}
	return cachedDir, nil
}

func (j jsonFileLoader) Unload(s SearchableStorage) error {
	if dir, ok := s.(*dirEntry); ok {
		m := make(map[string]Item)
		flatRecursive(m, dir)
		data, err := json.Marshal(m)
		if err != nil {
			return err
		}
		return os.WriteFile(j.localPath, data, 0644)
	}
	return errors.New("can not unload storage")
}
