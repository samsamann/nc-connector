package operator

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/samsamann/nc-connector/internal/config"
	"github.com/samsamann/nc-connector/internal/stream"
)

const apiOperatorName = "api"

func initAPIOperator(globalConfig *config.GlobalConfig, opConfig map[string]interface{}) (stream.Operator, error) {
	/*c, err := processConfig(util.NewConfigMap(config))
	if err != nil {
		return nil, err
	}*/
	return newAPIOperator(), nil
}

type apiOperator struct {
	channel chan stream.SyncItem
}

func newAPIOperator() stream.Operator {
	return &apiOperator{
		channel: make(chan stream.SyncItem),
	}
}

func (ms apiOperator) In() chan<- stream.SyncItem {
	return ms.channel
}

func (ms apiOperator) Out() <-chan stream.SyncItem {
	channel := make(chan stream.SyncItem)
	go func() {
		client := &http.Client{}
		defer func() {
			close(channel)
			client.CloseIdleConnections()
		}()
		for file := range ms.channel {
			req := prepareRequest(file)
			res, err := client.Do(req)
			if err == nil && res.StatusCode >= 200 && res.StatusCode < 300 {
				dir, fileName := path.Split(file.Path())
				fileName = strings.TrimRight(fileName, path.Ext(fileName)) + ".pdf"

				file.SetPath(path.Join(dir, fileName))
				file.SetData(res.Body)
			}
			channel <- file
		}
	}()
	return channel
}

func prepareRequest(file stream.SyncItem) *http.Request {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("files", filepath.Base(file.Path()))
	io.Copy(part, file.Data())
	writer.Close()

	r, _ := http.NewRequest("POST", "http://localhost:3000/forms/libreoffice/convert", body)
	r.Header.Add("Content-Type", writer.FormDataContentType())
	return r
}
