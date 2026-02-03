import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

const STORAGE_KEY = 'vue-go-session-auth'

export const useAuthStore = defineStore('auth', () => {
  const sessionId = ref<string>('')
  const isAuthenticated = ref<boolean>(false)

  // 从 localStorage 初始化
  function init() {
    const stored = localStorage.getItem(STORAGE_KEY)
    if (stored) {
      try {
        const data = JSON.parse(stored)
        sessionId.value = data.sessionId || ''
        isAuthenticated.value = data.isAuthenticated || false
      } catch (error) {
        console.error('Failed to parse auth data:', error)
        clear()
      }
    }
  }

  // 模拟登录 (仅在开发环境使用)
  function mockLogin() {
    const mockSessionId = `mock-session-${Date.now()}`
    sessionId.value = mockSessionId
    isAuthenticated.value = true
    saveToStorage()
    console.log('[MOCK AUTH] 模拟登录成功:', mockSessionId)
  }

  // 登录成功后设置认证状态
  function setAuthenticated(newSessionId: string) {
    sessionId.value = newSessionId
    isAuthenticated.value = true
    saveToStorage()
  }

  // 清除认证状态
  function clear() {
    sessionId.value = ''
    isAuthenticated.value = false
    localStorage.removeItem(STORAGE_KEY)
  }

  // 保存到 localStorage
  function saveToStorage() {
    const data = {
      sessionId: sessionId.value,
      isAuthenticated: isAuthenticated.value,
    }
    localStorage.setItem(STORAGE_KEY, JSON.stringify(data))
  }

  // 验证 session 是否有效
  async function validateSession(): Promise<boolean> {
    if (!sessionId.value) {
      clear()
      return false
    }

    try {
      const response = await fetch('/api/validate-session', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          session_id: sessionId.value,
        }),
      })

      if (!response.ok) {
        clear()
        return false
      }

      const data = await response.json()
      if (!data.valid) {
        clear()
        return false
      }

      // session 有效，保持认证状态
      isAuthenticated.value = true
      saveToStorage()
      return true
    } catch (error) {
      console.error('Failed to validate session:', error)
      clear()
      return false
    }
  }

  // 登出
  async function logout(): Promise<void> {
    try {
      await fetch('/api/logout', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
      })
    } catch (error) {
      console.error('Failed to logout:', error)
    } finally {
      clear()
    }
  }

  return {
    sessionId: computed(() => sessionId.value),
    isAuthenticated: computed(() => isAuthenticated.value),
    init,
    mockLogin,
    setAuthenticated,
    clear,
    validateSession,
    logout,
  }
})
