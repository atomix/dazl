// SPDX-FileCopyrightText: 2023-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package dazl

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestLoggerNames(t *testing.T) {
	assert.Equal(t, "", root.Name())
	assert.Equal(t, "foo", GetLogger("foo").Name())
	assert.Equal(t, "foo/bar", GetLogger("foo/bar").Name())
	assert.Equal(t, "github.com/atomix/dazl", GetPackageLogger().Name())
}

func TestSetLevel(t *testing.T) {
	parent := GetLogger("parent")
	child := parent.GetLogger("child")

	assert.Equal(t, EmptyLevel, root.Level())
	assert.Equal(t, EmptyLevel, parent.Level())
	assert.Equal(t, EmptyLevel, child.Level())

	SetLevel(ErrorLevel)
	assert.Equal(t, ErrorLevel, root.Level())
	assert.Equal(t, ErrorLevel, parent.Level())
	assert.Equal(t, ErrorLevel, child.Level())

	parent.SetLevel(WarnLevel)
	assert.Equal(t, ErrorLevel, root.Level())
	assert.Equal(t, WarnLevel, parent.Level())
	assert.Equal(t, WarnLevel, child.Level())

	child.SetLevel(ErrorLevel)
	assert.Equal(t, ErrorLevel, root.Level())
	assert.Equal(t, WarnLevel, parent.Level())
	assert.Equal(t, ErrorLevel, child.Level())
}

const testLoggerConfig = `
level: debug
sample:
  random:
    interval: 10
    level: debug
outputs:
  - stdout:
      level: info
  - stderr:
      level: error
  - file
`

func TestLoggerConfig(t *testing.T) {
	var config loggerConfig
	assert.NoError(t, yaml.Unmarshal([]byte(testLoggerConfig), &config))

	assert.Equal(t, DebugLevel, config.Level.Level())
	assert.Len(t, config.Outputs, 3)
	assert.Equal(t, "stdout", config.Outputs[0].Writer)
	assert.Equal(t, InfoLevel, config.Outputs[0].Level.Level())
	assert.Equal(t, "stderr", config.Outputs[1].Writer)
	assert.Equal(t, ErrorLevel, config.Outputs[1].Level.Level())
	assert.Equal(t, "file", config.Outputs[2].Writer)
	assert.Equal(t, EmptyLevel, config.Outputs[2].Level.Level())
}

const testLevelConfig = `
encoders:
  console:
    fields:
      - message
      - level:
          format: uppercase
      - time:
          format: iso8601
      - caller:
          format: short

writers:
  stdout:
    encoder: console

rootLogger:
  level: debug
  outputs:
    - stdout
`

func TestLoggerLevels(t *testing.T) {
	var config loggingConfig
	assert.NoError(t, yaml.Unmarshal([]byte(testLevelConfig), &config))
	assert.NoError(t, configure(config))
}
