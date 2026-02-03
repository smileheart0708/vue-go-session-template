<template>
  <div class="main-layout">
    <AppSidebar :is-open="sidebarOpen" />
    <div class="main-content">
      <AppHeader @toggle-sidebar="toggleSidebar" />
      <main class="content">
        <RouterView />
      </main>
    </div>
    <div v-if="sidebarOpen" class="sidebar-overlay" @click="toggleSidebar" />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import AppSidebar from './AppSidebar.vue'
import AppHeader from './AppHeader.vue'

const sidebarOpen = ref(false)

function toggleSidebar() {
  sidebarOpen.value = !sidebarOpen.value
}
</script>

<style scoped>
.main-layout {
  display: flex;
  min-height: 100vh;
}

.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  margin-left: var(--sidebar-width);
  min-height: 100vh;
}

.content {
  flex: 1;
  padding: 1.5rem;
  overflow-y: auto;
}

.sidebar-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 99;
}

@media (max-width: 767px) {
  .main-content {
    margin-left: 0;
  }
}
</style>
