<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import ThemeToggle from '@/components/common/ThemeToggle.vue'
import BaseButton from '@/components/common/BaseButton.vue'
import { useToast } from '@/composables'
import { useAuthStore } from '@/stores/auth'
import { HttpError, http, resolveRedirectPath } from '@/utils'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const authKey = ref('')
const isLoading = ref(false)
const { success, error: toastError } = useToast()

interface LoginResponse {
  success: boolean
  message: string
  session_id?: string
}

function extractErrorMessage(payload: unknown): string | null {
  if (typeof payload === 'string') {
    return payload.trim() || null
  }

  if (!payload || typeof payload !== 'object') {
    return null
  }

  const message = (payload as Record<string, unknown>).message
  if (typeof message !== 'string') {
    return null
  }

  return message.trim() || null
}

const loginRedirectPath = computed(() => resolveRedirectPath(route.query.redirect) ?? '/dashboard')

const handleLogin = async () => {
  if (!authKey.value.trim()) {
    toastError('请输入认证令牌')
    return
  }

  isLoading.value = true

  try {
    const data = await http<LoginResponse>('/login', {
      method: 'POST',
      body: { auth_key: authKey.value },
      skipUnauthorizedHandler: true,
    })

    if (!data.success || !data.session_id) {
      toastError(data.message || '认证失败，请重试')
      return
    }

    // 登录成功
    authStore.setAuthenticated(data.session_id)
    success('登录成功！')

    // 登录后回跳
    await router.replace(loginRedirectPath.value)
  } catch (error) {
    if (error instanceof HttpError) {
      toastError(extractErrorMessage(error.data) || '认证失败，请重试')
      return
    }

    console.error('Login error:', error)
    toastError('网络错误，请重试')
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <div class="login-container">
    <div class="login-card">
      <div class="theme-toggle-wrapper">
        <ThemeToggle />
      </div>
      <h1 class="login-title">身份认证</h1>

      <form @submit.prevent="handleLogin" class="login-form">
        <div class="form-group">
          <label for="authKey" class="form-label">认证令牌</label>
          <input
            id="authKey"
            v-model="authKey"
            type="password"
            class="form-input"
            placeholder="请输入 AUTH_KEY"
            :disabled="isLoading"
          />
        </div>

        <BaseButton
          type="submit"
          width="100%"
          :height="44"
          :text="isLoading ? '认证中...' : '登录'"
          :primary="true"
          :disabled="isLoading"
        />
      </form>
    </div>
  </div>
</template>

<style scoped>
.login-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  padding: 20px;
  position: relative;
}

.theme-toggle-wrapper {
  position: absolute;
  top: 12px;
  right: 12px;
}

.login-card {
  position: relative;
  width: 100%;
  max-width: 400px;
  padding: 40px;
  background: var(--sys-color-bg-surface);
  border: 1px solid var(--sys-color-border);
  border-radius: 8px;
  box-shadow: var(--sys-shadow-card);
}

.login-title {
  margin-bottom: 32px;
  font-size: 24px;
  font-weight: 600;
  text-align: center;
  color: var(--sys-color-text-primary);
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-label {
  font-size: 14px;
  font-weight: 500;
  color: var(--sys-color-text-primary);
}

.form-input {
  padding: 10px 12px;
  border-radius: 4px;
  font-size: 14px;
  transition: all 0.2s;
}

.form-input:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>
