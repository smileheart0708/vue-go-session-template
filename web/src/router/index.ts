import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '@/views/LoginView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/', name: 'login', component: LoginView },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: () => import('@/views/dashboard/DashboardView.vue'),
      // 路由级代码分割
      // 为该路由生成独立的代码块 (DashboardView.[hash].js)
      // 访问该路由时才进行懒加载
    },
  ],
})

export default router
