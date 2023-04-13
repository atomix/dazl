// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package dazl

// Level :
type Level int32

const (
	// EmptyLevel :
	EmptyLevel Level = iota
	// DebugLevel logs a message at debug level
	DebugLevel
	// InfoLevel logs a message at info level
	InfoLevel
	// WarnLevel logs a message at warning level
	WarnLevel
	// ErrorLevel logs a message at error level
	ErrorLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel
	// PanicLevel logs a message, then panics.
	PanicLevel
)

// Enabled indicates whether the log level is enabled
func (l Level) Enabled(level Level) bool {
	return l <= level
}

// String :
func (l Level) String() string {
	return [...]string{"", "debug", "info", "warn", "error", "fatal", "panic"}[l]
}

type levelConfig Level

func (c levelConfig) Level() Level {
	return Level(c)
}

func (c *levelConfig) UnmarshalText(text []byte) error {
	switch string(text) {
	case DebugLevel.String():
		*c = levelConfig(DebugLevel)
	case InfoLevel.String():
		*c = levelConfig(InfoLevel)
	case WarnLevel.String():
		*c = levelConfig(WarnLevel)
	case ErrorLevel.String():
		*c = levelConfig(ErrorLevel)
	case FatalLevel.String():
		*c = levelConfig(FatalLevel)
	case PanicLevel.String():
		*c = levelConfig(PanicLevel)
	default:
		*c = levelConfig(EmptyLevel)
	}
	return nil
}
