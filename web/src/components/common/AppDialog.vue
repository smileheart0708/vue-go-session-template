<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition-opacity duration-200 ease-[ease] [&_.dialog-panel]:transition-transform [&_.dialog-panel]:duration-200 [&_.dialog-panel]:ease-[ease]"
      enter-from-class="opacity-0 [&_.dialog-panel]:scale-95"
      enter-to-class="opacity-100 [&_.dialog-panel]:scale-100"
      leave-active-class="transition-opacity duration-200 ease-[ease] [&_.dialog-panel]:transition-transform [&_.dialog-panel]:duration-200 [&_.dialog-panel]:ease-[ease]"
      leave-from-class="opacity-100 [&_.dialog-panel]:scale-100"
      leave-to-class="opacity-0 [&_.dialog-panel]:scale-95"
    >
      <div
        v-if="modelValue"
        class="fixed inset-0 z-1000 flex items-center justify-center bg-black/50 p-4 backdrop-blur-xs"
        @click="handleOverlayClick"
      >
        <div
          class="dialog-panel flex max-h-[90vh] w-full max-w-120 flex-col overflow-hidden rounded-xl bg-bg-surface shadow-[0_25px_50px_-12px_rgba(0,0,0,0.25)]"
          @click.stop
        >
          <div
            class="flex items-center justify-between px-6 py-5"
            :class="showHeaderBorder ? 'border-b border-border' : 'border-b-0'"
          >
            <h3 class="m-0 text-lg font-semibold text-text-primary">{{ title }}</h3>
            <button
              class="flex size-8 items-center justify-center rounded-md border-0 bg-transparent text-text-secondary transition-all duration-200 ease-[ease] hover:bg-bg-subtle hover:text-text-primary"
              @click="handleClose"
              :aria-label="closeAriaLabel"
            >
              <X :size="20" />
            </button>
          </div>
          <div class="flex-1 overflow-y-auto p-6">
            <slot />
          </div>
          <div
            v-if="$slots['footer']"
            class="flex justify-end gap-3 px-6 py-4"
            :class="showFooterBorder ? 'border-t border-border' : 'border-t-0'"
          >
            <slot name="footer" />
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { X } from 'lucide-vue-next'

defineOptions({ name: 'AppDialog' })

interface Props {
  title: string
  closeOnOverlay?: boolean
  showHeaderBorder?: boolean
  showFooterBorder?: boolean
  closeAriaLabel?: string
}

const {
  title,
  closeOnOverlay = true,
  showHeaderBorder = true,
  showFooterBorder = true,
  closeAriaLabel = 'Close dialog',
} = defineProps<Props>()

const emit = defineEmits<{ close: [] }>()

const modelValue = defineModel<boolean>({ required: true })

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
