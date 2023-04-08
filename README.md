<!--
SPDX-FileCopyrightText: 2023-present Intel Corporation
SPDX-License-Identifier: Apache-2.0
-->

# dazl

## A simple Go framework for granular logging controls

[![Docs](https://pkg.go.dev/badge/github.com/atomix/dazl)](https://pkg.go.dev/github.com/atomix/dazl)
![License](https://img.shields.io/github/license/atomix/dazl)
[![Build](https://img.shields.io/github/actions/workflow/status/atomix/dazl/test.yml)](https://github.com/atomix/dazl/actions/workflows/test.yml)

dazl is a lightweight wrapper around [zap](https://github.com/uber-go/zap) loggers that adds a logger hierarchy
with log level inheritence to enable fine-grained control and configuration of log levels, encoding, and formats.

## Installation

To add dazl to your Go module:

```bash
go get github.com/atomix/dazl
```

## Usage

The typical usage of the framework is to create a `Logger` for each Go package:

```go
var log = dazl.GetLogger()
```

By default, the logger will be assigned the package path as its name.

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

But a custom name can optionally be assigned to loggers:

```go
var log = dazl.GetLogger("github.com/my/project/test")
```

Logger names must be formatted in path format, with each element separated by a `/`. This format is used
to establish a hierarchy for inheritence of logger configurations.

### Log levels

dazl supports a fairly standard set of log levels for loggers:

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

### Inheritance

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

## Logging configuration

Loggers can be configured via a YAML configuration file. The configuration files may be in one of many
locations on the file system:

* `dazl.yaml`
* `.atomix/dazl.yaml`
* `~/.atomix/dazl.yaml`
* `/etc/atomix/dazl.yaml`

The configuration file contains a set of `loggers` which specifies the level and outputs of each logger,
and `sinks` which specify where and how to write log messages.

## The root logger

The default logging configuration is configured via the `rootLogger` key:

```yaml
rootLogger:
  level: info
```

All loggers inherit their default configuration from the root logger.

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

Logs will be written to the sink using the specified `endoding`:
* `console` - plain text output suitable for human consumption
* `json` - structured JSON logging suitable for machine consumption

Sinks also support additional formatting options discussed in the [formatting](#formatting) section.

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

## Formatting

As with the loggers, the format of the log outputs can be configured either globally or on a per-sink basis.

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
