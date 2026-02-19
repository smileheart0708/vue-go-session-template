<script setup lang="ts">
import { onMounted } from 'vue'
import { RouterView } from 'vue-router'
import { ToastMessage } from '@/components/common'
import { useToast } from '@/composables'
import { useThemeStore } from '@/stores/theme'

const themeStore = useThemeStore()
const { toasts, removeToast } = useToast()

onMounted(() => {
  themeStore.init()
})
</script>

<template>
  <RouterView />
  <ToastMessage :toasts="toasts" close-aria-label="关闭通知" @close="removeToast" />
</template>
<style scoped></style>
