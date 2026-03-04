<template>
  <Teleport
    defer
    :to="teleportTo"
  >
    <Transition
      enter-active-class="transition duration-200 motion-reduce:transition-none"
      enter-from-class="-translate-y-2 opacity-0"
      enter-to-class="translate-y-0 opacity-100"
      leave-active-class="transition duration-200 motion-reduce:transition-none"
      leave-from-class="translate-y-0 opacity-100"
      leave-to-class="-translate-y-2 opacity-0"
    >
      <div
        v-if="open"
        ref="drawerRef"
        class="fixed z-(--dropdown-z-index,1000) min-w-(--dropdown-min-width,140px) rounded-lg border border-border bg-bg-glass p-1 shadow-floating backdrop-blur-md backdrop-saturate-140 [&_.dropdown-icon]:size-4 [&_.dropdown-icon]:shrink-0 [&_.dropdown-item]:flex [&_.dropdown-item]:w-full [&_.dropdown-item]:items-center [&_.dropdown-item]:gap-2 [&_.dropdown-item]:rounded-md [&_.dropdown-item]:border-0 [&_.dropdown-item]:bg-transparent [&_.dropdown-item]:px-3 [&_.dropdown-item]:py-2 [&_.dropdown-item]:text-left [&_.dropdown-item]:text-sm [&_.dropdown-item]:text-text-primary [&_.dropdown-item]:transition-colors [&_.dropdown-item]:duration-200 [&_.dropdown-item]:cursor-pointer [&_.dropdown-item:hover]:bg-bg-component-muted [&_.dropdown-item.active]:bg-accent [&_.dropdown-item.active]:text-on-accent"
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
import { nextTick, ref, useAttrs, useTemplateRef, watch } from 'vue'

defineOptions({ name: 'DropdownDrawer', inheritAttrs: false })

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

const open = defineModel<boolean>({ required: true })

const attrs = useAttrs()
const drawerRef = useTemplateRef<HTMLElement>('drawerRef')
const drawerStyle = ref<Record<string, string>>({})

function updatePosition(): void {
  const anchor = anchorEl
  if (!anchor) return
  const rect = anchor.getBoundingClientRect()
  const style: Record<string, string> = {
    top: `${rect.bottom + offset}px`,
    '--dropdown-min-width': `${minWidth}px`,
  }

  if (align === 'start') {
    style['left'] = `${Math.max(8, rect.left)}px`
  } else {
    style['right'] = `${Math.max(8, window.innerWidth - rect.right)}px`
  }

  drawerStyle.value = style
}

function handleWindowChange(): void {
  if (!open.value) return
  updatePosition()
}

function handleClickOutside(event: MouseEvent): void {
  if (!open.value || !closeOnOutside) return
  const target = event.target
  if (!(target instanceof Node)) return
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

function addGlobalListeners(shouldCloseOnOutside: boolean, shouldCloseOnEscape: boolean): void {
  if (shouldCloseOnOutside) {
    document.addEventListener('click', handleClickOutside)
  }
  if (shouldCloseOnEscape) {
    window.addEventListener('keydown', handleKeydown)
  }
  window.addEventListener('resize', handleWindowChange)
  window.addEventListener('scroll', handleWindowChange, true)
}

function removeGlobalListeners(shouldCloseOnOutside: boolean, shouldCloseOnEscape: boolean): void {
  if (shouldCloseOnOutside) {
    document.removeEventListener('click', handleClickOutside)
  }
  if (shouldCloseOnEscape) {
    window.removeEventListener('keydown', handleKeydown)
  }
  window.removeEventListener('resize', handleWindowChange)
  window.removeEventListener('scroll', handleWindowChange, true)
}

watch(
  (): [boolean, boolean, boolean] => [open.value, closeOnOutside, closeOnEscape],
  async ([isOpen, shouldCloseOnOutside, shouldCloseOnEscape]) => {
    if (!isOpen) return

    await nextTick()
    if (!anchorEl) {
      open.value = false
      return
    }

    updatePosition()
    addGlobalListeners(shouldCloseOnOutside, shouldCloseOnEscape)

    return () => {
      removeGlobalListeners(shouldCloseOnOutside, shouldCloseOnEscape)
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
</script>
