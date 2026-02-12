<template>
  <header class="mobile-header">
    <IconButton
      class="menu-toggle"
      title="切换侧边栏"
      @click="emit('toggle-sidebar')"
    >
      <Menu />
    </IconButton>

    <h1 class="header-title">{{ currentTitle }}</h1>

    <div class="header-actions">
      <ThemeToggle />
    </div>
  </header>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { Menu } from 'lucide-vue-next'
import IconButton from '@/components/common/IconButton.vue'
import ThemeToggle from '@/components/common/ThemeToggle.vue'

const emit = defineEmits<{ 'toggle-sidebar': [] }>()

const route = useRoute()

const currentTitle = computed(() => {
  // 从路由 meta 中获取标题
  const title = route.meta.title as string | undefined
  return title || '仪表板'
})
</script>

<style scoped>
.mobile-header {
  display: flex;
  align-items: center;
  gap: 1rem;
  height: var(--header-height);
  padding: 0 1rem;
  background: var(--color-background-elevated);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-bottom: 1px solid var(--color-border);
  position: sticky;
  top: 0;
  z-index: 90;
}

.header-title {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--color-text);
  margin: 0;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

@media (min-width: 768px) {
  .menu-toggle {
    display: none;
  }
}
</style>
