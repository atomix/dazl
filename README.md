<!--
SPDX-FileCopyrightText: 2023-present Intel Corporation
SPDX-License-Identifier: Apache-2.0
-->

# dazl

## A simple Go framework for granular logging controls

[![Docs](https://pkg.go.dev/badge/github.com/atomix/dazl)](https://pkg.go.dev/github.com/atomix/dazl)
![License](https://img.shields.io/github/license/atomix/dazl)
[![Build](https://img.shields.io/github/actions/workflow/status/atomix/dazl/test.yml)](https://github.com/atomix/dazl/actions/workflows/test.yml)

Dazl is not just another Go logging framework. We're not here to reinvent Go logging for the nth time. Instead,
dazl wraps [a logging library that already works](https://github.com/uber-go/zap), enriching it and the existing
ecosystem with fine-grained logging controls and configuration files, making your applications easier to configure 
and debug, both for you and your users.

### A brief history

This framework was originally developed at the [Open Networking Foundation](https://opennetworking.org) and is still 
used extensively in open source cloud software projects at [Intel](https://www.intel.com), providing users and
operators the ability to easily configure logging for specific systems and subsystems, control logging targets and
formats, and enable rapid debugging and analysis without the need to modify code.

## User Guide

* [Installation](#installation)
* [Getting started](#getting-started)
  * [Log levels](#log-levels)
  * [Structured logging](#structured-logging)
  * [Inheritance](#inheritance)
* [Configuration](#configuration)
  * [The root logger](#the-root-logger)
  * [Setting up sinks](#sinks)
    * [Encoding messages](#encodings)
    * [Field-level encodings](#formatting-outputs)
  * [Logger outputs](#outputs)
    * [Output level filters](#output-levels)
* [Writing unit tests](#writing-unit-tests)

## Installation

To add dazl to your Go module:

```bash
go get github.com/atomix/dazl
```

## Getting started

The typical usage of the framework is to create a `Logger` once at the top of each Go package:

```go
var log = dazl.GetLogger()
```

By default, loggers will be assigned the path of the package calling the `GetLogger()` function. So, if you 
call `dazl.GetLogger()` from the `github.com/atomix/atomix/runtime` package, the logger will be assigned the
name `github.com/atomix/atomix/runtime`. The naming strategy becomes important for 
[logger configuration](#configuration) and, in particular, [inheritance](#inheritance).

```go
const author = "kuujo"

var log = dazl.GetLogger()

func main() {
    log.Infof("The author of dazl is %s", author)
}
```

```
2023-03-31T01:13:30.607Z	INFO	main.go:12	My name is Jordan Halterman
```

A custom name may also be assigned to loggers:

```go
var log = dazl.GetLogger("test")
```

Logger names must be formatted in path format, with each element separated by a `/`. This format is used
to establish a hierarchy for [inheritence](#inheritance) of logger configurations.

## Log levels

Dazl supports a fairly standard set of log levels for loggers:

* `debug`
* `info`
* `warn`
* `error`
* `fatal`
* `panic`

The levels for each logger can be configured individually via their configuration:

```yaml
loggers:
  github.com/atomix/atomix/runtime:
    level: warn
```

The `Logger` interface exposes methods for simple logging, formatted logging, and structured logging with typed
fields for each log level:
* `Debug(args ...any)`
* `Debugf(msg string, args ...any)`
* `Debugw(msg string, fields ...Field)`
* `Info(args ...any)`
* `Infof(msg string, args ...any)`
* `Infow(msg string, fields ...Field)`
* `Warn(args ...any)`
* `Warnf(msg string, args ...any)`
* `Warnw(msg string, fields ...Field)`
* ...

## Structured logging

Structured logging is supported for the JSON [encoding](#encodings), and JSON fields are configurable via
the `Logger` API.

The simplest way to add fields to your structured logs is to call one of the `*w` methods on the `Logger` interface.
These methods accept an arbitrary number of `...Field`s to write to the logs. Fields are typed and named values
that can be constructed via functions in the `dazl` package:

```go
log.Warnw("Something went wrong!", 
	  dazl.String("user", user.Name), 
	  dazl.Uint64("user-id", user.ID))
```

Alternatively, you can create a structured logger with a fixed set of fields using the `WithFields` method:

```go
var log = dazl.GetLogger().WithFields(
    dazl.String("user", user.Name),
    dazl.Uint64("id", user.ID))
log.Warn("Something went wrong!")
```

When the logger is output to a JSON encoded sink, the above code will log the fields as part of the JSON object:

```
{"timestamp":"2023-04-07T19:24:09-07:00","logger":"2/4","message":"Something went wrong!","user":"Jordan Halterman","id":5678}
```

## Inheritance

The path-like format used for logger names is used to establish a hierarchy of loggers. The dazl configuration
enables developers and their users to configure individual loggers at runtime. Log levels are inherited by
descendants of a logger. This enables users to easily enable logging for specific Go packages, their subpackages,
or entire Go modules with a single configuration change:

```yaml
# Enable debug logging for all Atomix code
loggers:
  github.com/atomix:
    level: debug
```

You can set the `dazl.Level` for a logger at startup time via configuration files or at runtime
via the `Logger` API to control the granularity of a logger's output:

```go
dazl.GetLogger("github.com/atomix").SetLevel(dazl.DebugLevel)
```

If the level for a logger is not explicitly set, it will inherit its level from its nearest ancestor in
the logger hierarchy. For example, setting the `github.com/atomix/atomix/runtime` logger to the `debug`
level will change the loggers for all loggers in the `github.com/atomix/atomix/runtime/...` packages
to the `debug` level.

## Configuration

Loggers can be configured via a YAML configuration file. The configuration files may be in one of many
locations on the file system:

* `dazl.yaml`
* `.atomix/dazl.yaml`
* `~/.atomix/dazl.yaml`
* `/etc/atomix/dazl.yaml`

The configuration file contains a set of `loggers` which specifies the level and outputs of each logger,
and `sinks` which specify where and how to write log messages.

## Sinks

In order to write logs to some destination -- be it a log file or stdout -- the dazl configuration must 
define one or more `sinks`:

```yaml
sinks:
  stdout:
    path: stdout
    encoding: console
```

Each sink must be keyed by a unique name. This is used to reference sinks in [outputs](#outputs). Additionally, the
sink must specify a `path` to write logs. Paths are URLs indicating the sink target (with a couple exceptions):
* `stdout` - write to stdout
* `stderr` - write to stderr
* `file://...` - write to the given file

### Encodings

Logs will be written to the sink using the specified `endoding`:
* `console` - plain text output suitable for human consumption
* `json` - structured JSON logging suitable for machine consumption

Sinks also support additional formatting options discussed in the [formatting](#formatting) section.

## The root logger

The default logging configuration is configured via the `rootLogger` key:

```yaml
rootLogger:
  level: info
```

All loggers inherit their default configuration from the root logger. The root logger configuration should at
least specify a minimum log `level` for all loggers, and at least one `outputs` to a [sink](#sinks).

```yaml
rootLogger:
  level: info
  outputs:
    stdout:
      sink: stdout
```

## Named loggers

Once the `rootLogger` has been defined, all other loggers that are descendants of the root logger may be defined and 
configured in the `loggers` section of the configuration file:

```yaml
loggers:
  github.com/atomix/atomix/runtime:
    level: info
  github.com/atomix/atomix/sidecar:
    level: warn
```

## Outputs

Once you've defined the set of sinks to which to write your application logs, named loggers can be directed to
those sinks via their configured `outputs`:

```yaml
rootLogger:
  level: info
  outputs:
    stdout:
      sink: stdout
```

Each output must specify a `sink` to write to. As with log levels, loggers inherit the outputs of their ancestors, so
with an output to the `stdout` sink in the root logger configuration, all loggers will have an output to `stdout`.

### Output levels

In some cases, you may want to restrict the verbosity of logs to one output without restricting the verbosity of all
messages for a logger. For example, you may want to write `info` and higher messages to the console for 
human readability, and `debug` and higher messages to a file for later debugging. Each output supports a `level`
that can be used to filter the logs to that output:

```yaml
loggers:
  github.com/atomix/atomix:
    # Pass through all messages over 'debug' level to the outputs
    level: debug
    outputs:
      # Limit the messages from this logger to stdout to 'info' level
      stdout:
        sink: stdout
        level: info
      # Pass through 'debug' level messages to the file sink
      file:
        sink: file

sinks:
  stdout:
    path: stdout
    encoding: console
  file:
    path: file://app.log
    encoding: json
```

### Output overrides

Descendants may override their ancestor loggers' output configurations. This can be done by simply specifying the
same output name. For example, to override the log level for the `rootLogger`'s `stdout` output:

```yaml
sinks:
  # Define an stdout sink with console encoding
  stdout:
    path: stdout
    encoding: console

rootLogger:
  level: info
  outputs:
    # Output the root logger to the stdout sink
    stdout:
      sink: stdout

loggers:
  github.com/atomix/atomix/runtime:
    outputs:
      # Filter outputs from this logger to the stdout sink to minimum of 'warn' level
      stdout:
        level: warn
```

## Formatting outputs

As with the loggers, the format of the log outputs can be configured either globally or on a per-sink basis. In
addition to the sinks and loggers, the top-level configuration has the following fields which act as global 
settings for all the sinks and loggers:
* `messageKey` - the key to which to write the log message in structured logs
* `levelKey` - the key to which to write the log level in structured logs
* `timeKey` - the key to which to write the timestamp in structured logs
* `nameKey` - the key to which to write the logger name in structured logs
* `callerKey` - the key to which to write the caller in structured logs
* `stacktraceKey` - the key to which to write the stack trace in structured logs
* `skipLineEnding` - whether to skip the line endings in outputs
* `lineEnding` - the line ending to use in outputs
* `levelEncoder` - the [level encoder](#formatting-log-levels) to use for all logs
* `timeEncoder` - the [time encoder](#formatting-timestamps) to use for all logs
* `durationEncoder` - the [duration encoder](#formatting-durations) to use for all logs
* `callerEncoder` - the encoder to use to encode the caller info in log outputs
* `nameEncoder` - the encoder to use to encode the logger name in log outputs

By default, the following values are set:

```yaml
messageKey: message
levelKey: level
timeKey: timestamp
nameKey: logger
```

The above fields can also be configured for each individual sink, if necessary.

```yaml
sinks:
  stdout:
    path: stdout
    encoding: console
    levelEncoder: capitalColor
```

To include the [log level](#log-levels), [logger names](#named-loggers), [timestamps](#formatting-timestamps), and
other fields in the log output, you must configure the encoders for those fields:

```yaml
levelEncoder: capital
nameEncoder: full
timeEncoder: rfc3339
```

### Formatting log levels

The format of the log level written to sinks can be configured via the `levelEncoder` key:

```yaml
levelEncoder: capital
```

Level encoders can be configured on individual sinks, as well:

```yaml
sinks:
  stdout:
    path: stdout
    encoding: console
    levelEncoder: capitalColor
  file:
    path: file://app.log
    encoding: json
    levelEncoder: capital
```

Valid level formats are:
* `capital` - upper case level name
* `capitalColor` - upper case, color-coded level name
* `color` - lower case color-coded level name
* `` - defaults to lower case level name

### Formatting timestamps

The timestamp format can be configured via the `timeEncoder` key:

```yaml
timeEncoder: RFC3339
```

As with other formatting options, time encoders can be configured for individual sinks, as well:

```yaml
sinks:
  stdout:
    path: stdout
    encoding: console
    timeEncoder: rfc3339
  file:
    path: file://app.log
    encoding: json
    timeEncoder: nanos
```

Valid timestamp formats are:
* `rfc3339nano`
* `RFC3339Nano`
* `rfc3339`
* `RFC3339`
* `iso8601`
* `ISO8601`
* `millis`
* `nanos`

### Formatting durations

The format of directions written to the log can be configured via the `durationEncoder` key:

```yaml
durationEncoder: string
```

As with other formatting options, time encoders can be configured for individual sinks, as well:

```yaml
sinks:
  stdout:
    path: stdout
    encoding: console
    durationEncoder: string
  file:
    path: file://app.log
    encoding: json
    durationEncoder: nanos
```

Valid duration formats include:
* `string`
* `nanos`
* `ms`

## Writing unit tests

For unit testing, tests can configure custom sinks to arbitrary `io.Writer`s for verifying logging output.

```go
func TestLogger(t *testing.T) {
    var buf bytes.Buffer
    sink, err := dazl.NewSink(&buf,
        dazl.WithEncoding(dazl.ConsoleEncoding),
        dazl.WithNameKey("name"),
        dazl.WithMessageKey("message"),
        dazl.WithLevelKey("level"),
        dazl.WithNameEncoder(zapcore.FullNameEncoder),
        dazl.WithLevelEncoder(zapcore.CapitalLevelEncoder))
	  assert.NoError(t, err)
    
    var log = dazl.GetLogger().WithOutputs(dazl.NewOutput(sink))
    
    log.Info("Hello world!")
    assert.Equal(t, "INFO\tgithub.com/my/package\tHello world!\n", buf.String())
}
```
