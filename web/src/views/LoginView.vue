<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { BaseButton, ThemeToggle } from '@/components/common'
import { useTheme, useToast } from '@/composables'
import { useAuthStore } from '@/stores/auth'
import { loginResponseSchema } from '@/types/api'
import { HttpError, HttpResponseValidationError, http, resolveRedirectPath } from '@/utils'
import type { ThemeMode } from '@/composables'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const authKey = ref('')
const isLoading = ref(false)
const { success, error: toastError } = useToast()
const { mode, setTheme } = useTheme()

function hasMessageField(payload: unknown): payload is { message: unknown } {
  return typeof payload === 'object' && payload !== null && 'message' in payload
}

function extractErrorMessage(payload: unknown): string | null {
  if (typeof payload === 'string') {
    return payload.trim() || null
  }

  if (!hasMessageField(payload)) {
    return null
  }

  const { message } = payload
  if (typeof message !== 'string') {
    return null
  }

  return message.trim() || null
}

const loginRedirectPath = computed(
  () => resolveRedirectPath(route.query['redirect']) ?? '/dashboard',
)

async function handleThemeChange(nextMode: ThemeMode, event?: MouseEvent): Promise<void> {
  await setTheme(nextMode, event)
}

const handleLogin = async () => {
  if (!authKey.value.trim()) {
    toastError('请输入认证令牌')
    return
  }

  isLoading.value = true

  try {
    const data = await http('/login', {
      method: 'POST',
      body: { auth_key: authKey.value },
      skipUnauthorizedHandler: true,
      schema: loginResponseSchema,
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
    if (error instanceof HttpResponseValidationError) {
      console.error('Invalid login response payload:', error)
      toastError('服务端响应格式异常，请稍后重试')
      return
    }

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
  <div class="relative flex min-h-screen items-center justify-center p-5">
    <div
      class="relative w-full max-w-100 rounded-lg border border-border bg-bg-surface p-10 shadow-card"
    >
      <div class="absolute right-3 top-3">
        <ThemeToggle
          v-model="mode"
          light-label="浅色"
          dark-label="深色"
          auto-label="自动"
          tooltip-prefix="当前："
          tooltip-suffix="（长按切换）"
          @change="handleThemeChange"
        />
      </div>
      <h1 class="mb-8 text-center text-2xl font-semibold text-text-primary">身份认证</h1>

      <form
        @submit.prevent="handleLogin"
        class="flex flex-col gap-5"
      >
        <div class="flex flex-col gap-2">
          <label
            for="authKey"
            class="text-sm font-medium text-text-primary"
            >认证令牌</label
          >
          <input
            id="authKey"
            v-model="authKey"
            type="password"
            class="rounded-sm px-3 py-2.5 text-sm transition-all disabled:cursor-not-allowed disabled:opacity-60"
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
