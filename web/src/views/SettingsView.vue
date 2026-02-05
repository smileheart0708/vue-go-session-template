<template>
  <div class="settings">
    <div ref="tabsRef" class="settings-tabs">
      <span v-show="showIndicator" class="settings-indicator" :style="indicatorStyle" />
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

const activeTabName = computed(() => (typeof route.name === 'string' ? route.name : ''))
const routeKey = computed(() => route.fullPath)

const transitionName = ref('settings-slide-forward')

const tabsRef = ref<HTMLElement | null>(null)
const indicatorOffset = ref(0)
const indicatorWidth = ref(0)
const showIndicator = ref(false)

const indicatorStyle = computed(() => ({
  width: `${indicatorWidth.value}px`,
  transform: `translate3d(${indicatorOffset.value}px, 0, 0)`,
}))

const activeIndex = computed(() => tabs.findIndex((tab) => tab.name === activeTabName.value))

function updateIndicatorPosition() {
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
  (nextIndex, prevIndex) => {
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
