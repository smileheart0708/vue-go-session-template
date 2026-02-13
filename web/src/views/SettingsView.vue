<template>
  <div class="settings">
    <div ref="tabsRef" class="settings-tabs">
      <span v-show="showIndicator" class="settings-indicator" />
      <RouterLink
        v-for="tab in tabs"
        :key="tab.name"
        :to="{ name: tab.name }"
        class="settings-tab"
        :class="{ active: tab.name === activeTabName }"
      >
        <span class="settings-tab__label">{{ tab.label }}</span>
      </RouterLink>
    </div>

    <section class="settings-content">
      <RouterView v-slot="{ Component }">
        <Transition :name="transitionName" mode="out-in">
          <component :is="Component" :key="routeKey" />
        </Transition>
      </RouterView>
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'

type SettingsRouteName = 'settings-general' | 'settings-upstream' | 'settings-proxy'

interface SettingsTab {
  name: SettingsRouteName
  label: string
}

const tabs: SettingsTab[] = [
  { name: 'settings-general', label: '基本设置' },
  { name: 'settings-upstream', label: '上游服务' },
  { name: 'settings-proxy', label: '下游代理' },
]

const route = useRoute()

const activeTabName = computed<string>(() => (typeof route.name === 'string' ? route.name : ''))
const routeKey = computed<string>(() => route.fullPath)

const transitionName = ref<string>('settings-slide-forward')

const tabsRef = ref<HTMLElement | null>(null)
const indicatorOffset = ref<number>(0)
const indicatorWidth = ref<number>(0)
const showIndicator = ref<boolean>(false)

const indicatorOffsetPx = computed<string>(() => `${indicatorOffset.value}px`)
const indicatorWidthPx = computed<string>(() => `${indicatorWidth.value}px`)

const activeIndex = computed<number>(() => tabs.findIndex((tab) => tab.name === activeTabName.value))

function updateIndicatorPosition(): void {
  const container = tabsRef.value
  if (!container) return

  const activeElement = container.querySelector('.settings-tab.active') as HTMLElement | null
  if (!activeElement) {
    showIndicator.value = false
    return
  }

  indicatorOffset.value = activeElement.offsetLeft
  indicatorWidth.value = activeElement.offsetWidth
  showIndicator.value = true
}

watch(
  () => activeIndex.value,
  (nextIndex: number, prevIndex: number) => {
    if (nextIndex === -1 || prevIndex === -1 || nextIndex === prevIndex) return
    transitionName.value = nextIndex > prevIndex ? 'settings-slide-forward' : 'settings-slide-back'
  },
)

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
  window.addEventListener('resize', updateIndicatorPosition)
})

onUnmounted(() => {
  window.removeEventListener('resize', updateIndicatorPosition)
})
</script>

<style scoped>
/* Settings 页面布局 */
.settings {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

/* Tabs 导航 */
.settings-tabs {
  position: relative;
  display: flex;
  gap: 1.5rem;
  align-items: flex-end;
  padding-bottom: 0.75rem;
  border-bottom: 1px solid var(--sys-color-border);
  overflow-x: auto;
  scrollbar-width: none;
}

.settings-tabs::-webkit-scrollbar {
  height: 0;
}

/* Tab 指示器 - 使用 CSS v-bind */
.settings-indicator {
  position: absolute;
  left: 0;
  bottom: 0;
  height: 3px;
  width: v-bind(indicatorWidthPx);
  transform: translate3d(v-bind(indicatorOffsetPx), 0, 0);
  background: var(--sys-color-accent);
  border-radius: 999px;
  box-shadow: var(--sys-shadow-accent-glow);
  transition:
    transform 0.35s cubic-bezier(0.2, 0.8, 0.2, 1),
    width 0.35s cubic-bezier(0.2, 0.8, 0.2, 1);
  will-change: transform, width;
}

/* Tab 项 */
.settings-tab {
  display: inline-flex;
  flex-direction: column;
  gap: 0.2rem;
  padding: 0.2rem 0.5rem 0.55rem;
  color: var(--sys-color-text-secondary);
  font-weight: 500;
  text-decoration: none;
  transition: color 0.2s ease;
  white-space: nowrap;
}

.settings-tab:hover {
  color: var(--sys-color-text-primary);
}

.settings-tab.active {
  color: var(--sys-color-accent);
}

.settings-tab__label {
  font-size: 1rem;
}

.settings-tab__desc {
  font-size: 0.8rem;
  color: var(--sys-color-text-tertiary);
}

.settings-tab.active .settings-tab__desc {
  color: var(--sys-color-accent-hover);
}

/* 内容区域 */
.settings-content {
  min-height: 240px;
}

/* RouterView 过渡动画 */
.settings-slide-forward-enter-active,
.settings-slide-forward-leave-active,
.settings-slide-back-enter-active,
.settings-slide-back-leave-active {
  transition:
    transform 0.35s cubic-bezier(0.2, 0.8, 0.2, 1),
    opacity 0.25s ease;
  will-change: transform, opacity;
}

.settings-slide-forward-enter-from,
.settings-slide-back-leave-to {
  opacity: 0;
  transform: translate3d(12px, 0, 0) scale(0.98);
}

.settings-slide-forward-leave-to,
.settings-slide-back-enter-from {
  opacity: 0;
  transform: translate3d(-12px, 0, 0) scale(0.98);
}

/* 响应式设计 */
@media (width <= 640px) {
  .settings-tab {
    padding-inline: 0.25rem;
  }
}

/* 减少动画偏好 */
@media (prefers-reduced-motion: reduce) {
  .settings-slide-forward-enter-active,
  .settings-slide-forward-leave-active,
  .settings-slide-back-enter-active,
  .settings-slide-back-leave-active {
    transition: none;
  }
}
</style>
