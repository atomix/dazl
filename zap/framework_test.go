// SPDX-FileCopyrightText: 2023-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package zap

import (
	"bytes"
	"github.com/atomix/dazl"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJSONEncoder(t *testing.T) {
	framework := &Framework{}
	encoder := framework.JSONEncoder()
	buf := &bytes.Buffer{}
	writer, err := encoder.NewWriter(buf)
	assert.NoError(t, err)
	writer.Debug("Hello world!")
	assert.Equal(t, "{\"message\":\"Hello world!\"}\n", buf.String())
	buf.Reset()

	encoder, err = encoder.(dazl.MessageKeyEncoder).WithMessageKey("msg")
	assert.NoError(t, err)
	writer, err = encoder.NewWriter(buf)
	assert.NoError(t, err)
	writer.Info("Hello world!")
	assert.Equal(t, "{\"msg\":\"Hello world!\"}\n", buf.String())
	buf.Reset()

	encoder, err = encoder.(dazl.NameEncoder).WithNameEnabled()
	assert.NoError(t, err)
	writer, err = encoder.NewWriter(buf)
	assert.NoError(t, err)
	writer.Warn("Hello world!")
	assert.Equal(t, "{\"msg\":\"Hello world!\"}\n", buf.String())
	buf.Reset()

	writer = writer.WithName("test")
	writer.Error("Hello world!")
	assert.Equal(t, "{\"logger\":\"test\",\"msg\":\"Hello world!\"}\n", buf.String())
	buf.Reset()

	encoder, err = encoder.(dazl.NameKeyEncoder).WithNameKey("name")
	assert.NoError(t, err)
	writer, err = encoder.NewWriter(buf)
	assert.NoError(t, err)
	writer = writer.WithName("test")
	writer.Info("Hello world!")
	assert.Equal(t, "{\"name\":\"test\",\"msg\":\"Hello world!\"}\n", buf.String())
	buf.Reset()

	encoder, err = encoder.(dazl.LevelEncoder).WithLevelEnabled()
	assert.NoError(t, err)
	writer, err = encoder.NewWriter(buf)
	assert.NoError(t, err)
	writer.Info("Hello world!")
	assert.Equal(t, "{\"level\":\"info\",\"msg\":\"Hello world!\"}\n", buf.String())
	buf.Reset()

	encoder, err = encoder.(dazl.LevelKeyEncoder).WithLevelKey("lvl")
	assert.NoError(t, err)
	writer, err = encoder.NewWriter(buf)
	assert.NoError(t, err)
	writer.Warn("Hello world!")
	assert.Equal(t, "{\"lvl\":\"warn\",\"msg\":\"Hello world!\"}\n", buf.String())
	buf.Reset()

	encoder, err = encoder.(dazl.LevelFormattingEncoder).WithLevelFormat(dazl.UpperCaseLevelFormat)
	assert.NoError(t, err)
	writer, err = encoder.NewWriter(buf)
	assert.NoError(t, err)
	writer.Error("Hello world!")
	assert.Equal(t, "{\"lvl\":\"ERROR\",\"msg\":\"Hello world!\"}\n", buf.String())
	buf.Reset()

	encoder, err = encoder.(dazl.CallerEncoder).WithCallerEnabled()
	assert.NoError(t, err)
	writer, err = encoder.NewWriter(buf)
	assert.NoError(t, err)
	writer.Info("Hello world!")
	assert.Equal(t, "{\"lvl\":\"INFO\",\"caller\":\"zap/framework_test.go:82\",\"msg\":\"Hello world!\"}\n", buf.String())
	buf.Reset()

	encoder, err = encoder.(dazl.CallerKeyEncoder).WithCallerKey("call")
	assert.NoError(t, err)
	writer, err = encoder.NewWriter(buf)
	assert.NoError(t, err)
	writer.Info("Hello world!")
	assert.Equal(t, "{\"lvl\":\"INFO\",\"call\":\"zap/framework_test.go:90\",\"msg\":\"Hello world!\"}\n", buf.String())
	buf.Reset()
}

func TestConsoleEncoder(t *testing.T) {
	framework := &Framework{}
	encoder := framework.ConsoleEncoder()
	buf := &bytes.Buffer{}
	writer, err := encoder.NewWriter(buf)
	assert.NoError(t, err)
	writer.Debug("Hello world!")
	assert.Equal(t, "Hello world!\n", buf.String())
	buf.Reset()

	encoder, err = encoder.(dazl.NameEncoder).WithNameEnabled()
	assert.NoError(t, err)
	writer, err = encoder.NewWriter(buf)
	assert.NoError(t, err)
	writer.Info("Hello world!")
	assert.Equal(t, "Hello world!\n", buf.String())
	buf.Reset()

	writer = writer.WithName("test")
	writer.Warn("Hello world!")
	assert.Equal(t, "test\tHello world!\n", buf.String())
	buf.Reset()

	encoder, err = encoder.(dazl.LevelEncoder).WithLevelEnabled()
	assert.NoError(t, err)
	writer, err = encoder.NewWriter(buf)
	assert.NoError(t, err)
	writer.Error("Hello world!")
	assert.Equal(t, "error\tHello world!\n", buf.String())
	buf.Reset()

	encoder, err = encoder.(dazl.LevelFormattingEncoder).WithLevelFormat(dazl.UpperCaseLevelFormat)
	assert.NoError(t, err)
	writer, err = encoder.NewWriter(buf)
	assert.NoError(t, err)
	writer.Info("Hello world!")
	assert.Equal(t, "INFO\tHello world!\n", buf.String())
	buf.Reset()

	encoder, err = encoder.(dazl.CallerEncoder).WithCallerEnabled()
	assert.NoError(t, err)
	writer, err = encoder.NewWriter(buf)
	assert.NoError(t, err)
	writer.Warn("Hello world!")
	assert.Equal(t, "WARN\tzap/framework_test.go:138\tHello world!\n", buf.String())
	buf.Reset()
}
