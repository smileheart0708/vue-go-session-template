package stream

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"sync"
	"time"
)

// LogEntry 定义发送给前端的日志结构
type LogEntry struct {
	Time    string         `json:"time"`
	Level   string         `json:"level"`
	Message string         `json:"msg"`
	Attrs   map[string]any `json:"attrs,omitempty"`
}

const maxHistoryLogs = 100

// LogBroadcaster 管理 SSE 连接和日志分发
type LogBroadcaster struct {
	clients map[chan LogEntry]struct{}
	history []LogEntry
	mu      sync.RWMutex
}

func NewLogBroadcaster() *LogBroadcaster {
	return &LogBroadcaster{
		clients: make(map[chan LogEntry]struct{}),
		history: make([]LogEntry, 0, maxHistoryLogs),
	}
}

// Subscribe 添加一个新的客户端连接
func (b *LogBroadcaster) Subscribe() chan LogEntry {
	b.mu.Lock()
	defer b.mu.Unlock()
	ch := make(chan LogEntry, 100) // 带缓冲，防止阻塞
	b.clients[ch] = struct{}{}
	return ch
}

// Unsubscribe 移除客户端连接
func (b *LogBroadcaster) Unsubscribe(ch chan LogEntry) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.clients[ch]; ok {
		delete(b.clients, ch)
		close(ch)
	}
}

// Broadcast 发送日志给所有连接的客户端
func (b *LogBroadcaster) Broadcast(entry LogEntry) {
	b.mu.Lock()
	// 保存到历史记录
	b.history = append(b.history, entry)
	if len(b.history) > maxHistoryLogs {
		b.history = b.history[len(b.history)-maxHistoryLogs:]
	}
	b.mu.Unlock()

	b.mu.RLock()
	defer b.mu.RUnlock()

	for ch := range b.clients {
		select {
		case ch <- entry:
		default:
			// 如果客户端阻塞，跳过该消息，避免影响其他客户端
		}
	}
}

// JSONLogWriter 接收 slog.JSONHandler 输出并广播给前端
type JSONLogWriter struct {
	broadcaster *LogBroadcaster
	buffer      bytes.Buffer
	mu          sync.Mutex
}

func NewJSONLogWriter(broadcaster *LogBroadcaster) io.Writer {
	return &JSONLogWriter{
		broadcaster: broadcaster,
	}
}

func (w *JSONLogWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if _, err := w.buffer.Write(p); err != nil {
		return 0, err
	}

	for {
		line, ok := w.nextLine()
		if !ok {
			break
		}
		entry, parsed := parseJSONLogLine(line)
		if !parsed {
			continue
		}
		w.broadcaster.Broadcast(entry)
	}

	return len(p), nil
}

func (w *JSONLogWriter) nextLine() ([]byte, bool) {
	bufferBytes := w.buffer.Bytes()
	newlineIndex := bytes.IndexByte(bufferBytes, '\n')
	if newlineIndex < 0 {
		return nil, false
	}

	line := strings.TrimSpace(string(bufferBytes[:newlineIndex]))
	w.buffer.Next(newlineIndex + 1)
	return []byte(line), true
}

func parseJSONLogLine(line []byte) (LogEntry, bool) {
	if len(line) == 0 {
		return LogEntry{}, false
	}

	var payload map[string]any
	if err := json.Unmarshal(line, &payload); err != nil {
		return LogEntry{}, false
	}

	entry := LogEntry{
		Time:    normalizeEntryTime(payload["time"]),
		Level:   normalizeStringField(payload["level"], "INFO"),
		Message: normalizeStringField(payload["msg"], ""),
	}

	delete(payload, "time")
	delete(payload, "level")
	delete(payload, "msg")
	if len(payload) > 0 {
		entry.Attrs = payload
	}

	return entry, true
}

func normalizeStringField(raw any, fallback string) string {
	switch value := raw.(type) {
	case string:
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			return fallback
		}
		return trimmed
	case nil:
		return fallback
	default:
		return fmt.Sprint(value)
	}
}

func normalizeEntryTime(raw any) string {
	timeString := normalizeStringField(raw, "")
	if timeString == "" {
		return time.Now().Format(time.DateTime)
	}

	parsed, err := time.Parse(time.RFC3339Nano, timeString)
	if err == nil {
		return parsed.Local().Format(time.DateTime)
	}

	parsed, err = time.Parse(time.DateTime, timeString)
	if err == nil {
		return parsed.Format(time.DateTime)
	}

	return timeString
}

// GetHistory 获取历史日志记录
func (b *LogBroadcaster) GetHistory() []LogEntry {
	b.mu.RLock()
	defer b.mu.RUnlock()

	history := make([]LogEntry, len(b.history))
	copy(history, b.history)
	return history
}

func buildJSONStreamHandler(level slog.Level, writer io.Writer) slog.Handler {
	return slog.NewJSONHandler(writer, &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(_ []string, attr slog.Attr) slog.Attr {
			if attr.Key == slog.TimeKey {
				if timestamp, ok := attr.Value.Any().(time.Time); ok {
					return slog.String(slog.TimeKey, timestamp.Local().Format(time.DateTime))
				}
			}
			if attr.Value.Kind() == slog.KindAny {
				if err, ok := attr.Value.Any().(error); ok {
					return slog.String(attr.Key, err.Error())
				}
			}
			return attr
		},
	})
}

// NewSSELogHandler 创建用于 SSE 日志流的 slog.Handler
func NewSSELogHandler(level slog.Level, broadcaster *LogBroadcaster) slog.Handler {
	return buildJSONStreamHandler(level, NewJSONLogWriter(broadcaster))
}
