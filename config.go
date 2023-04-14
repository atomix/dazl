// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package dazl

import (
	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

const configFile = "logging.yaml"

type loggingConfig struct {
	Encoders   encodersConfig          `json:"encoders" yaml:"encoders"`
	Writers    writersConfig           `json:"writers" yaml:"writers"`
	RootLogger loggerConfig            `json:"rootLogger" yaml:"rootLogger"`
	Loggers    map[string]loggerConfig `json:"loggers" yaml:"loggers"`
}

func (c *loggingConfig) getLoggers() map[string]loggerConfig {
	if c.Loggers == nil {
		return map[string]loggerConfig{}
	}
	return c.Loggers
}

func (c *loggingConfig) getLogger(name string) (loggerConfig, bool) {
	config, ok := c.getLoggers()[name]
	return config, ok
}

func init() {
	var config loggingConfig
	if err := load(&config); err != nil {
		panic(err)
	} else if err := configure(config); err != nil {
		panic(err)
	}
}

// load the dazl configuration
func load(config *loggingConfig) error {
	bytes, err := os.ReadFile(configFile)
	if err == nil {
		return yaml.Unmarshal(bytes, config)
	} else if !os.IsNotExist(err) {
		return err
	}

	if home, err := homedir.Dir(); err == nil {
		bytes, err = os.ReadFile(filepath.Join(home, configFile))
		if err == nil {
			return yaml.Unmarshal(bytes, config)
		} else if !os.IsNotExist(err) {
			return err
		}
	}

	bytes, err = os.ReadFile(filepath.Join("/etc/dazl", configFile))
	if err == nil {
		return yaml.Unmarshal(bytes, config)
	} else if !os.IsNotExist(err) {
		return err
	}
	return nil
}
