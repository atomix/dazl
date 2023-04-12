// SPDX-FileCopyrightText: 2023-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package dazl

import (
	"fmt"
	"strconv"
)

type samplerName string

const (
	countingSamplerName samplerName = "counting"
	randomSamplerName   samplerName = "random"
)

type samplingConfig struct {
	Counting *countingSamplerConfig `json:"counting" yaml:"counting"`
	Random   *randomSamplerConfig   `json:"random" yaml:"random"`
}

func (c *samplingConfig) UnmarshalText(text []byte) error {
	name := samplerName(text)
	switch name {
	case countingSamplerName:
		c.Counting = &countingSamplerConfig{
			Count: 1,
		}
	case randomSamplerName:
		c.Random = &randomSamplerConfig{}
	default:
		return fmt.Errorf("unknown sampler '%s'", name)
	}
	return nil
}

type samplerConfig struct {
	Level *levelConfig `json:"level" yaml:"level"`
}

type countingSamplerConfig struct {
	samplerConfig `json:",inline" yaml:",inline"`
	Count         int `json:"count" yaml:"count"`
}

func (c *countingSamplerConfig) UnmarshalText(text []byte) error {
	i, err := strconv.Atoi(string(text))
	if err != nil {
		return err
	}
	c.Count = i
	return nil
}

type randomSamplerConfig struct {
	samplerConfig `json:",inline" yaml:",inline"`
}
