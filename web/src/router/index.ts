import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '@/views/LoginView.vue'
import MainLayout from '@/components/layout/MainLayout.vue'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: LoginView,
      meta: { requiresGuest: true }
    },
    {
      path: '/',
      component: MainLayout,
      meta: { requiresAuth: true },
      children: [
        {
          path: '',
          redirect: '/dashboard'
        },
        {
          path: 'dashboard',
          name: 'dashboard',
          component: () => import('@/views/DashboardView.vue'),
          meta: { title: '仪表板', requiresAuth: true },
        },
        {
          path: 'logs',
          name: 'logs',
          component: () => import('@/views/LogsView.vue'),
          meta: { title: '日志', requiresAuth: true },
        },
        {
          path: 'settings',
          name: 'settings',
          component: () => import('@/views/SettingsView.vue'),
          meta: { title: '设置', requiresAuth: true },
        },
      ],
    },
  ],
})

// 全局前置守卫
router.beforeEach(async (to, _from, next) => {
  const authStore = useAuthStore()

  // 检查路由是否需要认证
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth)
  const requiresGuest = to.matched.some(record => record.meta.requiresGuest)

  // 如果需要认证
  if (requiresAuth) {
    if (!authStore.isAuthenticated) {
      // 未登录，重定向到登录页
      next({ name: 'login' })
      return
    }

    // 已登录，验证 session 是否有效（非模拟模式）
    const isMockAuth = import.meta.env.VITE_MOCK_AUTH === 'true'
    if (!isMockAuth) {
      const isValid = await authStore.validateSession()
      if (!isValid) {
        // session 无效，重定向到登录页
        next({ name: 'login' })
        return
      }
    }
  }

  // 如果是登录页，但已经登录
  if (requiresGuest && authStore.isAuthenticated) {
    // 已登录，重定向到仪表板
    next({ name: 'dashboard' })
    return
  }

  next()
})

export default router
