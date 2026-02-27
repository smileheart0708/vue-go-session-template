<template>
  <div class="flex w-full flex-col gap-6">
    <div
      ref="tabsRef"
      class="relative flex items-end gap-6 overflow-x-auto border-b border-border pb-3 [scrollbar-width:none] [&::-webkit-scrollbar]:h-0 max-sm:gap-4"
    >
      <span
        v-show="showIndicator"
        class="absolute bottom-0 left-0 h-[3px] rounded-full bg-accent shadow-accent-glow transition-[transform,width] duration-350 ease-[cubic-bezier(0.2,0.8,0.2,1)]"
        :style="{ width: indicatorWidthPx, transform: `translate3d(${indicatorOffsetPx}, 0, 0)` }"
      />

      <RouterLink
        v-for="tab in tabs"
        :key="tab.name"
        :to="{ name: tab.name }"
        class="inline-flex whitespace-nowrap px-2 pb-[0.55rem] pt-[0.2rem] text-base font-medium text-text-secondary no-underline transition-colors duration-200 hover:text-text-primary"
        :class="{ 'text-accent': tab.name === activeTabName }"
      >
        <span>{{ tab.label }}</span>
      </RouterLink>
    </div>

    <section class="min-h-[240px] [overflow-anchor:none]">
      <RouterView v-slot="{ Component }">
        <Transition
          mode="out-in"
          enter-active-class="transition-[transform,opacity] duration-350 ease-[cubic-bezier(0.2,0.8,0.2,1)] motion-reduce:transition-none"
          leave-active-class="transition-[transform,opacity] duration-350 ease-[cubic-bezier(0.2,0.8,0.2,1)] motion-reduce:transition-none"
          :enter-from-class="transitionEnterFromClass"
          leave-from-class="opacity-100 translate-x-0 scale-100"
          :leave-to-class="transitionLeaveToClass"
          enter-to-class="opacity-100 translate-x-0 scale-100"
        >
          <component :is="Component" :key="routeKey" />
        </Transition>
      </RouterView>
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onMounted, onUnmounted, ref, useTemplateRef, watch } from 'vue'
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

const tabsRef = useTemplateRef<HTMLElement>('tabsRef')
const indicatorOffset = ref<number>(0)
const indicatorWidth = ref<number>(0)
const showIndicator = ref<boolean>(false)

const indicatorOffsetPx = computed<string>(() => `${indicatorOffset.value}px`)
const indicatorWidthPx = computed<string>(() => `${indicatorWidth.value}px`)

const transitionEnterFromClass = computed<string>(() => {
  if (transitionName.value === 'settings-slide-back') {
    return 'opacity-0 -translate-x-3 scale-[0.98]'
  }
  return 'opacity-0 translate-x-3 scale-[0.98]'
})

const transitionLeaveToClass = computed<string>(() => {
  if (transitionName.value === 'settings-slide-back') {
    return 'opacity-0 translate-x-3 scale-[0.98]'
  }
  return 'opacity-0 -translate-x-3 scale-[0.98]'
})

const activeIndex = computed<number>(() =>
  tabs.findIndex((tab) => tab.name === activeTabName.value),
)

function updateIndicatorPosition(): void {
  const container = tabsRef.value
  if (!container) return

  const activeElement = container.querySelector<HTMLElement>('a.text-accent')
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
