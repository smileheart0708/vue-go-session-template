import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import { useAuthStore } from './stores/auth'

const app = createApp(App)

const pinia = createPinia()
app.use(pinia)
app.use(router)

// 初始化 auth store
const authStore = useAuthStore()
authStore.init()

// 如果开启了模拟认证模式，且当前未登录，则自动模拟登录
const isMockAuth = import.meta.env.VITE_MOCK_AUTH === 'true'
if (isMockAuth && !authStore.isAuthenticated) {
  authStore.mockLogin()
}

app.mount('#app')
