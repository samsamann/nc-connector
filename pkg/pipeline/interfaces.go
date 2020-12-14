package pipeline

import (
	"context"

	"github.com/samsamann/nc-connector/internal/log"
)

// Pipeline is the interface that wraps the concret pipeline struct.
type Pipeline interface {
	Run()
}

// FileImporter is the interface that wraps the Import method.
type FileImporter interface {
	Connect() error
	Import(ImportContext, chan<- FileData) error
}

// FileManipulator is the interface that encapsulate the entire middleware processing.
type FileManipulator interface {
	CanManipulate(Context, FileData) bool
	Manipulate(Context, FileData) FileData
}

// FileExporter is the interface that wraps the Export method.
type FileExporter interface {
	Export(ExportContext, <-chan FileData) error
}

// ManipulatorPrio defines processing priority.
type ManipulatorPrio uint

const (
	// LowPrio is the lowest processing priority.
	LowPrio ManipulatorPrio = 10 ^ (iota + 1)
	// MediumPrio is the default processing priority.
	MediumPrio
	// HighPrio is the highest processing priority.
	HighPrio
)

// Context wraps the context
type Context interface {
	context.Context
	log.ReducedLogger
	Report()
}

type ImportContext interface {
	Context
	SetTotalEntities(uint)
}

type ExportContext interface {
	Context
	UpdateProcess(uint)
}
