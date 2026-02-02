<script setup lang="ts">
import { onMounted } from 'vue'
import { useThemeStore } from '@/stores/theme'
import LoginView from '@/views/LoginView.vue'

const themeStore = useThemeStore()

onMounted(() => {
  themeStore.init()
})
</script>
<template>
  <LoginView />
</template>
<style scoped></style>
