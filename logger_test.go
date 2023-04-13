// SPDX-FileCopyrightText: 2023-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package dazl

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoggerNames(t *testing.T) {
	assert.Equal(t, "", GetRootLogger().Name())
	assert.Equal(t, "", GetLogger("").Name())
	assert.Equal(t, "foo", GetLogger("foo").Name())
	assert.Equal(t, "foo/bar", GetLogger("foo/bar").Name())
	assert.Equal(t, "github.com/atomix/dazl", GetPackageLogger().Name())
}

func TestSetLevel(t *testing.T) {
	root := GetRootLogger()
	parent := GetLogger("parent")
	child := parent.GetLogger("child")

	assert.Equal(t, EmptyLevel, root.Level())
	assert.Equal(t, EmptyLevel, parent.Level())
	assert.Equal(t, EmptyLevel, child.Level())

	SetLevel(ErrorLevel)
	assert.Equal(t, ErrorLevel, root.Level())
	assert.Equal(t, ErrorLevel, parent.Level())
	assert.Equal(t, ErrorLevel, child.Level())

	parent.SetLevel(WarnLevel)
	assert.Equal(t, ErrorLevel, root.Level())
	assert.Equal(t, WarnLevel, parent.Level())
	assert.Equal(t, WarnLevel, child.Level())

	child.SetLevel(ErrorLevel)
	assert.Equal(t, ErrorLevel, root.Level())
	assert.Equal(t, WarnLevel, parent.Level())
	assert.Equal(t, ErrorLevel, child.Level())
}
