// SPDX-FileCopyrightText: 2023-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package zap

import (
	"fmt"
	"github.com/atomix/dazl"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"time"
)

func newWriter(writer io.Writer, encoding dazl.Encoding) (dazl.Writer, error) {
	encoder, err := newEncoder(encoding, zapcore.EncoderConfig{})
	if err != nil {
		return nil, err
	}

	var config zap.Config
	config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	config.Encoding = string(encoding)

	ws := &writeSyncer{
		Writer: writer,
	}

	logger, err := config.Build(
		zap.AddCallerSkip(5),
		zap.WrapCore(func(zapcore.Core) zapcore.Core {
			return zapcore.NewCore(encoder, ws, zap.DebugLevel)
		}))
	if err != nil {
		return nil, err
	}

	return &Writer{
		root:     logger,
		logger:   logger,
		writer:   ws,
		encoding: encoding,
	}, nil
}

func newEncoder(encoding dazl.Encoding, config zapcore.EncoderConfig) (zapcore.Encoder, error) {
	switch encoding {
	case dazl.ConsoleEncoding:
		return zapcore.NewConsoleEncoder(config), nil
	case dazl.JSONEncoding:
		return zapcore.NewJSONEncoder(config), nil
	default:
		return nil, fmt.Errorf("unknown encoding %s", encoding)
	}
}

// Writer is a dazl output implementation
type Writer struct {
	root     *zap.Logger
	logger   *zap.Logger
	writer   zapcore.WriteSyncer
	encoding dazl.Encoding
	config   zapcore.EncoderConfig
}

func (w *Writer) WithName(name string) dazl.Writer {
	return &Writer{
		root:     w.root,
		logger:   w.root.Named(name),
		writer:   w.writer,
		encoding: w.encoding,
		config:   w.config,
	}
}

func (w *Writer) withField(field zap.Field) dazl.Writer {
	return &Writer{
		root:     w.root,
		logger:   w.logger.With(field),
		writer:   w.writer,
		encoding: w.encoding,
		config:   w.config,
	}
}

func (w *Writer) withOptions(options ...zap.Option) dazl.Writer {
	return &Writer{
		root:     w.root,
		logger:   w.logger.WithOptions(options...),
		writer:   w.writer,
		encoding: w.encoding,
		config:   w.config,
	}
}

func (w *Writer) withEncoder(config zapcore.EncoderConfig) (dazl.Writer, error) {
	encoder, err := newEncoder(w.encoding, config)
	if err != nil {
		return nil, err
	}
	return w.withOptions(zap.WrapCore(func(zapcore.Core) zapcore.Core {
		return zapcore.NewCore(encoder, w.writer, zap.DebugLevel)
	})), nil
}

func (w *Writer) WithMessageKey(key string) (dazl.Writer, error) {
	config := w.config
	config.MessageKey = key
	return w.withEncoder(config)
}

func (w *Writer) WithNameEnabled() (dazl.Writer, error) {
	return w.WithNameKey("logger")
}

func (w *Writer) WithNameKey(key string) (dazl.Writer, error) {
	config := w.config
	config.NameKey = key
	return w.withEncoder(config)
}

func (w *Writer) WithLevelEnabled() (dazl.Writer, error) {
	config := w.config
	config.LevelKey = "level"
	config.EncodeLevel = zapcore.CapitalLevelEncoder
	return w.withEncoder(config)
}

func (w *Writer) WithLevelKey(key string) (dazl.Writer, error) {
	config := w.config
	config.LevelKey = key
	return w.withEncoder(config)
}

func (w *Writer) WithLevelFormat(format dazl.LevelFormat) (dazl.Writer, error) {
	config := w.config
	switch format {
	case dazl.LowerCaseLevelFormat:
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	case dazl.UpperCaseLevelFormat:
		config.EncodeLevel = zapcore.CapitalLevelEncoder
	default:
		return nil, fmt.Errorf("unsupported level format '%s'", format)
	}
	return w.withEncoder(config)
}

func (w *Writer) WithTimestampEnabled() (dazl.Writer, error) {
	config := w.config
	config.TimeKey = "time"
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	return w.withEncoder(config)
}

func (w *Writer) WithTimestampKey(key string) (dazl.Writer, error) {
	config := w.config
	config.TimeKey = key
	return w.withEncoder(config)
}

func (w *Writer) WithTimestampFormat(format dazl.TimestampFormat) (dazl.Writer, error) {
	config := w.config
	switch format {
	case dazl.ISO8601TimestampFormat:
		config.EncodeTime = zapcore.ISO8601TimeEncoder
	case dazl.UnixTimestampFormat:
		config.EncodeTime = zapcore.EpochTimeEncoder
	default:
		return nil, fmt.Errorf("unsupported time format '%s'", format)
	}
	return w.withEncoder(config)
}

func (w *Writer) WithCallerEnabled() (dazl.Writer, error) {
	config := w.config
	config.CallerKey = "caller"
	config.EncodeCaller = zapcore.ShortCallerEncoder
	return w.withEncoder(config)
}

func (w *Writer) WithCallerKey(key string) (dazl.Writer, error) {
	config := w.config
	config.CallerKey = key
	return w.withEncoder(config)
}

func (w *Writer) WithCallerFormat(format dazl.CallerFormat) (dazl.Writer, error) {
	config := w.config
	switch format {
	case dazl.ShortCallerFormat:
		config.EncodeCaller = zapcore.ShortCallerEncoder
	case dazl.FullCallerFormat:
		config.EncodeCaller = zapcore.FullCallerEncoder
	default:
		return nil, fmt.Errorf("unsupported caller format '%s'", format)
	}
	return w.withEncoder(config)
}

func (w *Writer) WithStacktraceEnabled() (dazl.Writer, error) {
	return w.WithStacktraceKey("trace")
}

func (w *Writer) WithStacktraceKey(key string) (dazl.Writer, error) {
	config := w.config
	config.StacktraceKey = key
	return w.withEncoder(config)
}

func (w *Writer) WithStringField(name string, value string) dazl.Writer {
	return w.withField(zap.String(name, value))
}

func (w *Writer) WithBoolField(name string, value bool) dazl.Writer {
	return w.withField(zap.Bool(name, value))
}

func (w *Writer) WithIntField(name string, value int) dazl.Writer {
	return w.withField(zap.Int(name, value))
}

func (w *Writer) WithInt32Field(name string, value int32) dazl.Writer {
	return w.withField(zap.Int32(name, value))
}

func (w *Writer) WithInt64Field(name string, value int64) dazl.Writer {
	return w.withField(zap.Int64(name, value))
}

func (w *Writer) WithUintField(name string, value uint) dazl.Writer {
	return w.withField(zap.Uint(name, value))
}

func (w *Writer) WithUint32Field(name string, value uint32) dazl.Writer {
	return w.withField(zap.Uint32(name, value))
}

func (w *Writer) WithUint64Field(name string, value uint64) dazl.Writer {
	return w.withField(zap.Uint64(name, value))
}

func (w *Writer) WithFloat32Field(name string, value float32) dazl.Writer {
	return w.withField(zap.Float32(name, value))
}

func (w *Writer) WithFloat64Field(name string, value float64) dazl.Writer {
	return w.withField(zap.Float64(name, value))
}

func (w *Writer) WithTimeField(name string, value time.Time) dazl.Writer {
	return w.withField(zap.Time(name, value))
}

func (w *Writer) WithDurationField(name string, value time.Duration) dazl.Writer {
	return w.withField(zap.Duration(name, value))
}

func (w *Writer) WithBinaryField(name string, value []byte) dazl.Writer {
	return w.withField(zap.Binary(name, value))
}

func (w *Writer) WithBytesField(name string, value []byte) dazl.Writer {
	return w.withField(zap.ByteString(name, value))
}

func (w *Writer) WithStringPointerField(name string, value *string) dazl.Writer {
	return w.withField(zap.Stringp(name, value))
}

func (w *Writer) WithBoolPointerField(name string, value *bool) dazl.Writer {
	return w.withField(zap.Boolp(name, value))
}

func (w *Writer) WithIntPointerField(name string, value *int) dazl.Writer {
	return w.withField(zap.Intp(name, value))
}

func (w *Writer) WithInt32PointerField(name string, value *int32) dazl.Writer {
	return w.withField(zap.Int32p(name, value))
}

func (w *Writer) WithInt64PointerField(name string, value *int64) dazl.Writer {
	return w.withField(zap.Int64p(name, value))
}

func (w *Writer) WithUintPointerField(name string, value *uint) dazl.Writer {
	return w.withField(zap.Uintp(name, value))
}

func (w *Writer) WithUint32PointerField(name string, value *uint32) dazl.Writer {
	return w.withField(zap.Uint32p(name, value))
}

func (w *Writer) WithUint64PointerField(name string, value *uint64) dazl.Writer {
	return w.withField(zap.Uint64p(name, value))
}

func (w *Writer) WithFloat32PointerField(name string, value *float32) dazl.Writer {
	return w.withField(zap.Float32p(name, value))
}

func (w *Writer) WithFloat64PointerField(name string, value *float64) dazl.Writer {
	return w.withField(zap.Float64p(name, value))
}

func (w *Writer) WithTimePointerField(name string, value *time.Time) dazl.Writer {
	return w.withField(zap.Timep(name, value))
}

func (w *Writer) WithDurationPointerField(name string, value *time.Duration) dazl.Writer {
	return w.withField(zap.Durationp(name, value))
}

func (w *Writer) WithStringSliceField(name string, values []string) dazl.Writer {
	return w.withField(zap.Strings(name, values))
}

func (w *Writer) WithBoolSliceField(name string, values []bool) dazl.Writer {
	return w.withField(zap.Bools(name, values))
}

func (w *Writer) WithIntSliceField(name string, values []int) dazl.Writer {
	return w.withField(zap.Ints(name, values))
}

func (w *Writer) WithInt32SliceField(name string, values []int32) dazl.Writer {
	return w.withField(zap.Int32s(name, values))
}

func (w *Writer) WithInt64SliceField(name string, values []int64) dazl.Writer {
	return w.withField(zap.Int64s(name, values))
}

func (w *Writer) WithUintSliceField(name string, values []uint) dazl.Writer {
	return w.withField(zap.Uints(name, values))
}

func (w *Writer) WithUint32SliceField(name string, values []uint32) dazl.Writer {
	return w.withField(zap.Uint32s(name, values))
}

func (w *Writer) WithUint64SliceField(name string, values []uint64) dazl.Writer {
	return w.withField(zap.Uint64s(name, values))
}

func (w *Writer) WithFloat32SliceField(name string, values []float32) dazl.Writer {
	return w.withField(zap.Float32s(name, values))
}

func (w *Writer) WithFloat64SliceField(name string, values []float64) dazl.Writer {
	return w.withField(zap.Float64s(name, values))
}

func (w *Writer) WithTimeSliceField(name string, values []time.Time) dazl.Writer {
	return w.withField(zap.Times(name, values))
}

func (w *Writer) WithDurationSliceField(name string, values []time.Duration) dazl.Writer {
	return w.withField(zap.Durations(name, values))
}

func (w *Writer) WithSkipCalls(calls int) dazl.Writer {
	return w.withOptions(zap.AddCallerSkip(calls))
}

func (w *Writer) Debug(msg string) {
	w.logger.Debug(msg)
}

func (w *Writer) Info(msg string) {
	w.logger.Info(msg)
}

func (w *Writer) Error(msg string) {
	w.logger.Error(msg)
}

func (w *Writer) Fatal(msg string) {
	w.logger.Fatal(msg)
}

func (w *Writer) Panic(msg string) {
	w.logger.Panic(msg)
}

func (w *Writer) Warn(msg string) {
	w.logger.Warn(msg)
}

func (w *Writer) Sync() error {
	return w.logger.Sync()
}

var _ dazl.Writer = (*Writer)(nil)
var _ dazl.FieldWriter = (*Writer)(nil)
var _ dazl.MessageKeyWriter = (*Writer)(nil)
var _ dazl.NameWriter = (*Writer)(nil)
var _ dazl.NameKeyWriter = (*Writer)(nil)
var _ dazl.LevelWriter = (*Writer)(nil)
var _ dazl.LevelKeyWriter = (*Writer)(nil)
var _ dazl.LevelFormattingWriter = (*Writer)(nil)
var _ dazl.TimestampWriter = (*Writer)(nil)
var _ dazl.TimestampKeyWriter = (*Writer)(nil)
var _ dazl.TimestampFormattingWriter = (*Writer)(nil)
var _ dazl.CallerWriter = (*Writer)(nil)
var _ dazl.CallerKeyWriter = (*Writer)(nil)
var _ dazl.CallerFormattingWriter = (*Writer)(nil)
var _ dazl.StacktraceWriter = (*Writer)(nil)
var _ dazl.StacktraceKeyWriter = (*Writer)(nil)

type writeSyncer struct {
	io.Writer
}

func (w *writeSyncer) Sync() error {
	return nil
}
