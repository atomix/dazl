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
	"sync/atomic"
)

var root Logger

const pathSep = "/"

func init() {
	var config loggingConfig
	if err := load(&config); err != nil {
		panic(err)
	} else if err := configure(config); err != nil {
		panic(err)
	}
}

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

	// Level returns the logger's level
	Level() Level

	// SetLevel sets the logger's level
	SetLevel(level Level)

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

// SetLevel sets the root logger level
func SetLevel(level Level) {
	root.SetLevel(level)
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
	context := newLoggingContext(config)
	logger, err := newLogger(context, nil, "")
	if err != nil {
		return err
	}
	root = logger
	return nil
}

func newLogger(context *loggingContext, parent *dazlLogger, name string) (*dazlLogger, error) {
	var path []string
	if parent != nil {
		path = append(parent.path, name)
	}
	loggerName := strings.Join(path, pathSep)

	// Create the child logger.
	logger := &dazlLogger{
		loggerContext: &loggerContext{
			loggingContext: context,
			name:           loggerName,
			path:           path,
		},
	}

	// Configure the logger's outputs from the configuration
	loggerConfig, _ := context.config.getLogger(loggerName)
	parentOutputs := make(map[string]int)
	if parent != nil {
		for i, output := range parent.outputs {
			parentOutputs[output.name] = i
		}

		defaultLevel := Level(parent.defaultLevel.Load())
		parentLevel := Level(parent.level.Load())
		if parentLevel != EmptyLevel {
			defaultLevel = parentLevel
		}
		logger.defaultLevel.Store(int32(defaultLevel))
	}
	if loggerConfig.Level != nil {
		logger.SetLevel(loggerConfig.Level.Level())
	}

	for outputName, outputConfig := range loggerConfig.getOutputs() {
		// If the configured output already exists, create a child output that inherits the
		// parent's level and sampling configuration. Otherwise, create a new output.
		var output *dazlOutput
		if i, ok := parentOutputs[outputName]; ok {
			parentOutput := parent.outputs[i]
			if outputConfig.Writer == "" {
				output = parentOutput.WithWriter(parentOutput.writer.WithName(loggerName))
			} else {
				writer, err := context.getWriter(outputConfig.Writer)
				if err != nil {
					return nil, err
				}
				output = parentOutput.WithWriter(writer.WithName(loggerName))
			}
			level := outputConfig.Level.Level()
			if level != EmptyLevel {
				output = output.WithLevel(level)
			}
		} else {
			writer, err := context.getWriter(outputConfig.Writer)
			if err != nil {
				return nil, err
			}
			output = newOutput(outputName, writer.WithName(loggerName), outputConfig.Level.Level(), &allSampler{})
		}

		// Configure sampling for the output, first adding the logger sampling and then
		// output sampling configurations.
		var err error
		output, err = configureSampling(output, loggerConfig.Sample)
		if err != nil {
			return nil, err
		}
		output, err = configureSampling(output, outputConfig.Sample)
		if err != nil {
			return nil, err
		}
		logger.outputs = append(logger.outputs, output)
	}
	return logger, nil
}

func newLoggingContext(config loggingConfig) *loggingContext {
	return &loggingContext{
		config: config,
	}
}

type loggingContext struct {
	config  loggingConfig
	writers sync.Map
	mu      sync.Mutex
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

	switch name {
	case "stdout":
		if c.config.Writers.Stdout == nil {
			return nil, fmt.Errorf("'%s' writer is not configured", name)
		}
		writer, err := getFramework().NewWriter(os.Stdout, c.config.Writers.Stdout.Encoder)
		if err != nil {
			return nil, err
		}
		c.writers.Store(name, writer)
		return writer, nil
	case "stderr":
		if c.config.Writers.Stderr == nil {
			return nil, fmt.Errorf("'%s' writer is not configured", name)
		}
		writer, err := getFramework().NewWriter(os.Stderr, c.config.Writers.Stderr.Encoder)
		if err != nil {
			return nil, err
		}
		c.writers.Store(name, writer)
		return writer, nil
	default:
		config, ok := c.config.Writers.getFile(name)
		if !ok {
			return nil, fmt.Errorf("'%s' writer is not configured", name)
		}
		file, err := os.OpenFile(config.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return nil, err
		}
		writer, err := getFramework().NewWriter(file, config.Encoder)
		if err != nil {
			return nil, err
		}
		c.writers.Store(name, writer)
		return writer, nil
	}
}

type loggerContext struct {
	*loggingContext
	name         string
	path         []string
	children     sync.Map
	mu           sync.Mutex
	level        atomic.Int32
	defaultLevel atomic.Int32
}

func (l *loggerContext) Level() Level {
	level := Level(l.level.Load())
	if level != EmptyLevel {
		return level
	}
	return Level(l.defaultLevel.Load())
}

func (l *loggerContext) SetLevel(level Level) {
	l.level.Store(int32(level))
	l.children.Range(func(key, value any) bool {
		value.(*dazlLogger).setDefaultLevel(level)
		return true
	})
}

func (l *loggerContext) setDefaultLevel(level Level) {
	l.defaultLevel.Store(int32(level))
	if Level(l.level.Load()) == EmptyLevel {
		l.children.Range(func(key, value any) bool {
			value.(*dazlLogger).setDefaultLevel(level)
			return true
		})
	}
}

type dazlLogger struct {
	*loggerContext
	outputs []*dazlOutput
}

func (l *dazlLogger) Name() string {
	return l.name
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
	outputs := make([]*dazlOutput, len(l.outputs))
	for i, output := range l.outputs {
		writer := output.Writer()
		var err error
		for _, field := range fields {
			if writer, err = field(writer); err != nil {
				panic(err)
			}
		}
		outputs[i] = output.WithWriter(writer)
	}
	return &dazlLogger{
		loggerContext: l.loggerContext,
		outputs:       outputs,
	}
}

func (l *dazlLogger) WithSkipCalls(calls int) Logger {
	outputs := make([]*dazlOutput, len(l.outputs))
	for i, output := range l.outputs {
		if writer, ok := output.Writer().(CallSkippingWriter); ok {
			outputs[i] = output.WithWriter(writer.WithSkipCalls(calls))
		} else {
			outputs[i] = output
		}
	}
	return &dazlLogger{
		loggerContext: l.loggerContext,
		outputs:       outputs,
	}
}

func (l *dazlLogger) Debug(args ...interface{}) {
	if l.Level().Enabled(DebugLevel) {
		for _, output := range l.outputs {
			output.Debug(fmt.Sprint(args...))
		}
	}
}

func (l *dazlLogger) Debugf(format string, args ...interface{}) {
	if l.Level().Enabled(DebugLevel) {
		for _, output := range l.outputs {
			output.Debug(fmt.Sprintf(format, args...))
		}
	}
}

func (l *dazlLogger) Debugw(msg string, fields ...Field) {
	l.WithFields(fields...).Debug(msg)
}

func (l *dazlLogger) Info(args ...interface{}) {
	if l.Level().Enabled(InfoLevel) {
		for _, output := range l.outputs {
			output.Info(fmt.Sprint(args...))
		}
	}
}

func (l *dazlLogger) Infof(format string, args ...interface{}) {
	if l.Level().Enabled(InfoLevel) {
		for _, output := range l.outputs {
			output.Info(fmt.Sprintf(format, args...))
		}
	}
}

func (l *dazlLogger) Infow(msg string, fields ...Field) {
	l.WithFields(fields...).Info(msg)
}

func (l *dazlLogger) Warn(args ...interface{}) {
	if l.Level().Enabled(WarnLevel) {
		for _, output := range l.outputs {
			output.Warn(fmt.Sprint(args...))
		}
	}
}

func (l *dazlLogger) Warnf(format string, args ...interface{}) {
	if l.Level().Enabled(WarnLevel) {
		for _, output := range l.outputs {
			output.Warn(fmt.Sprintf(format, args...))
		}
	}
}

func (l *dazlLogger) Warnw(msg string, fields ...Field) {
	l.WithFields(fields...).Warn(msg)
}

func (l *dazlLogger) Error(args ...interface{}) {
	if l.Level().Enabled(ErrorLevel) {
		for _, output := range l.outputs {
			output.Error(fmt.Sprint(args...))
		}
	}
}

func (l *dazlLogger) Errorf(format string, args ...interface{}) {
	if l.Level().Enabled(ErrorLevel) {
		for _, output := range l.outputs {
			output.Error(fmt.Sprintf(format, args...))
		}
	}
}

func (l *dazlLogger) Errorw(msg string, fields ...Field) {
	l.WithFields(fields...).Error(msg)
}

func (l *dazlLogger) Fatal(args ...interface{}) {
	if l.Level().Enabled(FatalLevel) {
		for _, output := range l.outputs {
			output.Fatal(fmt.Sprint(args...))
		}
	}
}

func (l *dazlLogger) Fatalf(format string, args ...interface{}) {
	if l.Level().Enabled(FatalLevel) {
		for _, output := range l.outputs {
			output.Fatal(fmt.Sprintf(format, args...))
		}
	}
}

func (l *dazlLogger) Fatalw(msg string, fields ...Field) {
	l.WithFields(fields...).Fatal(msg)
}

func (l *dazlLogger) Panic(args ...interface{}) {
	if l.Level().Enabled(PanicLevel) {
		for _, output := range l.outputs {
			output.Panic(fmt.Sprint(args...))
		}
	}
}

func (l *dazlLogger) Panicf(format string, args ...interface{}) {
	if l.Level().Enabled(PanicLevel) {
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
	Level   *levelConfig            `json:"level" yaml:"level"`
	Sample  *samplingConfig         `json:"sample" yaml:"sample"`
	Outputs map[string]outputConfig `json:"outputs" yaml:"outputs"`
}

func (c *loggerConfig) getOutputs() map[string]outputConfig {
	if c.Outputs == nil {
		return map[string]outputConfig{}
	}
	return c.Outputs
}

func (c *loggerConfig) getOutput(name string) (outputConfig, bool) {
	config, ok := c.getOutputs()[name]
	return config, ok
}

func configureSampling(output *dazlOutput, config *samplingConfig) (*dazlOutput, error) {
	if config == nil {
		return output, nil
	}

	// If sampling is configured for the output, add the sampler to the output's underlying writer
	// and remove any configured samplers from the output.
	// If the writer does not support the configured sampling strategy, add a sampler to the output.
	if config.Basic != nil {
		if samplingWriter, ok := output.Writer().(BasicSamplingWriter); ok {
			writer, err := samplingWriter.WithBasicSampler(config.Basic.Interval, config.Basic.Level.Level())
			if err != nil {
				return nil, err
			}
			output = output.WithWriter(writer).WithSampler(&allSampler{})
		} else {
			output = output.WithSampler(&basicSampler{
				Interval: uint32(config.Basic.Interval),
				Level:    config.Basic.Level.Level(),
			})
		}
	} else if config.Random != nil {
		if samplingWriter, ok := output.Writer().(RandomSamplingWriter); ok {
			writer, err := samplingWriter.WithRandomSampler(config.Random.Interval, config.Random.Level.Level())
			if err != nil {
				return nil, err
			}
			output = output.WithWriter(writer).WithSampler(&allSampler{})
		} else {
			output = output.WithSampler(randomSampler{
				Interval: config.Random.Interval,
				Level:    config.Random.Level.Level(),
			})
		}
	}
	return output, nil
}
