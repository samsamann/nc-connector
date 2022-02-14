package stream

import (
	ctx "context"
	"sync"

	"github.com/sirupsen/logrus"
)

type Context interface {
	Info(string)
}

type context struct {
	cx     ctx.Context
	mu     *sync.Mutex
	logger *logrus.Logger
}

func newContext(logger *logrus.Logger) Context {
	c := &context{
		cx:     ctx.Background(),
		mu:     &sync.Mutex{},
		logger: logger,
	}
	return c
}

func (c *context) Info(msg string) {
	c.logger.Info(msg)
}
