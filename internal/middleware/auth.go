package middleware

import (
	"log/slog"
	"net/http"
	"time"

	ginsessions "github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"main/internal/session"
)

// AuthMiddleware 认证中间件
func AuthMiddleware(cookieSecure bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		sess := ginsessions.Default(c)
		authenticated, ok := sess.Get("authenticated").(bool)
		if !ok || !authenticated {
			clearInvalidSessionCookie(sess, cookieSecure)
			slog.Warn("unauthorized request", "remote_addr", c.ClientIP())
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "未授权，请先登录",
			})
			c.Abort()
			return
		}

		sess.Set("last_seen_at", time.Now().Unix())
		session.SetCookieOptions(sess, cookieSecure, session.SessionMaxAgeSeconds)
		if err := sess.Save(); err != nil {
			slog.Warn("failed to refresh session", "error", err)
		}

		sessionID, _ := sess.Get("session_id").(string)
		if sessionID == "" {
			sessionID = "unknown"
		}

		c.Set("session_id", sessionID)
		c.Next()
	}
}

func clearInvalidSessionCookie(sess ginsessions.Session, secure bool) {
	session.ExpireCookie(sess, secure)
	if err := sess.Save(); err != nil {
		slog.Warn("failed to clear invalid session", "error", err)
	}
}
