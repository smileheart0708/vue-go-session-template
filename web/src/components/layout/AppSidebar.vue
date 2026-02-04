<template>
  <aside class="sidebar" :class="{ 'sidebar-open': isOpen }">
    <div class="sidebar-drawer">
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

withDefaults(defineProps<Props>(), { isOpen: true })
</script>

<style scoped>
.sidebar {
  display: flex;
  flex-direction: column;
  width: var(--sidebar-width);
  height: 100vh;
  background: transparent;
  transition: transform 0.35s cubic-bezier(0.2, 0.8, 0.2, 1);
  position: fixed;
  left: 0;
  top: 0;
  z-index: 100;
}

.sidebar-drawer {
  display: flex;
  flex-direction: column;
  width: 100%;
  height: 100%;
  background: var(--color-background-glass);
  border-right: 1px solid var(--color-border);
  border-radius: 0;
  box-shadow: none;
  backdrop-filter: blur(28px) saturate(1.1);
  -webkit-backdrop-filter: blur(28px) saturate(1.1);
  overflow: hidden;
  transition:
    transform 0.35s cubic-bezier(0.2, 0.8, 0.2, 1),
    opacity 0.25s ease,
    box-shadow 0.3s ease;
  will-change: transform, opacity;
}

@media (max-width: 767px) {
  .sidebar {
    transform: translate3d(-100%, 0, 0);
  }

  .sidebar .sidebar-drawer {
    opacity: 0;
    transform: translate3d(-12px, 0, 0) scale(0.98);
    border-radius: 0 16px 16px 0;
  }

  .sidebar-open {
    transform: translate3d(0, 0, 0);
  }

  .sidebar-open .sidebar-drawer {
    opacity: 1;
    transform: translate3d(0, 0, 0) scale(1);
  }
}
</style>
