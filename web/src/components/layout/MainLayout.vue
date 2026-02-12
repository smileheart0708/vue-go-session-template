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
import { ref, watch, onUnmounted } from 'vue'
import AppSidebar from './AppSidebar.vue'
import AppHeader from './AppHeader.vue'

const sidebarOpen = ref(false)

function toggleSidebar() {
  sidebarOpen.value = !sidebarOpen.value
}

// 监听侧边栏状态，在移动端控制 body 滚动
watch(sidebarOpen, (isOpen) => {
  // 只在移动端（宽度小于 768px）禁止滚动
  if (window.innerWidth < 768) {
    if (isOpen) {
      document.body.style.overflow = 'hidden'
    } else {
      document.body.style.overflow = ''
    }
  }
})

// 组件卸载时恢复滚动
onUnmounted(() => {
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
  padding: 1.5rem;
  overflow-y: auto;
  overflow-x: hidden;
  overscroll-behavior-y: contain;
  min-width: 0;
  min-height: 0;
}

.sidebar-overlay {
  position: fixed;
  inset: 0;
  background: rgba(10, 12, 16, 0.35);
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

@media (max-width: 767px) {
  .main-content {
    margin-left: 0;
  }
}
</style>
