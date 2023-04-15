// SPDX-FileCopyrightText: 2023-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package dazl

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync/atomic"
)

type samplingStintervalgy string

const (
	basicSamplingStintervalgy  samplingStintervalgy = "basic"
	randomSamplingStintervalgy samplingStintervalgy = "random"
)

type samplingConfig struct {
	Basic  *basicSamplerConfig  `json:"basic" yaml:"basic"`
	Random *randomSamplerConfig `json:"random" yaml:"random"`
}

func (c *samplingConfig) UnmarshalText(text []byte) error {
	name := samplingStintervalgy(text)
	switch name {
	case basicSamplingStintervalgy:
		c.Basic = &basicSamplerConfig{
			Interval: 1,
		}
	case randomSamplingStintervalgy:
		c.Random = &randomSamplerConfig{
			Interval: 10,
		}
	default:
		return fmt.Errorf("unknown sampler '%s'", name)
	}
	return nil
}

type samplerConfig struct {
	MinLevel levelConfig `json:"minLevel" yaml:"minLevel"`
}

type basicSamplerConfig struct {
	samplerConfig `json:",inline" yaml:",inline"`
	Interval      int `json:"interval" yaml:"interval"`
}

func (c *basicSamplerConfig) UnmarshalText(text []byte) error {
	i, err := strconv.Atoi(string(text))
	if err != nil {
		return err
	}
	c.Interval = i
	return nil
}

type randomSamplerConfig struct {
	samplerConfig `json:",inline" yaml:",inline"`
	Interval      int `json:"interval" yaml:"interval"`
}

type Sampler interface {
	Sample(level Level) bool
}

type allSampler struct{}

func (s allSampler) Sample(level Level) bool {
	return true
}

type basicSampler struct {
	Interval uint32
	MinLevel Level
	counter  atomic.Uint32
}

func (s *basicSampler) Sample(level Level) bool {
	if s.MinLevel == EmptyLevel || level.Enabled(s.MinLevel) {
		if s.Interval == 1 {
			return true
		}
		n := s.counter.Add(1)
		return n%s.Interval == 1
	}
	return true
}

type randomSampler struct {
	Interval int
	MinLevel Level
}

func (s randomSampler) Sample(level Level) bool {
	if s.MinLevel == EmptyLevel || level.Enabled(s.MinLevel) {
		if s.Interval <= 0 {
			return false
		}
		if rand.Intn(s.Interval) != 0 {
			return false
		}
		return true
	}
	return true
}
