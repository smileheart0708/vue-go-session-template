package middleware

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"time"

	"main/internal/stream"

	"github.com/lmittmann/tint"
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
	broadcastHandler := stream.NewBroadcastHandler(broadcaster)

	// 3. 组合处理器
	handler := NewFanoutHandler(consoleHandler, broadcastHandler)

	logger := slog.New(handler)
	return logger
}

// FanoutHandler 分发日志到多个处理器
type FanoutHandler struct {
	handlers []slog.Handler
}

func NewFanoutHandler(handlers ...slog.Handler) *FanoutHandler {
	return &FanoutHandler{handlers: handlers}
}

func (h *FanoutHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (h *FanoutHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, r.Level) {
			_ = handler.Handle(ctx, r.Clone())
		}
	}
	return nil
}

func (h *FanoutHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		handlers[i] = handler.WithAttrs(attrs)
	}
	return &FanoutHandler{handlers: handlers}
}

func (h *FanoutHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		handlers[i] = handler.WithGroup(name)
	}
	return &FanoutHandler{handlers: handlers}
}
