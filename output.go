// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package dazl

type outputConfig struct {
	Writer string          `json:"writer" yaml:"writer"`
	Level  *levelConfig    `json:"level" yaml:"level"`
	Sample *samplingConfig `json:"sample" yaml:"sample"`
}

func newOutput(name string, writer Writer, level Level, sampler Sampler) *dazlOutput {
	return &dazlOutput{
		name:    name,
		writer:  writer,
		level:   level,
		sampler: sampler,
	}
}

// dazlOutput is a dazl output implementation
type dazlOutput struct {
	name    string
	writer  Writer
	level   Level
	sampler Sampler
}

func (o *dazlOutput) Writer() Writer {
	return o.writer
}

func (o *dazlOutput) WithWriter(writer Writer) *dazlOutput {
	return &dazlOutput{
		name:    o.name,
		writer:  writer,
		level:   o.level,
		sampler: o.sampler,
	}
}

func (o *dazlOutput) Level() Level {
	return o.level
}

func (o *dazlOutput) WithLevel(level Level) *dazlOutput {
	return &dazlOutput{
		name:    o.name,
		writer:  o.writer,
		level:   level,
		sampler: o.sampler,
	}
}

func (o *dazlOutput) WithSampler(sampler Sampler) *dazlOutput {
	return &dazlOutput{
		name:    o.name,
		writer:  o.writer,
		level:   o.level,
		sampler: sampler,
	}
}

func (o *dazlOutput) Debug(msg string) {
	if o.level.Enabled(DebugLevel) && o.sampler.Sample(DebugLevel) {
		o.writer.Debug(msg)
	}
}

func (o *dazlOutput) Info(msg string) {
	if o.level.Enabled(InfoLevel) && o.sampler.Sample(InfoLevel) {
		o.writer.Info(msg)
	}
}

func (o *dazlOutput) Warn(msg string) {
	if o.level.Enabled(WarnLevel) && o.sampler.Sample(WarnLevel) {
		o.writer.Warn(msg)
	}
}

func (o *dazlOutput) Error(msg string) {
	if o.level.Enabled(ErrorLevel) && o.sampler.Sample(ErrorLevel) {
		o.writer.Error(msg)
	}
}

func (o *dazlOutput) Fatal(msg string) {
	if o.level.Enabled(FatalLevel) && o.sampler.Sample(FatalLevel) {
		o.writer.Fatal(msg)
	}
}

func (o *dazlOutput) Panic(msg string) {
	if o.level.Enabled(PanicLevel) && o.sampler.Sample(PanicLevel) {
		o.writer.Panic(msg)
	}
}
