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

type testBoolWriter struct {
	testWriter
}

func (w *testBoolWriter) WithBoolField(name string, value bool) Writer {
	assert.Equal(w.T, "foo", name)
	assert.True(w.T, value)
	return w
}

func TestBoolField(t *testing.T) {
	writer := testWriter{T: t}
	_, err := Bool("foo", true)(&testBoolWriter{writer})
	assert.NoError(t, err)
	_, err = Bool("foo", true)(&testStringWriter{testWriter: writer, Value: "true"})
	assert.NoError(t, err)
	_, err = Bool("foo", true)(writer)
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

type testInt32Writer struct {
	testWriter
}

func (w *testInt32Writer) WithInt32Field(name string, value int32) Writer {
	assert.Equal(w.T, "foo", name)
	assert.Equal(w.T, int32(1), value)
	return w
}

func TestInt32Field(t *testing.T) {
	writer := testWriter{T: t}
	_, err := Int32("foo", 1)(&testInt32Writer{writer})
	assert.NoError(t, err)
	_, err = Int32("foo", 1)(&testStringWriter{testWriter: writer, Value: "1"})
	assert.NoError(t, err)
	_, err = Int32("foo", 1)(writer)
	assert.Error(t, err)
}

type testInt64Writer struct {
	testWriter
}

func (w *testInt64Writer) WithInt64Field(name string, value int64) Writer {
	assert.Equal(w.T, "foo", name)
	assert.Equal(w.T, int64(1), value)
	return w
}

func TestInt64Field(t *testing.T) {
	writer := testWriter{T: t}
	_, err := Int64("foo", 1)(&testInt64Writer{writer})
	assert.NoError(t, err)
	_, err = Int64("foo", 1)(&testStringWriter{testWriter: writer, Value: "1"})
	assert.NoError(t, err)
	_, err = Int64("foo", 1)(writer)
	assert.Error(t, err)
}

type testUintWriter struct {
	testWriter
}

func (w *testUintWriter) WithUintField(name string, value uint) Writer {
	assert.Equal(w.T, "foo", name)
	assert.Equal(w.T, uint(1), value)
	return w
}

func TestUintField(t *testing.T) {
	writer := testWriter{T: t}
	_, err := Uint("foo", 1)(&testUintWriter{writer})
	assert.NoError(t, err)
	_, err = Uint("foo", 1)(&testStringWriter{testWriter: writer, Value: "1"})
	assert.NoError(t, err)
	_, err = Uint("foo", 1)(writer)
	assert.Error(t, err)
}

type testUint32Writer struct {
	testWriter
}

func (w *testUint32Writer) WithUint32Field(name string, value uint32) Writer {
	assert.Equal(w.T, "foo", name)
	assert.Equal(w.T, uint32(1), value)
	return w
}

func TestUint32Field(t *testing.T) {
	writer := testWriter{T: t}
	_, err := Uint32("foo", 1)(&testUint32Writer{writer})
	assert.NoError(t, err)
	_, err = Uint32("foo", 1)(&testStringWriter{testWriter: writer, Value: "1"})
	assert.NoError(t, err)
	_, err = Uint32("foo", 1)(writer)
	assert.Error(t, err)
}

type testUint64Writer struct {
	testWriter
}

func (w *testUint64Writer) WithUint64Field(name string, value uint64) Writer {
	assert.Equal(w.T, "foo", name)
	assert.Equal(w.T, uint64(1), value)
	return w
}

func TestUint64Field(t *testing.T) {
	writer := testWriter{T: t}
	_, err := Uint64("foo", 1)(&testUint64Writer{writer})
	assert.NoError(t, err)
	_, err = Uint64("foo", 1)(&testStringWriter{testWriter: writer, Value: "1"})
	assert.NoError(t, err)
	_, err = Uint64("foo", 1)(writer)
	assert.Error(t, err)
}

type testIntSliceWriter struct {
	testWriter
}

func (w *testIntSliceWriter) WithIntSliceField(name string, value []int) Writer {
	assert.Equal(w.T, "foo", name)
	assert.Equal(w.T, []int{1}, value)
	return w
}

func TestIntSliceField(t *testing.T) {
	writer := testWriter{T: t}
	_, err := Ints("foo", []int{1})(&testIntSliceWriter{writer})
	assert.NoError(t, err)
	_, err = Ints("foo", []int{1})(&testStringWriter{testWriter: writer, Value: "[1]"})
	assert.NoError(t, err)
	_, err = Ints("foo", []int{1})(writer)
	assert.Error(t, err)
}

type testBoolSliceWriter struct {
	testWriter
}

func (w *testBoolSliceWriter) WithBoolSliceField(name string, value []bool) Writer {
	assert.Equal(w.T, "foo", name)
	assert.Equal(w.T, []bool{true}, value)
	return w
}

func TestBoolSliceField(t *testing.T) {
	writer := testWriter{T: t}
	_, err := Bools("foo", []bool{true})(&testBoolSliceWriter{writer})
	assert.NoError(t, err)
	_, err = Bools("foo", []bool{true})(&testStringWriter{testWriter: writer, Value: "[true]"})
	assert.NoError(t, err)
	_, err = Bools("foo", []bool{true})(writer)
	assert.Error(t, err)
}

type testInt32SliceWriter struct {
	testWriter
}

func (w *testInt32SliceWriter) WithInt32SliceField(name string, value []int32) Writer {
	assert.Equal(w.T, "foo", name)
	assert.Equal(w.T, []int32{1}, value)
	return w
}

func TestInt32SliceField(t *testing.T) {
	writer := testWriter{T: t}
	_, err := Int32s("foo", []int32{1})(&testInt32SliceWriter{writer})
	assert.NoError(t, err)
	_, err = Int32s("foo", []int32{1})(&testStringWriter{testWriter: writer, Value: "[1]"})
	assert.NoError(t, err)
	_, err = Int32s("foo", []int32{1})(writer)
	assert.Error(t, err)
}

type testInt64SliceWriter struct {
	testWriter
}

func (w *testInt64SliceWriter) WithInt64SliceField(name string, value []int64) Writer {
	assert.Equal(w.T, "foo", name)
	assert.Equal(w.T, []int64{1}, value)
	return w
}

func TestInt64SliceField(t *testing.T) {
	writer := testWriter{T: t}
	_, err := Int64s("foo", []int64{1})(&testInt64SliceWriter{writer})
	assert.NoError(t, err)
	_, err = Int64s("foo", []int64{1})(&testStringWriter{testWriter: writer, Value: "[1]"})
	assert.NoError(t, err)
	_, err = Int64s("foo", []int64{1})(writer)
	assert.Error(t, err)
}

type testUintSliceWriter struct {
	testWriter
}

func (w *testUintSliceWriter) WithUintSliceField(name string, value []uint) Writer {
	assert.Equal(w.T, "foo", name)
	assert.Equal(w.T, []uint{1}, value)
	return w
}

func TestUintSliceField(t *testing.T) {
	writer := testWriter{T: t}
	_, err := Uints("foo", []uint{1})(&testUintSliceWriter{writer})
	assert.NoError(t, err)
	_, err = Uints("foo", []uint{1})(&testStringWriter{testWriter: writer, Value: "[1]"})
	assert.NoError(t, err)
	_, err = Uints("foo", []uint{1})(writer)
	assert.Error(t, err)
}

type testUint32SliceWriter struct {
	testWriter
}

func (w *testUint32SliceWriter) WithUint32SliceField(name string, value []uint32) Writer {
	assert.Equal(w.T, "foo", name)
	assert.Equal(w.T, []uint32{1}, value)
	return w
}

func TestUint32SliceField(t *testing.T) {
	writer := testWriter{T: t}
	_, err := Uint32s("foo", []uint32{1})(&testUint32SliceWriter{writer})
	assert.NoError(t, err)
	_, err = Uint32s("foo", []uint32{1})(&testStringWriter{testWriter: writer, Value: "[1]"})
	assert.NoError(t, err)
	_, err = Uint32s("foo", []uint32{1})(writer)
	assert.Error(t, err)
}

type testUint64SliceWriter struct {
	testWriter
}

func (w *testUint64SliceWriter) WithUint64SliceField(name string, value []uint64) Writer {
	assert.Equal(w.T, "foo", name)
	assert.Equal(w.T, []uint64{1}, value)
	return w
}

func TestUint64SliceField(t *testing.T) {
	writer := testWriter{T: t}
	_, err := Uint64s("foo", []uint64{1})(&testUint64SliceWriter{writer})
	assert.NoError(t, err)
	_, err = Uint64s("foo", []uint64{1})(&testStringWriter{testWriter: writer, Value: "[1]"})
	assert.NoError(t, err)
	_, err = Uint64s("foo", []uint64{1})(writer)
	assert.Error(t, err)
}
