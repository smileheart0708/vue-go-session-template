import { ref } from 'vue'
import { defineStore } from 'pinia'

export interface DashboardData {
  totalRequests: number
  successCount: number
  failedCount: number
  totalTokens: number
  promptTokens: number
  completionTokens: number
  memoryPercent: number
  memoryUsed: number
  memoryTotal: number
  uptime: string
  startTime: string
}

export const useDashboardStore = defineStore('dashboard', () => {
  // 数据是否可用
  const dataAvailable = ref(false)

  // 请求数据
  const totalRequests = ref(0)
  const successCount = ref(0)
  const failedCount = ref(0)

  // Token 数据
  const totalTokens = ref(0)
  const promptTokens = ref(0)
  const completionTokens = ref(0)

  // 内存数据
  const memoryPercent = ref(0)
  const memoryUsed = ref(0)
  const memoryTotal = ref(0)

  // 运行时间
  const uptime = ref('0天 0小时 0分钟')
  const startTime = ref('')

  // 设置数据
  function setData(data: Partial<DashboardData>) {
    if (data.totalRequests !== undefined) totalRequests.value = data.totalRequests
    if (data.successCount !== undefined) successCount.value = data.successCount
    if (data.failedCount !== undefined) failedCount.value = data.failedCount
    if (data.totalTokens !== undefined) totalTokens.value = data.totalTokens
    if (data.promptTokens !== undefined) promptTokens.value = data.promptTokens
    if (data.completionTokens !== undefined) completionTokens.value = data.completionTokens
    if (data.memoryPercent !== undefined) memoryPercent.value = data.memoryPercent
    if (data.memoryUsed !== undefined) memoryUsed.value = data.memoryUsed
    if (data.memoryTotal !== undefined) memoryTotal.value = data.memoryTotal
    if (data.uptime !== undefined) uptime.value = data.uptime
    if (data.startTime !== undefined) startTime.value = data.startTime

    dataAvailable.value = true
  }

  // 重置数据
  function reset() {
    dataAvailable.value = false
    totalRequests.value = 0
    successCount.value = 0
    failedCount.value = 0
    totalTokens.value = 0
    promptTokens.value = 0
    completionTokens.value = 0
    memoryPercent.value = 0
    memoryUsed.value = 0
    memoryTotal.value = 0
    uptime.value = '0天 0小时 0分钟'
    startTime.value = ''
  }

  // 模拟数据（用于开发测试）
  function loadMockData() {
    setData({
      totalRequests: 128456,
      successCount: 127890,
      failedCount: 566,
      totalTokens: 4567890,
      promptTokens: 2345678,
      completionTokens: 2222212,
      memoryPercent: 45.6,
      memoryUsed: 7.4 * 1024 * 1024 * 1024, // 7.4 GB
      memoryTotal: 16 * 1024 * 1024 * 1024, // 16 GB
      uptime: '15天 8小时 32分钟',
      startTime: '2025-01-18 04:18:14',
    })
  }

  return {
    dataAvailable,
    totalRequests,
    successCount,
    failedCount,
    totalTokens,
    promptTokens,
    completionTokens,
    memoryPercent,
    memoryUsed,
    memoryTotal,
    uptime,
    startTime,
    setData,
    reset,
    loadMockData,
  }
})
