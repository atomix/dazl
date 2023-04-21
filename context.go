// SPDX-FileCopyrightText: 2023-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package dazl

type ContextualLogger struct {
	logger Logger
}

func (l ContextualLogger) With() Context {
	return Context{
		logger: l.logger,
	}
}

func (l ContextualLogger) Debug() Event {
	return Event{
		logger: l.logger,
		msgFunc: func(logger Logger, msg string) {
			logger.Debug(msg)
		},
		msgfFunc: func(logger Logger, msg string, args ...any) {
			logger.Debugf(msg, args...)
		},
	}
}

func (l ContextualLogger) Info() Event {
	return Event{
		logger: l.logger,
		msgFunc: func(logger Logger, msg string) {
			logger.Info(msg)
		},
		msgfFunc: func(logger Logger, msg string, args ...any) {
			logger.Infof(msg, args...)
		},
	}
}

func (l ContextualLogger) Warn() Event {
	return Event{
		logger: l.logger,
		msgFunc: func(logger Logger, msg string) {
			logger.Warn(msg)
		},
		msgfFunc: func(logger Logger, msg string, args ...any) {
			logger.Warnf(msg, args...)
		},
	}
}

func (l ContextualLogger) Error() Event {
	return Event{
		logger: l.logger,
		msgFunc: func(logger Logger, msg string) {
			logger.Error(msg)
		},
		msgfFunc: func(logger Logger, msg string, args ...any) {
			logger.Errorf(msg, args...)
		},
	}
}

func (l ContextualLogger) Panic() Event {
	return Event{
		logger: l.logger,
		msgFunc: func(logger Logger, msg string) {
			logger.Panic(msg)
		},
		msgfFunc: func(logger Logger, msg string, args ...any) {
			logger.Panicf(msg, args...)
		},
	}
}

func (l ContextualLogger) Fatal() Event {
	return Event{
		logger: l.logger,
		msgFunc: func(logger Logger, msg string) {
			logger.Fatal(msg)
		},
		msgfFunc: func(logger Logger, msg string, args ...any) {
			logger.Fatalf(msg, args...)
		},
	}
}

type Event struct {
	logger   Logger
	msgFunc  func(Logger, string)
	msgfFunc func(Logger, string, ...any)
}

func (e Event) Str(name string, value string) Event {
	return Event{
		logger: e.logger.WithFields(String(name, value)),
	}
}

func (e Event) Msg(msg string) {
	e.msgFunc(e.logger, msg)
}

type Context struct {
	logger Logger
}

func (c Context) Logger() ContextualLogger {
	return ContextualLogger{
		logger: c.logger,
	}
}

func (c Context) Str(name string, value string) Context {
	return Context{
		logger: c.logger.WithFields(String(name, value)),
	}
}

func (c Context) Int(name string, value int) Context {
	return Context{
		logger: c.logger.WithFields(Int(name, value)),
	}
}

func (c Context) Int32(name string, value int32) Context {
	return Context{
		logger: c.logger.WithFields(Int32(name, value)),
	}
}

func (c Context) Int64(name string, value int64) Context {
	return Context{
		logger: c.logger.WithFields(Int64(name, value)),
	}
}

func (c Context) Uint(name string, value uint) Context {
	return Context{
		logger: c.logger.WithFields(Uint(name, value)),
	}
}

func (c Context) Uint32(name string, value uint32) Context {
	return Context{
		logger: c.logger.WithFields(Uint32(name, value)),
	}
}

func (c Context) Uint64(name string, value uint64) Context {
	return Context{
		logger: c.logger.WithFields(Uint64(name, value)),
	}
}

func (c Context) Float32(name string, value float32) Context {
	return Context{
		logger: c.logger.WithFields(Float32(name, value)),
	}
}

func (c Context) Float64(name string, value float64) Context {
	return Context{
		logger: c.logger.WithFields(Float64(name, value)),
	}
}

func (c Context) Bool(name string, value bool) Context {
	return Context{
		logger: c.logger.WithFields(Bool(name, value)),
	}
}
