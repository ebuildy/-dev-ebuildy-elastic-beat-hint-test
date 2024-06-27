package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/elastic/elastic-agent-libs/mapstr"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestBuildConfig(t *testing.T) {
	hints := mapstr.M{
		"logs": mapstr.M{
			"processors": mapstr.M{
				"add_fields": ` {"fields": {"foo": "bar"}}`,
			},
		},
	}

	config := buildConfig(hints)

	b := bytes.Buffer{}
	yamlEncoder := yaml.NewEncoder(&b)
	yamlEncoder.SetIndent(2)
	yamlEncoder.Encode(&config)

	configYAML := strings.TrimSpace(b.String())

	assert.Equal(t, strings.TrimSpace(`
enabled: false
excludeLines: []
includeLines: []
json: {}
processors:
  - add_fields:
      fields:
        foo: bar
`), configYAML)
}
