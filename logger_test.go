// SPDX-FileCopyrightText: 2023-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package dazl

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	io "io"
	"testing"
)

func TestLoggerNames(t *testing.T) {
	assert.Equal(t, "", root.Name())
	assert.Equal(t, "foo", GetLogger("foo").Name())
	assert.Equal(t, "foo/bar", GetLogger("foo/bar").Name())
	assert.Equal(t, "github.com/atomix/dazl", GetPackageLogger().Name())
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

const testConfig = `
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
  file:
    encoder: json
    path: test.log

rootLogger:
  level: info
  outputs:
    - stdout

loggers:
  test/level:
    level: warn
  test/sample:
    sample:
      basic: 
        interval: 2
        minLevel: warn
  test/sample/outputs:
    outputs:
      - file
  test/outputs:
    outputs:
      - file
  test/outputs/level:
    outputs:
      - stdout:
          level: warn
  test/outputs/sample:
    outputs:
      - stdout:
          sample:
            basic:
              interval: 2
              minLevel: warn
`

func TestLogger(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	Register(&testFramework{ctrl})

	var config loggingConfig
	assert.NoError(t, yaml.Unmarshal([]byte(testConfig), &config))
	assert.NoError(t, configure(config))

	var log = root
	writer, _ := log.(*dazlLogger).writers.Load("stdout")
	stdout := writer.(*MockWriter)

	stdout.EXPECT().WithName(gomock.Eq("test")).Return(stdout)
	log = GetLogger("test")

	log.Debug("debug")
	stdout.EXPECT().Info(gomock.Eq("info"))
	log.Info("info")
	stdout.EXPECT().Warn(gomock.Eq("warn"))
	log.Warn("warn")
	stdout.EXPECT().Error(gomock.Eq("error"))
	log.Error("error")

	stdout.EXPECT().WithName(gomock.Eq("test/level")).Return(stdout)
	log = GetLogger("test/level")

	log.Debug("debug")
	log.Info("info")
	stdout.EXPECT().Warn(gomock.Eq("warn"))
	log.Warn("warn")
	stdout.EXPECT().Error(gomock.Eq("error"))
	log.Error("error")

	stdout.EXPECT().WithName(gomock.Eq("test/sample")).Return(stdout)
	log = GetLogger("test/sample")

	stdout.EXPECT().Error(gomock.Eq("error"))
	log.Error("error")
	stdout.EXPECT().Error(gomock.Eq("error"))
	log.Error("error")
	stdout.EXPECT().Warn(gomock.Eq("warn"))
	log.Warn("warn")
	log.Warn("warn")
	stdout.EXPECT().Warn(gomock.Eq("warn"))
	log.Warn("warn")
	log.Warn("warn")

	stdout.EXPECT().WithName(gomock.Eq("test/sample/outputs")).Return(stdout)
	log = GetLogger("test/sample/outputs")

	writer, _ = log.(*dazlLogger).writers.Load("file")
	file := writer.(*MockWriter)

	stdout.EXPECT().Error(gomock.Eq("error"))
	file.EXPECT().Error(gomock.Eq("error"))
	log.Error("error")
	stdout.EXPECT().Error(gomock.Eq("error"))
	file.EXPECT().Error(gomock.Eq("error"))
	log.Error("error")
	stdout.EXPECT().Warn(gomock.Eq("warn"))
	file.EXPECT().Warn(gomock.Eq("warn"))
	log.Warn("warn")
	log.Warn("warn")
	stdout.EXPECT().Warn(gomock.Eq("warn"))
	file.EXPECT().Warn(gomock.Eq("warn"))
	log.Warn("warn")
	log.Warn("warn")

	stdout.EXPECT().WithName(gomock.Eq("test/outputs")).Return(stdout)
	file.EXPECT().WithName(gomock.Eq("test/outputs")).Return(file)
	log = GetLogger("test/outputs")

	log.Debug("debug")
	stdout.EXPECT().Info(gomock.Eq("info"))
	file.EXPECT().Info(gomock.Eq("info"))
	log.Info("info")
	stdout.EXPECT().Warn(gomock.Eq("warn"))
	file.EXPECT().Warn(gomock.Eq("warn"))
	log.Warn("warn")
	stdout.EXPECT().Error(gomock.Eq("error"))
	file.EXPECT().Error(gomock.Eq("error"))
	log.Error("error")

	stdout.EXPECT().WithName(gomock.Eq("test/outputs/level")).Return(stdout)
	file.EXPECT().WithName(gomock.Eq("test/outputs/level")).Return(file)
	log = GetLogger("test/outputs/level")

	log.Debug("debug")
	file.EXPECT().Info(gomock.Eq("info"))
	log.Info("info")
	stdout.EXPECT().Warn(gomock.Eq("warn"))
	file.EXPECT().Warn(gomock.Eq("warn"))
	log.Warn("warn")
	stdout.EXPECT().Error(gomock.Eq("error"))
	file.EXPECT().Error(gomock.Eq("error"))
	log.Error("error")

	stdout.EXPECT().WithName(gomock.Eq("test/outputs/sample")).Return(stdout)
	file.EXPECT().WithName(gomock.Eq("test/outputs/sample")).Return(file)
	log = GetLogger("test/outputs/sample")

	stdout.EXPECT().Error(gomock.Eq("error"))
	file.EXPECT().Error(gomock.Eq("error"))
	log.Error("error")
	stdout.EXPECT().Error(gomock.Eq("error"))
	file.EXPECT().Error(gomock.Eq("error"))
	log.Error("error")
	stdout.EXPECT().Warn(gomock.Eq("warn"))
	file.EXPECT().Warn(gomock.Eq("warn"))
	log.Warn("warn")
	file.EXPECT().Warn(gomock.Eq("warn"))
	log.Warn("warn")
	stdout.EXPECT().Warn(gomock.Eq("warn"))
	file.EXPECT().Warn(gomock.Eq("warn"))
	log.Warn("warn")
	file.EXPECT().Warn(gomock.Eq("warn"))
	log.Warn("warn")
}

type testFramework struct {
	ctrl *gomock.Controller
}

func (f *testFramework) NewWriter(io.Writer, Encoding) (Writer, error) {
	writer := NewMockWriter(f.ctrl)
	writer.EXPECT().WithName(gomock.Any()).DoAndReturn(func(name string) Writer {
		return writer
	})
	return writer, nil
}
