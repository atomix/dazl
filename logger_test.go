// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package dazl

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
	"os"
	"testing"
	"time"
)

func TestLoggerPackage(t *testing.T) {
	assert.Equal(t, "github.com/atomix/dazl", GetLogger().Name())
	SetLevel(InfoLevel)
	assert.Equal(t, InfoLevel, GetLogger().Level())
}

func TestLoggerConfig(t *testing.T) {
	config := Config{}
	data, err := os.ReadFile("test-data/test-config.yaml")
	assert.NoError(t, err)
	err = yaml.Unmarshal(data, &config)
	assert.NoError(t, err)
	err = configure(config)
	assert.NoError(t, err)

	buf := &bytes.Buffer{}
	sink, err := NewSink(buf,
		WithEncoding(ConsoleEncoding),
		WithNameKey("name"),
		WithMessageKey("message"),
		WithLevelKey("level"),
		WithCallerKey("caller"),
		WithNameEncoder(zapcore.FullNameEncoder),
		WithCallerEncoder(zapcore.ShortCallerEncoder),
		WithLevelEncoder(zapcore.CapitalLevelEncoder))
	assert.NoError(t, err)
	assert.NotNil(t, sink)

	root, err := newLogger(config)
	assert.NoError(t, err)
	root = root.WithOutputs(NewOutput(sink))

	// The root logger should be configured with INFO level
	logger := root
	assert.Equal(t, "", logger.Name())
	assert.Equal(t, InfoLevel, logger.Level())
	logger.Debug("should not be printed")
	assert.Equal(t, "", buf.String())
	logger.Infof("should be %s", "printed")
	assert.Equal(t, "INFO\tdazl/logger_test.go:56\tshould be printed\n", buf.String())
	logger.WithSkipCalls(2).Infof("should be %s again", "printed")
	assert.NoError(t, logger.Sync())
	assert.NotEqual(t, "INFO\tdazl/logger_test.go:56\tshould be printed\n", buf.String())

	// The "test" logger should inherit the INFO level from the root logger
	buf.Reset()
	logger = root.GetLogger("test").WithFields(Bool("printed", true))
	assert.Equal(t, "test", logger.Name())
	assert.Equal(t, InfoLevel, logger.Level())
	logger.Debugf("should %s be", "not")
	assert.Equal(t, "", buf.String())
	logger.Info("should be")
	assert.Equal(t, "INFO\ttest\tdazl/logger_test.go:69\tshould be\t{\"printed\": true}\n", buf.String())

	// The "test/1" logger should be configured with DEBUG level
	buf.Reset()
	logger = root.GetLogger("test/1")
	assert.Equal(t, "test/1", logger.Name())
	assert.Equal(t, DebugLevel, logger.Level())
	logger.Debugw("should be", Bool("printed", true))
	assert.Equal(t, "DEBUG\ttest/1\tdazl/logger_test.go:77\tshould be\t{\"printed\": true}\n", buf.String())
	logger.Infow("should be", Bool("printed", true))
	assert.Equal(t, "DEBUG\ttest/1\tdazl/logger_test.go:77\tshould be\t{\"printed\": true}\nINFO\ttest/1\tdazl/logger_test.go:79\tshould be\t{\"printed\": true}\n", buf.String())

	// The "test/1/2" logger should inherit the DEBUG level from "test/1"
	buf.Reset()
	logger = root.GetLogger("test/1/2").WithFields(Bool("printed", true))
	assert.Equal(t, "test/1/2", logger.Name())
	assert.Equal(t, DebugLevel, logger.Level())
	logger.Debugw("printed", String("should", "be"))
	assert.Equal(t, "DEBUG\ttest/1/2\tdazl/logger_test.go:87\tprinted\t{\"printed\": true, \"should\": \"be\"}\n", buf.String())
	logger.Infow("printed", String("should", "be"))
	assert.Equal(t, "DEBUG\ttest/1/2\tdazl/logger_test.go:87\tprinted\t{\"printed\": true, \"should\": \"be\"}\nINFO\ttest/1/2\tdazl/logger_test.go:89\tprinted\t{\"printed\": true, \"should\": \"be\"}\n", buf.String())

	// The "test" logger should still inherit the INFO level from the root logger
	buf.Reset()
	logger = root.GetLogger("test")
	assert.Equal(t, "test", logger.Name())
	assert.Equal(t, InfoLevel, logger.Level())
	logger.Debug("should not be printed")
	assert.Equal(t, "", buf.String())
	logger.Info("should be printed")
	assert.Equal(t, "INFO\ttest\tdazl/logger_test.go:99\tshould be printed\n", buf.String())

	// The "test/2" logger should be configured with WARN level
	buf.Reset()
	logger = root.GetLogger("test/2")
	assert.Equal(t, "test/2", logger.Name())
	assert.Equal(t, WarnLevel, logger.Level())
	logger.Debug("should not be printed")
	assert.Equal(t, "", buf.String())
	logger.Infow("should not be", Bool("printed", true))
	assert.Equal(t, "", buf.String())
	logger.Warnw("should be printed", Int("times", 2))
	assert.Equal(t, "WARN\ttest/2\tdazl/logger_test.go:111\tshould be printed\t{\"times\": 2}\n", buf.String())

	// The "test/2/3" logger should be configured with INFO level
	buf.Reset()
	logger = root.GetLogger("test/2/3")
	assert.Equal(t, "test/2/3", logger.Name())
	assert.Equal(t, InfoLevel, logger.Level())
	logger.Debug("should not be printed")
	assert.Equal(t, "", buf.String())
	logger.Infow("should be printed", Int("times", 2))
	assert.Equal(t, "INFO\ttest/2/3\tdazl/logger_test.go:121\tshould be printed\t{\"times\": 2}\n", buf.String())
	logger.Warn("should be printed twice")
	assert.Equal(t, "INFO\ttest/2/3\tdazl/logger_test.go:121\tshould be printed\t{\"times\": 2}\nWARN\ttest/2/3\tdazl/logger_test.go:123\tshould be printed twice\n", buf.String())

	// The "test/2/4" logger should inherit the WARN level from "test/2"
	buf.Reset()
	logger = root.GetLogger("test/2/4")
	assert.Equal(t, "test/2/4", logger.Name())
	assert.Equal(t, WarnLevel, logger.Level())
	logger.Debug("should not be printed")
	assert.Equal(t, "", buf.String())
	logger.Info("should not be printed")
	assert.Equal(t, "", buf.String())
	logger.Warn("should be printed twice")
	assert.Equal(t, "WARN\ttest/2/4\tdazl/logger_test.go:135\tshould be printed twice\n", buf.String())

	// The "test/2" logger level should be changed to DEBUG
	buf.Reset()
	logger = root.GetLogger("test/2")
	assert.Equal(t, "test/2", logger.Name())
	logger.SetLevel(DebugLevel)
	assert.Equal(t, DebugLevel, logger.Level())
	logger.Debugw("should be", Bool("printed", true))
	assert.Equal(t, "DEBUG\ttest/2\tdazl/logger_test.go:144\tshould be\t{\"printed\": true}\n", buf.String())
	logger.Infow("should be printed", Int("times", 2))
	assert.Equal(t, "DEBUG\ttest/2\tdazl/logger_test.go:144\tshould be\t{\"printed\": true}\nINFO\ttest/2\tdazl/logger_test.go:146\tshould be printed\t{\"times\": 2}\n", buf.String())
	logger.Warn("should be printed twice")
	assert.Equal(t, "DEBUG\ttest/2\tdazl/logger_test.go:144\tshould be\t{\"printed\": true}\nINFO\ttest/2\tdazl/logger_test.go:146\tshould be printed\t{\"times\": 2}\nWARN\ttest/2\tdazl/logger_test.go:148\tshould be printed twice\n", buf.String())

	// The "test/2/3" logger should not inherit the change to the "test/2" logger since its level has been explicitly set
	buf.Reset()
	logger = root.GetLogger("test/2/3")
	assert.Equal(t, "test/2/3", logger.Name())
	assert.Equal(t, InfoLevel, logger.Level())
	logger.Debug("should not be printed")
	assert.Equal(t, "", buf.String())
	logger.Info("should be printed twice")
	assert.Equal(t, "INFO\ttest/2/3\tdazl/logger_test.go:158\tshould be printed twice\n", buf.String())
	logger.Warn("should be printed twice")
	assert.Equal(t, "INFO\ttest/2/3\tdazl/logger_test.go:158\tshould be printed twice\nWARN\ttest/2/3\tdazl/logger_test.go:160\tshould be printed twice\n", buf.String())

	// The "test/2/4" logger should inherit the change to the "test/2" logger since its level has not been explicitly set
	// The "test/2/4" logger should not output DEBUG messages since the output level is explicitly set to WARN
	buf.Reset()
	logger = root.GetLogger("test/2/4")
	assert.Equal(t, "test/2/4", logger.Name())
	assert.Equal(t, DebugLevel, logger.Level())
	logger.Debug("should be printed")
	assert.Equal(t, "DEBUG\ttest/2/4\tdazl/logger_test.go:169\tshould be printed\n", buf.String())
	logger.Info("should be printed twice")
	assert.Equal(t, "DEBUG\ttest/2/4\tdazl/logger_test.go:169\tshould be printed\nINFO\ttest/2/4\tdazl/logger_test.go:171\tshould be printed twice\n", buf.String())
	logger.Warn("should be printed twice")
	assert.Equal(t, "DEBUG\ttest/2/4\tdazl/logger_test.go:169\tshould be printed\nINFO\ttest/2/4\tdazl/logger_test.go:171\tshould be printed twice\nWARN\ttest/2/4\tdazl/logger_test.go:173\tshould be printed twice\n", buf.String())

	// The "test/3" logger should be configured with INFO level
	// The "test/3" logger should write to multiple outputs
	buf.Reset()
	logger = root.GetLogger("test/3")
	assert.Equal(t, "test/3", logger.Name())
	assert.Equal(t, InfoLevel, logger.Level())
	logger.Debug("should not be printed")
	assert.Equal(t, "", buf.String())
	logger.Info("should be printed")
	assert.Equal(t, "INFO\ttest/3\tdazl/logger_test.go:184\tshould be printed\n", buf.String())
	logger.Warn("should be printed twice")
	assert.Equal(t, "INFO\ttest/3\tdazl/logger_test.go:184\tshould be printed\nWARN\ttest/3\tdazl/logger_test.go:186\tshould be printed twice\n", buf.String())

	// The "test/3/4" logger should inherit INFO level from "test/3"
	// The "test/3/4" logger should inherit multiple outputs from "test/3"
	buf.Reset()
	logger = root.GetLogger("test/3/4")
	assert.Equal(t, "test/3/4", logger.Name())
	assert.Equal(t, InfoLevel, logger.Level())
	logger.Debug("should not be printed")
	assert.Equal(t, "", buf.String())
	logger.Info("should be printed")
	assert.Equal(t, "INFO\ttest/3/4\tdazl/logger_test.go:197\tshould be printed\n", buf.String())
	logger.Warn("should be printed twice")
	assert.Equal(t, "INFO\ttest/3/4\tdazl/logger_test.go:197\tshould be printed\nWARN\ttest/3/4\tdazl/logger_test.go:199\tshould be printed twice\n", buf.String())

	//logger = GetLogger("test/kafka")
	//assert.Equal(t, InfoLevel, logger.Level())
}

func TestLoggerFuncs(t *testing.T) {
	buf := &bytes.Buffer{}
	sink, err := NewSink(buf,
		WithEncoding(ConsoleEncoding),
		WithMessageKey("message"),
		WithLevelKey("level"),
		WithLevelEncoder(zapcore.LowercaseLevelEncoder),
		WithTimeEncoder(zapcore.EpochNanosTimeEncoder),
		WithDurationEncoder(zapcore.StringDurationEncoder),
		WithLineEnding(""))

	assert.NoError(t, err)
	assert.NotNil(t, sink)

	now := time.Now()
	duration := time.Second

	str := "foo"
	var i = 1
	var i32 int32 = 2
	var i64 int64 = 3
	var u uint = 4
	var u32 uint32 = 5
	var u64 uint64 = 6
	var b = true

	fields := []Field{
		String("string", str),
		Stringp("stringp", &str),
		Strings("strings", []string{str}),
		Int("int", i),
		Intp("intp", &i),
		Ints("ints", []int{i}),
		Int32("int32", i32),
		Int32p("int32p", &i32),
		Int32s("int32s", []int32{i32}),
		Int64("int64", i64),
		Int64p("int64p", &i64),
		Int64s("int64s", []int64{i64}),
		Uint("uint", u),
		Uintp("uintp", &u),
		Uints("uints", []uint{u}),
		Uint32("uint32", u32),
		Uint32p("uint32p", &u32),
		Uint32s("uint32s", []uint32{u32}),
		Uint64("uint64", u64),
		Uint64p("uint64p", &u64),
		Uint64s("uint64s", []uint64{u64}),
		Bool("bool", b),
		Boolp("boolp", &b),
		Bools("bools", []bool{b}),
		Time("time", now),
		Timep("timep", &now),
		Times("times", []time.Time{now}),
		Duration("duration", duration),
		Durationp("durationp", &duration),
		Durations("durations", []time.Duration{duration}),
	}
	fieldsString := fmt.Sprintf("{\"string\": \"foo\", \"stringp\": \"foo\", \"strings\": [\"foo\"], \"int\": 1, \"intp\": 1, \"ints\": [1], \"int32\": 2, \"int32p\": 2, \"int32s\": [2], \"int64\": 3, \"int64p\": 3, \"int64s\": [3], \"uint\": 4, \"uintp\": 4, \"uints\": [4], \"uint32\": 5, \"uint32p\": 5, \"uint32s\": [5], \"uint64\": 6, \"uint64p\": 6, \"uint64s\": [6], \"bool\": true, \"boolp\": true, \"bools\": [true], \"time\": %d, \"timep\": %d, \"times\": [%d], \"duration\": \"%s\", \"durationp\": \"%s\", \"durations\": [\"%s\"]}", now.UnixNano(), now.UnixNano(), now.UnixNano(), duration, duration, duration)

	log := GetLogger().WithOutputs(NewOutput(sink))
	log.SetLevel(DebugLevel)
	assert.Equal(t, DebugLevel, log.Level())

	buf.Reset()
	log.Debug("foo")
	assert.Equal(t, "debug\tfoo\n", buf.String())
	buf.Reset()
	log.Debugf("foo %d", 1)
	assert.Equal(t, "debug\tfoo 1\n", buf.String())
	buf.Reset()
	log.Debugw("foo", fields...)
	assert.Equal(t, fmt.Sprintf("debug\tfoo\t%s\n", fieldsString), buf.String())
	buf.Reset()
	log.WithFields(fields...).Debug("foo")
	assert.Equal(t, fmt.Sprintf("debug\tfoo\t%s\n", fieldsString), buf.String())

	buf.Reset()
	log.Info("foo")
	assert.Equal(t, "info\tfoo\n", buf.String())
	buf.Reset()
	log.Infof("foo %d", 1)
	assert.Equal(t, "info\tfoo 1\n", buf.String())
	buf.Reset()
	log.Infow("foo", fields...)
	assert.Equal(t, fmt.Sprintf("info\tfoo\t%s\n", fieldsString), buf.String())
	buf.Reset()
	log.WithFields(fields...).Info("foo")
	assert.Equal(t, fmt.Sprintf("info\tfoo\t%s\n", fieldsString), buf.String())

	buf.Reset()
	log.Warn("foo")
	assert.Equal(t, "warn\tfoo\n", buf.String())
	buf.Reset()
	log.Warnf("foo %d", 1)
	assert.Equal(t, "warn\tfoo 1\n", buf.String())
	buf.Reset()
	log.Warnw("foo", fields...)
	assert.Equal(t, fmt.Sprintf("warn\tfoo\t%s\n", fieldsString), buf.String())
	buf.Reset()
	log.WithFields(fields...).Warn("foo")
	assert.Equal(t, fmt.Sprintf("warn\tfoo\t%s\n", fieldsString), buf.String())

	buf.Reset()
	log.Error("foo")
	assert.Equal(t, "error\tfoo\n", buf.String())
	buf.Reset()
	log.Errorf("foo %d", 1)
	assert.Equal(t, "error\tfoo 1\n", buf.String())
	buf.Reset()
	log.Errorw("foo", fields...)
	assert.Equal(t, fmt.Sprintf("error\tfoo\t%s\n", fieldsString), buf.String())
	buf.Reset()
	log.WithFields(fields...).Error("foo")
	assert.Equal(t, fmt.Sprintf("error\tfoo\t%s\n", fieldsString), buf.String())
}

func TestSetLevel(t *testing.T) {
	root := GetLogger("")
	parent := GetLogger("parent")
	child := parent.GetLogger("child")
	clone := parent.WithOutputs()

	assert.Equal(t, EmptyLevel, root.Level())
	assert.Equal(t, EmptyLevel, parent.Level())
	assert.Equal(t, EmptyLevel, child.Level())
	assert.Equal(t, EmptyLevel, clone.Level())

	SetLevel(ErrorLevel)
	assert.Equal(t, ErrorLevel, root.Level())
	assert.Equal(t, ErrorLevel, parent.Level())
	assert.Equal(t, ErrorLevel, child.Level())
	assert.Equal(t, ErrorLevel, clone.Level())

	parent.SetLevel(WarnLevel)
	assert.Equal(t, ErrorLevel, root.Level())
	assert.Equal(t, WarnLevel, parent.Level())
	assert.Equal(t, WarnLevel, child.Level())
	assert.Equal(t, WarnLevel, clone.Level())

	child.SetLevel(ErrorLevel)
	assert.Equal(t, ErrorLevel, root.Level())
	assert.Equal(t, WarnLevel, parent.Level())
	assert.Equal(t, ErrorLevel, child.Level())
	assert.Equal(t, WarnLevel, clone.Level())

	clone.SetLevel(DebugLevel)
	assert.Equal(t, ErrorLevel, root.Level())
	assert.Equal(t, DebugLevel, parent.Level())
	assert.Equal(t, ErrorLevel, child.Level())
	assert.Equal(t, DebugLevel, clone.Level())
}
