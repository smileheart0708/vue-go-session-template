package session

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	// SessionDuration session 有效期为 7 天
	SessionDuration = 7 * 24 * time.Hour
	// SessionIDLength session ID 长度（字节数）
	SessionIDLength = 32
)

// Session 会话数据结构
type Session struct {
	ID        string `json:"id"`
	CreatedAt int64  `json:"created_at"` // Unix 秒时间戳
	ExpiresAt int64  `json:"expires_at"` // Unix 秒时间戳
}

// Manager session 管理器
type Manager struct {
	dataDir  string
	sessions map[string]*Session
	mu       sync.RWMutex
}

// NewManager 创建 session 管理器
func NewManager(dataDir string) (*Manager, error) {
	// 确保数据目录存在
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	m := &Manager{
		dataDir:  dataDir,
		sessions: make(map[string]*Session),
	}

	// 加载已有的 session
	count, err := m.loadSessions()
	if err != nil {
		slog.Warn("加载 session 失败", "error", err)
	} else {
		slog.Info("已加载 session", "count", count)
	}

	return m, nil
}

// CreateSession 创建新的 session
func (m *Manager) CreateSession() (*Session, error) {
	sessionID, err := generateSessionID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session ID: %w", err)
	}

	now := time.Now()
	session := &Session{
		ID:        sessionID,
		CreatedAt: now.Unix(),
		ExpiresAt: now.Add(SessionDuration).Unix(),
	}

	m.mu.Lock()
	m.sessions[sessionID] = session
	m.mu.Unlock()

	// 持久化到文件
	if err := m.saveSession(session); err != nil {
		return nil, fmt.Errorf("failed to save session: %w", err)
	}

	slog.Info("session created", "session_id", sessionID, "expires_at", session.ExpiresAt)
	return session, nil
}

// ValidateSession 验证 session 是否有效
func (m *Manager) ValidateSession(sessionID string) bool {
	m.mu.RLock()
	session, exists := m.sessions[sessionID]
	m.mu.RUnlock()

	if !exists {
		return false
	}

	// 检查是否过期
	if time.Now().Unix() > session.ExpiresAt {
		// 删除过期的 session
		m.DeleteSession(sessionID)
		return false
	}

	return true
}

// RefreshSession 刷新 session 有效期
func (m *Manager) RefreshSession(sessionID string) error {
	m.mu.Lock()
	session, exists := m.sessions[sessionID]
	if !exists {
		m.mu.Unlock()
		return fmt.Errorf("session not found")
	}

	// 更新过期时间
	session.ExpiresAt = time.Now().Add(SessionDuration).Unix()
	m.mu.Unlock()

	// 持久化更新
	if err := m.saveSession(session); err != nil {
		return fmt.Errorf("failed to refresh session: %w", err)
	}

	return nil
}

// DeleteSession 删除 session
func (m *Manager) DeleteSession(sessionID string) {
	m.mu.Lock()
	delete(m.sessions, sessionID)
	m.mu.Unlock()

	// 删除文件
	filePath := m.getSessionFilePath(sessionID)
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		slog.Warn("failed to delete session file", "session_id", sessionID, "error", err)
	}

	slog.Info("session deleted", "session_id", sessionID)
}

// CleanExpiredSessions 清理过期的 session
func (m *Manager) CleanExpiredSessions() {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now().Unix()
	for id, session := range m.sessions {
		if now > session.ExpiresAt {
			delete(m.sessions, id)
			filePath := m.getSessionFilePath(id)
			if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
				slog.Warn("failed to delete expired session file", "session_id", id, "error", err)
			}
			slog.Info("expired session cleaned", "session_id", id)
		}
	}
}

// saveSession 保存 session 到文件
func (m *Manager) saveSession(session *Session) error {
	filePath := m.getSessionFilePath(session.ID)

	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0600); err != nil {
		return fmt.Errorf("failed to write session file: %w", err)
	}

	return nil
}

// loadSessions 从文件加载所有 session
func (m *Manager) loadSessions() (int, error) {
	files, err := os.ReadDir(m.dataDir)
	if err != nil {
		return 0, fmt.Errorf("failed to read data directory: %w", err)
	}

	now := time.Now().Unix()
	count := 0
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}

		filePath := filepath.Join(m.dataDir, file.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			slog.Warn("failed to read session file", "file", file.Name(), "error", err)
			continue
		}

		var session Session
		if err := json.Unmarshal(data, &session); err != nil {
			slog.Warn("failed to unmarshal session", "file", file.Name(), "error", err)
			continue
		}

		// 跳过已过期的 session
		if now > session.ExpiresAt {
			os.Remove(filePath)
			slog.Info("expired session file removed", "session_id", session.ID)
			continue
		}

		m.sessions[session.ID] = &session
		count++
	}

	return count, nil
}

// getSessionFilePath 获取 session 文件路径
func (m *Manager) getSessionFilePath(sessionID string) string {
	return filepath.Join(m.dataDir, fmt.Sprintf("%s.json", sessionID))
}

// generateSessionID 生成随机 session ID
func generateSessionID() (string, error) {
	bytes := make([]byte, SessionIDLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
