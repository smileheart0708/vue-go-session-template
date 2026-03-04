import { computed } from 'vue'
import { z } from 'zod'
import { defineStore } from 'pinia'
import { useValidatedLocalStorage } from '@/composables/useValidatedLocalStorage'
import { logoutResponseSchema, sessionStatusResponseSchema } from '@/types/api'
import { useDashboardStore } from './dashboard'
import {
  api,
  ApiResponseValidationError,
  HTTPError,
  normalizeApiEndpoint,
  parseWithSchema,
  setUnauthorizedHandler,
  withUnauthorizedHandlerSkipped,
} from '@/utils'
import { isMockAuthEnabled } from '@/utils/env'

const STORAGE_KEY = 'vue-go-session-auth'

const storedAuthStateSchema = z
  .object({
    isAuthenticated: z.boolean().catch(false),
  })
  .transform((value) => ({
    isAuthenticated: value.isAuthenticated === true,
  }))

type StoredAuthState = z.infer<typeof storedAuthStateSchema>

const DEFAULT_AUTH_STATE: StoredAuthState = { isAuthenticated: false }

export const useAuthStore = defineStore('auth', () => {
  const storedState = useValidatedLocalStorage(
    STORAGE_KEY,
    storedAuthStateSchema,
    DEFAULT_AUTH_STATE,
  )
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
    applyState(storedState.value)
    setUnauthorizedHandler(handleUnauthorized)
  }

  function mockLogin(): void {
    applyState({ isAuthenticated: true })
    console.log('[MOCK AUTH] 模拟登录成功')
  }

  function setAuthenticated(): void {
    applyState({ isAuthenticated: true })
  }

  async function validateSession(): Promise<boolean> {
    if (isMockAuthEnabled) {
      return isAuthenticated.value
    }

    try {
      const response = await api.get(
        normalizeApiEndpoint('/session'),
        withUnauthorizedHandlerSkipped(),
      )
      const payload = await response.json<unknown>()
      const data = parseWithSchema(payload, sessionStatusResponseSchema, response.url)

      if (!data.authenticated) {
        handleUnauthorized()
        return false
      }

      applyState({ isAuthenticated: true })
      return true
    } catch (error) {
      if (error instanceof ApiResponseValidationError) {
        console.error('Invalid session status response payload:', error)
      } else if (!(error instanceof HTTPError && error.response.status === 401)) {
        console.error('Failed to validate session:', error)
      }
      handleUnauthorized()
      return false
    }
  }

  async function logout(): Promise<void> {
    try {
      const response = await api.post(
        normalizeApiEndpoint('/logout'),
        withUnauthorizedHandlerSkipped(),
      )
      const payload = await response.json<unknown>()
      const data = parseWithSchema(payload, logoutResponseSchema, response.url)
      if (!data.success) {
        console.warn('Logout response indicates failure:', data)
      }
    } catch (error) {
      if (error instanceof ApiResponseValidationError) {
        console.error('Invalid logout response payload:', error)
        return
      }
      if (!(error instanceof HTTPError && error.response.status === 401)) {
        console.error('Failed to logout:', error)
      }
    } finally {
      clear()
    }
  }

  return {
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
