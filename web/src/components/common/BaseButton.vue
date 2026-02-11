<template>
  <button
    class="base-button"
    :class="{ 'is-primary': primary }"
    :style="buttonStyle"
    :disabled="disabled"
    @click="handleClick"
  >
    <span v-if="$slots.icon || icon" class="button-icon">
      <slot name="icon">
        <component :is="icon" v-if="icon" :size="16" />
      </slot>
    </span>
    <span class="button-text">{{ text }}</span>
  </button>
</template>

<script setup lang="ts">
import { computed, type Component } from 'vue'

defineOptions({
  name: 'BaseButton',
})

type CssSize = string | number

interface Props {
  width: CssSize
  height: CssSize
  text: string
  primary?: boolean
  icon?: string | Component | null
  disabled?: boolean
}

const {
  width,
  height,
  text,
  primary = false,
  icon = null,
  disabled = false,
} = defineProps<Props>()

const emit = defineEmits<{
  click: [event: MouseEvent]
}>()

// 处理尺寸单位
const buttonStyle = computed<Record<string, string>>(() => {
  const formatSize = (size: CssSize): string => (typeof size === 'number' ? `${size}px` : size)

  return { width: formatSize(width), height: formatSize(height) }
})

function handleClick(event: MouseEvent): void {
  if (!disabled) {
    emit('click', event)
  }
}
</script>

<style scoped>
.base-button {
  /* Win11 风格固定圆角 */
  border-radius: 6px;

  /* 布局 - 默认居中 */
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;

  /* 基础样式 */
  padding: 0 16px;
  font-size: 14px;
  font-weight: 500;
  font-family: inherit;
  cursor: pointer;
  outline: none;
  user-select: none;

  /* 过渡动画 */
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);

  /* 默认模式：背景使用 elevated 背景色，边框和文字使用主题色 */
  background-color: var(--color-background-elevated);
  border: 1px solid var(--color-primary);
  color: var(--color-primary);
}

/* 图标容器 */
.button-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  line-height: 1;
}

/* 文字容器 */
.button-text {
  line-height: 1;
  white-space: nowrap;
}

/* 默认模式 Hover：背景变为主题色，文字变白 */
.base-button:not(.is-primary):not(:disabled):hover {
  background-color: var(--color-primary);
  color: #ffffff;
  border-color: var(--color-primary);
}

/* 默认模式 Active：使用更深主题色 */
.base-button:not(.is-primary):not(:disabled):active {
  background-color: var(--color-primary-active);
  border-color: var(--color-primary-active);
}

/* 主题色模式：背景使用主题色，文字白色，无边框或边框同色 */
.base-button.is-primary {
  background-color: var(--color-primary);
  border: 1px solid var(--color-primary);
  color: #ffffff;
}

/* 主题色模式 Hover */
.base-button.is-primary:not(:disabled):hover {
  background-color: var(--color-primary-hover);
  border-color: var(--color-primary-hover);
}

/* 主题色模式 Active */
.base-button.is-primary:not(:disabled):active {
  background-color: var(--color-primary-active);
  border-color: var(--color-primary-active);
}

/* 禁用状态 */
.base-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* 深色主题适配 - 确保对比度 */
:root[data-theme='dark'] .base-button:not(.is-primary) {
  background-color: var(--color-background-elevated);
}

/* 焦点状态 - 添加轮廓 */
.base-button:focus-visible {
  outline: 2px solid var(--color-primary);
  outline-offset: 2px;
}
</style>
