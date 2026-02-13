<template>
  <div class="stats-grid">
    <StatsCard
      label="24小时总请求"
      :value="operationCount"
      :change="operationChange"
      color="var(--toast-error)"
    >
      <template #icon>
        <TrendingUp />
      </template>
    </StatsCard>

    <StatsCard
      label="额度消耗"
      :value="tokenUsage"
      :change="tokenDetail"
      color="var(--toast-info)"
    >
      <template #icon>
        <Coins />
      </template>
    </StatsCard>

    <StatsCard
      label="内存占用"
      :value="memoryUsage"
      :change="memoryDetail"
      color="var(--toast-warning)"
    >
      <template #icon>
        <IconMemory />
      </template>
    </StatsCard>

    <StatsCard
      label="运行时间"
      :value="dashboardStore.uptime"
      :change="dashboardStore.startTime"
      color="var(--toast-success)"
    >
      <template #icon>
        <Clock />
      </template>
    </StatsCard>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { TrendingUp, Coins, Clock } from 'lucide-vue-next'
import IconMemory from '@/components/icons/IconMemory.vue'
import StatsCard from '@/components/common/StatsCard.vue'
import { useDashboardStore } from '@/stores/dashboard'
import { formatNumber, formatBytes, formatTokens } from '@/utils'

const dashboardStore = useDashboardStore()

// 计算属性：24小时总请求
const operationCount = computed(() => {
  if (!dashboardStore.dataAvailable) return '0'
  return formatNumber(dashboardStore.totalRequests)
})

const operationChange = computed(() => {
  if (!dashboardStore.dataAvailable) return '-'
  const success = dashboardStore.successCount
  const failed = dashboardStore.failedCount
  return `成功 ${formatNumber(success)} / 失败 ${formatNumber(failed)}`
})

// 计算属性：Token 消耗
const tokenUsage = computed(() => {
  if (!dashboardStore.dataAvailable) return '0'
  return formatTokens(dashboardStore.totalTokens)
})

const tokenDetail = computed(() => {
  if (!dashboardStore.dataAvailable) return '-'
  const prompt = dashboardStore.promptTokens
  const completion = dashboardStore.completionTokens
  return `输入 ${formatTokens(prompt)} / 输出 ${formatTokens(completion)}`
})

// 计算属性：内存占用
const memoryUsage = computed(() => {
  if (!dashboardStore.dataAvailable) return '0%'
  return `${dashboardStore.memoryPercent.toFixed(1)}%`
})

const memoryDetail = computed(() => {
  if (!dashboardStore.dataAvailable) return '-'
  const used = formatBytes(dashboardStore.memoryUsed)
  const total = formatBytes(dashboardStore.memoryTotal)
  return `${used} / ${total}`
})
</script>

<style scoped>
.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 1.5rem;
}

@media (max-width: 1280px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 640px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }
}
</style>
