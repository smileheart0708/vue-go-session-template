<script setup lang="ts">
import { onMounted } from 'vue'
import { RouterView } from 'vue-router'
import { useThemeStore } from '@/stores/theme'
import ToastMessage from '@/components/common/ToastMessage.vue'

const themeStore = useThemeStore()

onMounted(() => {
  themeStore.init()
})
</script>

<template>
  <RouterView />
  <ToastMessage />
</template>
<style scoped></style>
