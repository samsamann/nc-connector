package config

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestGlobalYamlParsing(t *testing.T) {
	tests := []struct {
		content        string
		expectErrType  interface{}
		expectErrMsg   []string
		expectedConfig *Config
	}{
		{
			content:        "",
			expectedConfig: new(Config)},
		{
			content:       "foo: bar",
			expectErrType: new(yaml.TypeError),
			expectErrMsg: []string{
				"yaml: unmarshal errors:",
				"line 1: field foo not found in type config.Config",
			},
		},
		{
			content:       "global: 123",
			expectErrType: new(yaml.TypeError),
			expectErrMsg: []string{
				"yaml: unmarshal errors:",
				"line 1: cannot unmarshal !!int `123` into config.GlobalConfig",
			},
		},
		{
			content:       "global: 123\nnext: bar",
			expectErrType: new(yaml.TypeError),
			expectErrMsg: []string{
				"yaml: unmarshal errors:",
				"line 1: cannot unmarshal !!int `123` into config.GlobalConfig",
				"line 2: field next not found in type config.Config",
			},
		},
	}

	for _, test := range tests {
		config, err := Load(test.content)
		if err != nil {
			assert.IsType(t, test.expectErrType, err)
			assert.EqualError(t, err, strings.Join(test.expectErrMsg, "\n  "))
			continue
		}
		assert.Equal(t, test.expectedConfig, config)
	}
}

func TestPipelineConfigParsing(t *testing.T) {
	baseContent := `
pipeline:
  producer:
    test:
      conf: 123
      test: foo
`
	config, err := Load(baseContent)
	assert.NoError(t, err)
	assert.NotEmpty(t, config)
}
