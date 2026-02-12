import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import { useAuthStore } from './stores/auth'
import { isMockApiEnabled, isMockAuthEnabled } from './utils'

async function enableApiMocking(): Promise<void> {
  if (!import.meta.env.DEV || !isMockApiEnabled) {
    return
  }

  const { worker } = await import('./mocks/browser')
  await worker.start({
    quiet: true,
    onUnhandledRequest(request, print) {
      const requestUrl = new URL(request.url)
      if (requestUrl.pathname === '/api' || requestUrl.pathname.startsWith('/api/')) {
        print.error()
      }
    },
  })
}

async function bootstrap(): Promise<void> {
  await enableApiMocking()

  const app = createApp(App)
  const pinia = createPinia()
  app.use(pinia)

  // 初始化 auth store
  const authStore = useAuthStore()
  authStore.init()

  // 如果开启了模拟认证模式，且当前未登录，则自动模拟登录
  if (isMockAuthEnabled && !authStore.isAuthenticated) {
    authStore.mockLogin()
  }

  app.use(router)
  await router.isReady()
  app.mount('#app')
}

void bootstrap()
