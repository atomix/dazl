// SPDX-FileCopyrightText: 2023-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package dazl

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"time"
)

// Writer is a dazl writer
type Writer interface {
	WithName(name string) Writer
	Debug(msg string)
	Info(msg string)
	Error(msg string)
	Fatal(msg string)
	Panic(msg string)
	Warn(msg string)
}

type CallSkippingWriter interface {
	WithSkipCalls(calls int) Writer
}

type NameWriter interface {
	WithNameEnabled() (Writer, error)
}

type NameKeyWriter interface {
	WithNameKey(key string) (Writer, error)
}

type MessageKeyWriter interface {
	WithMessageKey(key string) (Writer, error)
}

type LevelWriter interface {
	WithLevelEnabled() (Writer, error)
}

type LevelKeyWriter interface {
	WithLevelKey(key string) (Writer, error)
}

type LevelFormattingWriter interface {
	WithLevelFormat(format LevelFormat) (Writer, error)
}

type TimestampWriter interface {
	WithTimestampEnabled() (Writer, error)
}

type TimestampKeyWriter interface {
	WithTimestampKey(key string) (Writer, error)
}

type TimestampFormattingWriter interface {
	WithTimestampFormat(format TimestampFormat) (Writer, error)
}

type CallerWriter interface {
	WithCallerEnabled() (Writer, error)
}

type CallerKeyWriter interface {
	WithCallerKey(key string) (Writer, error)
}

type CallerFormattingWriter interface {
	WithCallerFormat(format CallerFormat) (Writer, error)
}

type StacktraceWriter interface {
	WithStacktraceEnabled() (Writer, error)
}

type StacktraceKeyWriter interface {
	WithStacktraceKey(key string) (Writer, error)
}

type BasicSamplingWriter interface {
	WithBasicSampler(interval int, minLevel Level) (Writer, error)
}

type RandomSamplingWriter interface {
	WithRandomSampler(interval int, minLevel Level) (Writer, error)
}

type FieldWriter interface {
	StringFieldWriter
	BoolFieldWriter
	IntFieldWriter
	Int32FieldWriter
	Int64FieldWriter
	UintFieldWriter
	Uint32FieldWriter
	Uint64FieldWriter
	Float32FieldWriter
	Float64FieldWriter
	TimeFieldWriter
	DurationFieldWriter
	BinaryFieldWriter
	BytesFieldWriter
	StringPointerFieldWriter
	BoolPointerFieldWriter
	IntPointerFieldWriter
	Int32PointerFieldWriter
	Int64PointerFieldWriter
	UintPointerFieldWriter
	Uint32PointerFieldWriter
	Uint64PointerFieldWriter
	Float32PointerFieldWriter
	Float64PointerFieldWriter
	TimePointerFieldWriter
	DurationPointerFieldWriter
	StringSliceFieldWriter
	BoolSliceFieldWriter
	IntSliceFieldWriter
	Int32SliceFieldWriter
	Int64SliceFieldWriter
	UintSliceFieldWriter
	Uint32SliceFieldWriter
	Uint64SliceFieldWriter
	Float32SliceFieldWriter
	Float64SliceFieldWriter
	TimeSliceFieldWriter
	DurationSliceFieldWriter
}

type ErrorFieldWriter interface {
	WithErrorField(name string, err error) Writer
}

type StringerFieldWriter interface {
	WithStringerField(name string, value fmt.Stringer) Writer
}

type StringFieldWriter interface {
	WithStringField(name string, value string) Writer
}

type BoolFieldWriter interface {
	WithBoolField(name string, value bool) Writer
}

type IntFieldWriter interface {
	WithIntField(name string, value int) Writer
}

type Int32FieldWriter interface {
	WithInt32Field(name string, value int32) Writer
}

type Int64FieldWriter interface {
	WithInt64Field(name string, value int64) Writer
}

type UintFieldWriter interface {
	WithUintField(name string, value uint) Writer
}

type Uint32FieldWriter interface {
	WithUint32Field(name string, value uint32) Writer
}

type Uint64FieldWriter interface {
	WithUint64Field(name string, value uint64) Writer
}

type Float32FieldWriter interface {
	WithFloat32Field(name string, value float32) Writer
}

type Float64FieldWriter interface {
	WithFloat64Field(name string, value float64) Writer
}

type TimeFieldWriter interface {
	WithTimeField(name string, value time.Time) Writer
}

type DurationFieldWriter interface {
	WithDurationField(name string, value time.Duration) Writer
}

type BinaryFieldWriter interface {
	WithBinaryField(name string, value []byte) Writer
}

type BytesFieldWriter interface {
	WithBytesField(name string, value []byte) Writer
}

type StringPointerFieldWriter interface {
	WithStringPointerField(name string, value *string) Writer
}

type BoolPointerFieldWriter interface {
	WithBoolPointerField(name string, value *bool) Writer
}

type IntPointerFieldWriter interface {
	WithIntPointerField(name string, value *int) Writer
}

type Int32PointerFieldWriter interface {
	WithInt32PointerField(name string, value *int32) Writer
}

type Int64PointerFieldWriter interface {
	WithInt64PointerField(name string, value *int64) Writer
}

type UintPointerFieldWriter interface {
	WithUintPointerField(name string, value *uint) Writer
}

type Uint32PointerFieldWriter interface {
	WithUint32PointerField(name string, value *uint32) Writer
}

type Uint64PointerFieldWriter interface {
	WithUint64PointerField(name string, value *uint64) Writer
}

type Float32PointerFieldWriter interface {
	WithFloat32PointerField(name string, value *float32) Writer
}

type Float64PointerFieldWriter interface {
	WithFloat64PointerField(name string, value *float64) Writer
}

type TimePointerFieldWriter interface {
	WithTimePointerField(name string, value *time.Time) Writer
}

type DurationPointerFieldWriter interface {
	WithDurationPointerField(name string, value *time.Duration) Writer
}

type StringSliceFieldWriter interface {
	WithStringSliceField(name string, values []string) Writer
}

type BoolSliceFieldWriter interface {
	WithBoolSliceField(name string, values []bool) Writer
}

type IntSliceFieldWriter interface {
	WithIntSliceField(name string, values []int) Writer
}

type Int32SliceFieldWriter interface {
	WithInt32SliceField(name string, values []int32) Writer
}

type Int64SliceFieldWriter interface {
	WithInt64SliceField(name string, values []int64) Writer
}

type UintSliceFieldWriter interface {
	WithUintSliceField(name string, values []uint) Writer
}

type Uint32SliceFieldWriter interface {
	WithUint32SliceField(name string, values []uint32) Writer
}

type Uint64SliceFieldWriter interface {
	WithUint64SliceField(name string, values []uint64) Writer
}

type Float32SliceFieldWriter interface {
	WithFloat32SliceField(name string, values []float32) Writer
}

type Float64SliceFieldWriter interface {
	WithFloat64SliceField(name string, values []float64) Writer
}

type TimeSliceFieldWriter interface {
	WithTimeSliceField(name string, values []time.Time) Writer
}

type DurationSliceFieldWriter interface {
	WithDurationSliceField(name string, values []time.Duration) Writer
}

func newWriter(config loggingConfig, name string) (Writer, error) {
	var ioWriter io.Writer
	var encoding Encoding
	switch name {
	case "stdout":
		if config.Writers.Stdout == nil {
			return nil, fmt.Errorf("'%s' writer is not configured", name)
		}
		ioWriter = os.Stdout
		encoding = config.Writers.Stdout.Encoder
	case "stderr":
		if config.Writers.Stderr == nil {
			return nil, fmt.Errorf("'%s' writer is not configured", name)
		}
		ioWriter = os.Stderr
		encoding = config.Writers.Stderr.Encoder
	default:
		config, ok := config.Writers.getFile(name)
		if !ok {
			return nil, fmt.Errorf("'%s' writer is not configured", name)
		}
		writer, err := os.OpenFile(config.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return nil, err
		}
		ioWriter = writer
		encoding = config.Encoder
	}

	writer, err := getFramework().NewWriter(ioWriter, encoding)
	if err != nil {
		return nil, err
	}

	encoder := newEncoder(encoding, config.Encoders)
	writer, err = encoder.configure(writer)
	if err != nil {
		return nil, err
	}
	return writer, nil
}

type writersConfig struct {
	Stdout *stdoutWriterConfig         `json:"stdout" yaml:"stdout"`
	Stderr *stderrWriterConfig         `json:"stderr" yaml:"stderr"`
	Files  map[string]fileWriterConfig `json:"files" yaml:"files"`
}

func (c *writersConfig) getFiles() map[string]fileWriterConfig {
	if c.Files == nil {
		return map[string]fileWriterConfig{}
	}
	return c.Files
}

func (c *writersConfig) getFile(name string) (fileWriterConfig, bool) {
	config, ok := c.getFiles()[name]
	return config, ok
}

func (c *writersConfig) UnmarshalYAML(unmarshal func(any) error) error {
	writers := make(map[string]any)
	if err := unmarshal(writers); err != nil {
		return err
	}

	c.Files = make(map[string]fileWriterConfig)
	for name, config := range writers {
		switch name {
		case "stdout":
			text, err := yaml.Marshal(config)
			if err != nil {
				return err
			}
			writer := &stdoutWriterConfig{}
			if err := yaml.Unmarshal(text, writer); err != nil {
				return err
			}
			if writer.Encoder == "" {
				return fmt.Errorf("writer '%s' is missing required encoder", name)
			}
			c.Stdout = writer
		case "stderr":
			text, err := yaml.Marshal(config)
			if err != nil {
				return err
			}
			writer := &stderrWriterConfig{}
			if err := yaml.Unmarshal(text, writer); err != nil {
				return err
			}
			if writer.Encoder == "" {
				return fmt.Errorf("writer '%s' is missing required encoder", name)
			}
			c.Stderr = writer
		default:
			text, err := yaml.Marshal(config)
			if err != nil {
				return err
			}
			var writer fileWriterConfig
			if err := yaml.Unmarshal(text, &writer); err != nil {
				return err
			}
			if writer.Encoder == "" {
				return fmt.Errorf("writer '%s' is missing required encoder", name)
			}
			if writer.Path == "" {
				return fmt.Errorf("writer '%s' is missing required path", name)
			}
			c.Files[name] = writer
		}
	}
	return nil
}

type writerConfig struct {
	Encoder Encoding `json:"encoder" yaml:"encoder"`
}

type stdoutWriterConfig struct {
	writerConfig `json:",inline" yaml:",inline"`
}

type stderrWriterConfig struct {
	writerConfig `json:",inline" yaml:",inline"`
}

type fileWriterConfig struct {
	writerConfig `json:",inline" yaml:",inline"`
	Path         string `json:"path" yaml:"path"`
}
