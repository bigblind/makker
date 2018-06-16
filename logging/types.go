package logging

import (
	"context"
	"fmt"
	"strings"
)

type Logger interface{
	Debugf(ctx context.Context, format string, args... interface{})
	Infof(ctx context.Context, format string, args... interface{})
	Warnf(ctx context.Context, format string, args... interface{})
	Error(ctx context.Context, format string, args... interface{})
}

type field struct{
	name string
	value interface{}
}

type StructuredLogger struct{
	logger Logger
	fields []field
}

func NewStructuredLogger(l Logger) *StructuredLogger {
	sl := StructuredLogger{l, make([]field, 0)}
	return &sl
}

func (sl *StructuredLogger) Debugf(ctx context.Context, format string, args... interface{}) {
	sl.logger.Debugf(ctx, sl.makeFormat(format), args...)
}

func (sl *StructuredLogger) Infof(ctx context.Context, format string, args ... interface{}) {
	sl.logger.Infof(ctx, sl.makeFormat(format), args...)
}

func (sl *StructuredLogger) Warnf(ctx context.Context, format string, args ... interface{}) {
	sl.Warnf(ctx, sl.makeFormat(format), args...)
}

func (sl *StructuredLogger) Error(ctx context.Context, format string, args ... interface{}) {
	sl.Error(ctx, sl.makeFormat(format), args...)
}

func (sl *StructuredLogger) WithField(name string, value interface{}) *StructuredLogger {
	sl2 := StructuredLogger{
		logger: sl.logger,
		fields: append(sl.fields, field{name, value}),
	}

	return &sl2
}

func (sl *StructuredLogger) formatFields() string {
	parts := make([]string, len(sl.fields))
	for i, f := range sl.fields {
		parts[i] = fmt.Sprintf("%v=%#v", f.name, f.value)
	}

	return strings.Join(parts, " ")
}

func (sl *StructuredLogger) makeFormat(format string) string {
	return format + "\n" + sl.formatFields()
}

