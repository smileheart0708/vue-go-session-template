<template>
  <button
    ref="buttonRef"
    class="icon-button"
    :class="[size, { active }]"
    :title="title"
    :aria-pressed="toggle ? active : undefined"
    :disabled="disabled"
    @click="handleClick"
  >
    <slot />
  </button>
</template>

<script setup lang="ts">
import { ref } from 'vue'

interface Props {
  title?: string
  disabled?: boolean
  size?: 'small' | 'medium' | 'large'
  toggle?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  title: '',
  disabled: false,
  size: 'medium',
  toggle: false,
})

const active = defineModel<boolean>('active', { default: false })

const emit = defineEmits<{
  click: [event: MouseEvent]
}>()

const buttonRef = ref<HTMLButtonElement | null>(null)

function handleClick(event: MouseEvent) {
  if (props.disabled) return
  if (props.toggle) {
    active.value = !active.value
  }
  emit('click', event)
}

defineExpose({
  $el: buttonRef,
})
</script>

<style scoped>
.icon-button {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0;
  background: transparent;
  border: 1px solid transparent;
  border-radius: 8px;
  cursor: pointer;
  color: var(--color-text);
  transition:
    all 0.2s,
    transform 0.1s;
}

.icon-button:hover:not(:disabled) {
  background: var(--color-background-secondary);
  border-color: var(--color-border);
}

.icon-button:active:not(:disabled) {
  transform: scale(0.95);
}

.icon-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.icon-button.active {
  background: var(--color-primary);
  color: white;
}

/* 尺寸变体 */
.icon-button.small {
  width: 32px;
  height: 32px;
}

.icon-button.medium {
  width: 40px;
  height: 40px;
}

.icon-button.large {
  width: 48px;
  height: 48px;
}

/* 图标样式 */
.icon-button :deep(svg) {
  flex-shrink: 0;
}

.icon-button.small :deep(svg) {
  width: 16px;
  height: 16px;
}

.icon-button.medium :deep(svg) {
  width: 20px;
  height: 20px;
}

.icon-button.large :deep(svg) {
  width: 24px;
  height: 24px;
}
</style>
