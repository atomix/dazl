// SPDX-FileCopyrightText: 2023-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package zap

import (
	"github.com/atomix/dazl"
	"io"
)

func init() {
	dazl.Register(&Framework{})
}

type Framework struct{}

func (f *Framework) NewWriter(writer io.Writer, encoding dazl.Encoding) (dazl.Writer, error) {
	return newWriter(writer, encoding)
}
