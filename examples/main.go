// SPDX-FileCopyrightText: 2023-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package main

import "github.com/atomix/dazl"

// To configure the zap logger, import the github.com/atomix/dazl/zap package
//import _ "github.com/atomix/dazl/zap"

// To configure the zerolog logger, import the github.com/atomix/dazl/zerolog package
import _ "github.com/atomix/dazl/zerolog"

var log = dazl.GetPackageLogger()

const projectName = "dazl"

func main() {
	log.Infof("%s is a logging framework for Go", projectName)
}
