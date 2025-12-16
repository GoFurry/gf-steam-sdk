// Package log 提供基于 Zap 的高性能日志封装
// 支持开发/生产模式切换、结构化日志、日志轮转和自定义格式等功能
// Package log provides high-performance log wrapper based on Zap
// Supports dev/prod mode switch, structured logging, log rotation and custom format

package log

import (
	"os"
	"time"

	"github.com/GoFurry/gf-steam-sdk/pkg/util"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 全局日志实例 | Global log instances
var (
	GlobalLogger *zap.Logger        // 标准 Zap Logger | Standard Zap Logger
	SugarLogger  *zap.SugaredLogger // 简化版日志（易用性优先）| Simplified logger (ease of use first)
)

// Config 日志配置结构体
// 整合日志级别、输出模式、文件轮转等所有可配置项
// Config is the log configuration structure
// Integrates all configurable items such as log level, output mode, file rotation
type Config struct {
	Level      string // 日志级别 debug/info/warn/error | Log level debug/info/warn/error
	Mode       string // 运行模式 dev(控制台)/prod(文件) | Run mode dev(console)/prod(file)
	FilePath   string // 日志文件路径(prod必填) | Log file path (required for prod)
	MaxSize    int    // 单个日志文件大小(MB) | Single log file size (MB)
	MaxBackups int    // 最大备份数 | Max backup count
	MaxAge     int    // 最大保留天数 | Max retention days
	Compress   bool   // 是否压缩备份 | Whether to compress backups
	ShowLine   bool   // 是否显示代码行号 | Whether to show code line numbers
	EncodeJson bool   // 是否 JSON 格式输出 | Whether to output in JSON format
	TimeFormat string // 时间格式 | Time format
}

// defaultConfig 获取默认日志配置
// 返回值:
//   - Config: 默认配置 | Default configuration
func defaultConfig() Config {
	return Config{
		Level:      "info",
		Mode:       "dev", // 默认在控制台输出 | Default output to console
		FilePath:   "./logs/gf-steam-sdk.log",
		MaxSize:    100,
		MaxBackups: 7,
		MaxAge:     30,
		Compress:   true,
		ShowLine:   true,
		EncodeJson: false,
		TimeFormat: util.TIME_FORMAT_DATE,
	}
}

// InitLogger 初始化日志
// 参数:
//   - cfg: 日志配置（nil 则使用默认配置）| Log config (use default if nil)
//
// 返回值:
//   - error: 初始化失败时返回错误 | Error if initialization failed
func InitLogger(cfg *Config) error {
	if cfg == nil {
		cfg = &Config{}
	}
	// 合并默认配置 | Merge default config
	defaultCfg := defaultConfig()
	if cfg.Level == "" {
		cfg.Level = defaultCfg.Level
	}
	if cfg.Mode == "" {
		cfg.Mode = defaultCfg.Mode
	}
	if cfg.FilePath == "" {
		cfg.FilePath = defaultCfg.FilePath
	}
	if cfg.MaxSize == 0 {
		cfg.MaxSize = defaultCfg.MaxSize
	}
	if cfg.MaxBackups == 0 {
		cfg.MaxBackups = defaultCfg.MaxBackups
	}
	if cfg.MaxAge == 0 {
		cfg.MaxAge = defaultCfg.MaxAge
	}
	if cfg.TimeFormat == "" {
		cfg.TimeFormat = defaultCfg.TimeFormat
	}

	// 设置日志级别 | Set log level
	level := zapcore.InfoLevel
	switch cfg.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	case "dpanic":
		level = zapcore.DPanicLevel
	case "panic":
		level = zapcore.PanicLevel
	case "fatal":
		level = zapcore.FatalLevel
	}

	// 配置编码器 | Configure encoder
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     customTimeEncoder(cfg.TimeFormat), // 格式化时间 | Format time
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 配置输出 | Configure output
	var writeSyncer zapcore.WriteSyncer
	if cfg.Mode == "prod" {
		// 生产模式: 输出到文件 | Prod mode: output to file
		lumberjackLogger := &lumberjack.Logger{
			Filename:   cfg.FilePath,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		}
		writeSyncer = zapcore.AddSync(lumberjackLogger)
		// 生产模式默认 JSON 格式 | Prod mode default to JSON format
		if !cfg.EncodeJson {
			cfg.EncodeJson = true
		}
	} else {
		// 开发模式: 输出到控制台 | Dev mode: output to console
		writeSyncer = zapcore.AddSync(os.Stdout)
	}

	// 选择编码器 | Select encoder
	var encoder zapcore.Encoder
	if cfg.EncodeJson {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 构建 Logger 核心 | Build Logger core
	core := zapcore.NewCore(encoder, writeSyncer, level)

	// 配置选项 | Configure options
	options := []zap.Option{}
	if cfg.ShowLine {
		options = append(options, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	}

	// 初始化全局 Logger | Initialize global Logger
	GlobalLogger = zap.New(core, options...)
	SugarLogger = GlobalLogger.Sugar()

	// 测试日志 | Test log
	GlobalLogger.Info("gf-steam-sdk logger init success",
		String("mode", cfg.Mode),
		String("level", cfg.Level),
	)
	return nil
}

// customTimeEncoder 自定义时间编码器
// 参数:
//   - timeFormat: 时间格式字符串 | Time format string
//
// 返回值:
//   - zapcore.TimeEncoder: 时间编码器 | Time encoder
func customTimeEncoder(timeFormat string) zapcore.TimeEncoder {
	return func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(timeFormat))
	}
}

// ============================ 结构化字段 ============================

// String 结构化日志字符串字段
// 参数:
//   - key: 字段名 | Field name
//   - value: 字段值 | Field value
//
// 返回值:
//   - zap.Field: Zap 字段 | Zap field
func String(key, value string) zap.Field {
	return zap.String(key, value)
}

// Int 结构化日志整型字段
// 参数:
//   - key: 字段名 | Field name
//   - value: 字段值 | Field value
//
// 返回值:
//   - zap.Field: Zap 字段 | Zap field
func Int(key string, value int) zap.Field {
	return zap.Int(key, value)
}

// Uint64 结构化日志无符号整型字段
// 参数:
//   - key: 字段名 | Field name
//   - value: 字段值 | Field value
//
// 返回值:
//   - zap.Field: Zap 字段 | Zap field
func Uint64(key string, value uint64) zap.Field {
	return zap.Uint64(key, value)
}

// Duration 结构化日志时长字段
// 参数:
//   - key: 字段名 | Field name
//   - value: 字段值 | Field value
//
// 返回值:
//   - zap.Field: Zap 字段 | Zap field
func Duration(key string, value time.Duration) zap.Field {
	return zap.Duration(key, value)
}

// ============================ 日志调用方法 ============================

// Debug 输出 Debug 级别日志
// 参数:
//   - args: 日志参数 | Log arguments
func Debug(args ...interface{}) {
	SugarLogger.Debug(args...)
}

// Info 输出 Info 级别日志
// 参数:
//   - args: 日志参数 | Log arguments
func Info(args ...interface{}) {
	SugarLogger.Info(args...)
}

// Warn 输出 Warn 级别日志
// 参数:
//   - args: 日志参数 | Log arguments
func Warn(args ...interface{}) {
	SugarLogger.Warn(args...)
}

// Error 输出 Error 级别日志
// 参数:
//   - args: 日志参数 | Log arguments
func Error(args ...interface{}) {
	SugarLogger.Error(args...)
}

// Fatal 输出 Fatal 级别日志
// 参数:
//   - args: 日志参数 | Log arguments
func Fatal(args ...interface{}) {
	SugarLogger.Fatal(args...)
}

// ============================ 格式化日志 ============================

// Debugf 输出格式化 Debug 级别日志
// 参数:
//   - template: 格式化模板 | Format template
//   - args: 模板参数 | Template arguments
func Debugf(template string, args ...interface{}) {
	SugarLogger.Debugf(template, args...)
}

// Infof 输出格式化 Info 级别日志
// 参数:
//   - template: 格式化模板 | Format template
//   - args: 模板参数 | Template arguments
func Infof(template string, args ...interface{}) {
	SugarLogger.Infof(template, args...)
}

// Warnf 输出格式化 Warn 级别日志
// 参数:
//   - template: 格式化模板 | Format template
//   - args: 模板参数 | Template arguments
func Warnf(template string, args ...interface{}) {
	SugarLogger.Warnf(template, args...)
}

// Errorf 输出格式化 Error 级别日志
// 参数:
//   - template: 格式化模板 | Format template
//   - args: 模板参数 | Template arguments
func Errorf(template string, args ...interface{}) {
	SugarLogger.Errorf(template, args...)
}

// Fatalf 输出格式化 Fatal 级别日志
// 参数:
//   - template: 格式化模板 | Format template
//   - args: 模板参数 | Template arguments
func Fatalf(template string, args ...interface{}) {
	SugarLogger.Fatalf(template, args...)
}

// ============================ zap.Field 结构化日志 ============================

// DebugWithFields 输出带字段的 Debug 级别日志
// 参数:
//   - msg: 日志消息 | Log message
//   - fields: 结构化字段 | Structured fields
func DebugWithFields(msg string, fields ...zap.Field) {
	GlobalLogger.Debug(msg, fields...)
}

// InfoWithFields 输出带字段的 Info 级别日志
// 参数:
//   - msg: 日志消息 | Log message
//   - fields: 结构化字段 | Structured fields
func InfoWithFields(msg string, fields ...zap.Field) {
	GlobalLogger.Info(msg, fields...)
}

// WarnWithFields 输出带字段的 Warn 级别日志
// 参数:
//   - msg: 日志消息 | Log message
//   - fields: 结构化字段 | Structured fields
func WarnWithFields(msg string, fields ...zap.Field) {
	GlobalLogger.Warn(msg, fields...)
}

// ErrorWithFields 输出带字段的 Error 级别日志
// 参数:
//   - msg: 日志消息 | Log message
//   - fields: 结构化字段 | Structured fields
func ErrorWithFields(msg string, fields ...zap.Field) {
	GlobalLogger.Error(msg, fields...)
}

// FatalWithFields 输出带字段的 Fatal 级别日志
// 参数:
//   - msg: 日志消息 | Log message
//   - fields: 结构化字段 | Structured fields
func FatalWithFields(msg string, fields ...zap.Field) {
	GlobalLogger.Fatal(msg, fields...)
}
