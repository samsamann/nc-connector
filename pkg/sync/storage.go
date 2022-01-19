package sync

import "time"

type Item interface {
	eTag() string
	modifiedDate() time.Time
	name() string
	path() string
	setParent(*dirEntry)
	isRemovable() bool
	markAsDurable()
}

type SearchableStorage interface {
	Get(string) Item
	Add(string, Item)
	Delete(string)
	remove() []Item
}

func extendStorageTree(current *dirEntry, pathParts []string, item Item) {
	if len(pathParts) == 1 {
		item.setParent(current)

		if _, ok := current.files[item.name()]; !ok {
		current.files[item.name()] = item
		}
		return
	}
	currentPathSeg := pathParts[0]
	if dir, ok := current.dirs[currentPathSeg]; ok {
		extendStorageTree(dir, pathParts[1:], item)
	} else {
		newDir := newDirEntry(current, currentPathSeg)
		current.dirs[currentPathSeg] = newDir
		extendStorageTree(newDir, pathParts[1:], item)
	}
}
