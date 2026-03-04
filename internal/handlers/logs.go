package handlers

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"

	"main/internal/stream"
)

// LogsHandler 日志处理器
type LogsHandler struct {
	broadcaster *stream.LogBroadcaster
}

// NewLogsHandler 创建日志处理器
func NewLogsHandler(broadcaster *stream.LogBroadcaster) *LogsHandler {
	return &LogsHandler{
		broadcaster: broadcaster,
	}
}

// StreamLogs 处理 SSE 日志流请求
func (h *LogsHandler) StreamLogs(c *gin.Context) {
	// 设置 SSE 响应头
	c.Writer.Header().Set("Content-Type", sse.ContentType)
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no") // 禁用 Nginx 缓冲

	// 订阅日志广播
	ch := h.broadcaster.Subscribe()
	defer h.broadcaster.Unsubscribe(ch)

	sessionID, _ := c.Get("session_id")
	slog.Info("客户端连接日志流", "session_id", sessionID, "remote_addr", c.ClientIP())

	sendHistory := c.DefaultQuery("history", "1") != "0"
	if sendHistory {
		// 发送历史日志
		history := h.broadcaster.GetHistory()
		for _, entry := range history {
			if err := sse.Encode(c.Writer, sse.Event{Data: entry}); err != nil {
				return
			}
		}
		c.Writer.Flush()
	}

	// 持续推送新日志
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.Request.Context().Done():
			return

		case logEntry, ok := <-ch:
			if !ok {
				// 通道已关闭
				return
			}

			if err := sse.Encode(c.Writer, sse.Event{Data: logEntry}); err != nil {
				return
			}
			c.Writer.Flush()

		case <-ticker.C:
			// 发送心跳，保持连接
			if err := sse.Encode(c.Writer, sse.Event{Event: "heartbeat", Data: "ping"}); err != nil {
				return
			}
			c.Writer.Flush()
		}
	}
}

// GetHistory 获取历史日志
func (h *LogsHandler) GetHistory(c *gin.Context) {
	history := h.broadcaster.GetHistory()
	c.JSON(http.StatusOK, gin.H{
		"logs":  history,
		"count": len(history),
	})
}
