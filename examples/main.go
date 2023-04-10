// SPDX-FileCopyrightText: 2023-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package main

import "github.com/atomix/dazl"

var log = dazl.GetPackageLogger()

func main() {
	log.Info("Hello world!")
}
