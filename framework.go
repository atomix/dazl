// SPDX-FileCopyrightText: 2023-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package dazl

import (
	"io"
	"sync/atomic"
)

var globalFramework = &atomic.Pointer[Framework]{}
var defaultFramework = &noopFramework{}

func Register(framework Framework) {
	globalFramework.Store(&framework)
}

func getFramework() Framework {
	framework := globalFramework.Load()
	if framework == nil {
		return defaultFramework
	}
	return *framework
}

type Framework interface {
	NewWriter(writer io.Writer, encoding Encoding) (Writer, error)
}

type noopFramework struct{}

func (f *noopFramework) NewWriter(writer io.Writer, encoding Encoding) (Writer, error) {
	return &noopWriter{}, nil
}

type noopWriter struct{}

func (w *noopWriter) WithName(name string) Writer {
	return &noopWriter{}
}

func (w *noopWriter) Debug(msg string) {

}

func (w *noopWriter) Info(msg string) {

}

func (w *noopWriter) Error(msg string) {

}

func (w *noopWriter) Fatal(msg string) {

}

func (w *noopWriter) Panic(msg string) {

}

func (w *noopWriter) Warn(msg string) {

}
