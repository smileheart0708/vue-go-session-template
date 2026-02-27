<template>
  <div class="flex h-screen h-dvh w-full overflow-hidden">
    <AppSidebar :is-open="sidebarOpen" />

    <div class="ml-0 flex min-h-0 min-w-0 flex-1 flex-col md:ml-sidebar">
      <AppHeader @toggle-sidebar="toggleSidebar" />

      <main
        class="flex min-h-0 min-w-0 flex-1 flex-col overflow-y-auto overflow-x-hidden px-6 pb-6 pt-[calc(var(--sys-layout-header-height)+1.5rem)] overscroll-y-contain [scrollbar-gutter:stable_both-edges] max-md:px-4 max-md:pb-4 max-md:pt-[calc(var(--sys-layout-header-height)+1rem)]"
      >
        <RouterView v-slot="{ Component }">
          <Transition
            mode="out-in"
            enter-active-class="transition-opacity duration-300"
            enter-from-class="opacity-0"
            enter-to-class="opacity-100"
            leave-active-class="transition-opacity duration-300"
            leave-from-class="opacity-100"
            leave-to-class="opacity-0"
          >
            <component :is="Component" />
          </Transition>
        </RouterView>
      </main>
    </div>

    <Transition
      enter-active-class="transition-opacity duration-300"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition-opacity duration-300"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="sidebarOpen"
        class="fixed inset-0 z-[99] bg-[var(--sys-color-overlay)] backdrop-blur-[6px] md:hidden"
        @click="toggleSidebar"
      />
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import AppSidebar from './AppSidebar.vue'
import AppHeader from './AppHeader.vue'

const sidebarOpen = ref(false)

function toggleSidebar(): void {
  sidebarOpen.value = !sidebarOpen.value
}

onMounted(() => {
  document.body.style.overflow = 'hidden'
})

onUnmounted(() => {
  document.body.style.overflow = ''
})
</script>
