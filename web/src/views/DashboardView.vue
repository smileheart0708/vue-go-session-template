<template>
  <div class="dashboard">
    <DashboardStats />
    <ModelDistribution />
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import DashboardStats from '@/components/dashboard/DashboardStats.vue'
import ModelDistribution from '@/components/dashboard/ModelDistribution.vue'
import { useDashboardStore } from '@/stores/dashboard'

const dashboardStore = useDashboardStore()

let uptimeTimer: number | null = null
let statsTimer: number | null = null

async function loadStats() {
  await dashboardStore.fetchSystemStats()
}

onMounted(async () => {
  await loadStats()

  uptimeTimer = window.setInterval(() => {
    dashboardStore.refreshUptime()
  }, 1000)

  statsTimer = window.setInterval(() => {
    void loadStats()
  }, 10000)
})

onUnmounted(() => {
  if (uptimeTimer !== null) {
    window.clearInterval(uptimeTimer)
  }

  if (statsTimer !== null) {
    window.clearInterval(statsTimer)
  }
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
