package main

import (
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/samsamann/nc-connector/internal/config"
	"github.com/samsamann/nc-connector/internal/stream"
	"github.com/samsamann/nc-connector/internal/stream/consumer"
	"github.com/samsamann/nc-connector/internal/stream/operator"
	"github.com/samsamann/nc-connector/internal/stream/producer"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	configPath, err := parseCliArgs(os.Args[1:])
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	c, err := config.LoadFile(configPath)
	if err != nil || c == nil {
		fmt.Printf("can not process config file: %s", err)
		os.Exit(1)
	}
	logger := createLogger(c.GlobalConfig.Logger)
	buildPipeline(c, logger)
}

func parseCliArgs(args []string) (string, error) {
	app := kingpin.New("nc-connector", "A command-line program to connect Nexcloud with other systems.")
	configPath := app.Flag("config", "path to the config file (yaml)").
		Default("./config.yaml").
		Short('c').
		String()
	_, err := app.Parse(args)
	if err != nil {
		return "", err
	}
	return *configPath, nil
}

func buildPipeline(config *config.Config, logger *logrus.Logger) {
	if config == nil {
		return
	}
	// TODO: check config first
	p, err := producer.CreateProducer(
		config.Pipeline.Producer.Name,
		config.Pipeline.Producer.Config,
	)
	if err != nil {
		return
	}
	flow := stream.NewStream(p)
	for _, m := range config.Pipeline.Middleware {
		o, err := operator.CreateOperator(m, &config.GlobalConfig)
		if err != nil {
			return
		}
		flow = flow.Via(o)
	}
	c, err := consumer.CreateConsumer(
		config.Pipeline.Consumer,
		&config.GlobalConfig,
	)
	if err != nil {
		return
	}
	flow.To(c).Start(logger)
}

func createLogger(logger config.LoggerConfig) *logrus.Logger {
	log := logrus.New()
	log.Out = os.Stdout

	if logger.LogLevel != "" {
		if lvl, err := logrus.ParseLevel(logger.LogLevel); err == nil {
			log.Level = lvl
		} else {
			log.Warn("failed to set log level", err)
		}
	}
	if logger.LogPath != "" {
		t, err := template.New("logfile").Parse(logger.LogPath)
		if err != nil {
			log.Error("failed to parse 'logfile' template, using default stdout")
			return log
		}
		b := strings.Builder{}
		if err = t.Execute(&b, nil); err != nil {
			log.Error("failed to parse 'logfile' template, using default stdout")
			return log
		}

		file, err := os.OpenFile(b.String(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			log.Out = file
		} else {
			log.Error("failed to open log file, using default stdout")
		}
	}
	return log
}
