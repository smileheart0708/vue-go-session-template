<template>
  <button :class="buttonClasses" :disabled="disabled || loading" :type="type" @click="handleClick">
    <!-- Loading Spinner -->
    <span v-if="loading" class="btn-spinner">
      <Loader2 />
    </span>

    <!-- Prefix Icon -->
    <span v-if="(prefixIcon || $slots['prefix-icon']) && !loading" class="btn-icon btn-icon-prefix">
      <slot name="prefix-icon">
        <component :is="prefixIcon" />
      </slot>
    </span>

    <!-- Button Content -->
    <span v-if="$slots.default" class="btn-content">
      <slot />
    </span>

    <!-- Suffix Icon -->
    <span v-if="(suffixIcon || $slots['suffix-icon']) && !loading" class="btn-icon btn-icon-suffix">
      <slot name="suffix-icon">
        <component :is="suffixIcon" />
      </slot>
    </span>
  </button>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Loader2 } from 'lucide-vue-next'

interface Props {
  variant?: 'primary' | 'secondary' | 'danger' | 'ghost' | 'text'
  size?: 'small' | 'medium' | 'large'
  type?: 'button' | 'submit' | 'reset'
  disabled?: boolean
  loading?: boolean
  block?: boolean
  iconOnly?: boolean
  prefixIcon?: unknown
  suffixIcon?: unknown
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'primary',
  size: 'medium',
  type: 'button',
  disabled: false,
  loading: false,
  block: false,
  iconOnly: false,
  prefixIcon: undefined,
  suffixIcon: undefined,
})

const emit = defineEmits<{ click: [event: MouseEvent] }>()

const buttonClasses = computed(() => {
  return [
    'app-button',
    `app-button--${props.variant}`,
    `app-button--${props.size}`,
    {
      'app-button--disabled': props.disabled || props.loading,
      'app-button--loading': props.loading,
      'app-button--block': props.block,
      'app-button--icon-only': props.iconOnly,
      'app-button--has-prefix': props.prefixIcon,
      'app-button--has-suffix': props.suffixIcon,
    },
  ]
})

function handleClick(event: MouseEvent) {
  if (!props.disabled && !props.loading) {
    emit('click', event)
  }
}
</script>

<style scoped>
.app-button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  border: none;
  border-radius: 8px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  white-space: nowrap;
  position: relative;
  outline: none;
}

.app-button:focus-visible {
  box-shadow: 0 0 0 3px var(--primary-light);
}

/* Sizes */
.app-button--small {
  padding: 0.5rem 0.875rem;
  font-size: 0.8125rem;
}

.app-button--medium {
  padding: 0.625rem 1rem;
  font-size: 0.875rem;
}

.app-button--large {
  padding: 0.875rem 1.5rem;
  font-size: 1rem;
}

/* Icon Only */
.app-button--icon-only {
  padding: 0.5rem;
  width: 36px;
  height: 36px;
}

.app-button--small.app-button--icon-only {
  width: 28px;
  height: 28px;
  padding: 0.375rem;
}

.app-button--large.app-button--icon-only {
  width: 44px;
  height: 44px;
  padding: 0.625rem;
}

/* Block */
.app-button--block {
  width: 100%;
}

/* Variants - Primary */
.app-button--primary {
  background: var(--primary-color);
  color: var(--text-white);
}

.app-button--primary:hover:not(:disabled) {
  background: var(--primary-hover);
  transform: translateY(-1px);
  box-shadow: var(--shadow-md);
}

.app-button--primary:active:not(:disabled) {
  transform: translateY(0);
}

/* Variants - Secondary */
.app-button--secondary {
  background: var(--bg-primary);
  color: var(--text-secondary);
  border: 1px solid var(--border-color);
}

.app-button--secondary:hover:not(:disabled) {
  background: var(--hover-bg);
  border-color: var(--primary-color);
  color: var(--primary-color);
}

/* Variants - Danger */
.app-button--danger {
  background: var(--error-color);
  color: var(--text-white);
}

.app-button--danger:hover:not(:disabled) {
  background: var(--error-hover);
  transform: translateY(-1px);
  box-shadow: var(--shadow-md);
}

.app-button--danger:active:not(:disabled) {
  transform: translateY(0);
}

/* Variants - Ghost */
.app-button--ghost {
  background: var(--bg-secondary);
  color: var(--text-secondary);
  border: 1px solid var(--border-color);
}

.app-button--ghost:hover:not(:disabled) {
  background: var(--hover-bg);
  color: var(--primary-color);
  border-color: var(--primary-color);
}

/* Variants - Text */
.app-button--text {
  background: transparent;
  color: var(--primary-color);
  padding: 0.375rem 0.75rem;
  border-radius: 6px;
}

.app-button--text:hover:not(:disabled) {
  background: var(--primary-light);
}

.app-button--text.app-button--small {
  padding: 0.25rem 0.5rem;
}

.app-button--text.app-button--large {
  padding: 0.5rem 1rem;
}

/* Disabled State */
.app-button--disabled {
  opacity: 0.6;
  cursor: not-allowed;
  pointer-events: none;
}

/* Loading State */
.app-button--loading {
  pointer-events: none;
}

.btn-spinner {
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-spinner svg {
  width: 16px;
  height: 16px;
  animation: spin 1s linear infinite;
}

.app-button--small .btn-spinner svg {
  width: 14px;
  height: 14px;
}

.app-button--large .btn-spinner svg {
  width: 18px;
  height: 18px;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

/* Icons */
.btn-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.btn-icon svg {
  width: 16px;
  height: 16px;
}

.app-button--small .btn-icon svg {
  width: 14px;
  height: 14px;
}

.app-button--large .btn-icon svg {
  width: 18px;
  height: 18px;
}

.app-button--icon-only .btn-icon svg {
  width: 16px;
  height: 16px;
}

.app-button--small.app-button--icon-only .btn-icon svg {
  width: 14px;
  height: 14px;
}

.app-button--large.app-button--icon-only .btn-icon svg {
  width: 18px;
  height: 18px;
}

/* Content */
.btn-content {
  display: inline-flex;
  align-items: center;
}

/* Hide content when loading */
.app-button--loading .btn-content {
  opacity: 0.7;
}
</style>
