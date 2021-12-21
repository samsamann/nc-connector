package config

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// LoadFile parses the given YAMl file.
func LoadFile(filename string) (*Config, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return Load(string(content))

}

// Load parses the given string s into a Config.
func Load(s string) (*Config, error) {
	cfg := &Config{}
	err := yaml.UnmarshalStrict([]byte(s), cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// Config is the top-level configuration.
type Config struct {
	GlobalConfig GlobalConfig   `yaml:"global"`
	Pipeline     PipelineConfig `yaml:"pipeline"`
}

// GlobalConfig defines global variables that are used everywhere.
type GlobalConfig struct {
	NCClient NCClientConfig `yaml:"ncClient"`
}

type NCClientConfig struct {
	Host     string `yaml:"host"`
	BasePath string `yaml:"basePath"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type PipelineConfig struct {
	Producer   StreamElem   `yaml:"producer"`
	Middleware []StreamElem `yaml:"middleware"`
	Consumer   StreamElem   `yaml:"consumer"`
}

type StreamElem struct {
	Name   string           `yaml:"name"`
	Config StreamElemConfig `yaml:",inline"`
}

type StreamElemConfig map[string]interface{}
