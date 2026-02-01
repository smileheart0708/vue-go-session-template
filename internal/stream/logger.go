package stream

import (
	"context"
	"encoding/json"
	"log/slog"
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
	clients map[chan []byte]struct{}
	history []LogEntry
	mu      sync.RWMutex
}

func NewLogBroadcaster() *LogBroadcaster {
	return &LogBroadcaster{
		clients: make(map[chan []byte]struct{}),
		history: make([]LogEntry, 0, maxHistoryLogs),
	}
}

// Subscribe 添加一个新的客户端连接
func (b *LogBroadcaster) Subscribe() chan []byte {
	b.mu.Lock()
	defer b.mu.Unlock()
	ch := make(chan []byte, 100) // 带缓冲，防止阻塞
	b.clients[ch] = struct{}{}
	return ch
}

// Unsubscribe 移除客户端连接
func (b *LogBroadcaster) Unsubscribe(ch chan []byte) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.clients[ch]; ok {
		delete(b.clients, ch)
		close(ch)
	}
}

// Broadcast 发送日志给所有连接的客户端
func (b *LogBroadcaster) Broadcast(record slog.Record) {
	// 转换为 JSON
	entry := LogEntry{
		Time:    record.Time.Format(time.DateTime),
		Level:   record.Level.String(),
		Message: record.Message,
		Attrs:   make(map[string]any),
	}

	record.Attrs(func(a slog.Attr) bool {
		entry.Attrs[a.Key] = a.Value.Any()
		return true
	})

	data, err := json.Marshal(entry)
	if err != nil {
		return
	}

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
		case ch <- data:
		default:
			// 如果客户端阻塞，跳过该消息，避免影响其他客户端
		}
	}
}

// BroadcastHandler 实现 slog.Handler 接口
type BroadcastHandler struct {
	broadcaster *LogBroadcaster
	attrs       []slog.Attr
	group       string
}

func NewBroadcastHandler(b *LogBroadcaster) *BroadcastHandler {
	return &BroadcastHandler{broadcaster: b}
}

func (h *BroadcastHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

func (h *BroadcastHandler) Handle(_ context.Context, r slog.Record) error {
	// 复制 record 并添加预设属性（如果有）
	// 注意：为了简化，这里直接广播原始 record，属性处理在 Broadcast 方法中进行
	// 实际生产中可能需要更复杂的属性合并逻辑
	h.broadcaster.Broadcast(r)
	return nil
}

func (h *BroadcastHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &BroadcastHandler{
		broadcaster: h.broadcaster,
		attrs:       append(h.attrs, attrs...),
	}
}

func (h *BroadcastHandler) WithGroup(name string) slog.Handler {
	// 简单实现，暂不支持 Group
	return h
}

// GetHistory 获取历史日志记录
func (b *LogBroadcaster) GetHistory() []LogEntry {
	b.mu.RLock()
	defer b.mu.RUnlock()

	history := make([]LogEntry, len(b.history))
	copy(history, b.history)
	return history
}
