// SPDX-FileCopyrightText: 2023-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package dazl

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"testing"
)

const testSampler = `

`

func TestUnmarshalSampler(t *testing.T) {
	text := "counting: 10"
	config := &samplingConfig{}
	assert.NoError(t, yaml.Unmarshal([]byte(text), config))
	assert.NotNil(t, config.Counting)
	assert.Nil(t, config.Random)
	assert.Equal(t, 10, config.Counting.Count)
	assert.Nil(t, config.Counting.Level)

	text = `counting:
  count: 10
  level: info
`
	config = &samplingConfig{}
	assert.NoError(t, yaml.Unmarshal([]byte(text), config))
	assert.NotNil(t, config.Counting)
	assert.Nil(t, config.Random)
	assert.Equal(t, 10, config.Counting.Count)
	assert.Equal(t, InfoLevel, config.Counting.Level.Level())

	text = "random"
	config = &samplingConfig{}
	assert.NoError(t, yaml.Unmarshal([]byte(text), config))
	assert.Nil(t, config.Counting)
	assert.NotNil(t, config.Random)
	assert.Nil(t, config.Random.Level)

	text = `random:
  level: info`
	config = &samplingConfig{}
	assert.NoError(t, yaml.Unmarshal([]byte(text), config))
	assert.Nil(t, config.Counting)
	assert.NotNil(t, config.Random)
	assert.Equal(t, InfoLevel, config.Random.Level.Level())
}
