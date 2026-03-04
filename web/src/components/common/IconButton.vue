<template>
  <button
    ref="buttonRef"
    class="inline-flex items-center justify-center rounded-md border-0 bg-transparent text-text-primary transition-all duration-200 active:scale-95 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-(--sys-color-focus-ring) disabled:cursor-not-allowed disabled:opacity-50"
    :class="[sizeClass, activeClass]"
    :title="title"
    :aria-pressed="toggle ? active : undefined"
    :disabled="disabled"
    @click="handleClick"
  >
    <slot />
  </button>
</template>

<script setup lang="ts">
import { computed, useTemplateRef } from 'vue'

defineOptions({ name: 'IconButton' })

interface Props {
  title?: string
  disabled?: boolean
  size?: 'small' | 'medium' | 'large'
  toggle?: boolean
}

const { title = '', disabled = false, size = 'medium', toggle = false } = defineProps<Props>()

const active = defineModel<boolean>('active', { default: false })

const emit = defineEmits<{ click: [event: MouseEvent] }>()

const buttonRef = useTemplateRef<HTMLButtonElement>('buttonRef')

const sizeClass = computed<string>(() => {
  if (size === 'small') return 'size-8 [&>svg]:size-4'
  if (size === 'large') return 'size-12 [&>svg]:size-6'
  return 'size-10 [&>svg]:size-5'
})

const activeClass = computed<string>(() => {
  if (active.value) {
    return 'bg-bg-component text-accent'
  }
  return 'enabled:hover:bg-bg-component-muted'
})

function handleClick(event: MouseEvent): void {
  if (disabled) return
  if (toggle) {
    active.value = !active.value
  }
  emit('click', event)
}

defineExpose({ $el: buttonRef })
</script>
