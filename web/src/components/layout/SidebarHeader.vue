<template>
  <div class="sidebar-header">
    <div class="logo-container">
      <img :src="logoSvg" alt="logo" class="logo-icon" />
      <h2 class="sidebar-title">CodeHub</h2>
    </div>
  </div>
</template>

<script setup lang="ts">
import logoSvg from '@/assets/logo.svg'
</script>

<style scoped>
.sidebar-header {
  padding: 1.25rem 1.25rem 0.75rem;
  background: transparent;
}

.logo-container {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.logo-icon {
  width: 28px;
  height: 28px;
  flex-shrink: 0;
}

.sidebar-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--color-primary);
  margin: 0;
}
</style>
