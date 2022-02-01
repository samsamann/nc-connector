package main

import (
	"fmt"
	"os"

	"github.com/samsamann/nc-connector/internal/config"
	"github.com/samsamann/nc-connector/internal/stream"
	"github.com/samsamann/nc-connector/internal/stream/consumer"
	"github.com/samsamann/nc-connector/internal/stream/operator"
	"github.com/samsamann/nc-connector/internal/stream/producer"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	configPath, err := parseCliArgs(os.Args[1:])
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	c, err := config.LoadFile(configPath)
	if err != nil {
		fmt.Printf("can not process config file: %s", err)
	}
	buildPipeline(c)
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

func buildPipeline(config *config.Config) {
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
	flow.To(c).Start()
}
