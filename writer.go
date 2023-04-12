// SPDX-FileCopyrightText: 2023-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package dazl

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

type writersConfig struct {
	Stdout *stdoutWriterConfig         `json:"stdout" yaml:"stdout"`
	Stderr *stderrWriterConfig         `json:"stderr" yaml:"stderr"`
	Files  map[string]fileWriterConfig `json:"files" yaml:"files"`
}

func (c *writersConfig) UnmarshalYAML(unmarshal func(any) error) error {
	writers := make(map[string]any)
	if err := unmarshal(writers); err != nil {
		return err
	}

	c.Files = make(map[string]fileWriterConfig)
	for name, config := range writers {
		switch name {
		case "stdout":
			text, err := yaml.Marshal(config)
			if err != nil {
				return err
			}
			writer := &stdoutWriterConfig{}
			if err := yaml.Unmarshal(text, writer); err != nil {
				return err
			}
			if writer.Encoder == "" {
				return fmt.Errorf("writer '%s' is missing required encoder", name)
			}
			c.Stdout = writer
		case "stderr":
			text, err := yaml.Marshal(config)
			if err != nil {
				return err
			}
			writer := &stderrWriterConfig{}
			if err := yaml.Unmarshal(text, writer); err != nil {
				return err
			}
			if writer.Encoder == "" {
				return fmt.Errorf("writer '%s' is missing required encoder", name)
			}
			c.Stderr = writer
		default:
			text, err := yaml.Marshal(config)
			if err != nil {
				return err
			}
			var writer fileWriterConfig
			if err := yaml.Unmarshal(text, &writer); err != nil {
				return err
			}
			if writer.Encoder == "" {
				return fmt.Errorf("writer '%s' is missing required encoder", name)
			}
			if writer.Path == "" {
				return fmt.Errorf("writer '%s' is missing required path", name)
			}
			c.Files[name] = writer
		}
	}
	return nil
}

type writerConfig struct {
	Encoder encoderName `json:"encoder" yaml:"encoder"`
}

type stdoutWriterConfig struct {
	writerConfig `json:",inline" yaml:",inline"`
}

type stderrWriterConfig struct {
	writerConfig `json:",inline" yaml:",inline"`
}

type fileWriterConfig struct {
	writerConfig `json:",inline" yaml:",inline"`
	Path         string `json:"path" yaml:"path"`
}
