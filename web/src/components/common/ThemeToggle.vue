<template>
  <div ref="containerRef" class="theme-toggle-container">
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

    <Teleport to="body">
      <Transition name="dropdown">
        <div
          v-if="showDropdown"
          ref="dropdownRef"
          class="theme-dropdown"
          :style="dropdownStyle"
        >
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
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { Sun, Moon, Monitor } from 'lucide-vue-next'
import { useTheme, type ThemeMode } from '@/composables'
import IconButton from './IconButton.vue'

const { mode, setTheme } = useTheme()

const containerRef = ref<HTMLElement | null>(null)
const buttonRef = ref<InstanceType<typeof IconButton> | null>(null)
const dropdownRef = ref<HTMLElement | null>(null)

const showDropdown = ref(false)
let longPressTimer: ReturnType<typeof setTimeout> | null = null
let isLongPress = false
const LONG_PRESS_DURATION = 500
const dropdownStyle = ref<Record<string, string>>({})

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
  if (buttonRef.value?.$el) {
    const rect = buttonRef.value.$el.getBoundingClientRect()
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

function handleClickOutside(event: MouseEvent) {
  const target = event.target as Node
  if (containerRef.value?.contains(target)) return
  if (dropdownRef.value?.contains(target)) return
  showDropdown.value = false
}

function updateDropdownPosition() {
  const buttonEl = buttonRef.value?.$el
  if (!buttonEl) return
  const rect = buttonEl.getBoundingClientRect()
  const offset = 8
  dropdownStyle.value = {
    top: `${rect.bottom + offset}px`,
    right: `${Math.max(8, window.innerWidth - rect.right)}px`,
  }
}

function handleWindowChange() {
  if (!showDropdown.value) return
  updateDropdownPosition()
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
  window.removeEventListener('resize', handleWindowChange)
  window.removeEventListener('scroll', handleWindowChange, true)
  if (longPressTimer) {
    clearTimeout(longPressTimer)
  }
})

watch(showDropdown, async (isOpen) => {
  if (isOpen) {
    await nextTick()
    updateDropdownPosition()
    window.addEventListener('resize', handleWindowChange)
    window.addEventListener('scroll', handleWindowChange, true)
  } else {
    window.removeEventListener('resize', handleWindowChange)
    window.removeEventListener('scroll', handleWindowChange, true)
  }
})
</script>

<style scoped>
.theme-toggle-container {
  position: relative;
  display: inline-block;
}

.theme-dropdown {
  position: fixed;
  min-width: 140px;
  padding: 4px;
  background-color: var(--color-background-glass);
  backdrop-filter: blur(12px) saturate(140%);
  -webkit-backdrop-filter: blur(12px) saturate(140%);
  border: 1px solid var(--color-border);
  border-radius: 8px;
  box-shadow:
    0 12px 28px rgba(0, 0, 0, 0.18),
    0 6px 12px rgba(0, 0, 0, 0.1),
    0 0 0 1px rgba(255, 255, 255, 0.08) inset;
  z-index: 1000;
}

.dropdown-item {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
  padding: 8px 12px;
  background: transparent;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  color: var(--color-text);
  font-size: 14px;
  text-align: left;
  transition: background-color 0.2s;
}

.dropdown-item:hover {
  background: var(--color-background-secondary);
}

.dropdown-item.active {
  background: var(--color-primary);
  color: wheat;
}

.dropdown-icon {
  width: 16px;
  height: 16px;
  flex-shrink: 0;
}

.dropdown-enter-active,
.dropdown-leave-active {
  transition:
    opacity 0.2s,
    transform 0.2s;
}

.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}
</style>
