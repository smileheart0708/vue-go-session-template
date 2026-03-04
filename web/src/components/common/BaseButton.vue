<template>
  <button
    class="inline-flex select-none items-center justify-center gap-2 rounded-md border px-4 text-sm font-medium transition-all duration-200 ease-out focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-(--sys-color-focus-ring) disabled:cursor-not-allowed disabled:opacity-50"
    :class="buttonClass"
    :style="buttonStyle"
    :disabled="disabled"
    @click="handleClick"
  >
    <span
      v-if="$slots['icon'] || icon"
      class="inline-flex items-center justify-center text-base leading-none"
    >
      <slot name="icon">
        <component
          :is="icon"
          v-if="icon"
          :size="16"
        />
      </slot>
    </span>
    <span class="whitespace-nowrap leading-none">{{ text }}</span>
  </button>
</template>

<script setup lang="ts">
import { computed, type Component } from 'vue'

defineOptions({ name: 'BaseButton' })

type CssSize = string | number

interface Props {
  width?: CssSize
  height?: CssSize
  text: string
  primary?: boolean
  icon?: string | Component | null
  disabled?: boolean
}

const {
  width = 'auto',
  height = 40,
  text,
  primary = false,
  icon = null,
  disabled = false,
} = defineProps<Props>()

const emit = defineEmits<{ click: [event: MouseEvent] }>()

const buttonStyle = computed<Record<string, string>>(() => {
  const formatSize = (size: CssSize): string => (typeof size === 'number' ? `${size}px` : size)
  return { width: formatSize(width), height: formatSize(height) }
})

const buttonClass = computed<string>(() => {
  if (primary) {
    return 'border-accent bg-accent text-on-accent enabled:hover:border-accent-hover enabled:hover:bg-accent-hover enabled:active:border-accent-active enabled:active:bg-accent-active'
  }

  return 'border-accent bg-bg-surface text-accent enabled:hover:bg-accent enabled:hover:text-on-accent enabled:active:border-accent-active enabled:active:bg-accent-active'
})

function handleClick(event: MouseEvent): void {
  if (!disabled) {
    emit('click', event)
  }
}
</script>
