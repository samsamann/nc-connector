package operator

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strconv"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/samsamann/nc-connector/internal/config"
	"github.com/samsamann/nc-connector/internal/stream"
	"github.com/samsamann/nc-connector/internal/stream/util"
)

const extractPageOperatorrName = "extract-pdf-page"

const (
	extractIndicatorCName = "indicator"
	extractSpanCName      = "span"
)

func initPageSplitOperator(global *config.GlobalConfig, config map[string]interface{}) (stream.Operator, error) {
	cMap := util.NewConfigMap(config)
	indicator := cMap.Get(extractIndicatorCName).String()
	span := cMap.Get(extractSpanCName).String()
	if err := cMap.Error(); err != nil {
		return nil, err
	}
	return newExtractPageOperator(indicator, span), nil
}

type extractPageOperator struct {
	channel   chan stream.SyncItem
	indicator string
	span      string
}

func newExtractPageOperator(indicator, span string) *extractPageOperator {
	return &extractPageOperator{
		channel:   make(chan stream.SyncItem),
		indicator: indicator,
		span:      span,
	}
}

func (ps extractPageOperator) In(ctx stream.Context) chan<- stream.SyncItem {
	return ps.channel
}

func (ps extractPageOperator) Out(ctx stream.Context) <-chan stream.SyncItem {
	channel := make(chan stream.SyncItem)
	go func() {
		defer close(channel)
		for file := range ps.channel {
			if ps.indicator != "" {
				extractPageWithIndicator(channel, file, ps.indicator)
			}
		}
	}()
	return channel
}

func extractPageWithIndicator(channel chan<- stream.SyncItem, file stream.SyncItem, indicator string) {
	data, _ := ioutil.ReadAll(file.Data())
	ctx, _ := pdfcpu.Read(bytes.NewReader(data), nil)
	ctx.EnsurePageCount()
	pages := ctx.PageCount
	regex := regexp.MustCompile(fmt.Sprintf("%s\\s+(\\d+)", indicator))
	pageDiff := 0
	for p := 1; p <= pages; p++ {
		reader, _ := ctx.ExtractPageContent(p)
		page, _ := ioutil.ReadAll(reader)
		results := regex.FindStringSubmatch(string(page))
		if len(results) == 2 {
			i, _ := strconv.Atoi(results[1])
			if i > pageDiff {
				if p == pages {
					item := copySyncItem(file)
					reader := pdfcpuSplit(ctx, p-pageDiff, p)
					item.SetData(reader)
					channel <- item
				}
				pageDiff += 1
			} else if pageDiff > 0 {
				item := copySyncItem(file)
				reader := pdfcpuSplit(ctx, p-pageDiff, p-1)
				item.SetData(reader)
				pageDiff = 1
				channel <- item
			} else {
				item := copySyncItem(file)
				reader := pdfcpuSplit(ctx, p, p)
				item.SetData(reader)
				channel <- item
			}
		} else {
			if pageDiff > 0 {
				item := copySyncItem(file)
				reader := pdfcpuSplit(ctx, p-pageDiff, p-1)
				item.SetData(reader)
				pageDiff = 0
				channel <- item
			}
			pageDiff = 0
			item := copySyncItem(file)
			reader := pdfcpuSplit(ctx, p, p)
			item.SetData(reader)
			channel <- item
		}
	}
}

func copySyncItem(item stream.SyncItem) stream.SyncItem {
	switch t := item.(type) {
	case *stream.File:
		concreteFile := *t
		return &concreteFile
	}
	return item
}

func pdfcpuSplit(ctx *pdfcpu.Context, from, thru int) io.ReadCloser {
	selectedPages := pagesForPageRange(from, thru)
	ctxNew, err := ctx.ExtractPages(selectedPages, false)
	if err != nil {
	}
	buffer := bytes.Buffer{}
	ctxNew.Write.Writer = bufio.NewWriter(&buffer)
	defer ctxNew.Write.Flush()
	err = pdfcpu.Write(ctxNew)
	if err != nil {
	}
	return io.NopCloser(&buffer)
}

func pagesForPageRange(from, thru int) []int {
	s := make([]int, thru-from+1)
	for i := 0; i < len(s); i++ {
		s[i] = from + i
	}
	return s
}
