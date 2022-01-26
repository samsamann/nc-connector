package sync

import (
	"fmt"
	"strings"
	"time"
)

type Entry struct {
	parent    *dirEntry
	Name      string `json:"name"`
	removable bool
}

func (e *Entry) path() string {
	if e.parent == nil {
		return ""
	}
	return fmt.Sprintf("%s/%s", e.parent.path(), e.Name)
}

type dirEntry struct {
	*Entry
	dirs  map[string]*dirEntry
	files map[string]Item
}

func newDirEntry(parent *dirEntry, name string) *dirEntry {
	return &dirEntry{
		Entry: &Entry{parent: parent, Name: name},
		dirs:  make(map[string]*dirEntry),
		files: make(map[string]Item),
	}
}

func (d *dirEntry) isRemovable() bool {
	return len(d.dirs) == 0 && len(d.files) == 0
}

func (d *dirEntry) markAsDurable() {
	if len(d.dirs) > 0 && len(d.files) > 0 {
		d.Entry.removable = false
	}
}

func (d *dirEntry) Get(path string) Item {
	pathParts := splitPath(path)
	var cachedFile Item
	if len(pathParts) > 1 {
		if sub, ok := d.dirs[pathParts[0]]; ok {
			cachedFile = sub.Get(strings.Join(pathParts[1:], "/"))
		}
	} else {
		if file, ok := d.files[pathParts[0]]; ok {
			cachedFile = file
		}
	}
	return cachedFile
}

func (d *dirEntry) Add(path string, item Item) {
	extendStorageTree(d, splitPath(path), item)
}

func (d *dirEntry) Delete(path string) {
	pathParts := splitPath(path)
	if len(pathParts) > 1 {
		if sub, ok := d.dirs[pathParts[0]]; ok {
			sub.Delete(strings.Join(pathParts[1:], "/"))
		}
	} else {
		delete(d.files, pathParts[0])
	}
}

func (d *dirEntry) removable() []Entry {
	removedItems := make([]Entry, 0)
	for _, dir := range d.dirs {
		removableSubItems := dir.removable()
		if dir.isRemovable() {
			removedItems = append(removedItems, *dir.Entry)
			delete(d.dirs, dir.Name)
		} else {
			// Add removable sub-items only if current directory contains other
			// non-removable items otherwise remove entire directory
			removedItems = append(removedItems, removableSubItems...)
		}
	}
	for _, file := range d.files {
		if f, ok := file.(*fileEntry); ok && file.isRemovable() {
			removedItems = append(removedItems, *f.Entry)
			delete(d.files, f.Name)
		}
	}

	return removedItems
}

type fileEntry struct {
	*Entry
	ModifiedDate time.Time `json:"modDate"`
	ETag         string    `json:"etag"`
}

func newFileEntry(name string) *fileEntry {
	return &fileEntry{
		Entry: &Entry{Name: name},
	}
}

func (f *fileEntry) modifiedDate() time.Time {
	return f.ModifiedDate
}

func (f *fileEntry) eTag() string {
	return f.ETag
}

func (f *fileEntry) name() string {
	return f.Name
}

func (f *fileEntry) setParent(parent *dirEntry) {
	f.Entry.parent = parent
}

func (f *fileEntry) isRemovable() bool {
	return f.Entry.removable
}

func (f *fileEntry) markAsDurable() {
	f.Entry.removable = false
}
