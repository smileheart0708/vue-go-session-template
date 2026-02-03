<template>
  <nav class="sidebar-nav">
    <div
      v-show="showIndicator"
      class="nav-indicator"
      :style="{ transform: `translate3d(0, ${indicatorOffset}px, 0)` }"
    ></div>
    <RouterLink
      v-for="item in menuItems"
      :key="item.path"
      :to="item.path"
      class="nav-item"
      :exact-active-class="item.path === '/dashboard' ? 'active' : ''"
      :active-class="item.path !== '/dashboard' ? 'active' : ''"
    >
      <component :is="item.icon" class="nav-icon" />
      <span>{{ item.label }}</span>
    </RouterLink>
  </nav>
</template>

<script setup lang="ts">
import { ref, watch, nextTick, onMounted, type Component } from 'vue'
import { useRoute } from 'vue-router'
import { LayoutGrid, FileText, Settings } from 'lucide-vue-next'

interface MenuItem {
  path: string
  label: string
  icon: Component
}

const menuItems: MenuItem[] = [
  { path: '/dashboard', label: '仪表板', icon: LayoutGrid },
  { path: '/logs', label: '日志', icon: FileText },
  { path: '/settings', label: '设置', icon: Settings },
]

const route = useRoute()

const indicatorOffset = ref(0)
const showIndicator = ref(false)

function findActiveIndex(path: string): number {
  return menuItems.findIndex((item) => {
    if (item.path === '/dashboard') {
      return path === '/dashboard'
    }
    return path.startsWith(item.path)
  })
}

function updateIndicatorPosition() {
  const index = findActiveIndex(route.path)
  if (index >= 0) {
    // 使用 getBoundingClientRect 获取实际位置
    const navElement = document.querySelector('.sidebar-nav')
    const activeElement = document.querySelectorAll('.nav-item')[index] as HTMLElement

    if (navElement && activeElement) {
      const navRect = navElement.getBoundingClientRect()
      const activeRect = activeElement.getBoundingClientRect()

      // 计算相对于导航容器的偏移量
      indicatorOffset.value = activeRect.top - navRect.top
      showIndicator.value = true
    }
  } else {
    showIndicator.value = false
  }
}

watch(
  () => route.fullPath,
  async () => {
    await nextTick()
    updateIndicatorPosition()
  },
  { immediate: true },
)

onMounted(() => {
  // 初始化时更新指示器位置
  updateIndicatorPosition()
})
</script>

<style scoped>
.sidebar-nav {
  flex: 1;
  padding: 1rem 0;
  overflow-y: auto;
  position: relative;
  background: var(--color-background-elevated);
}

.nav-indicator {
  position: absolute;
  top: 0;
  left: 0;
  width: 3px;
  height: 48px;
  background: var(--color-primary);
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  will-change: transform;
  z-index: 10;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0 1.25rem;
  height: 48px;
  color: var(--color-text-secondary);
  text-decoration: none;
  transition: all 0.2s ease;
  position: relative;
  z-index: 1;
}

.nav-item:hover {
  background: var(--color-background-secondary);
  color: var(--color-text);
}

.nav-item.active {
  background: var(--color-background-secondary);
  color: var(--color-primary);
}

.nav-icon {
  width: 20px;
  height: 20px;
  flex-shrink: 0;
}
</style>
