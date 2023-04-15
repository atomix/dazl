// SPDX-FileCopyrightText: 2023-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package zerolog

import (
	"fmt"
	"github.com/atomix/dazl"
	"github.com/rs/zerolog"
	"io"
	"time"
)

type consoleEncoder struct{}

func (e *consoleEncoder) NewWriter(writer io.Writer) (dazl.Writer, error) {
	return &Writer{
		logger: zerolog.New(zerolog.ConsoleWriter{
			Out: writer,
		}),
	}, nil
}

type jsonEncoder struct {
	nameKey string
}

func (e *jsonEncoder) NewWriter(writer io.Writer) (dazl.Writer, error) {
	return &Writer{
		logger: zerolog.New(writer),
	}, nil
}

func (e *jsonEncoder) WithMessageKey(key string) (dazl.Encoder, error) {
	zerolog.MessageFieldName = key
	return e, nil
}

func (e *jsonEncoder) WithNameEnabled() (dazl.Encoder, error) {
	return &jsonEncoder{
		nameKey: "logger",
	}, nil
}

func (e *jsonEncoder) WithNameKey(key string) (dazl.Encoder, error) {
	return &jsonEncoder{
		nameKey: key,
	}, nil
}

func (e *jsonEncoder) WithLevelEnabled() (dazl.Encoder, error) {
	return e, nil
}

func (e *jsonEncoder) WithLevelKey(key string) (dazl.Encoder, error) {
	zerolog.LevelFieldName = key
	return e, nil
}

func (e *jsonEncoder) WithLevelFormat(format dazl.LevelFormat) (dazl.Encoder, error) {
	switch format {
	case dazl.LowerCaseLevelFormat:
		zerolog.LevelTraceValue = "trace"
		zerolog.LevelDebugValue = "debug"
		zerolog.LevelInfoValue = "info"
		zerolog.LevelWarnValue = "warn"
		zerolog.LevelErrorValue = "error"
		zerolog.LevelFatalValue = "fatal"
		zerolog.LevelPanicValue = "panic"
	case dazl.UpperCaseLevelFormat:
		zerolog.LevelTraceValue = "TRACE"
		zerolog.LevelDebugValue = "DEBUG"
		zerolog.LevelInfoValue = "INFO"
		zerolog.LevelWarnValue = "WARN"
		zerolog.LevelErrorValue = "ERROR"
		zerolog.LevelFatalValue = "FATAL"
		zerolog.LevelPanicValue = "PANIC"
	}
	return e, nil
}

func (e *jsonEncoder) WithTimestampEnabled() (dazl.Encoder, error) {
	return e, nil
}

func (e *jsonEncoder) WithTimestampKey(key string) (dazl.Encoder, error) {
	zerolog.TimestampFieldName = key
	return e, nil
}

func (e *jsonEncoder) WithTimestampFormat(format dazl.TimestampFormat) (dazl.Encoder, error) {
	switch format {
	case dazl.UnixTimestampFormat:
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	case dazl.ISO8601TimestampFormat:
		zerolog.TimeFieldFormat = time.RFC3339
	default:
		return nil, fmt.Errorf("unsupoorted timestamp format %s", format)
	}
	return e, nil
}

func (e *jsonEncoder) WithCallerEnabled() (dazl.Encoder, error) {
	return e, nil
}

func (e *jsonEncoder) WithCallerKey(key string) (dazl.Encoder, error) {
	zerolog.CallerFieldName = key
	return e, nil
}

func (e *jsonEncoder) WithCallerFormat(format dazl.CallerFormat) (dazl.Encoder, error) {
	switch format {
	case dazl.ShortCallerFormat:
	default:
		return nil, fmt.Errorf("unsupoorted caller format %s", format)
	}
	return e, nil
}

func (e *jsonEncoder) WithStacktraceEnabled() (dazl.Encoder, error) {
	return e, nil
}

func (e *jsonEncoder) WithStacktraceKey(key string) (dazl.Encoder, error) {
	zerolog.ErrorStackFieldName = key
	return e, nil
}

var _ dazl.MessageKeyEncoder = (*jsonEncoder)(nil)
var _ dazl.NameEncoder = (*jsonEncoder)(nil)
var _ dazl.NameKeyEncoder = (*jsonEncoder)(nil)
var _ dazl.LevelEncoder = (*jsonEncoder)(nil)
var _ dazl.LevelKeyEncoder = (*jsonEncoder)(nil)
var _ dazl.LevelFormattingEncoder = (*jsonEncoder)(nil)
var _ dazl.TimestampEncoder = (*jsonEncoder)(nil)
var _ dazl.TimestampKeyEncoder = (*jsonEncoder)(nil)
var _ dazl.TimestampFormattingEncoder = (*jsonEncoder)(nil)
var _ dazl.CallerEncoder = (*jsonEncoder)(nil)
var _ dazl.CallerKeyEncoder = (*jsonEncoder)(nil)
var _ dazl.CallerFormattingEncoder = (*jsonEncoder)(nil)
var _ dazl.StacktraceEncoder = (*jsonEncoder)(nil)
var _ dazl.StacktraceKeyEncoder = (*jsonEncoder)(nil)
