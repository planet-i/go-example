package main

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewErrStackCore(c zapcore.Core) zapcore.Core {
	return &errStackCore{c}
}

type errStackCore struct {
	zapcore.Core
}

func (c *errStackCore) With(fields []zapcore.Field) zapcore.Core {
	return &errStackCore{
		c.Core.With(fields),
	}
}

func (c *errStackCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	// 判断fields里有没有error字段
	if !hasStackedErr(fields) {
		return c.Core.Write(ent, fields)
	}
	// 这里是重点，从fields里取出error字段，把内容放到ent.Stack里，逻辑就是这样，具体代码就不给出了
	ent.Stack, fields = getStacks(fields)

	return c.Core.Write(ent, fields)
}

func (c *errStackCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	return c.Core.Check(ent, ce)
}
func hasStackedErr([]zapcore.Field) bool {
	return true
}

func getStacks([]zapcore.Field) bool {
	return true
}

var l *zap.Logger

func a() error {
	err := b()
	if err != nil {
		return err
	}
	return nil
}

func b() error {
	err := c()
	if err != nil {
		return err
	}
	return nil
}

func c() error {
	return errors.New("do c fail")
}

func main() {
	l, _ = zap.NewDevelopment()
	err := a()
	if err != nil {
		l.Error("main error", zap.Error(err))
	}
}
