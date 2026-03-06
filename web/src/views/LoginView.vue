<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { z } from 'zod'
import { BaseButton, ThemeToggle } from '@/components/common'
import { useTheme, useToast } from '@/composables'
import { useAuthStore } from '@/stores/auth'
import { loginResponseSchema } from '@/types/api'
import {
  api,
  ApiResponseValidationError,
  HTTPError,
  normalizeApiEndpoint,
  parseWithSchema,
  readHttpErrorData,
  resolveRedirectPath,
  withUnauthorizedHandlerSkipped,
} from '@/utils'
import type { ThemeMode } from '@/composables'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const authKey = ref('')
const isLoading = ref(false)
const { toast } = useToast()
const { mode, setTheme } = useTheme()

const errorMessageSchema = z.union([z.string(), z.object({ message: z.string() })])

function extractErrorMessage(payload: unknown): string | null {
  const parsedPayload = errorMessageSchema.safeParse(payload)
  if (!parsedPayload.success) {
    return null
  }

  const message =
    typeof parsedPayload.data === 'string' ? parsedPayload.data : parsedPayload.data.message

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
    toast.error('请输入认证令牌')
    return
  }

  isLoading.value = true

  try {
    const response = await api.post(
      normalizeApiEndpoint('/login'),
      withUnauthorizedHandlerSkipped({ json: { auth_key: authKey.value } }),
    )
    const payload = await response.json<unknown>()
    const data = parseWithSchema(payload, loginResponseSchema, response.url)

    if (!data.success) {
      toast.error(data.message || '认证失败，请重试')
      return
    }

    // 登录成功
    authStore.setAuthenticated()
    toast.success('登录成功！')

    // 登录后回跳
    await router.replace(loginRedirectPath.value)
  } catch (error) {
    if (error instanceof ApiResponseValidationError) {
      console.error('Invalid login response payload:', error)
      toast.error('服务端响应格式异常，请稍后重试')
      return
    }

    if (error instanceof HTTPError) {
      const errorPayload = await readHttpErrorData(error)
      toast.error(extractErrorMessage(errorPayload) || '认证失败，请重试')
      return
    }

    console.error('Login error:', error)
    toast.error('网络错误，请重试')
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
