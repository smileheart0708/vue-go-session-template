<template>
  <Teleport to="body">
    <Transition name="dialog-fade">
      <div v-if="modelValue" class="dialog-overlay" @click="handleOverlayClick">
        <div class="dialog-container" @click.stop>
          <div class="dialog-header" :class="{ 'no-border': !showHeaderBorder }">
            <h3 class="dialog-title">{{ title }}</h3>
            <button class="dialog-close" @click="handleClose" aria-label="关闭">
              <svg viewBox="0 0 24 24" width="20" height="20">
                <path
                  fill="currentColor"
                  d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z"
                />
              </svg>
            </button>
          </div>
          <div class="dialog-body">
            <slot />
          </div>
          <div
            v-if="$slots.footer"
            class="dialog-footer"
            :class="{ 'no-border': !showFooterBorder }"
          >
            <slot name="footer" />
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { onWatcherCleanup, watch } from 'vue'

defineOptions({ name: 'AppDialog' })

interface Props {
  title: string
  closeOnOverlay?: boolean
  showHeaderBorder?: boolean
  showFooterBorder?: boolean
}

const {
  title,
  closeOnOverlay = true,
  showHeaderBorder = true,
  showFooterBorder = true,
} = defineProps<Props>()

const emit = defineEmits<{ close: [] }>()

const modelValue = defineModel<boolean>({ default: false })

watch(
  modelValue,
  (visible) => {
    if (!visible) return

    const previousOverflow = document.body.style.overflow
    document.body.style.overflow = 'hidden'

    onWatcherCleanup(() => {
      document.body.style.overflow = previousOverflow
    })
  },
  { immediate: true },
)

function handleClose(): void {
  modelValue.value = false
  emit('close')
}

function handleOverlayClick(): void {
  if (closeOnOverlay) {
    handleClose()
  }
}
</script>

<style scoped>
.dialog-overlay {
  position: fixed;
  inset: 0;
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
  padding: 1rem;
}

.dialog-container {
  background-color: var(--sys-color-bg-surface);
  border-radius: 12px;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
  width: 100%;
  max-width: 480px;
  max-height: 90vh;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.dialog-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1.25rem 1.5rem;
  border-bottom: 1px solid var(--sys-color-border);
}

.dialog-title {
  margin: 0;
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--sys-color-text-primary);
}

.dialog-close {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: none;
  background: transparent;
  color: var(--sys-color-text-secondary);
  cursor: pointer;
  border-radius: 6px;
  transition: all 0.2s ease;
}

.dialog-close:hover {
  background-color: var(--sys-color-bg-subtle);
  color: var(--sys-color-text-primary);
}

.dialog-body {
  padding: 1.5rem;
  overflow-y: auto;
  flex: 1;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  padding: 1rem 1.5rem;
  border-top: 1px solid var(--sys-color-border);
}

.dialog-header.no-border {
  border-bottom: none;
}

.dialog-footer.no-border {
  border-top: none;
}

/* 动画 */
.dialog-fade-enter-active,
.dialog-fade-leave-active {
  transition: opacity 0.2s ease;
}

.dialog-fade-enter-from,
.dialog-fade-leave-to {
  opacity: 0;
}

.dialog-fade-enter-active .dialog-container,
.dialog-fade-leave-active .dialog-container {
  transition: transform 0.2s ease;
}

.dialog-fade-enter-from .dialog-container,
.dialog-fade-leave-to .dialog-container {
  transform: scale(0.95);
}
</style>
