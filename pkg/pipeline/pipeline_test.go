package pipeline

import (
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
	ctx, cancel := newFilePipelineContext(nil, 5*time.Second)

	wg.Add(1)
	reporter(ctx, &wg, 6*time.Second)
	cancel()
	wg.Wait()

	ctx, cancel = newFilePipelineContext(nil, 1*time.Second)
	reporter(ctx, &wg, 6*time.Second)
	wg.Wait()
}

type FileExporterMock struct {
	mock.Mock
}

func (f *FileExporterMock) Export(ctx ExportContext, fileChan <-chan FileData) Done {
	args := f.Called(ctx, fileChan)

	return args.Get(0).(Done)
}

func TestExportData(t *testing.T) {

}
