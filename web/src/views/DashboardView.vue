<template>
  <div class="dashboard">
    <DashboardStats />
    <ModelDistribution />
  </div>
</template>

<script setup lang="ts">
import DashboardStats from '@/components/dashboard/DashboardStats.vue'
import ModelDistribution from '@/components/dashboard/ModelDistribution.vue'
</script>

<style scoped>
.dashboard {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}
</style>
