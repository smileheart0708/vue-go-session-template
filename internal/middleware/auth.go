package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const sessionMaxAgeSeconds = 7 * 24 * 60 * 60

// AuthMiddleware 认证中间件
func AuthMiddleware(cookieSecure bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		sess := sessions.Default(c)
		authenticated, ok := sess.Get("authenticated").(bool)
		if !ok || !authenticated {
			slog.Warn("unauthorized request", "remote_addr", c.ClientIP())
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "未授权，请先登录",
			})
			c.Abort()
			return
		}

		sess.Set("last_seen_at", time.Now().Unix())
		sess.Options(sessions.Options{
			Path:     "/",
			MaxAge:   sessionMaxAgeSeconds,
			HttpOnly: true,
			Secure:   cookieSecure,
			SameSite: http.SameSiteLaxMode,
		})
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
