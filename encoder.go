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
	if e.config.Fields.Name != nil {
		if nameWriter, ok := writer.(NameWriter); ok {
			writer, err = nameWriter.WithNameEnabled()
			if err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New("level is not an optional field for the configured writer")
		}
	}
	if e.config.Fields.Level != nil {
		if levelWriter, ok := writer.(LevelWriter); ok {
			writer, err = levelWriter.WithLevelEnabled()
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
		if timestampWriter, ok := writer.(TimestampWriter); ok {
			writer, err = timestampWriter.WithTimestampEnabled()
			if err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New("time is not an optional field for the configured writer")
		}
		if e.config.Fields.Time.Format != nil {
			if timestampFormattingWriter, ok := writer.(TimestampFormattingWriter); ok {
				writer, err = timestampFormattingWriter.WithTimestampFormat(*e.config.Fields.Time.Format)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, errors.New("the configured writer does not support configuring timestamp formats")
			}
		}
	}
	if e.config.Fields.Caller != nil {
		if callerWriter, ok := writer.(CallerWriter); ok {
			writer, err = callerWriter.WithCallerEnabled()
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
			writer, err = stacktraceWriter.WithStacktraceEnabled()
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
	if e.config.Fields.Name != nil {
		if nameWriter, ok := writer.(NameWriter); ok {
			writer, err = nameWriter.WithNameEnabled()
			if err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New("name is not an optional field for the configured writer")
		}
		if e.config.Fields.Name.Key != "" {
			if nameKeyWriter, ok := writer.(NameKeyWriter); ok {
				writer, err = nameKeyWriter.WithNameKey(e.config.Fields.Name.Key)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, errors.New("the configured writer does not support configuring name keys")
			}
		}
	}
	if e.config.Fields.Level != nil {
		if levelWriter, ok := writer.(LevelWriter); ok {
			writer, err = levelWriter.WithLevelEnabled()
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
		if timestampWriter, ok := writer.(TimestampWriter); ok {
			writer, err = timestampWriter.WithTimestampEnabled()
			if err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New("timestamp is not an optional field for the configured writer")
		}
		if e.config.Fields.Time.Key != "" {
			if timestampKeyWriter, ok := writer.(TimestampKeyWriter); ok {
				writer, err = timestampKeyWriter.WithTimestampKey(e.config.Fields.Time.Key)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, errors.New("the configured writer does not support configuring timestamp keys")
			}
		}
		if e.config.Fields.Time.Format != nil {
			if timestampFormattingWriter, ok := writer.(TimestampFormattingWriter); ok {
				writer, err = timestampFormattingWriter.WithTimestampFormat(*e.config.Fields.Time.Format)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, errors.New("the configured writer does not support configuring timestamp formats")
			}
		}
	}
	if e.config.Fields.Caller != nil {
		if callerWriter, ok := writer.(CallerWriter); ok {
			writer, err = callerWriter.WithCallerEnabled()
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
			writer, err = stacktraceWriter.WithStacktraceEnabled()
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
	Name       *nameEncoderConfig       `json:"name" yaml:"name"`
	Message    *messageEncoderConfig    `json:"message" yaml:"message"`
	Level      *levelEncoderConfig      `json:"level" yaml:"level"`
	Time       *timestampEncoderConfig  `json:"timestamp" yaml:"timestamp"`
	Caller     *callerEncoderConfig     `json:"caller" yaml:"caller"`
	Stacktrace *stacktraceEncoderConfig `json:"stacktrace" yaml:"stacktrace"`
}

func (c *encoderFieldsConfig) UnmarshalYAML(unmarshal func(any) error) error {
	var fields []encoderFieldSchema
	if err := unmarshal(&fields); err != nil {
		return err
	}
	for _, field := range fields {
		if field.Name != nil {
			c.Name = field.Name
		}
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
	Name       *nameEncoderConfig       `json:"name" yaml:"name"`
	Message    *messageEncoderConfig    `json:"message" yaml:"message"`
	Level      *levelEncoderConfig      `json:"level" yaml:"level"`
	Time       *timestampEncoderConfig  `json:"timestamp" yaml:"timestamp"`
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
		case nameFieldName:
			c.Name = &nameEncoderConfig{}
			return yaml.Unmarshal(text, c.Name)
		case messageFieldName:
			c.Message = &messageEncoderConfig{}
			return yaml.Unmarshal(text, c.Message)
		case levelFieldName:
			c.Level = &levelEncoderConfig{}
			return yaml.Unmarshal(text, c.Level)
		case timestampFieldName:
			c.Time = &timestampEncoderConfig{}
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
	case nameFieldName:
		c.Name = &nameEncoderConfig{}
	case messageFieldName:
		c.Message = &messageEncoderConfig{}
	case levelFieldName:
		c.Level = &levelEncoderConfig{}
	case timestampFieldName:
		c.Time = &timestampEncoderConfig{}
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
	nameFieldName       fieldEncoderName = "name"
	messageFieldName    fieldEncoderName = "message"
	levelFieldName      fieldEncoderName = "level"
	timestampFieldName  fieldEncoderName = "timestamp"
	callerFieldName     fieldEncoderName = "caller"
	stacktraceFieldName fieldEncoderName = "stacktrace"
)

type fieldEncoderConfig struct {
	Key string `json:"key" yaml:"key"`
}

type nameEncoderConfig struct {
	fieldEncoderConfig `json:",inline" yaml:",inline"`
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

type TimestampFormat string

const (
	ISO8601TimestampFormat TimestampFormat = "iso8601"
	UnixTimestampFormat    TimestampFormat = "unix"
)

type timestampEncoderConfig struct {
	fieldEncoderConfig `json:",inline" yaml:",inline"`
	Format             *TimestampFormat `json:"format" yaml:"format"`
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
