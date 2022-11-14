package main

import (
	"fmt"
	"os"

	"github.com/natefinch/lumberjack"
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
	fmt.Println("这里----------------------------------")
	return true
}

func getStacks([]zapcore.Field) (string, []zapcore.Field) {
	fmt.Println("这里呀----------------------------------")
	return "", nil
}

type ALL_LOGGER_LEVEL string

var (
	appLogger *zap.SugaredLogger
)
var (
	DEBUG ALL_LOGGER_LEVEL = "DEBUG"
	INFO  ALL_LOGGER_LEVEL = "INFO"
	WARN  ALL_LOGGER_LEVEL = "WARN"
	ERROR ALL_LOGGER_LEVEL = "ERROR"
	FATAL ALL_LOGGER_LEVEL = "FATAL"
	PANIC ALL_LOGGER_LEVEL = "PANIC"
)

var (
	LOGGER_LEVEL      ALL_LOGGER_LEVEL = INFO
	LOGGER_MAX_SIZE   int              = 100
	LOGGER_MAX_BACKUP int              = 7
	LOGGER_MAX_AGE    int              = 30
)

// GetLogger 获取logger对象
func GetLogger() *zap.SugaredLogger {
	if appLogger == nil {
		InitLog()
	}
	return appLogger
}

func InitLog() *zap.SugaredLogger {
	var level zapcore.Level
	switch LOGGER_LEVEL {
	case DEBUG:
		level = zapcore.DebugLevel
	case INFO:
		level = zapcore.InfoLevel
	case WARN:
		level = zapcore.WarnLevel
	case ERROR:
		level = zapcore.ErrorLevel
	case FATAL:
		level = zapcore.FatalLevel
	case PANIC:
		level = zapcore.PanicLevel
	default:
		level = zapcore.InfoLevel
	}

	logger := NewLogger(level, LOGGER_MAX_SIZE, LOGGER_MAX_BACKUP, LOGGER_MAX_AGE, true).Sugar()
	appLogger = logger
	return logger
}

/**
 * 获取日志
 * filePath 日志文件路径
 * level 日志级别
 * maxSize 每个日志文件保存的最大尺寸 单位：M
 * maxBackups 日志文件最多保存多少个备份
 * maxAge 文件最多保存多少天
 * compress 是否压缩
 * serviceName 服务名
 */
func NewLogger(level zapcore.Level, maxSize int, maxBackups int, maxAge int, compress bool) *zap.Logger {
	core := newCore(level, maxSize, maxBackups, maxAge, compress)
	errCore := NewErrStackCore(core)
	return zap.New(errCore, zap.AddCaller(), zap.Development())
}

/**
 * zapcore构造
 */
func newCore(level zapcore.Level, maxSize int, maxBackups int, maxAge int, compress bool) zapcore.Core {
	//日志文件路径配置
	hook := lumberjack.Logger{
		MaxSize:    maxSize,    // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: maxBackups, // 日志文件最多保存多少个备份
		MaxAge:     maxAge,     // 文件最多保存多少天
		Compress:   compress,   // 是否压缩
	}
	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(level)
	//公用编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,    // 大写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 包路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	return zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		// 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
		atomicLevel, // 日志级别
	)
}

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
	InitLog()
	err := a()
	if err != nil {
		fmt.Printf("%+v\n", err)
		appLogger.Errorf("%+v", err)
	}
}
