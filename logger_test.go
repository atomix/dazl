// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package dazl

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
	"os"
	"testing"
)

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
	assert.Equal(t, "INFO\tdazl/logger_test.go:48\tshould be printed\n", buf.String())

	// The "test" logger should inherit the INFO level from the root logger
	buf.Reset()
	logger = root.GetLogger("test").WithFields(Bool("printed", true))
	assert.Equal(t, "test", logger.Name())
	assert.Equal(t, InfoLevel, logger.Level())
	logger.Debugf("should %s be", "not")
	assert.Equal(t, "", buf.String())
	logger.Info("should be")
	assert.Equal(t, "INFO\ttest\tdazl/logger_test.go:58\tshould be\t{\"printed\": true}\n", buf.String())

	// The "test/1" logger should be configured with DEBUG level
	buf.Reset()
	logger = root.GetLogger("test/1")
	assert.Equal(t, "test/1", logger.Name())
	assert.Equal(t, DebugLevel, logger.Level())
	logger.Debugw("should be", Bool("printed", true))
	assert.Equal(t, "DEBUG\ttest/1\tdazl/logger_test.go:66\tshould be\t{\"printed\": true}\n", buf.String())
	logger.Infow("should be", Bool("printed", true))
	assert.Equal(t, "DEBUG\ttest/1\tdazl/logger_test.go:66\tshould be\t{\"printed\": true}\nINFO\ttest/1\tdazl/logger_test.go:68\tshould be\t{\"printed\": true}\n", buf.String())

	// The "test/1/2" logger should inherit the DEBUG level from "test/1"
	buf.Reset()
	logger = root.GetLogger("test/1/2").WithFields(Bool("printed", true))
	assert.Equal(t, "test/1/2", logger.Name())
	assert.Equal(t, DebugLevel, logger.Level())
	logger.Debugw("printed", String("should", "be"))
	assert.Equal(t, "DEBUG\ttest/1/2\tdazl/logger_test.go:76\tprinted\t{\"printed\": true, \"should\": \"be\"}\n", buf.String())
	logger.Infow("printed", String("should", "be"))
	assert.Equal(t, "DEBUG\ttest/1/2\tdazl/logger_test.go:76\tprinted\t{\"printed\": true, \"should\": \"be\"}\nINFO\ttest/1/2\tdazl/logger_test.go:78\tprinted\t{\"printed\": true, \"should\": \"be\"}\n", buf.String())

	// The "test" logger should still inherit the INFO level from the root logger
	buf.Reset()
	logger = root.GetLogger("test")
	assert.Equal(t, "test", logger.Name())
	assert.Equal(t, InfoLevel, logger.Level())
	logger.Debug("should not be printed")
	assert.Equal(t, "", buf.String())
	logger.Info("should be printed")
	assert.Equal(t, "INFO\ttest\tdazl/logger_test.go:88\tshould be printed\n", buf.String())

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
	assert.Equal(t, "WARN\ttest/2\tdazl/logger_test.go:100\tshould be printed\t{\"times\": 2}\n", buf.String())

	// The "test/2/3" logger should be configured with INFO level
	buf.Reset()
	logger = root.GetLogger("test/2/3")
	assert.Equal(t, "test/2/3", logger.Name())
	assert.Equal(t, InfoLevel, logger.Level())
	logger.Debug("should not be printed")
	assert.Equal(t, "", buf.String())
	logger.Infow("should be printed", Int("times", 2))
	assert.Equal(t, "INFO\ttest/2/3\tdazl/logger_test.go:110\tshould be printed\t{\"times\": 2}\n", buf.String())
	logger.Warn("should be printed twice")
	assert.Equal(t, "INFO\ttest/2/3\tdazl/logger_test.go:110\tshould be printed\t{\"times\": 2}\nWARN\ttest/2/3\tdazl/logger_test.go:112\tshould be printed twice\n", buf.String())

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
	assert.Equal(t, "WARN\ttest/2/4\tdazl/logger_test.go:124\tshould be printed twice\n", buf.String())

	// The "test/2" logger level should be changed to DEBUG
	buf.Reset()
	logger = root.GetLogger("test/2")
	assert.Equal(t, "test/2", logger.Name())
	logger.SetLevel(DebugLevel)
	assert.Equal(t, DebugLevel, logger.Level())
	logger.Debugw("should be", Bool("printed", true))
	assert.Equal(t, "DEBUG\ttest/2\tdazl/logger_test.go:133\tshould be\t{\"printed\": true}\n", buf.String())
	logger.Infow("should be printed", Int("times", 2))
	assert.Equal(t, "DEBUG\ttest/2\tdazl/logger_test.go:133\tshould be\t{\"printed\": true}\nINFO\ttest/2\tdazl/logger_test.go:135\tshould be printed\t{\"times\": 2}\n", buf.String())
	logger.Warn("should be printed twice")
	assert.Equal(t, "DEBUG\ttest/2\tdazl/logger_test.go:133\tshould be\t{\"printed\": true}\nINFO\ttest/2\tdazl/logger_test.go:135\tshould be printed\t{\"times\": 2}\nWARN\ttest/2\tdazl/logger_test.go:137\tshould be printed twice\n", buf.String())

	// The "test/2/3" logger should not inherit the change to the "test/2" logger since its level has been explicitly set
	buf.Reset()
	logger = root.GetLogger("test/2/3")
	assert.Equal(t, "test/2/3", logger.Name())
	assert.Equal(t, InfoLevel, logger.Level())
	logger.Debug("should not be printed")
	assert.Equal(t, "", buf.String())
	logger.Info("should be printed twice")
	assert.Equal(t, "INFO\ttest/2/3\tdazl/logger_test.go:147\tshould be printed twice\n", buf.String())
	logger.Warn("should be printed twice")
	assert.Equal(t, "INFO\ttest/2/3\tdazl/logger_test.go:147\tshould be printed twice\nWARN\ttest/2/3\tdazl/logger_test.go:149\tshould be printed twice\n", buf.String())

	// The "test/2/4" logger should inherit the change to the "test/2" logger since its level has not been explicitly set
	// The "test/2/4" logger should not output DEBUG messages since the output level is explicitly set to WARN
	buf.Reset()
	logger = root.GetLogger("test/2/4")
	assert.Equal(t, "test/2/4", logger.Name())
	assert.Equal(t, DebugLevel, logger.Level())
	logger.Debug("should be printed")
	assert.Equal(t, "DEBUG\ttest/2/4\tdazl/logger_test.go:158\tshould be printed\n", buf.String())
	logger.Info("should be printed twice")
	assert.Equal(t, "DEBUG\ttest/2/4\tdazl/logger_test.go:158\tshould be printed\nINFO\ttest/2/4\tdazl/logger_test.go:160\tshould be printed twice\n", buf.String())
	logger.Warn("should be printed twice")
	assert.Equal(t, "DEBUG\ttest/2/4\tdazl/logger_test.go:158\tshould be printed\nINFO\ttest/2/4\tdazl/logger_test.go:160\tshould be printed twice\nWARN\ttest/2/4\tdazl/logger_test.go:162\tshould be printed twice\n", buf.String())

	// The "test/3" logger should be configured with INFO level
	// The "test/3" logger should write to multiple outputs
	buf.Reset()
	logger = root.GetLogger("test/3")
	assert.Equal(t, "test/3", logger.Name())
	assert.Equal(t, InfoLevel, logger.Level())
	logger.Debug("should not be printed")
	assert.Equal(t, "", buf.String())
	logger.Info("should be printed")
	assert.Equal(t, "INFO\ttest/3\tdazl/logger_test.go:173\tshould be printed\n", buf.String())
	logger.Warn("should be printed twice")
	assert.Equal(t, "INFO\ttest/3\tdazl/logger_test.go:173\tshould be printed\nWARN\ttest/3\tdazl/logger_test.go:175\tshould be printed twice\n", buf.String())

	// The "test/3/4" logger should inherit INFO level from "test/3"
	// The "test/3/4" logger should inherit multiple outputs from "test/3"
	buf.Reset()
	logger = root.GetLogger("test/3/4")
	assert.Equal(t, "test/3/4", logger.Name())
	assert.Equal(t, InfoLevel, logger.Level())
	logger.Debug("should not be printed")
	assert.Equal(t, "", buf.String())
	logger.Info("should be printed")
	assert.Equal(t, "INFO\ttest/3/4\tdazl/logger_test.go:186\tshould be printed\n", buf.String())
	logger.Warn("should be printed twice")
	assert.Equal(t, "INFO\ttest/3/4\tdazl/logger_test.go:186\tshould be printed\nWARN\ttest/3/4\tdazl/logger_test.go:188\tshould be printed twice\n", buf.String())

	//logger = GetLogger("test/kafka")
	//assert.Equal(t, InfoLevel, logger.Level())
}
