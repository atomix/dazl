// SPDX-FileCopyrightText: 2023-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package zerolog

import (
	"bytes"
	"github.com/atomix/dazl"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestWriter(t *testing.T) {
	zerolog.MessageFieldName = "message"
	zerolog.LevelFieldName = "level"
	zerolog.LevelFieldMarshalFunc = func(l zerolog.Level) string {
		return strings.ToLower(l.String())
	}

	buf := &bytes.Buffer{}
	var writer dazl.Writer = &Writer{
		logger: zerolog.New(buf),
	}

	writer.Debug("Hello world!")
	assert.Equal(t, "{\"level\":\"debug\",\"message\":\"Hello world!\"}\n", buf.String())
	buf.Reset()

	writer.Info("Hello world!")
	assert.Equal(t, "{\"level\":\"info\",\"message\":\"Hello world!\"}\n", buf.String())
	buf.Reset()

	writer.Warn("Hello world!")
	assert.Equal(t, "{\"level\":\"warn\",\"message\":\"Hello world!\"}\n", buf.String())
	buf.Reset()

	writer.Error("Hello world!")
	assert.Equal(t, "{\"level\":\"error\",\"message\":\"Hello world!\"}\n", buf.String())
	buf.Reset()

	writer.(dazl.StringFieldWriter).WithStringField("foo", "bar").Info("Hello world!")
	assert.Equal(t, "{\"level\":\"info\",\"foo\":\"bar\",\"message\":\"Hello world!\"}\n", buf.String())
	buf.Reset()

	writer.(dazl.IntFieldWriter).WithIntField("foo", 1).Info("Hello world!")
	assert.Equal(t, "{\"level\":\"info\",\"foo\":1,\"message\":\"Hello world!\"}\n", buf.String())
	buf.Reset()

	writer.(dazl.Int32FieldWriter).WithInt32Field("foo", 1).Info("Hello world!")
	assert.Equal(t, "{\"level\":\"info\",\"foo\":1,\"message\":\"Hello world!\"}\n", buf.String())
	buf.Reset()

	writer.(dazl.Int64FieldWriter).WithInt64Field("foo", 1).Info("Hello world!")
	assert.Equal(t, "{\"level\":\"info\",\"foo\":1,\"message\":\"Hello world!\"}\n", buf.String())
	buf.Reset()

	writer.(dazl.UintFieldWriter).WithUintField("foo", 1).Info("Hello world!")
	assert.Equal(t, "{\"level\":\"info\",\"foo\":1,\"message\":\"Hello world!\"}\n", buf.String())
	buf.Reset()

	writer.(dazl.Uint32FieldWriter).WithUint32Field("foo", 1).Info("Hello world!")
	assert.Equal(t, "{\"level\":\"info\",\"foo\":1,\"message\":\"Hello world!\"}\n", buf.String())
	buf.Reset()

	writer.(dazl.Uint64FieldWriter).WithUint64Field("foo", 1).Info("Hello world!")
	assert.Equal(t, "{\"level\":\"info\",\"foo\":1,\"message\":\"Hello world!\"}\n", buf.String())
	buf.Reset()

	writer.(dazl.Float32FieldWriter).WithFloat32Field("foo", 1).Info("Hello world!")
	assert.Equal(t, "{\"level\":\"info\",\"foo\":1,\"message\":\"Hello world!\"}\n", buf.String())
	buf.Reset()

	writer.(dazl.Float64FieldWriter).WithFloat64Field("foo", 1).Info("Hello world!")
	assert.Equal(t, "{\"level\":\"info\",\"foo\":1,\"message\":\"Hello world!\"}\n", buf.String())
	buf.Reset()

	writer.(dazl.BoolFieldWriter).WithBoolField("foo", true).Info("Hello world!")
	assert.Equal(t, "{\"level\":\"info\",\"foo\":true,\"message\":\"Hello world!\"}\n", buf.String())
	buf.Reset()

	writer.(dazl.StringSliceFieldWriter).WithStringSliceField("foo", []string{"bar"}).Info("Hello world!")
	assert.Equal(t, "{\"level\":\"info\",\"foo\":[\"bar\"],\"message\":\"Hello world!\"}\n", buf.String())
	buf.Reset()

	writer.(dazl.IntSliceFieldWriter).WithIntSliceField("foo", []int{1}).Info("Hello world!")
	assert.Equal(t, "{\"level\":\"info\",\"foo\":[1],\"message\":\"Hello world!\"}\n", buf.String())
	buf.Reset()

	writer.(dazl.Int32SliceFieldWriter).WithInt32SliceField("foo", []int32{1}).Info("Hello world!")
	assert.Equal(t, "{\"level\":\"info\",\"foo\":[1],\"message\":\"Hello world!\"}\n", buf.String())
	buf.Reset()

	writer.(dazl.Int64SliceFieldWriter).WithInt64SliceField("foo", []int64{1}).Info("Hello world!")
	assert.Equal(t, "{\"level\":\"info\",\"foo\":[1],\"message\":\"Hello world!\"}\n", buf.String())
	buf.Reset()

	writer.(dazl.UintSliceFieldWriter).WithUintSliceField("foo", []uint{1}).Info("Hello world!")
	assert.Equal(t, "{\"level\":\"info\",\"foo\":[1],\"message\":\"Hello world!\"}\n", buf.String())
	buf.Reset()

	writer.(dazl.Uint32SliceFieldWriter).WithUint32SliceField("foo", []uint32{1}).Info("Hello world!")
	assert.Equal(t, "{\"level\":\"info\",\"foo\":[1],\"message\":\"Hello world!\"}\n", buf.String())
	buf.Reset()

	writer.(dazl.Uint64SliceFieldWriter).WithUint64SliceField("foo", []uint64{1}).Info("Hello world!")
	assert.Equal(t, "{\"level\":\"info\",\"foo\":[1],\"message\":\"Hello world!\"}\n", buf.String())
	buf.Reset()

	writer.(dazl.Float32SliceFieldWriter).WithFloat32SliceField("foo", []float32{1}).Info("Hello world!")
	assert.Equal(t, "{\"level\":\"info\",\"foo\":[1],\"message\":\"Hello world!\"}\n", buf.String())
	buf.Reset()

	writer.(dazl.Float64SliceFieldWriter).WithFloat64SliceField("foo", []float64{1}).Info("Hello world!")
	assert.Equal(t, "{\"level\":\"info\",\"foo\":[1],\"message\":\"Hello world!\"}\n", buf.String())
	buf.Reset()

	writer.(dazl.BoolSliceFieldWriter).WithBoolSliceField("foo", []bool{true}).Info("Hello world!")
	assert.Equal(t, "{\"level\":\"info\",\"foo\":[true],\"message\":\"Hello world!\"}\n", buf.String())
	buf.Reset()
}
