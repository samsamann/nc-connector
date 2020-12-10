package pipeline

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/samsamann/nc-connector/internal/log"
)

type FilePipeline struct {
	source FileImporter
	dest   FileExporter
}

type filePipelineContext struct {
	context.Context

	mux    *sync.Mutex
	logger log.Logger

	totalEntites     uint
	processedEntites uint
}

func (c filePipelineContext) Error(args ...interface{}) {
	c.logger.Error(args...)
}

func (c *filePipelineContext) SetTotalEntities(total uint) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.totalEntites = total
	c.processedEntites = 0
}

func (c *filePipelineContext) UpdateProcess(processed uint) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.processedEntites = processed
}

func newFilePipelineContext(logger log.Logger, timeout time.Duration) (*filePipelineContext, context.CancelFunc) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
	return &filePipelineContext{
			Context: ctx,
			logger:  logger,
		},
		cancelFunc
}

// Run runs
func (p FilePipeline) Run() {
	ctx, cancel := newFilePipelineContext(nil, 15*time.Minute)
	defer func() {
		cancel()
		if r := recover(); r != nil {
			// TODO: Log msg
			fmt.Printf("Recovering from panic. error is: %v \n", r)
		}
	}()
	var wg sync.WaitGroup

	wg.Add(4)
	importChan := importData(ctx, &wg, p.source)
	exportChan := middleware(ctx, &wg, []FileManipulator{}, importChan)
	exportData(ctx, &wg, p.dest, exportChan)
	reporter(ctx, &wg, 5*time.Second)
	wg.Wait()
}

func importData(ctx ImportContext, wg *sync.WaitGroup, importer FileImporter) <-chan FileData {
	channel := make(chan FileData)
	go func() {
		defer func() {
			close(channel)
			wg.Done()
		}()
		done := importer.Import(ctx, channel)
		select {
		case <-ctx.Done():
		case <-done:
		}
	}()
	return channel
}

func middleware(ctx Context, wg *sync.WaitGroup, manipulators []FileManipulator, input <-chan FileData) <-chan FileData {
	channel := make(chan FileData)

	go func() {
		defer func() {
			close(channel)
			wg.Done()
		}()
		for {
			select {
			case <-ctx.Done():
				break
			case data, ok := <-input:
				if !ok {
					break
				}

				for _, manipulator := range manipulators {
					if manipulator.CanManipulate(ctx, data) {
						channel <- manipulator.Manipulate(ctx, data)
					}
				}
			}
		}
	}()

	return channel
}

func exportData(ctx ExportContext, wg *sync.WaitGroup, exporter FileExporter, input <-chan FileData) {
	go func() {
		defer wg.Done()
		done := exporter.Export(ctx, input)
		select {
		case <-ctx.Done():
		case <-done:
		}
	}()
}

func reporter(ctx Context, wg *sync.WaitGroup, interval time.Duration) {
	go func() {
		defer func() {
			wg.Done()
		}()
		ticker := time.NewTicker(interval)
		ok := true
		for ok {
			select {
			case <-ctx.Done():
				ticker.Stop()
				ok = false
			case <-ticker.C:
				// TODO: Log Info
			}
		}
	}()
}
