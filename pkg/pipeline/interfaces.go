package pipeline

import (
	"context"

	"github.com/samsamann/nc-connector/internal/log"
)

// Done is a read-only channel that indicates whether a pipeline element has done the job.
type Done <-chan struct{}

// FileImporter is the interface that wraps the Import method.
type FileImporter interface {
	Import(ImportContext, chan<- FileData) Done
}

// FileManipulator is the interface that encapsulate the entire middleware processing.
type FileManipulator interface {
	CanManipulate(Context, FileData) bool
	Manipulate(Context, FileData) FileData
}

// FileExporter is the interface that wraps the Export method.
type FileExporter interface {
	Export(ExportContext, <-chan FileData) Done
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
}

type ImportContext interface {
	Context
	SetTotalEntities(uint)
}

type ExportContext interface {
	Context
	UpdateProcess(uint)
}
