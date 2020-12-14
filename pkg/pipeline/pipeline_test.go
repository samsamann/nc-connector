package pipeline

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

/*type dummyManipulator struct {
}

func (m dummyManipulator) Manipulate(data FileData) {
	valOf := reflect.ValueOf(data)
	if valOf.CanSet() {
		s := struct {
			A string
		}{
			"Test",
		}

		unsafe.Pointer()
		valOf.SetPointer(reflect.ValueOf(s).Pointer())
	}
}

func TestDummyManipulator(t *testing.T) {
	manipulator := new(dummyManipulator)
	manipulator.Manipulate(nil)
}*/

func TestReporter(t *testing.T) {
	var wg sync.WaitGroup
	ctx, cancel := newFilePipelineContext(context.Background(), nil, 5*time.Second)

	wg.Add(1)
	reporter(ctx, &wg, 6*time.Second)
	cancel()
	wg.Wait()

	ctx, cancel = newFilePipelineContext(context.Background(), nil, 1*time.Second)
	reporter(ctx, &wg, 6*time.Second)
	wg.Wait()
}

type FileExporterMock struct {
	mock.Mock
}

func (f *FileExporterMock) Export(ctx ExportContext, fileChan <-chan FileData) error {
	args := f.Called(ctx, fileChan)

	retVal := args.Get(0)
	if retVal != nil {
		return retVal.(error)
	}
	return nil
}

func TestExportData(t *testing.T) {

}
