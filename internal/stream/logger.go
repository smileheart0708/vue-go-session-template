package stream

import (
	"context"
	"encoding/json"
	"fmt"
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
func (b *LogBroadcaster) Broadcast(entry LogEntry) {
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
	attrs       []scopedAttr
	groups      []string
}

type scopedAttr struct {
	attr   slog.Attr
	groups []string
}

func NewBroadcastHandler(b *LogBroadcaster) *BroadcastHandler {
	return &BroadcastHandler{broadcaster: b}
}

func (h *BroadcastHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

func (h *BroadcastHandler) Handle(_ context.Context, r slog.Record) error {
	entry := LogEntry{
		Time:    r.Time.Format(time.DateTime),
		Level:   r.Level.String(),
		Message: r.Message,
	}

	attrs := make(map[string]any)

	for _, scoped := range h.attrs {
		appendAttr(attrs, scoped.groups, scoped.attr)
	}

	r.Attrs(func(attr slog.Attr) bool {
		appendAttr(attrs, h.groups, attr)
		return true
	})

	if len(attrs) > 0 {
		entry.Attrs = attrs
	}

	h.broadcaster.Broadcast(entry)
	return nil
}

func (h *BroadcastHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	mergedAttrs := make([]scopedAttr, len(h.attrs), len(h.attrs)+len(attrs))
	copy(mergedAttrs, h.attrs)
	for _, attr := range attrs {
		attrGroups := make([]string, len(h.groups))
		copy(attrGroups, h.groups)
		mergedAttrs = append(mergedAttrs, scopedAttr{
			attr:   attr,
			groups: attrGroups,
		})
	}

	groups := make([]string, len(h.groups))
	copy(groups, h.groups)

	return &BroadcastHandler{
		broadcaster: h.broadcaster,
		attrs:       mergedAttrs,
		groups:      groups,
	}
}

func (h *BroadcastHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}

	attrs := make([]scopedAttr, len(h.attrs))
	copy(attrs, h.attrs)

	groups := make([]string, len(h.groups), len(h.groups)+1)
	copy(groups, h.groups)
	groups = append(groups, name)

	return &BroadcastHandler{
		broadcaster: h.broadcaster,
		attrs:       attrs,
		groups:      groups,
	}
}

// GetHistory 获取历史日志记录
func (b *LogBroadcaster) GetHistory() []LogEntry {
	b.mu.RLock()
	defer b.mu.RUnlock()

	history := make([]LogEntry, len(b.history))
	copy(history, b.history)
	return history
}

func appendAttr(target map[string]any, groups []string, attr slog.Attr) {
	value := attr.Value.Resolve()

	if value.Kind() == slog.KindGroup {
		nextGroups := groups
		if attr.Key != "" {
			nextGroups = appendPath(groups, attr.Key)
		}
		for _, nestedAttr := range value.Group() {
			appendAttr(target, nextGroups, nestedAttr)
		}
		return
	}

	if attr.Key == "" {
		return
	}

	setNestedValue(target, groups, attr.Key, normalizeSlogValue(value))
}

func appendPath(path []string, value string) []string {
	next := make([]string, len(path)+1)
	copy(next, path)
	next[len(path)] = value
	return next
}

func setNestedValue(target map[string]any, groups []string, key string, value any) {
	current := target
	for _, group := range groups {
		nested, ok := current[group].(map[string]any)
		if !ok {
			nested = make(map[string]any)
			current[group] = nested
		}
		current = nested
	}
	current[key] = value
}

func normalizeSlogValue(value slog.Value) any {
	switch value.Kind() {
	case slog.KindString:
		return value.String()
	case slog.KindInt64:
		return value.Int64()
	case slog.KindUint64:
		return value.Uint64()
	case slog.KindFloat64:
		return value.Float64()
	case slog.KindBool:
		return value.Bool()
	case slog.KindDuration:
		return value.Duration().String()
	case slog.KindTime:
		return value.Time().Format(time.RFC3339Nano)
	case slog.KindAny:
		return normalizeAny(value.Any())
	case slog.KindLogValuer:
		return normalizeSlogValue(value.Resolve())
	case slog.KindGroup:
		grouped := make(map[string]any)
		for _, nestedAttr := range value.Group() {
			appendAttr(grouped, nil, nestedAttr)
		}
		return grouped
	default:
		return value.String()
	}
}

func normalizeAny(raw any) any {
	if raw == nil {
		return nil
	}

	if err, ok := raw.(error); ok {
		return err.Error()
	}

	bytes, err := json.Marshal(raw)
	if err != nil {
		return fmt.Sprint(raw)
	}

	var normalized any
	if err := json.Unmarshal(bytes, &normalized); err != nil {
		return string(bytes)
	}

	return normalized
}
