<template>
  <Teleport to="body">
    <div class="toast-container">
      <TransitionGroup name="toast">
        <div
          v-for="toast in toasts"
          :key="toast.id"
          :class="['toast', `toast-${toast.type}`]"
          role="alert"
        >
          <div class="toast-glow"></div>
          <div class="toast-icon">
            <Check v-if="toast.type === 'success'" :size="20" />
            <XCircle v-else-if="toast.type === 'error'" :size="20" />
            <AlertTriangle v-else-if="toast.type === 'warning'" :size="20" />
            <Info v-else :size="20" />
          </div>
          <div class="toast-content">
            <div class="toast-message">{{ toast.message }}</div>
          </div>
          <button
            class="toast-close"
            type="button"
            :aria-label="closeAriaLabel"
            @click="removeToast(toast.id)"
          >
            <X :size="14" />
          </button>
          <div v-if="toast.duration > 0" class="toast-progress">
            <div
              class="toast-progress-bar"
              :style="{ animationDuration: `${toast.duration}ms` }"
            ></div>
          </div>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { Check, XCircle, AlertTriangle, Info, X } from 'lucide-vue-next'
import { useToast } from '@/composables'

defineOptions({
  name: 'ToastMessage',
})

interface Props {
  closeAriaLabel?: string
}

const { closeAriaLabel = 'Close notification' } = defineProps<Props>()

const { toasts, removeToast } = useToast()
</script>

<style scoped>
.toast-container {
  position: fixed;
  top: 1.5rem;
  right: 1.5rem;
  z-index: 9999;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  pointer-events: none;
  max-width: 90vw;
  width: 380px;
}

.toast {
  position: relative;
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem 1.25rem;
  background: var(--sys-color-bg-surface);
  backdrop-filter: blur(12px) saturate(180%);
  backdrop-filter: blur(12px) saturate(180%);
  border-radius: 14px;
  box-shadow: var(--sys-shadow-toast);
  pointer-events: auto;
  min-height: 64px;
  overflow: hidden;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

/* 状态样式 - 使用 CSS 变量统一管理 */
.toast-success {
  border: 1px solid rgba(var(--state-color-success-rgb), 0.3);
}

.toast-error {
  border: 1px solid rgba(var(--state-color-danger-rgb), 0.3);
}

.toast-warning {
  border: 1px solid rgba(var(--state-color-warning-rgb), 0.3);
}

.toast-info {
  border: 1px solid rgba(var(--state-color-info-rgb), 0.3);
}

/* 发光效果 */
.toast-glow {
  position: absolute;
  inset: 0;
  pointer-events: none;
}

.toast-success .toast-glow {
  background: radial-gradient(
    circle at 0% 0%,
    rgba(var(--state-color-success-rgb), 0.15) 0%,
    transparent 50%
  );
}

.toast-error .toast-glow {
  background: radial-gradient(
    circle at 0% 0%,
    rgba(var(--state-color-danger-rgb), 0.15) 0%,
    transparent 50%
  );
}

.toast-warning .toast-glow {
  background: radial-gradient(
    circle at 0% 0%,
    rgba(var(--state-color-warning-rgb), 0.15) 0%,
    transparent 50%
  );
}

.toast-info .toast-glow {
  background: radial-gradient(
    circle at 0% 0%,
    rgba(var(--state-color-info-rgb), 0.15) 0%,
    transparent 50%
  );
}

/* 图标样式 */
.toast-icon {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border-radius: 10px;
  transition: transform 0.2s ease;
}

.toast:hover .toast-icon {
  transform: scale(1.1) rotate(5deg);
}

.toast-success .toast-icon {
  color: var(--state-color-success);
  background: rgba(var(--state-color-success-rgb), 0.1);
}

.toast-error .toast-icon {
  color: var(--state-color-danger);
  background: rgba(var(--state-color-danger-rgb), 0.1);
}

.toast-warning .toast-icon {
  color: var(--state-color-warning);
  background: rgba(var(--state-color-warning-rgb), 0.1);
}

.toast-info .toast-icon {
  color: var(--state-color-info);
  background: rgba(var(--state-color-info-rgb), 0.1);
}

.toast-content {
  flex: 1;
  min-width: 0;
}

.toast-message {
  color: var(--sys-color-text-primary);
  font-size: 0.9375rem;
  font-weight: 500;
  line-height: 1.5;
  overflow-wrap: break-word;
}

.toast-close {
  flex-shrink: 0;
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  border-radius: 8px;
  color: var(--sys-color-text-tertiary);
  cursor: pointer;
  transition: all 0.2s ease;
  padding: 0;
  opacity: 0.6;
}

.toast-close:hover {
  background: var(--sys-color-bg-subtle);
  color: var(--sys-color-text-primary);
  opacity: 1;
  transform: rotate(90deg);
}

/* 进度条 */
.toast-progress {
  position: absolute;
  bottom: 0;
  left: 0;
  width: 100%;
  height: 3px;
  background: var(--sys-color-progress-track);
}

.toast-progress-bar {
  height: 100%;
  width: 100%;
  background: currentcolor;
  transform-origin: left;
  animation: progress linear forwards;
}

.toast-success .toast-progress-bar {
  color: var(--state-color-success);
}

.toast-error .toast-progress-bar {
  color: var(--state-color-danger);
}

.toast-warning .toast-progress-bar {
  color: var(--state-color-warning);
}

.toast-info .toast-progress-bar {
  color: var(--state-color-info);
}

@keyframes progress {
  from {
    transform: scaleX(1);
  }

  to {
    transform: scaleX(0);
  }
}

/* 动画效果 */
.toast-enter-active {
  transition: all 0.4s cubic-bezier(0.18, 0.89, 0.32, 1.28);
}

.toast-leave-active {
  transition: all 0.3s cubic-bezier(0.6, -0.28, 0.735, 0.045);
}

.toast-enter-from {
  opacity: 0;
  transform: translateX(30px) scale(0.9);
}

.toast-leave-to {
  opacity: 0;
  transform: translateX(100%) scale(0.9);
}

/* 响应式适配 */
@media (width <= 480px) {
  .toast-container {
    top: 1rem;
    right: 50%;
    transform: translateX(50%);
    width: calc(100vw - 2rem);
  }

  .toast {
    padding: 0.875rem 1rem;
    gap: 0.75rem;
    min-height: 56px;
  }

  .toast-enter-from {
    transform: translateY(-20px) scale(0.9);
  }

  .toast-leave-to {
    transform: translateY(-20px) scale(0.9);
  }
}

@media (prefers-reduced-motion: reduce) {
  .toast-enter-active,
  .toast-leave-active,
  .toast-progress-bar {
    transition: none;
    animation: none;
  }
}
</style>
