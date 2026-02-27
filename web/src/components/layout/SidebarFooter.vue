<template>
  <div class="min-h-20 bg-transparent px-5 pb-6 pt-5">
    <BaseButton width="100%" :height="40" text="退出登录" :icon="LogOut" @click="handleLogout" />
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { LogOut } from 'lucide-vue-next'
import { BaseButton } from '@/components/common'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

async function handleLogout(): Promise<void> {
  await authStore.logout()
  router.replace('/login')
}
</script>
