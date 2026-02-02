<script setup lang="ts">
import { onMounted } from 'vue'
import { useThemeStore } from '@/stores/theme'
import LoginView from '@/views/LoginView.vue'
import ToastMessage from '@/components/common/ToastMessage.vue'

const themeStore = useThemeStore()

onMounted(() => {
  themeStore.init()
})
</script>
<template>
  <LoginView />
  <ToastMessage />
</template>
<style scoped></style>
