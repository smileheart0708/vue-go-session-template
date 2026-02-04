<template>
  <div class="sidebar-footer">
    <BaseButton width="100%" :height="40" text="退出登录" :icon="LogOut" @click="handleLogout" />
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { LogOut } from 'lucide-vue-next'
import BaseButton from '@/components/common/BaseButton.vue'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

async function handleLogout() {
  await authStore.logout()
  router.replace('/login')
}
</script>

<style scoped>
.sidebar-footer {
  padding: 1.25rem 1.25rem 1.5rem;
  min-height: 80px;
  background: transparent;
}
</style>
