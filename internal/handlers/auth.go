package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"main/internal/session"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authKey        string
	sessionManager *session.Manager
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(authKey string, sessionManager *session.Manager) *AuthHandler {
	return &AuthHandler{
		authKey:        authKey,
		sessionManager: sessionManager,
	}
}

// LoginRequest 登录请求结构
type LoginRequest struct {
	AuthKey string `json:"auth_key" binding:"required"`
}

// LoginResponse 登录响应结构
type LoginResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	SessionID string `json:"session_id,omitempty"`
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

	// 创建 session
	sess, err := h.sessionManager.CreateSession()
	if err != nil {
		slog.Error("failed to create session", "error", err)
		c.JSON(http.StatusInternalServerError, LoginResponse{
			Success: false,
			Message: "服务器内部错误",
		})
		return
	}

	// 设置 cookie
	c.SetCookie(
		"session_id",
		sess.ID,
		int(session.SessionDuration.Seconds()),
		"/",
		"",
		false, // secure (生产环境应设为 true)
		true,  // httpOnly
	)

	slog.Info("user logged in", "session_id", sess.ID, "remote_addr", c.ClientIP())

	c.JSON(http.StatusOK, LoginResponse{
		Success:   true,
		Message:   "登录成功",
		SessionID: sess.ID,
	})
}

// ValidateSessionRequest 验证 session 请求结构
type ValidateSessionRequest struct {
	SessionID string `json:"session_id" binding:"required"`
}

// ValidateSessionResponse 验证 session 响应结构
type ValidateSessionResponse struct {
	Valid bool `json:"valid"`
}

// ValidateSession 验证 session 是否有效
func (h *AuthHandler) ValidateSession(c *gin.Context) {
	var req ValidateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ValidateSessionResponse{
			Valid: false,
		})
		return
	}

	valid := h.sessionManager.ValidateSession(req.SessionID)

	// 如果有效，刷新 session 有效期
	if valid {
		if err := h.sessionManager.RefreshSession(req.SessionID); err != nil {
			slog.Warn("failed to refresh session", "session_id", req.SessionID, "error", err)
		}
	}

	c.JSON(http.StatusOK, ValidateSessionResponse{
		Valid: valid,
	})
}

// Logout 处理登出请求
func (h *AuthHandler) Logout(c *gin.Context) {
	sessionID, err := c.Cookie("session_id")
	if err == nil && sessionID != "" {
		h.sessionManager.DeleteSession(sessionID)
	}

	// 清除 cookie
	c.SetCookie(
		"session_id",
		"",
		-1,
		"/",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "登出成功",
	})
}
