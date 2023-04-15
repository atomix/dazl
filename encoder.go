// SPDX-FileCopyrightText: 2023-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package dazl

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
)

type Encoding string

const (
	ConsoleEncoding Encoding = "console"
	JSONEncoding    Encoding = "json"
)

func newEncoder(encoding Encoding, config encodersConfig) encoder {
	switch encoding {
	case ConsoleEncoding:
		return &consoleEncoder{
			config: config.Console,
		}
	case JSONEncoding:
		return &jsonEncoder{
			config: config.JSON,
		}
	default:
		panic(fmt.Sprintf("unkown encoding '%s'", encoding))
	}
}

type encoder interface {
	configure(writer Writer) (Writer, error)
}

type consoleEncoder struct {
	config encoderConfig
}

func (e *consoleEncoder) configure(writer Writer) (Writer, error) {
	var err error
	if e.config.Fields.Level != nil {
		if levelWriter, ok := writer.(LevelWriter); ok {
			writer, err = levelWriter.WithLevel()
			if err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New("level is not an optional field for the configured writer")
		}
		if e.config.Fields.Level.Format != nil {
			if levelFormattingWriter, ok := writer.(LevelFormattingWriter); ok {
				writer, err = levelFormattingWriter.WithLevelFormat(*e.config.Fields.Level.Format)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, errors.New("the configured writer does not support configuring level formats")
			}
		}
	}
	if e.config.Fields.Time != nil {
		if timeWriter, ok := writer.(TimeWriter); ok {
			writer, err = timeWriter.WithTime()
			if err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New("time is not an optional field for the configured writer")
		}
		if e.config.Fields.Time.Format != nil {
			if timeFormattingWriter, ok := writer.(TimeFormattingWriter); ok {
				writer, err = timeFormattingWriter.WithTimeFormat(*e.config.Fields.Time.Format)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, errors.New("the configured writer does not support configuring time formats")
			}
		}
	}
	if e.config.Fields.Caller != nil {
		if callerWriter, ok := writer.(CallerWriter); ok {
			writer, err = callerWriter.WithCaller()
			if err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New("caller is not an optional field for the configured writer")
		}
		if e.config.Fields.Caller.Format != nil {
			if callerFormattingWriter, ok := writer.(CallerFormattingWriter); ok {
				writer, err = callerFormattingWriter.WithCallerFormat(*e.config.Fields.Caller.Format)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, errors.New("the configured writer does not support configuring caller formats")
			}
		}
	}
	if e.config.Fields.Stacktrace != nil {
		if stacktraceWriter, ok := writer.(StacktraceWriter); ok {
			writer, err = stacktraceWriter.WithStacktrace()
			if err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New("stacktrace is not an optional field for the configured writer")
		}
	}
	return writer, nil
}

type jsonEncoder struct {
	config encoderConfig
}

func (e *jsonEncoder) configure(writer Writer) (Writer, error) {
	var err error
	if e.config.Fields.Message != nil {
		if e.config.Fields.Message.Key != "" {
			if messageKeyWriter, ok := writer.(MessageKeyWriter); ok {
				writer, err = messageKeyWriter.WithMessageKey(e.config.Fields.Message.Key)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, errors.New("the configured writer does not support configuring message keys")
			}
		}
	}
	if e.config.Fields.Level != nil {
		if levelWriter, ok := writer.(LevelWriter); ok {
			writer, err = levelWriter.WithLevel()
			if err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New("level is not an optional field for the configured writer")
		}
		if e.config.Fields.Level.Key != "" {
			if levelKeyWriter, ok := writer.(LevelKeyWriter); ok {
				writer, err = levelKeyWriter.WithLevelKey(e.config.Fields.Level.Key)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, errors.New("the configured writer does not support configuring level keys")
			}
		}
		if e.config.Fields.Level.Format != nil {
			if levelFormattingWriter, ok := writer.(LevelFormattingWriter); ok {
				writer, err = levelFormattingWriter.WithLevelFormat(*e.config.Fields.Level.Format)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, errors.New("the configured writer does not support configuring level formats")
			}
		}
	}
	if e.config.Fields.Time != nil {
		if timeWriter, ok := writer.(TimeWriter); ok {
			writer, err = timeWriter.WithTime()
			if err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New("time is not an optional field for the configured writer")
		}
		if e.config.Fields.Time.Key != "" {
			if timeKeyWriter, ok := writer.(TimeKeyWriter); ok {
				writer, err = timeKeyWriter.WithTimeKey(e.config.Fields.Time.Key)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, errors.New("the configured writer does not support configuring time keys")
			}
		}
		if e.config.Fields.Time.Format != nil {
			if timeFormattingWriter, ok := writer.(TimeFormattingWriter); ok {
				writer, err = timeFormattingWriter.WithTimeFormat(*e.config.Fields.Time.Format)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, errors.New("the configured writer does not support configuring time formats")
			}
		}
	}
	if e.config.Fields.Caller != nil {
		if callerWriter, ok := writer.(CallerWriter); ok {
			writer, err = callerWriter.WithCaller()
			if err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New("caller is not an optional field for the configured writer")
		}
		if e.config.Fields.Caller.Key != "" {
			if callerKeyWriter, ok := writer.(CallerKeyWriter); ok {
				writer, err = callerKeyWriter.WithCallerKey(e.config.Fields.Caller.Key)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, errors.New("the configured writer does not support configuring caller keys")
			}
		}
		if e.config.Fields.Caller.Format != nil {
			if callerFormattingWriter, ok := writer.(CallerFormattingWriter); ok {
				writer, err = callerFormattingWriter.WithCallerFormat(*e.config.Fields.Caller.Format)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, errors.New("the configured writer does not support configuring caller formats")
			}
		}
	}
	if e.config.Fields.Stacktrace != nil {
		if stacktraceWriter, ok := writer.(StacktraceWriter); ok {
			writer, err = stacktraceWriter.WithStacktrace()
			if err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New("stacktrace is not an optional field for the configured writer")
		}
		if e.config.Fields.Stacktrace.Key != "" {
			if stacktraceKeyWriter, ok := writer.(StacktraceKeyWriter); ok {
				writer, err = stacktraceKeyWriter.WithStacktraceKey(e.config.Fields.Stacktrace.Key)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, errors.New("the configured writer does not support configuring stacktrace keys")
			}
		}
	}
	return writer, nil
}

type encodersConfig struct {
	Console encoderConfig `json:"console" yaml:"console"`
	JSON    encoderConfig `json:"json" yaml:"json"`
}

type encoderConfig struct {
	Fields encoderFieldsConfig `json:"fields" yaml:"fields"`
}

type encoderFieldsConfig struct {
	Message    *messageEncoderConfig    `json:"message" yaml:"message"`
	Level      *levelEncoderConfig      `json:"level" yaml:"level"`
	Time       *timeEncoderConfig       `json:"time" yaml:"time"`
	Caller     *callerEncoderConfig     `json:"caller" yaml:"caller"`
	Stacktrace *stacktraceEncoderConfig `json:"stacktrace" yaml:"stacktrace"`
}

func (c *encoderFieldsConfig) UnmarshalYAML(unmarshal func(any) error) error {
	var fields []encoderFieldSchema
	if err := unmarshal(&fields); err != nil {
		return err
	}
	for _, field := range fields {
		if field.Message != nil {
			c.Message = field.Message
		}
		if field.Level != nil {
			c.Level = field.Level
		}
		if field.Time != nil {
			c.Time = field.Time
		}
		if field.Caller != nil {
			c.Caller = field.Caller
		}
		if field.Stacktrace != nil {
			c.Stacktrace = field.Stacktrace
		}
	}
	return nil
}

type encoderFieldSchema struct {
	Message    *messageEncoderConfig    `json:"message" yaml:"message"`
	Level      *levelEncoderConfig      `json:"level" yaml:"level"`
	Time       *timeEncoderConfig       `json:"time" yaml:"time"`
	Caller     *callerEncoderConfig     `json:"caller" yaml:"caller"`
	Stacktrace *stacktraceEncoderConfig `json:"stacktrace" yaml:"stacktrace"`
}

func (c *encoderFieldSchema) UnmarshalYAML(unmarshal func(any) error) error {
	config := make(map[string]any)
	if err := unmarshal(&config); err != nil {
		var text string
		if err := unmarshal(&text); err != nil {
			return err
		}
		return c.UnmarshalText([]byte(text))
	}
	if len(config) > 1 {
		return fmt.Errorf("encoder fields must configure one encoder per list item")
	}
	for key, value := range config {
		text, err := yaml.Marshal(value)
		if err != nil {
			return err
		}
		name := fieldEncoderName(key)
		switch name {
		case messageFieldName:
			c.Message = &messageEncoderConfig{}
			return yaml.Unmarshal(text, c.Message)
		case levelFieldName:
			c.Level = &levelEncoderConfig{}
			return yaml.Unmarshal(text, c.Level)
		case timeFieldName:
			c.Time = &timeEncoderConfig{}
			return yaml.Unmarshal(text, c.Time)
		case callerFieldName:
			c.Caller = &callerEncoderConfig{}
			return yaml.Unmarshal(text, c.Caller)
		case stacktraceFieldName:
			c.Stacktrace = &stacktraceEncoderConfig{}
			return yaml.Unmarshal(text, c.Stacktrace)
		default:
			return fmt.Errorf("unknown field encoder '%s'", name)
		}
	}
	return nil
}

func (c *encoderFieldSchema) UnmarshalText(text []byte) error {
	name := fieldEncoderName(text)
	switch name {
	case messageFieldName:
		c.Message = &messageEncoderConfig{}
	case levelFieldName:
		c.Level = &levelEncoderConfig{}
	case timeFieldName:
		c.Time = &timeEncoderConfig{}
	case callerFieldName:
		c.Caller = &callerEncoderConfig{}
	case stacktraceFieldName:
		c.Stacktrace = &stacktraceEncoderConfig{}
	default:
		return fmt.Errorf("unknown field encoder '%s'", name)
	}
	return nil
}

type fieldEncoderName string

const (
	messageFieldName    fieldEncoderName = "message"
	levelFieldName      fieldEncoderName = "level"
	timeFieldName       fieldEncoderName = "time"
	callerFieldName     fieldEncoderName = "caller"
	stacktraceFieldName fieldEncoderName = "stacktrace"
)

type fieldEncoderConfig struct {
	Key string `json:"key" yaml:"key"`
}

type messageEncoderConfig struct {
	fieldEncoderConfig `json:",inline" yaml:",inline"`
}

type LevelFormat string

const (
	LowerCaseLevelFormat LevelFormat = "lowercase"
	UpperCaseLevelFormat LevelFormat = "uppercase"
)

type levelEncoderConfig struct {
	fieldEncoderConfig `json:",inline" yaml:",inline"`
	Format             *LevelFormat `json:"format" yaml:"format"`
}

type TimeFormat string

const (
	ISO8601TimeFormat TimeFormat = "iso8601"
	UnixTimeFormat    TimeFormat = "unix"
)

type timeEncoderConfig struct {
	fieldEncoderConfig `json:",inline" yaml:",inline"`
	Format             *TimeFormat `json:"format" yaml:"format"`
}

type CallerFormat string

const (
	ShortCallerFormat CallerFormat = "short"
	FullCallerFormat  CallerFormat = "full"
)

type callerEncoderConfig struct {
	fieldEncoderConfig `json:",inline" yaml:",inline"`
	Format             *CallerFormat `json:"format" yaml:"format"`
}

type stacktraceEncoderConfig struct {
	fieldEncoderConfig `json:",inline" yaml:",inline"`
}
