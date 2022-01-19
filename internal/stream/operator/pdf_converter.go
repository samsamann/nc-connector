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
	"github.com/samsamann/nc-connector/internal/stream/util"
)

const apiOperatorName = "to-pdf-converter-api"

const (
	toPdfApiMethodCName     = "method"
	toPdfApiUrlCName        = "url"
	toPdfApiExtensionsCName = "acceptedExtensions"
)

const httpContentTypeHeaderName = "Content-Type"

func initToPDFAPIOperator(globalConfig *config.GlobalConfig, opConfig map[string]interface{}) (stream.Operator, error) {
	c, err := processConfig(util.NewConfigMap(opConfig))
	if err != nil {
		return nil, err
	}
	return newAPIOperator(c), nil
}

type apiOperator struct {
	config  apiConfig
	channel chan stream.SyncItem
}

func newAPIOperator(c apiConfig) stream.Operator {
	return &apiOperator{
		config:  c,
		channel: make(chan stream.SyncItem),
	}
}

func (ao apiOperator) In() chan<- stream.SyncItem {
	return ao.channel
}

func (ao apiOperator) Out() <-chan stream.SyncItem {
	channel := make(chan stream.SyncItem)
	go func() {
		client := &http.Client{}
		defer func() {
			close(channel)
			client.CloseIdleConnections()
		}()
		for file := range ao.channel {
			dir, fileName := path.Split(file.Path())
			if ao.isAllowedToConvert(path.Ext(fileName)) {
				req := prepareRequest(ao.config.method, ao.config.url, file)
				res, err := client.Do(req)
				if err == nil && res.StatusCode >= 200 && res.StatusCode < 300 {
					fileName = strings.TrimRight(fileName, path.Ext(fileName)) + ".pdf"

					file.SetPath(path.Join(dir, fileName))
					file.SetData(res.Body)
				}
			}
			channel <- file
		}
	}()
	return channel
}

func (ao apiOperator) isAllowedToConvert(fileExtension string) bool {
	for _, e := range ao.config.allowedExts {
		if "."+strings.TrimLeft(e, ".") == fileExtension {
			return true
		}
	}
	return false
}

func prepareRequest(method, url string, file stream.SyncItem) *http.Request {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("files", filepath.Base(file.Path()))
	io.Copy(part, file.Data())
	writer.Close()

	r, _ := http.NewRequest(method, url, body)
	r.Header.Add(httpContentTypeHeaderName, writer.FormDataContentType())
	return r
}

type apiConfig struct {
	method      string
	url         string
	allowedExts []string
}

func processConfig(operatorConfig *util.ConfigMap) (apiConfig, error) {
	method := operatorConfig.Get(toPdfApiMethodCName).Required().String()
	url := operatorConfig.Get(toPdfApiUrlCName).Required().String()
	allowedExts := operatorConfig.Get(toPdfApiExtensionsCName).Slice()
	if err := operatorConfig.Error(); err != nil {
		return apiConfig{}, err
	}
	return apiConfig{
		method:      method,
		url:         url,
		allowedExts: allowedExts,
	}, nil
}
