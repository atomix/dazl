// SPDX-FileCopyrightText: 2023-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package dazl

func Register(framework Framework) {
	var config loggingConfig
	if err := load(&config); err != nil {
		panic(err)
	} else if err := configure(framework, config); err != nil {
		panic(err)
	}
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
