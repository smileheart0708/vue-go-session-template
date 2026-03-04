import { createRouter, createWebHistory, type RouteLocationNormalizedLoaded } from 'vue-router'
import { MainLayout } from '@/components'
import LoginView from '@/views/LoginView.vue'
import { useAuthStore } from '@/stores/auth'
import { isMockAuthEnabled, resolveRedirectPath } from '@/utils'

function createLoginRedirect(redirectPath: string) {
  const safeRedirectPath = resolveRedirectPath(redirectPath)
  if (!safeRedirectPath) {
    return { name: 'login' as const }
  }

  return { name: 'login' as const, query: { redirect: safeRedirectPath } }
}

const APP_NAME = 'web'

function resolveRouteTitle(to: RouteLocationNormalizedLoaded): string {
  const title = to.meta['title']
  if (typeof title !== 'string') return ''
  return title.trim()
}

function updateDocumentTitle(to: RouteLocationNormalizedLoaded): void {
  const routeTitle = resolveRouteTitle(to)
  document.title = routeTitle ? `${APP_NAME}-${routeTitle}` : APP_NAME
}

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: LoginView,
      meta: { title: '登录', requiresGuest: true },
    },
    {
      path: '/',
      component: MainLayout,
      meta: { requiresAuth: true },
      children: [
        { path: '', redirect: '/dashboard' },
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
          path: 'keys',
          name: 'keys',
          component: () => import('@/views/KeysView.vue'),
          meta: { title: '密钥', requiresAuth: true },
        },
        {
          path: 'settings',
          name: 'settings',
          component: () => import('@/views/SettingsView.vue'),
          meta: { title: '设置', requiresAuth: true },
          redirect: '/settings/general',
          children: [
            {
              path: 'general',
              name: 'settings-general',
              component: () => import('@/views/settings/GeneralSettings.vue'),
              meta: { title: '基本设置', requiresAuth: true },
            },
            {
              path: 'upstream',
              name: 'settings-upstream',
              component: () => import('@/views/settings/UpstreamSettings.vue'),
              meta: { title: '上游服务', requiresAuth: true },
            },
            {
              path: 'proxy',
              name: 'settings-proxy',
              component: () => import('@/views/settings/ProxySettings.vue'),
              meta: { title: '下游代理', requiresAuth: true },
            },
          ],
        },
      ],
    },
  ],
})

// 全局前置守卫
router.beforeEach(async (to) => {
  const authStore = useAuthStore()
  const isMockAuth = isMockAuthEnabled

  // 检查路由是否需要认证
  const requiresAuth = to.matched.some((record) => record.meta['requiresAuth'])
  const requiresGuest = to.matched.some((record) => record.meta['requiresGuest'])

  // 如果需要认证
  if (requiresAuth) {
    if (isMockAuth) {
      if (!authStore.isAuthenticated) {
        return createLoginRedirect(to.fullPath)
      }
    } else {
      const isValid = await authStore.validateSession()
      if (!isValid) {
        return createLoginRedirect(to.fullPath)
      }
    }
  }

  if (requiresGuest) {
    const hasActiveSession = isMockAuth ? authStore.isAuthenticated : await authStore.validateSession()
    if (!hasActiveSession) {
      return true
    }

    const redirectPath = resolveRedirectPath(to.query['redirect'])
    if (redirectPath) {
      return redirectPath
    }
    return { name: 'dashboard' }
  }

  return true
})

router.afterEach((to) => {
  updateDocumentTitle(to)
})

export default router
