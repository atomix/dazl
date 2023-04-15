// SPDX-FileCopyrightText: 2023-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package dazl

import (
	"sync/atomic"
)

var globalFramework = &atomic.Pointer[Framework]{}

func Register(framework Framework) {
	globalFramework.Store(&framework)
}

func getFramework() Framework {
	framework := globalFramework.Load()
	if framework == nil {
		return &defaultFramework{}
	}
	return *framework
}

type Framework interface {
	Name() string
}

type ConsoleEncodingFramework interface {
	ConsoleEncoder() Encoder
}

type JSONEncodingFramework interface {
	JSONEncoder() Encoder
}

type defaultFramework struct{}

func (f *defaultFramework) Name() string {
	return "default"
}
