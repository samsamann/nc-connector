package pipeline

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/samsamann/nc-connector/internal/config"
	"github.com/samsamann/nc-connector/internal/log"
	"github.com/samsamann/nc-connector/pkg/util"
)

type FilePipeline struct {
	ctx    context.Context
	logger log.Logger

	source FileImporter
	dest   FileExporter
}

func NewFilePipeline(rootCtx context.Context, logger log.Logger, options config.FilePipelineConfig) *FilePipeline {
	source, err := GetFileImporter(options.Import.Name, options.Import.Options)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return &FilePipeline{ctx: rootCtx, logger: logger, source: source}
}

type filePipelineContext struct {
	context.Context

	mux    *sync.Mutex
	logger log.Logger

	totalEntites     uint
	processedEntites uint
}

func (c filePipelineContext) Error(args ...interface{}) {
	c.logger.With("caller", util.GetFuncName()).Error(args...)
}

func (c filePipelineContext) Report() {

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

func newFilePipelineContext(
	ctx context.Context,
	logger log.Logger,
	timeout time.Duration,
) (*filePipelineContext, context.CancelFunc) {
	ctx, cancelFunc := context.WithTimeout(ctx, timeout)
	return &filePipelineContext{
			Context: ctx,
			logger:  logger,
		},
		cancelFunc
}

// Run executes the defined pipeline
func (p FilePipeline) Run() {
	//TODO: Make max. timeout configurable
	ctx, cancel := newFilePipelineContext(p.ctx, p.logger, 15*time.Minute)
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
		err := importer.Connect()
		if err != nil {
			ctx.Error(err)
			return
		}
		ctx.Error(importer.Import(ctx, channel))
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
			data, ok := <-input
			if !ok {
				break
			}

			for _, manipulator := range manipulators {
				if manipulator.CanManipulate(ctx, data) {
					channel <- manipulator.Manipulate(ctx, data)
				}
			}
		}
	}()

	return channel
}

func exportData(ctx ExportContext, wg *sync.WaitGroup, exporter FileExporter, input <-chan FileData) {
	go func() {
		defer wg.Done()
		ctx.Error(exporter.Export(ctx, input))
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
				ctx.Report()
			}
		}
	}()
}
