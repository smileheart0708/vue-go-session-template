<script setup lang="ts">
import { ref } from 'vue'
import ThemeToggle from '@/components/common/ThemeToggle.vue'

const authKey = ref('')
const isLoading = ref(false)
const errorMessage = ref('')

const handleLogin = async () => {
  if (!authKey.value.trim()) {
    errorMessage.value = '请输入认证令牌'
    return
  }

  isLoading.value = true
  errorMessage.value = ''

  try {
    // TODO: 实现实际的认证逻辑
    console.log('认证令牌:', authKey.value)
    // 模拟 API 调用
    await new Promise((resolve) => setTimeout(resolve, 1000))
  } catch (error) {
    console.error(error)
    errorMessage.value = '认证失败，请重试'
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <div class="login-container">
    <div class="theme-toggle-wrapper">
      <ThemeToggle />
    </div>
    <div class="login-card">
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

        <div v-if="errorMessage" class="error-message">
          {{ errorMessage }}
        </div>

        <button type="submit" class="login-button" :disabled="isLoading">
          {{ isLoading ? '认证中...' : '登录' }}
        </button>
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
  top: 20px;
  right: 20px;
}

.login-card {
  width: 100%;
  max-width: 400px;
  padding: 40px;
  background: var(--color-background-elevated);
  border: 1px solid var(--color-border);
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.login-title {
  margin-bottom: 32px;
  font-size: 24px;
  font-weight: 600;
  text-align: center;
  color: var(--color-text);
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
  color: var(--color-text);
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

.error-message {
  padding: 10px 12px;
  background: var(--color-error-bg);
  border: 1px solid var(--color-error-border);
  border-radius: 4px;
  color: var(--color-error-text);
  font-size: 14px;
}

.login-button {
  padding: 12px 24px;
  background: var(--color-primary);
  color: #ffffff;
  font-size: 14px;
  font-weight: 500;
  border-radius: 4px;
  transition: all 0.2s;
  cursor: pointer;
}

.login-button:hover:not(:disabled) {
  background: var(--color-primary-hover);
}

.login-button:active:not(:disabled) {
  background: var(--color-primary-active);
}

.login-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>
