<template>
  <aside class="sidebar" :class="{ 'sidebar-open': isOpen }">
    <SidebarHeader />
    <SidebarNav />
    <SidebarFooter />
  </aside>
</template>

<script setup lang="ts">
import SidebarHeader from './SidebarHeader.vue'
import SidebarNav from './SidebarNav.vue'
import SidebarFooter from './SidebarFooter.vue'

interface Props {
  isOpen?: boolean
}

withDefaults(defineProps<Props>(), { isOpen: true })
</script>

<style scoped>
.sidebar {
  display: flex;
  flex-direction: column;
  width: var(--sidebar-width);
  height: 100vh;
  background: var(--color-background-elevated);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-right: 1px solid var(--color-border);
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: fixed;
  left: 0;
  top: 0;
  z-index: 100;
}

@media (max-width: 767px) {
  .sidebar {
    transform: translateX(-100%);
  }

  .sidebar-open {
    transform: translateX(0);
  }
}
</style>
