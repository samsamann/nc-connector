package operator

import (
	"path"
	"strings"
	"text/template"

	"github.com/samsamann/nc-connector/internal/config"
	"github.com/samsamann/nc-connector/internal/stream"
	"github.com/samsamann/nc-connector/internal/stream/util"
)

const renameOperatorName = "rename"

const (
	renameNewFileNameCName = "newFileName"

	renameFileExtKeyName = "file_ext"
)

func initRenameOperator(gConfig *config.GlobalConfig, config map[string]interface{}) (stream.Operator, error) {
	cMap := util.NewConfigMap(config)
	fileNameTemplate := cMap.Get(renameNewFileNameCName).Required().String()

	if err := cMap.Error(); err != nil {
		return nil, err
	}
	t, err := template.New("filename").Parse(fileNameTemplate)
	if err != nil {
		return nil, err
	}

	return newRenameOperator(t), nil
}

type renameOperator struct {
	channel  chan stream.SyncItem
	template *template.Template
}

func newRenameOperator(t *template.Template) *renameOperator {
	return &renameOperator{
		channel:  make(chan stream.SyncItem),
		template: t,
	}
}

func (ms renameOperator) In(ctx stream.Context) chan<- stream.SyncItem {
	return ms.channel
}

func (ms renameOperator) Out(ctx stream.Context) <-chan stream.SyncItem {
	channel := make(chan stream.SyncItem)
	go func() {
		defer close(channel)
		for file := range ms.channel {
			renameFile(file, ms.template)
			channel <- file
		}
	}()
	return channel
}

func renameFile(file stream.SyncItem, t *template.Template) {
	attrs := make(map[string]interface{}, len(file.Attributes()))
	for k, v := range file.Attributes() {
		attrs[k] = v
	}
	attrs[renameFileExtKeyName] = path.Ext(file.Path())

	builder := strings.Builder{}
	if err := t.Execute(&builder, attrs); err != nil {
		return
	}
	newFileName := strings.ReplaceAll(builder.String(), "/", "")
	// prevent duplicated file extension
	newFileExt := path.Ext(newFileName)
	lastPos := strings.LastIndex(newFileName, newFileExt)
	if lastPos != -1 {
		newFileName = strings.ReplaceAll(
			newFileName[:lastPos],
			newFileExt,
			"",
		) + newFileExt
	}

	dir, _ := path.Split(file.Path())
	file.SetPath(
		path.Join(dir, "/", newFileName),
	)
}
