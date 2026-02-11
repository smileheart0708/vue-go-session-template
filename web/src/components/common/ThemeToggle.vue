<template>
  <div class="theme-toggle-container">
    <IconButton
      ref="buttonRef"
      :title="tooltipText"
      @click="handleClick"
      @mousedown="handleMouseDown"
      @mouseup="handleMouseUp"
      @mouseleave="handleMouseUp"
      @touchstart.passive="handleTouchStart"
      @touchend="handleTouchEnd"
    >
      <Sun v-if="mode === 'light'" />
      <Moon v-else-if="mode === 'dark'" />
      <Monitor v-else />
    </IconButton>

    <DropdownDrawer v-model="showDropdown" :anchor-el="anchorEl">
      <button
        class="dropdown-item"
        :class="{ active: mode === 'light' }"
        @click="selectMode('light')"
      >
        <Sun class="dropdown-icon" />
        <span>浅色</span>
      </button>
      <button
        class="dropdown-item"
        :class="{ active: mode === 'dark' }"
        @click="selectMode('dark')"
      >
        <Moon class="dropdown-icon" />
        <span>深色</span>
      </button>
      <button
        class="dropdown-item"
        :class="{ active: mode === 'auto' }"
        @click="selectMode('auto')"
      >
        <Monitor class="dropdown-icon" />
        <span>自动</span>
      </button>
    </DropdownDrawer>
  </div>
</template>

<script setup lang="ts">
import { computed, onUnmounted, ref, useTemplateRef } from 'vue'
import { Sun, Moon, Monitor } from 'lucide-vue-next'
import DropdownDrawer from './DropdownDrawer.vue'
import IconButton from './IconButton.vue'
import { useTheme } from '@/composables'
import type { ThemeMode } from '@/composables'

defineOptions({
  name: 'ThemeToggle',
})

const { mode, setTheme } = useTheme()

const buttonRef = useTemplateRef<InstanceType<typeof IconButton>>('buttonRef')
const anchorEl = computed<HTMLElement | null>(() => buttonRef.value?.$el ?? null)

const showDropdown = ref(false)
let longPressTimer: ReturnType<typeof setTimeout> | null = null
let isLongPress = false
const LONG_PRESS_DURATION = 500

const tooltipText = computed(() => {
  const modeText: Record<ThemeMode, string> = {
    light: '浅色模式',
    dark: '深色模式',
    auto: '自动模式',
  }
  return `当前：${modeText[mode.value]}（长按切换）`
})

async function handleClick(event: MouseEvent) {
  if (isLongPress) {
    isLongPress = false
    return
  }

  const modes: ThemeMode[] = ['light', 'dark', 'auto']
  const currentIndex = modes.indexOf(mode.value)
  const nextIndex = (currentIndex + 1) % modes.length
  const nextMode = modes[nextIndex]
  if (nextMode) {
    await setTheme(nextMode, event)
  }
}

async function selectMode(newMode: ThemeMode) {
  showDropdown.value = false
  const buttonEl = anchorEl.value
  if (buttonEl) {
    const rect = buttonEl.getBoundingClientRect()
    const event = new MouseEvent('click', {
      clientX: rect.left + rect.width / 2,
      clientY: rect.top + rect.height / 2,
    })
    await setTheme(newMode, event)
  }
}

function handleMouseDown() {
  isLongPress = false
  longPressTimer = setTimeout(() => {
    isLongPress = true
    showDropdown.value = true
  }, LONG_PRESS_DURATION)
}

function handleMouseUp() {
  if (longPressTimer) {
    clearTimeout(longPressTimer)
    longPressTimer = null
  }
}

function handleTouchStart() {
  isLongPress = false
  longPressTimer = setTimeout(() => {
    isLongPress = true
    showDropdown.value = true
  }, LONG_PRESS_DURATION)
}

function handleTouchEnd() {
  if (longPressTimer) {
    clearTimeout(longPressTimer)
    longPressTimer = null
  }
}

onUnmounted(() => {
  if (longPressTimer) {
    clearTimeout(longPressTimer)
  }
})
</script>

<style scoped>
.theme-toggle-container {
  position: relative;
  display: inline-block;
}
</style>
