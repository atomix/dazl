// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package dazl

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
)

var root Logger

const pathSep = "/"

// GetPackageLogger gets the logger for the current package.
func GetPackageLogger() Logger {
	pkg, ok := getCallerPackage()
	if !ok {
		panic("could not retrieve logger package")
	}
	return GetLogger(pkg)
}

// GetLogger gets the logger for the given path.
func GetLogger(path string) Logger {
	return root.GetLogger(path)
}

// Logger represents an abstract logging interface.
type Logger interface {
	// Name returns the logger name
	Name() string

	// GetLogger gets a descendant of this Logger
	GetLogger(path string) Logger

	// WithFields adds fields to the logger
	WithFields(fields ...Field) Logger

	// WithSkipCalls skipsthe given number of calls to the logger methods
	WithSkipCalls(calls int) Logger

	Debug(...interface{})
	Debugf(format string, args ...interface{})
	Debugw(msg string, fields ...Field)

	Info(...interface{})
	Infof(format string, args ...interface{})
	Infow(msg string, fields ...Field)

	Warn(...interface{})
	Warnf(format string, args ...interface{})
	Warnw(msg string, fields ...Field)

	Error(...interface{})
	Errorf(format string, args ...interface{})
	Errorw(msg string, fields ...Field)

	Fatal(...interface{})
	Fatalf(format string, args ...interface{})
	Fatalw(msg string, fields ...Field)

	Panic(...interface{})
	Panicf(format string, args ...interface{})
	Panicw(msg string, fields ...Field)
}

// getCallerPackage gets the package name of the calling function'ss caller
func getCallerPackage() (string, bool) {
	var pkg string
	pc, _, _, ok := runtime.Caller(2)
	if !ok {
		return pkg, false
	}
	parts := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	if parts[len(parts)-2][0] == '(' {
		pkg = strings.Join(parts[0:len(parts)-2], ".")
	} else {
		pkg = strings.Join(parts[0:len(parts)-1], ".")
	}
	return pkg, true
}

func configure(config loggingConfig) error {
	context, err := newLoggingContext(config)
	if err != nil {
		return err
	}
	logger, err := newLogger(context, nil, "")
	if err != nil {
		return err
	}
	root = logger
	return nil
}

func newLogger(context *loggingContext, parent *dazlLogger, name string) (*dazlLogger, error) {
	var config loggerConfig
	var logger *dazlLogger
	var path []string
	if parent != nil {
		path = append(parent.path, name)
		loggerName := strings.Join(path, pathSep)
		config, _ = context.config.getLogger(loggerName)
		defaultLevel := parent.defaultLevel
		if parent.level != EmptyLevel {
			defaultLevel = parent.level
		}
		logger = &dazlLogger{
			loggerContext: &loggerContext{
				loggingContext: context,
				name:           loggerName,
				path:           path,
				defaultLevel:   defaultLevel,
				sampler:        parent.sampler,
			},
			outputs: make(map[string]*dazlOutput),
		}

		for outputName, output := range parent.outputs {
			logger.outputs[outputName] = output.WithWriter(output.writer.WithName(loggerName))
		}
	} else {
		config = context.config.RootLogger
		logger = &dazlLogger{
			loggerContext: &loggerContext{
				loggingContext: context,
				sampler:        &allSampler{},
			},
			outputs: make(map[string]*dazlOutput),
		}
	}

	level := config.Level.Level()
	if level != EmptyLevel {
		logger.level = level
	}

	if config.Sample.Basic != nil {
		logger.sampler = &basicSampler{
			Interval: uint32(config.Sample.Basic.Interval),
			MinLevel: config.Sample.Basic.MinLevel.Level(),
		}
	} else if config.Sample.Random != nil {
		logger.sampler = randomSampler{
			Interval: config.Sample.Random.Interval,
			MinLevel: config.Sample.Random.MinLevel.Level(),
		}
	}

	for _, outputConfig := range config.Outputs {
		// If the configured output already exists, override the output configuration.
		// Otherwise, create a new output.
		output, ok := logger.outputs[outputConfig.Writer]
		if !ok {
			writer, err := context.getWriter(outputConfig.Writer)
			if err != nil {
				return nil, err
			}
			if logger.name != "" {
				writer = writer.WithName(logger.name)
			}
			output = newOutput(writer, EmptyLevel, &allSampler{})
		}

		// Add the level to the output if configured
		outputLevel := outputConfig.Level.Level()
		if outputLevel != EmptyLevel {
			output = output.WithLevel(outputLevel)
		}

		// Configure sampling for the output
		if outputConfig.Sample.Basic != nil {
			if samplingWriter, ok := output.Writer().(BasicSamplingWriter); ok {
				writer, err := samplingWriter.WithBasicSampler(
					outputConfig.Sample.Basic.Interval,
					outputConfig.Sample.Basic.MinLevel.Level())
				if err != nil {
					return nil, err
				}
				output = output.WithWriter(writer)
			} else {
				output = output.WithSampler(&basicSampler{
					Interval: uint32(outputConfig.Sample.Basic.Interval),
					MinLevel: outputConfig.Sample.Basic.MinLevel.Level(),
				})
			}
		} else if outputConfig.Sample.Random != nil {
			if samplingWriter, ok := output.Writer().(RandomSamplingWriter); ok {
				writer, err := samplingWriter.WithRandomSampler(outputConfig.Sample.Random.Interval, outputConfig.Sample.Random.MinLevel.Level())
				if err != nil {
					return nil, err
				}
				output = output.WithWriter(writer)
			} else {
				output = output.WithSampler(randomSampler{
					Interval: outputConfig.Sample.Random.Interval,
					MinLevel: outputConfig.Sample.Random.MinLevel.Level(),
				})
			}
		}
		logger.outputs[outputConfig.Writer] = output
	}
	return logger, nil
}

func newLoggingContext(config loggingConfig) (*loggingContext, error) {
	framework := getFramework()
	encoders := make(map[Encoding]Encoder)
	if consoleEncodingFramework, ok := framework.(ConsoleEncodingFramework); ok {
		encoder, err := configureConsoleEncoder(config.Encoders.Console, consoleEncodingFramework.ConsoleEncoder())
		if err != nil {
			return nil, err
		}
		encoders[ConsoleEncoding] = encoder
	}
	if jsonEncodingFramework, ok := framework.(JSONEncodingFramework); ok {
		encoder, err := configureJSONEncoder(config.Encoders.JSON, jsonEncodingFramework.JSONEncoder())
		if err != nil {
			return nil, err
		}
		encoders[JSONEncoding] = encoder
	}
	return &loggingContext{
		config:   config,
		encoders: encoders,
	}, nil
}

type loggingContext struct {
	config   loggingConfig
	encoders map[Encoding]Encoder
	writers  sync.Map
	mu       sync.Mutex
}

func (c *loggingContext) getWriter(name string) (Writer, error) {
	writer, ok := c.writers.Load(name)
	if ok {
		return writer.(Writer), nil
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	writer, ok = c.writers.Load(name)
	if ok {
		return writer.(Writer), nil
	}

	writer, err := c.newWriter(name)
	if err != nil {
		return nil, err
	}

	c.writers.Store(name, writer)
	return writer.(Writer), nil
}

func (c *loggingContext) newWriter(name string) (Writer, error) {
	switch name {
	case "stdout":
		if c.config.Writers.Stdout == nil {
			return nil, fmt.Errorf("'%s' writer is not configured", name)
		}
		encoder, ok := c.encoders[c.config.Writers.Stdout.Encoder]
		if !ok {
			return nil, fmt.Errorf("%s framework does not support %s encoding", getFramework().Name(), c.config.Writers.Stdout.Encoder)
		}
		return encoder.NewWriter(os.Stdout)
	case "stderr":
		if c.config.Writers.Stderr == nil {
			return nil, fmt.Errorf("'%s' writer is not configured", name)
		}
		encoder, ok := c.encoders[c.config.Writers.Stderr.Encoder]
		if !ok {
			return nil, fmt.Errorf("%s framework does not support %s encoding", getFramework().Name(), c.config.Writers.Stderr.Encoder)
		}
		return encoder.NewWriter(os.Stderr)
	default:
		config, ok := c.config.Writers.getFile(name)
		if !ok {
			return nil, fmt.Errorf("'%s' writer is not configured", name)
		}
		file, err := os.OpenFile(config.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return nil, err
		}
		encoder, ok := c.encoders[config.Encoder]
		if !ok {
			return nil, fmt.Errorf("%s framework does not support %s encoding", getFramework().Name(), config.Encoder)
		}
		return encoder.NewWriter(file)
	}
}

type loggerContext struct {
	*loggingContext
	name         string
	path         []string
	children     sync.Map
	mu           sync.Mutex
	level        Level
	defaultLevel Level
	sampler      Sampler
}

type dazlLogger struct {
	*loggerContext
	outputs map[string]*dazlOutput
}

func (l *dazlLogger) Name() string {
	return l.name
}

func (l *dazlLogger) Level() Level {
	if l.level != EmptyLevel {
		return l.level
	}
	return l.defaultLevel
}

func (l *dazlLogger) GetLogger(path string) Logger {
	if path == "" {
		panic("logger path must not be empty")
	}
	logger := l
	names := strings.Split(path, pathSep)
	for _, name := range names {
		child, err := logger.getChild(name)
		if err != nil {
			panic(err)
		}
		logger = child
	}
	return logger
}

func (l *dazlLogger) getChild(name string) (*dazlLogger, error) {
	child, ok := l.children.Load(name)
	if ok {
		return child.(*dazlLogger), nil
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	child, ok = l.children.Load(name)
	if ok {
		return child.(*dazlLogger), nil
	}

	logger, err := newLogger(l.loggingContext, l, name)
	if err != nil {
		return nil, err
	}

	l.children.Store(name, logger)
	return logger, nil
}

func (l *dazlLogger) WithFields(fields ...Field) Logger {
	outputs := make(map[string]*dazlOutput)
	for name, output := range l.outputs {
		writer := output.Writer()
		var err error
		for _, field := range fields {
			if writer, err = field(writer); err != nil {
				panic(err)
			}
		}
		outputs[name] = output.WithWriter(writer)
	}
	return &dazlLogger{
		loggerContext: l.loggerContext,
		outputs:       outputs,
	}
}

func (l *dazlLogger) WithSkipCalls(calls int) Logger {
	outputs := make(map[string]*dazlOutput)
	for name, output := range l.outputs {
		if writer, ok := output.Writer().(CallSkippingWriter); ok {
			outputs[name] = output.WithWriter(writer.WithSkipCalls(calls))
		} else {
			outputs[name] = output
		}
	}
	return &dazlLogger{
		loggerContext: l.loggerContext,
		outputs:       outputs,
	}
}

func (l *dazlLogger) Debug(args ...interface{}) {
	if l.Level().Enabled(DebugLevel) && l.sampler.Sample(DebugLevel) {
		for _, output := range l.outputs {
			output.Debug(fmt.Sprint(args...))
		}
	}
}

func (l *dazlLogger) Debugf(format string, args ...interface{}) {
	if l.Level().Enabled(DebugLevel) && l.sampler.Sample(DebugLevel) {
		for _, output := range l.outputs {
			output.Debug(fmt.Sprintf(format, args...))
		}
	}
}

func (l *dazlLogger) Debugw(msg string, fields ...Field) {
	l.WithFields(fields...).Debug(msg)
}

func (l *dazlLogger) Info(args ...interface{}) {
	if l.Level().Enabled(InfoLevel) && l.sampler.Sample(InfoLevel) {
		for _, output := range l.outputs {
			output.Info(fmt.Sprint(args...))
		}
	}
}

func (l *dazlLogger) Infof(format string, args ...interface{}) {
	if l.Level().Enabled(InfoLevel) && l.sampler.Sample(InfoLevel) {
		for _, output := range l.outputs {
			output.Info(fmt.Sprintf(format, args...))
		}
	}
}

func (l *dazlLogger) Infow(msg string, fields ...Field) {
	l.WithFields(fields...).Info(msg)
}

func (l *dazlLogger) Warn(args ...interface{}) {
	if l.Level().Enabled(WarnLevel) && l.sampler.Sample(WarnLevel) {
		for _, output := range l.outputs {
			output.Warn(fmt.Sprint(args...))
		}
	}
}

func (l *dazlLogger) Warnf(format string, args ...interface{}) {
	if l.Level().Enabled(WarnLevel) && l.sampler.Sample(WarnLevel) {
		for _, output := range l.outputs {
			output.Warn(fmt.Sprintf(format, args...))
		}
	}
}

func (l *dazlLogger) Warnw(msg string, fields ...Field) {
	l.WithFields(fields...).Warn(msg)
}

func (l *dazlLogger) Error(args ...interface{}) {
	if l.Level().Enabled(ErrorLevel) && l.sampler.Sample(ErrorLevel) {
		for _, output := range l.outputs {
			output.Error(fmt.Sprint(args...))
		}
	}
}

func (l *dazlLogger) Errorf(format string, args ...interface{}) {
	if l.Level().Enabled(ErrorLevel) && l.sampler.Sample(ErrorLevel) {
		for _, output := range l.outputs {
			output.Error(fmt.Sprintf(format, args...))
		}
	}
}

func (l *dazlLogger) Errorw(msg string, fields ...Field) {
	l.WithFields(fields...).Error(msg)
}

func (l *dazlLogger) Fatal(args ...interface{}) {
	if l.Level().Enabled(FatalLevel) && l.sampler.Sample(FatalLevel) {
		for _, output := range l.outputs {
			output.Fatal(fmt.Sprint(args...))
		}
	}
}

func (l *dazlLogger) Fatalf(format string, args ...interface{}) {
	if l.Level().Enabled(FatalLevel) && l.sampler.Sample(FatalLevel) {
		for _, output := range l.outputs {
			output.Fatal(fmt.Sprintf(format, args...))
		}
	}
}

func (l *dazlLogger) Fatalw(msg string, fields ...Field) {
	l.WithFields(fields...).Fatal(msg)
}

func (l *dazlLogger) Panic(args ...interface{}) {
	if l.Level().Enabled(PanicLevel) && l.sampler.Sample(PanicLevel) {
		for _, output := range l.outputs {
			output.Panic(fmt.Sprint(args...))
		}
	}
}

func (l *dazlLogger) Panicf(format string, args ...interface{}) {
	if l.Level().Enabled(PanicLevel) && l.sampler.Sample(PanicLevel) {
		for _, output := range l.outputs {
			output.Panic(fmt.Sprintf(format, args...))
		}
	}
}

func (l *dazlLogger) Panicw(msg string, fields ...Field) {
	l.WithFields(fields...).Panic(msg)
}

var _ Logger = &dazlLogger{}

type loggerConfig struct {
	Level   levelConfig    `json:"level" yaml:"level"`
	Sample  samplingConfig `json:"sample" yaml:"sample"`
	Outputs []outputConfig `json:"outputs" yaml:"outputs"`
}
