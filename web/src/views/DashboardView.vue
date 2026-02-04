<template>
  <div class="dashboard">
    <DashboardStats />
    <ModelDistribution />
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import { useIntervalFn } from '@vueuse/core'
import DashboardStats from '@/components/dashboard/DashboardStats.vue'
import ModelDistribution from '@/components/dashboard/ModelDistribution.vue'
import { useDashboardStore } from '@/stores/dashboard'
import { useRefreshStore } from '@/stores/refresh'

const dashboardStore = useDashboardStore()
const refreshStore = useRefreshStore()

async function loadStats() {
  await dashboardStore.fetchDashboardStats()
}

const { pause: pauseUptime, resume: resumeUptime } = useIntervalFn(
  () => {
    dashboardStore.refreshUptime()
  },
  1000,
  { immediate: false },
)

let unregisterStats: (() => void) | undefined

onMounted(async () => {
  unregisterStats = refreshStore.register('dashboard-stats', loadStats)
  await loadStats()
  resumeUptime()
})

onUnmounted(() => {
  unregisterStats?.()
  pauseUptime()
})
</script>

<style scoped>
.dashboard {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}
</style>
