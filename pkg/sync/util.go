package sync

import (
	"strings"
	"time"
)

func flatRecursive(m map[string]Item, currentDir *dirEntry) {
	for _, f := range currentDir.files {
		m[f.path()] = f
	}
	for _, sub := range currentDir.dirs {
		flatRecursive(m, sub)
	}
}

func convert(path, etag string, date time.Time) Item {
	pathParts := splitPath(path)
	if i := len(pathParts); i > 0 {
		path = pathParts[i-1]
	}
	item := newFileEntry(path)
	item.ETag = etag
	item.ModifiedDate = date
	return item
}

func splitPath(path string) []string {
	return strings.Split(strings.TrimLeft(path, "/"), "/")
}
