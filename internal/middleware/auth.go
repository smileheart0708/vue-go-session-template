package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"main/internal/session"
)

// AuthMiddleware 认证中间件
func AuthMiddleware(sessionManager *session.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Cookie 中获取 session_id
		sessionID, err := c.Cookie("session_id")
		if err != nil || sessionID == "" {
			slog.Warn("未提供 session_id", "remote_addr", c.ClientIP())
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "未授权，请先登录",
			})
			c.Abort()
			return
		}

		// 验证 session 是否有效
		if !sessionManager.ValidateSession(sessionID) {
			slog.Warn("无效的 session_id", "session_id", sessionID, "remote_addr", c.ClientIP())
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "会话已过期，请重新登录",
			})
			c.Abort()
			return
		}

		// 刷新 session 有效期
		if err := sessionManager.RefreshSession(sessionID); err != nil {
			slog.Warn("刷新 session 失败", "session_id", sessionID, "error", err)
		}

		// 将 session_id 存入上下文，供后续使用
		c.Set("session_id", sessionID)
		c.Next()
	}
}
