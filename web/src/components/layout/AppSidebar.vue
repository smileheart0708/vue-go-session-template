<template>
  <aside
    class="fixed left-0 top-0 z-[100] flex h-screen h-dvh w-sidebar flex-col bg-transparent transition-transform duration-300 ease-[cubic-bezier(0.2,0.8,0.2,1)] will-change-transform"
    :class="isOpen ? 'max-md:translate-x-0 max-md:pointer-events-auto' : 'max-md:-translate-x-full max-md:pointer-events-none'"
  >
    <div
      class="flex h-full w-full flex-col overflow-hidden border-r border-border bg-bg-glass backdrop-blur-[28px] backdrop-saturate-[1.1] transition-all duration-300 ease-[cubic-bezier(0.2,0.8,0.2,1)] will-change-transform max-md:rounded-r-2xl"
      :class="isOpen ? 'max-md:translate-x-0 max-md:scale-100 max-md:opacity-100' : 'max-md:-translate-x-3 max-md:scale-[0.98] max-md:opacity-0'"
    >
      <SidebarHeader />
      <SidebarNav />
      <SidebarFooter />
    </div>
  </aside>
</template>

<script setup lang="ts">
import SidebarHeader from './SidebarHeader.vue'
import SidebarNav from './SidebarNav.vue'
import SidebarFooter from './SidebarFooter.vue'

interface Props {
  isOpen?: boolean
}

const { isOpen = true } = defineProps<Props>()
</script>
