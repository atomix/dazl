// SPDX-FileCopyrightText: 2023-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package dazl

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

type Encoding string

const (
	ConsoleEncoding Encoding = "console"
	JSONEncoding    Encoding = "json"
)

type encodersConfig struct {
	Console *encoderConfig `json:"console" yaml:"console"`
	JSON    *encoderConfig `json:"json" yaml:"json"`
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
		field := fieldEncoderConfig{
			Key: key,
		}
		switch name {
		case messageFieldName:
			c.Message = &messageEncoderConfig{
				fieldEncoderConfig: field,
			}
			return yaml.Unmarshal(text, c.Message)
		case levelFieldName:
			c.Level = &levelEncoderConfig{
				fieldEncoderConfig: field,
				Format:             UpperCaseLevelFormat,
			}
			return yaml.Unmarshal(text, c.Level)
		case timeFieldName:
			c.Time = &timeEncoderConfig{
				fieldEncoderConfig: field,
				Format:             ISO8601TimeFormat,
			}
			return yaml.Unmarshal(text, c.Time)
		case callerFieldName:
			c.Caller = &callerEncoderConfig{
				fieldEncoderConfig: field,
				Format:             ShortCallerFormat,
			}
			return yaml.Unmarshal(text, c.Caller)
		case stacktraceFieldName:
			c.Stacktrace = &stacktraceEncoderConfig{
				fieldEncoderConfig: field,
			}
			return yaml.Unmarshal(text, c.Stacktrace)
		default:
			return fmt.Errorf("unknown field encoder '%s'", name)
		}
	}
	return nil
}

func (c *encoderFieldSchema) UnmarshalText(text []byte) error {
	name := fieldEncoderName(text)
	field := fieldEncoderConfig{
		Key: string(text),
	}
	switch name {
	case messageFieldName:
		c.Message = &messageEncoderConfig{
			fieldEncoderConfig: field,
		}
	case levelFieldName:
		c.Level = &levelEncoderConfig{
			fieldEncoderConfig: field,
			Format:             UpperCaseLevelFormat,
		}
	case timeFieldName:
		c.Time = &timeEncoderConfig{
			fieldEncoderConfig: field,
			Format:             ISO8601TimeFormat,
		}
	case callerFieldName:
		c.Caller = &callerEncoderConfig{
			fieldEncoderConfig: field,
			Format:             ShortCallerFormat,
		}
	case stacktraceFieldName:
		c.Stacktrace = &stacktraceEncoderConfig{
			fieldEncoderConfig: field,
		}
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

func (c *messageEncoderConfig) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		c.Key = string(messageFieldName)
	} else {
		c.Key = string(text)
	}
	return nil
}

type LevelFormat string

const (
	LowerCaseLevelFormat LevelFormat = "lowercase"
	UpperCaseLevelFormat LevelFormat = "uppercase"
)

func (f *LevelFormat) UnmarshalText(text []byte) error {
	format := LevelFormat(text)
	switch format {
	case LowerCaseLevelFormat, UpperCaseLevelFormat:
		*f = format
	case "":
		*f = UpperCaseLevelFormat
	default:
		return fmt.Errorf("unknown level format '%s'", format)
	}
	return nil
}

type levelEncoderConfig struct {
	fieldEncoderConfig `json:",inline" yaml:",inline"`
	Format             LevelFormat `json:"format" yaml:"format"`
}

func (c *levelEncoderConfig) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		c.Key = string(levelFieldName)
	} else {
		c.Key = string(text)
	}
	return nil
}

type TimeFormat string

const (
	ISO8601TimeFormat TimeFormat = "iso8601"
	UnixTimeFormat    TimeFormat = "unix"
)

func (f *TimeFormat) UnmarshalText(text []byte) error {
	format := TimeFormat(text)
	switch format {
	case ISO8601TimeFormat, UnixTimeFormat:
		*f = format
	case "":
		*f = ISO8601TimeFormat
	default:
		return fmt.Errorf("unknown time format '%s'", format)
	}
	return nil
}

type timeEncoderConfig struct {
	fieldEncoderConfig `json:",inline" yaml:",inline"`
	Format             TimeFormat `json:"format" yaml:"format"`
}

func (c *timeEncoderConfig) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		c.Key = string(timeFieldName)
	} else {
		c.Key = string(text)
	}
	return nil
}

type CallerFormat string

const (
	ShortCallerFormat CallerFormat = "short"
	FulCallerFormat   CallerFormat = "full"
)

func (f *CallerFormat) UnmarshalText(text []byte) error {
	format := CallerFormat(text)
	switch format {
	case ShortCallerFormat, FulCallerFormat:
		*f = format
	case "":
		*f = ShortCallerFormat
	default:
		return fmt.Errorf("unknown caller format '%s'", format)
	}
	return nil
}

type callerEncoderConfig struct {
	fieldEncoderConfig `json:",inline" yaml:",inline"`
	Format             CallerFormat `json:"format" yaml:"format"`
}

func (c *callerEncoderConfig) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		c.Key = string(callerFieldName)
	} else {
		c.Key = string(text)
	}
	return nil
}

type stacktraceEncoderConfig struct {
	fieldEncoderConfig `json:",inline" yaml:",inline"`
}

func (c *stacktraceEncoderConfig) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		c.Key = string(stacktraceFieldName)
	} else {
		c.Key = string(text)
	}
	return nil
}
