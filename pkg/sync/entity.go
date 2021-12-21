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
		return fmt.Sprintf("/%s", e.Name)
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
	return d.Entry.removable
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

func (d *dirEntry) remove() []Item {
	removedItems := make([]Item, 0)
	for name, dir := range d.dirs {
		removedItems = append(removedItems, dir.remove()...)
		if dir.isRemovable() {
			removedItems = append(removedItems, newFileEntry(""))
			delete(d.dirs, name)
		}
	}
	for name, file := range d.files {
		if file.isRemovable() {
			removedItems = append(removedItems, file)
			delete(d.files, name)
		}
	}
	d.markAsDurable()

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
