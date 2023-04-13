// SPDX-FileCopyrightText: 2023-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package dazl

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testWriter struct {
	Writer
	T *testing.T
}

type testStringWriter struct {
	testWriter
	Value string
}

func (w *testStringWriter) WithStringField(name string, value string) Writer {
	assert.Equal(w.T, "foo", name)
	assert.Equal(w.T, w.Value, value)
	return w
}

func TestStringField(t *testing.T) {
	writer := testWriter{T: t}
	_, err := String("foo", "bar")(&testStringWriter{testWriter: writer, Value: "bar"})
	assert.NoError(t, err)
	_, err = String("foo", "bar")(writer)
	assert.Error(t, err)
}

type testIntWriter struct {
	testWriter
}

func (w *testIntWriter) WithIntField(name string, value int) Writer {
	assert.Equal(w.T, "foo", name)
	assert.Equal(w.T, 1, value)
	return w
}

func TestIntField(t *testing.T) {
	writer := testWriter{T: t}
	_, err := Int("foo", 1)(&testIntWriter{writer})
	assert.NoError(t, err)
	_, err = Int("foo", 1)(&testStringWriter{testWriter: writer, Value: "1"})
	assert.NoError(t, err)
	_, err = Int("foo", 1)(writer)
	assert.Error(t, err)
}

// TODO: Add test cases for remaining fields
