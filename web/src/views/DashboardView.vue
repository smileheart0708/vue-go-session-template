<template>
  <div class="flex w-full flex-col gap-6">
    <DashboardStats />
    <div class="grid grid-cols-2 items-stretch gap-6 max-lg:grid-cols-1">
      <RequestChart />
      <ModelDistribution />
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import { useIntervalFn } from '@vueuse/core'
import { DashboardStats, ModelDistribution, RequestChart } from '@/components'
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
