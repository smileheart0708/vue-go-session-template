<template>
  <div class="main-layout">
    <AppSidebar :is-open="sidebarOpen" />
    <div class="main-content">
      <AppHeader @toggle-sidebar="toggleSidebar" />
      <main class="content">
        <RouterView v-slot="{ Component }">
          <Transition name="fade" mode="out-in">
            <component :is="Component" />
          </Transition>
        </RouterView>
      </main>
    </div>
    <Transition name="fade">
      <div v-if="sidebarOpen" class="sidebar-overlay" @click="toggleSidebar" />
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import AppSidebar from './AppSidebar.vue'
import AppHeader from './AppHeader.vue'

const BODY_SCROLL_LOCK_CLASS = 'main-layout-scroll-lock'

const sidebarOpen = ref(false)

function toggleSidebar() {
  sidebarOpen.value = !sidebarOpen.value
}

onMounted(() => {
  document.body.classList.add(BODY_SCROLL_LOCK_CLASS)
})

onUnmounted(() => {
  document.body.classList.remove(BODY_SCROLL_LOCK_CLASS)
  document.body.style.overflow = ''
})
</script>

<style scoped>
.main-layout {
  display: flex;
  width: 100%;
  height: 100vh;
  height: 100dvh;
  overflow: hidden;
}

.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  margin-left: var(--sidebar-width);
  min-width: 0;
  min-height: 0;
  height: 100%;
  overflow: hidden;
}

.content {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: calc(var(--header-height) + 1.5rem) 1.5rem 1.5rem;
  min-width: 0;
  min-height: 0;
  overflow-y: auto;
  overflow-x: hidden;
  overscroll-behavior-y: contain;
}

.sidebar-overlay {
  position: fixed;
  inset: 0;
  background: var(--color-overlay);
  backdrop-filter: blur(6px);
  -webkit-backdrop-filter: blur(6px);
  z-index: 99;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

:global(body.main-layout-scroll-lock) {
  overflow: hidden;
}

@media (max-width: 767px) {
  .main-content {
    margin-left: 0;
  }
}
</style>
