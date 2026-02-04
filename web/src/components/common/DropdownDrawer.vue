<template>
  <Teleport :to="teleportTo">
    <Transition name="dropdown-drawer">
      <div
        v-if="open"
        ref="drawerRef"
        class="dropdown-drawer"
        :style="drawerStyle"
        role="menu"
        v-bind="attrs"
      >
        <slot />
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, watch, nextTick, onUnmounted, useAttrs } from 'vue'

defineOptions({
  inheritAttrs: false,
})

interface Props {
  anchorEl: HTMLElement | null
  offset?: number
  align?: 'start' | 'end'
  teleportTo?: string
  minWidth?: number
  closeOnOutside?: boolean
  closeOnEscape?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  offset: 8,
  align: 'end',
  teleportTo: 'body',
  minWidth: 140,
  closeOnOutside: true,
  closeOnEscape: true,
})

const open = defineModel<boolean>({ default: false })

const attrs = useAttrs()
const drawerRef = ref<HTMLElement | null>(null)
const drawerStyle = ref<Record<string, string>>({})

let listenersActive = false

function updatePosition() {
  const anchor = props.anchorEl
  if (!anchor) return
  const rect = anchor.getBoundingClientRect()
  const offset = props.offset
  const style: Record<string, string> = {
    top: `${rect.bottom + offset}px`,
    '--dropdown-min-width': `${props.minWidth}px`,
  }

  if (props.align === 'start') {
    style.left = `${Math.max(8, rect.left)}px`
  } else {
    style.right = `${Math.max(8, window.innerWidth - rect.right)}px`
  }

  drawerStyle.value = style
}

function handleWindowChange() {
  if (!open.value) return
  updatePosition()
}

function handleClickOutside(event: MouseEvent) {
  if (!open.value || !props.closeOnOutside) return
  const target = event.target as Node
  if (drawerRef.value?.contains(target)) return
  if (props.anchorEl?.contains(target)) return
  open.value = false
}

function handleKeydown(event: KeyboardEvent) {
  if (!open.value || !props.closeOnEscape) return
  if (event.key === 'Escape') {
    open.value = false
  }
}

function addGlobalListeners() {
  if (listenersActive) return
  listenersActive = true
  if (props.closeOnOutside) {
    document.addEventListener('click', handleClickOutside)
  }
  if (props.closeOnEscape) {
    window.addEventListener('keydown', handleKeydown)
  }
  window.addEventListener('resize', handleWindowChange)
  window.addEventListener('scroll', handleWindowChange, true)
}

function removeGlobalListeners() {
  if (!listenersActive) return
  listenersActive = false
  document.removeEventListener('click', handleClickOutside)
  window.removeEventListener('keydown', handleKeydown)
  window.removeEventListener('resize', handleWindowChange)
  window.removeEventListener('scroll', handleWindowChange, true)
}

watch(
  open,
  async (isOpen) => {
    if (isOpen) {
      await nextTick()
      if (!props.anchorEl) {
        open.value = false
        return
      }
      updatePosition()
      addGlobalListeners()
    } else {
      removeGlobalListeners()
    }
  },
  { immediate: true },
)

watch(
  () => [props.anchorEl, props.align, props.offset, props.minWidth],
  () => {
    if (open.value) {
      updatePosition()
    }
  },
)

onUnmounted(() => {
  removeGlobalListeners()
})
</script>

<style scoped>
.dropdown-drawer {
  position: fixed;
  min-width: var(--dropdown-min-width, 140px);
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
  z-index: var(--dropdown-z-index, 1000);
}

:deep(.dropdown-item) {
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

:deep(.dropdown-item:hover) {
  background: var(--color-background-secondary);
}

:deep(.dropdown-item.active) {
  background: var(--dropdown-active-bg, var(--color-primary));
  color: var(--dropdown-active-color, wheat);
}

:deep(.dropdown-icon) {
  width: 16px;
  height: 16px;
  flex-shrink: 0;
}

.dropdown-drawer-enter-active,
.dropdown-drawer-leave-active {
  transition:
    opacity 0.2s,
    transform 0.2s;
}

.dropdown-drawer-enter-from,
.dropdown-drawer-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}

@media (prefers-reduced-motion: reduce) {
  .dropdown-drawer-enter-active,
  .dropdown-drawer-leave-active {
    transition: none;
  }
}
</style>
