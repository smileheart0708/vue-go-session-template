<template>
  <nav ref="navRef" class="relative flex flex-1 flex-col gap-1 overflow-y-auto bg-transparent px-0 py-2">
    <div
      v-show="showIndicator"
      class="absolute left-3 top-0 z-10 h-[22px] w-1 rounded-full bg-accent shadow-accent-glow transition-transform duration-300 ease-[cubic-bezier(0.2,0.8,0.2,1)] will-change-transform"
      :style="{ transform: `translate3d(0, ${indicatorOffset}px, 0)` }"
    />

    <RouterLink
      v-for="item in menuItems"
      :key="item.path"
      :to="item.path"
      class="relative z-[1] mx-3 flex h-11 items-center justify-start gap-3 rounded-xl px-4 pl-10 text-left text-text-secondary no-underline transition-all duration-200 ease-out hover:translate-x-0.5 hover:bg-bg-component-muted hover:text-text-primary focus-visible:outline-none focus-visible:shadow-[0_0_0_2px_var(--sys-color-focus-ring),inset_0_0_0_1px_var(--sys-color-border)] [&.active]:bg-bg-component [&.active]:font-semibold [&.active]:text-accent [&.active]:shadow-[inset_0_0_0_1px_var(--sys-color-border)]"
      :exact-active-class="item.path === '/dashboard' ? 'active' : ''"
      :active-class="item.path !== '/dashboard' ? 'active' : ''"
    >
      <component :is="item.icon" class="size-5 shrink-0" />
      <span>{{ item.label }}</span>
    </RouterLink>
  </nav>
</template>

<script setup lang="ts">
import { nextTick, onMounted, ref, useTemplateRef, watch, type Component } from 'vue'
import { useRoute } from 'vue-router'
import { LayoutGrid, FileText, Key, Settings } from 'lucide-vue-next'

interface MenuItem {
  path: string
  label: string
  icon: Component
}

const menuItems: MenuItem[] = [
  { path: '/dashboard', label: '仪表板', icon: LayoutGrid },
  { path: '/keys', label: '密钥', icon: Key },
  { path: '/logs', label: '日志', icon: FileText },
  { path: '/settings', label: '设置', icon: Settings },
]

const route = useRoute()
const navRef = useTemplateRef<HTMLElement>('navRef')

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

function updateIndicatorPosition(): void {
  const navElement = navRef.value
  if (!navElement) {
    showIndicator.value = false
    return
  }

  const index = findActiveIndex(route.path)
  if (index < 0) {
    showIndicator.value = false
    return
  }

  const activeElements = navElement.querySelectorAll<HTMLElement>('a')
  const activeElement = activeElements[index]
  if (!activeElement) {
    showIndicator.value = false
    return
  }

  const navRect = navElement.getBoundingClientRect()
  const activeRect = activeElement.getBoundingClientRect()
  const indicatorOffsetAdjustment = activeRect.height / 4

  indicatorOffset.value = activeRect.top - navRect.top + indicatorOffsetAdjustment
  showIndicator.value = true
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
  updateIndicatorPosition()
})
</script>
