import { computed } from 'vue'
import { useLocalStorage } from '@vueuse/core'
import { defineStore } from 'pinia'
import { useDashboardStore } from './dashboard'
import { HttpError, http, setUnauthorizedHandler } from '@/utils'

const STORAGE_KEY = 'vue-go-session-auth'

interface StoredAuthState {
  sessionId: string
  isAuthenticated: boolean
}

interface ValidateSessionResponse {
  valid: boolean
}

const DEFAULT_AUTH_STATE: StoredAuthState = {
  sessionId: '',
  isAuthenticated: false,
}

function isValidStoredAuthState(value: unknown): value is StoredAuthState {
  if (!value || typeof value !== 'object') {
    return false
  }

  const state = value as Partial<StoredAuthState>
  return typeof state.sessionId === 'string' && typeof state.isAuthenticated === 'boolean'
}

export const useAuthStore = defineStore('auth', () => {
  const storedState = useLocalStorage<StoredAuthState>(STORAGE_KEY, {
    ...DEFAULT_AUTH_STATE,
  })

  const sessionId = computed(() => storedState.value.sessionId)
  const isAuthenticated = computed(() => storedState.value.isAuthenticated)

  function applyState(nextState: StoredAuthState): void {
    storedState.value = { ...nextState }
  }

  function resetUiState(): void {
    const dashboardStore = useDashboardStore()
    dashboardStore.reset()
  }

  function clear(): void {
    applyState(DEFAULT_AUTH_STATE)
    resetUiState()
  }

  function handleUnauthorized(): void {
    clear()
  }

  function init(): void {
    if (!isValidStoredAuthState(storedState.value)) {
      applyState(DEFAULT_AUTH_STATE)
    }
    setUnauthorizedHandler(handleUnauthorized)
  }

  function mockLogin(): void {
    const mockSessionId = `mock-session-${Date.now()}`
    applyState({
      sessionId: mockSessionId,
      isAuthenticated: true,
    })
    console.log('[MOCK AUTH] 模拟登录成功:', mockSessionId)
  }

  function setAuthenticated(newSessionId: string): void {
    applyState({
      sessionId: newSessionId,
      isAuthenticated: true,
    })
  }

  async function validateSession(): Promise<boolean> {
    if (!sessionId.value) {
      handleUnauthorized()
      return false
    }

    try {
      const data = await http<ValidateSessionResponse>('/validate-session', {
        method: 'POST',
        body: {
          session_id: sessionId.value,
        },
        skipUnauthorizedHandler: true,
      })

      if (!data.valid) {
        handleUnauthorized()
        return false
      }

      applyState({
        sessionId: sessionId.value,
        isAuthenticated: true,
      })
      return true
    } catch (error) {
      if (!(error instanceof HttpError)) {
        console.error('Failed to validate session:', error)
      }
      handleUnauthorized()
      return false
    }
  }

  async function logout(): Promise<void> {
    try {
      await http('/logout', {
        method: 'POST',
        responseType: 'response',
        skipUnauthorizedHandler: true,
      })
    } catch (error) {
      if (!(error instanceof HttpError && error.status === 401)) {
        console.error('Failed to logout:', error)
      }
    } finally {
      clear()
    }
  }

  return {
    sessionId,
    isAuthenticated,
    init,
    mockLogin,
    setAuthenticated,
    clear,
    handleUnauthorized,
    validateSession,
    logout,
  }
})
