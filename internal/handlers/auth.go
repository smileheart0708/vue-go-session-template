package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const sessionMaxAgeSeconds = 7 * 24 * 60 * 60

// AuthHandler 认证处理器
type AuthHandler struct {
	authKey      string
	cookieSecure bool
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(authKey string, cookieSecure bool) *AuthHandler {
	return &AuthHandler{
		authKey:      authKey,
		cookieSecure: cookieSecure,
	}
}

// LoginRequest 登录请求结构
type LoginRequest struct {
	AuthKey string `json:"auth_key" binding:"required"`
}

// LoginResponse 登录响应结构
type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Login 处理登录请求
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("invalid login request", "error", err)
		c.JSON(http.StatusBadRequest, LoginResponse{
			Success: false,
			Message: "请求格式错误",
		})
		return
	}

	// 验证 AUTH_KEY
	if req.AuthKey != h.authKey {
		slog.Warn("login failed: invalid auth key", "remote_addr", c.ClientIP())
		c.JSON(http.StatusUnauthorized, LoginResponse{
			Success: false,
			Message: "认证失败，请检查令牌是否正确",
		})
		return
	}

	sessionID, err := generateSessionID()
	if err != nil {
		slog.Error("failed to generate session ID", "error", err)
		c.JSON(http.StatusInternalServerError, LoginResponse{
			Success: false,
			Message: "服务器内部错误",
		})
		return
	}

	sess := sessions.Default(c)
	now := time.Now().Unix()
	sess.Set("authenticated", true)
	sess.Set("session_id", sessionID)
	sess.Set("login_at", now)
	sess.Set("last_seen_at", now)
	setSessionCookieOptions(sess, h.cookieSecure, sessionMaxAgeSeconds)
	if err := sess.Save(); err != nil {
		slog.Error("failed to save session", "error", err)
		c.JSON(http.StatusInternalServerError, LoginResponse{
			Success: false,
			Message: "服务器内部错误",
		})
		return
	}

	slog.Info("user logged in", "session_id", sessionID, "remote_addr", c.ClientIP())

	c.JSON(http.StatusOK, LoginResponse{
		Success: true,
		Message: "登录成功",
	})
}

// SessionStatusResponse 会话状态响应
type SessionStatusResponse struct {
	Authenticated bool   `json:"authenticated"`
	Message       string `json:"message,omitempty"`
}

// Session 验证当前会话是否有效
func (h *AuthHandler) Session(c *gin.Context) {
	sess := sessions.Default(c)
	authenticated, ok := sess.Get("authenticated").(bool)
	if !ok || !authenticated {
		c.JSON(http.StatusUnauthorized, SessionStatusResponse{
			Authenticated: false,
			Message:       "未授权，请先登录",
		})
		return
	}

	sess.Set("last_seen_at", time.Now().Unix())
	setSessionCookieOptions(sess, h.cookieSecure, sessionMaxAgeSeconds)
	if err := sess.Save(); err != nil {
		slog.Warn("failed to refresh session", "error", err)
	}

	c.JSON(http.StatusOK, SessionStatusResponse{
		Authenticated: true,
	})
}

// Logout 处理登出请求
func (h *AuthHandler) Logout(c *gin.Context) {
	sess := sessions.Default(c)
	sessionID, _ := sess.Get("session_id").(string)

	sess.Clear()
	setSessionCookieOptions(sess, h.cookieSecure, -1)
	if err := sess.Save(); err != nil {
		slog.Error("failed to clear session", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "服务器内部错误",
		})
		return
	}

	slog.Info("user logged out", "session_id", sessionID, "remote_addr", c.ClientIP())

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "登出成功",
	})
}

func setSessionCookieOptions(sess sessions.Session, secure bool, maxAge int) {
	sess.Options(sessions.Options{
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})
}

func generateSessionID() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
