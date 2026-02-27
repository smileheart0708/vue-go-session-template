<template>
  <div class="relative inline-block">
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
        @click="selectMode('light', $event)"
      >
        <Sun class="dropdown-icon" />
        <span>{{ lightLabel }}</span>
      </button>
      <button
        class="dropdown-item"
        :class="{ active: mode === 'dark' }"
        @click="selectMode('dark', $event)"
      >
        <Moon class="dropdown-icon" />
        <span>{{ darkLabel }}</span>
      </button>
      <button
        class="dropdown-item"
        :class="{ active: mode === 'auto' }"
        @click="selectMode('auto', $event)"
      >
        <Monitor class="dropdown-icon" />
        <span>{{ autoLabel }}</span>
      </button>
    </DropdownDrawer>
  </div>
</template>

<script setup lang="ts">
import { computed, onUnmounted, ref, useTemplateRef } from 'vue'
import { Sun, Moon, Monitor } from 'lucide-vue-next'
import DropdownDrawer from './DropdownDrawer.vue'
import IconButton from './IconButton.vue'

defineOptions({ name: 'ThemeToggle' })

type ThemeMode = 'light' | 'dark' | 'auto'

interface TooltipContext {
  mode: ThemeMode
  modeLabel: string
}

interface Props {
  lightLabel?: string
  darkLabel?: string
  autoLabel?: string
  tooltipPrefix?: string
  tooltipSuffix?: string
  tooltipFormatter?: (context: TooltipContext) => string
}

const {
  lightLabel = 'Light',
  darkLabel = 'Dark',
  autoLabel = 'Auto',
  tooltipPrefix = 'Current:',
  tooltipSuffix = '(Long press to switch)',
  tooltipFormatter,
} = defineProps<Props>()

const mode = defineModel<ThemeMode>({ required: true })

const emit = defineEmits<{ change: [mode: ThemeMode, event: MouseEvent | undefined] }>()

const buttonRef = useTemplateRef<InstanceType<typeof IconButton>>('buttonRef')
const anchorEl = computed<HTMLElement | null>(() => buttonRef.value?.$el ?? null)

const showDropdown = ref(false)
let longPressTimer: ReturnType<typeof setTimeout> | null = null
let isLongPress = false
const LONG_PRESS_DURATION = 500

const tooltipText = computed(() => {
  const modeText: Record<ThemeMode, string> = {
    light: lightLabel,
    dark: darkLabel,
    auto: autoLabel,
  }
  const currentModeLabel = modeText[mode.value]
  if (tooltipFormatter) {
    return tooltipFormatter({ mode: mode.value, modeLabel: currentModeLabel })
  }
  return `${tooltipPrefix} ${currentModeLabel} ${tooltipSuffix}`.replace(/\s+/g, ' ').trim()
})

async function applyMode(nextMode: ThemeMode, event?: MouseEvent): Promise<void> {
  mode.value = nextMode
  emit('change', nextMode, event)
}

async function handleClick(event: MouseEvent): Promise<void> {
  if (isLongPress) {
    isLongPress = false
    return
  }

  const modes: ThemeMode[] = ['light', 'dark', 'auto']
  const currentIndex = modes.indexOf(mode.value)
  const nextIndex = (currentIndex + 1) % modes.length
  const nextMode = modes[nextIndex]
  if (nextMode) {
    await applyMode(nextMode, event)
  }
}

async function selectMode(newMode: ThemeMode, event: MouseEvent): Promise<void> {
  showDropdown.value = false
  await applyMode(newMode, event)
}

function handleMouseDown(): void {
  isLongPress = false
  longPressTimer = setTimeout(() => {
    isLongPress = true
    showDropdown.value = true
  }, LONG_PRESS_DURATION)
}

function handleMouseUp(): void {
  if (longPressTimer) {
    clearTimeout(longPressTimer)
    longPressTimer = null
  }
}

function handleTouchStart(): void {
  isLongPress = false
  longPressTimer = setTimeout(() => {
    isLongPress = true
    showDropdown.value = true
  }, LONG_PRESS_DURATION)
}

function handleTouchEnd(): void {
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
