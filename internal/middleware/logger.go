package middleware

import (
	"log/slog"
	"os"
	"strings"
	"time"

	"main/internal/stream"

	"github.com/lmittmann/tint"
	slogmulti "github.com/samber/slog-multi"
)

// ParseLogLevel 解析日志等级字符串
func ParseLogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// InitLogger 初始化日志系统
func InitLogger(logLevel string, broadcaster *stream.LogBroadcaster) *slog.Logger {
	level := ParseLogLevel(logLevel)

	// 1. CLI 日志处理器 (Tint)
	consoleHandler := tint.NewHandler(os.Stdout, &tint.Options{
		Level:      level,
		TimeFormat: time.DateTime,
	})

	// 2. SSE 广播处理器
	broadcastHandler := stream.NewSSELogHandler(level, broadcaster)

	// 3. 组合处理器
	handler := slogmulti.Fanout(consoleHandler, broadcastHandler)

	logger := slog.New(handler)
	return logger
}
