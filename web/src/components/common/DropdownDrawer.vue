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
import { nextTick, onUnmounted, ref, useAttrs, useTemplateRef, watch } from 'vue'

defineOptions({
  name: 'DropdownDrawer',
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

const {
  anchorEl,
  offset = 8,
  align = 'end',
  teleportTo = 'body',
  minWidth = 140,
  closeOnOutside = true,
  closeOnEscape = true,
} = defineProps<Props>()

const open = defineModel<boolean>({ default: false })

const attrs = useAttrs()
const drawerRef = useTemplateRef<HTMLElement>('drawerRef')
const drawerStyle = ref<Record<string, string>>({})

let listenersActive = false

function updatePosition(): void {
  const anchor = anchorEl
  if (!anchor) return
  const rect = anchor.getBoundingClientRect()
  const style: Record<string, string> = {
    top: `${rect.bottom + offset}px`,
    '--dropdown-min-width': `${minWidth}px`,
  }

  if (align === 'start') {
    style.left = `${Math.max(8, rect.left)}px`
  } else {
    style.right = `${Math.max(8, window.innerWidth - rect.right)}px`
  }

  drawerStyle.value = style
}

function handleWindowChange(): void {
  if (!open.value) return
  updatePosition()
}

function handleClickOutside(event: MouseEvent): void {
  if (!open.value || !closeOnOutside) return
  const target = event.target as Node
  if (drawerRef.value?.contains(target)) return
  if (anchorEl?.contains(target)) return
  open.value = false
}

function handleKeydown(event: KeyboardEvent): void {
  if (!open.value || !closeOnEscape) return
  if (event.key === 'Escape') {
    open.value = false
  }
}

function addGlobalListeners(): void {
  if (listenersActive) return
  listenersActive = true
  if (closeOnOutside) {
    document.addEventListener('click', handleClickOutside)
  }
  if (closeOnEscape) {
    window.addEventListener('keydown', handleKeydown)
  }
  window.addEventListener('resize', handleWindowChange)
  window.addEventListener('scroll', handleWindowChange, true)
}

function removeGlobalListeners(): void {
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
      if (!anchorEl) {
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
  () => [anchorEl, align, offset, minWidth],
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
  box-shadow: var(--shadow-floating);
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
  background: var(--color-component-muted);
}

:deep(.dropdown-item.active) {
  background: var(--dropdown-active-bg, var(--color-primary));
  color: var(--dropdown-active-color, var(--color-on-primary));
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
