// SPDX-FileCopyrightText: 2023-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package main

import "github.com/atomix/dazl"
import _ "github.com/atomix/dazl/zap"

var log = dazl.GetPackageLogger()

const projectName = "dazl"

func main() {
	log.Infof("%s is a logging framework for Go", projectName)
}
