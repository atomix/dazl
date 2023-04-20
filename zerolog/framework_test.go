// SPDX-FileCopyrightText: 2023-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package zerolog

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
	assert.Equal(t, "{\"level\":\"debug\",\"message\":\"Hello world!\"}\n", buf.String())
	buf.Reset()

	encoder, err = encoder.(dazl.MessageKeyEncoder).WithMessageKey("msg")
	assert.NoError(t, err)
	writer, err = encoder.NewWriter(buf)
	assert.NoError(t, err)
	writer.Info("Hello world!")
	assert.Equal(t, "{\"level\":\"info\",\"msg\":\"Hello world!\"}\n", buf.String())
	buf.Reset()

	encoder, err = encoder.(dazl.NameEncoder).WithNameEnabled()
	assert.NoError(t, err)
	writer, err = encoder.NewWriter(buf)
	assert.NoError(t, err)
	writer.Warn("Hello world!")
	assert.Equal(t, "{\"level\":\"warn\",\"msg\":\"Hello world!\"}\n", buf.String())
	buf.Reset()

	writer = writer.WithName("test")
	writer.Error("Hello world!")
	assert.Equal(t, "{\"level\":\"error\",\"logger\":\"test\",\"msg\":\"Hello world!\"}\n", buf.String())
	buf.Reset()

	encoder, err = encoder.(dazl.NameKeyEncoder).WithNameKey("name")
	assert.NoError(t, err)
	writer, err = encoder.NewWriter(buf)
	assert.NoError(t, err)
	writer = writer.WithName("test")
	writer.Info("Hello world!")
	assert.Equal(t, "{\"level\":\"info\",\"name\":\"test\",\"msg\":\"Hello world!\"}\n", buf.String())
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
	writer.Info("Hello world!")
	assert.Equal(t, "{\"lvl\":\"info\",\"msg\":\"Hello world!\"}\n", buf.String())
	buf.Reset()

	encoder, err = encoder.(dazl.LevelFormattingEncoder).WithLevelFormat(dazl.UpperCaseLevelFormat)
	assert.NoError(t, err)
	writer, err = encoder.NewWriter(buf)
	assert.NoError(t, err)
	writer.Info("Hello world!")
	assert.Equal(t, "{\"lvl\":\"INFO\",\"msg\":\"Hello world!\"}\n", buf.String())
	buf.Reset()

	encoder, err = encoder.(dazl.CallerEncoder).WithCallerEnabled()
	assert.NoError(t, err)
	writer, err = encoder.NewWriter(buf)
	assert.NoError(t, err)
	writer.Warn("Hello world!")
	assert.Equal(t, "{\"lvl\":\"WARN\",\"msg\":\"Hello world!\"}\n", buf.String())
	buf.Reset()

	encoder, err = encoder.(dazl.CallerKeyEncoder).WithCallerKey("call")
	assert.NoError(t, err)
	writer, err = encoder.NewWriter(buf)
	assert.NoError(t, err)
	writer.Info("Hello world!")
	assert.Equal(t, "{\"lvl\":\"INFO\",\"call\":\"framework_test.go:90\",\"msg\":\"Hello world!\"}\n", buf.String())
	buf.Reset()
}

func TestConsoleEncoder(t *testing.T) {
	framework := &Framework{}
	encoder := framework.ConsoleEncoder()
	buf := &bytes.Buffer{}
	writer, err := encoder.NewWriter(buf)
	assert.NoError(t, err)
	writer.Debug("Hello world!")
	assert.Equal(t, "\x1b[33mDBG\x1b[0m Hello world!\n", buf.String())
	buf.Reset()

	encoder, err = encoder.(dazl.NameEncoder).WithNameEnabled()
	assert.NoError(t, err)
	writer, err = encoder.NewWriter(buf)
	assert.NoError(t, err)
	writer.Info("Hello world!")
	assert.Equal(t, "\x1b[32mINF\x1b[0m Hello world!\n", buf.String())
	buf.Reset()

	writer = writer.WithName("test")
	writer.Info("Hello world!")
	assert.Equal(t, "\x1b[32mINF\x1b[0m Hello world! \x1b[36mlogger=\x1b[0mtest\n", buf.String())
	buf.Reset()

	encoder, err = encoder.(dazl.LevelEncoder).WithLevelEnabled()
	assert.NoError(t, err)
	writer, err = encoder.NewWriter(buf)
	assert.NoError(t, err)
	writer.Warn("Hello world!")
	assert.Equal(t, "warn Hello world!\n", buf.String())
	buf.Reset()

	encoder, err = encoder.(dazl.LevelFormattingEncoder).WithLevelFormat(dazl.UpperCaseLevelFormat)
	assert.NoError(t, err)
	writer, err = encoder.NewWriter(buf)
	assert.NoError(t, err)
	writer.Error("Hello world!")
	assert.Equal(t, "ERROR Hello world!\n", buf.String())
	buf.Reset()

	encoder, err = encoder.(dazl.CallerEncoder).WithCallerEnabled()
	assert.NoError(t, err)
	writer, err = encoder.NewWriter(buf)
	assert.NoError(t, err)
	writer.Info("Hello world!")
	assert.Equal(t, "INFO \x1b[1mframework_test.go:138\x1b[0m\x1b[36m >\x1b[0m Hello world!\n", buf.String())
	buf.Reset()
}
