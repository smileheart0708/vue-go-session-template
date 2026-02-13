<template>
  <section class="content-card">
    <div class="card-header">
      <h2 class="card-title">模型分布</h2>
    </div>
    <div class="chart-wrapper">
      <div ref="chartRef" class="chart-container"></div>
    </div>

    <!-- 自定义两栏图例 -->
    <div class="custom-legend">
      <div v-for="(item, index) in chartData" :key="item.name" class="legend-item">
        <span class="legend-dot" :style="{ backgroundColor: getChartColors()[index] }"></span>
        <span class="legend-text" :title="item.name">{{ item.name }}</span>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch, computed } from 'vue'
import * as echarts from 'echarts/core'
import { PieChart } from 'echarts/charts'
import { TitleComponent, TooltipComponent, LegendComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import { useTheme } from '@/composables/useTheme'

echarts.use([PieChart, TitleComponent, TooltipComponent, LegendComponent, CanvasRenderer])

const { isDark } = useTheme()
const chartRef = ref<HTMLElement>()
let chartInstance: echarts.ECharts | null = null
let resizeObserver: ResizeObserver | null = null

// 模拟数据 - 后续可替换为真实 API 数据
const mockData = [
  { model: 'GPT-4', count: 1250 },
  { model: 'GPT-3.5', count: 890 },
  { model: 'Claude-3', count: 640 },
  { model: 'Gemini Pro', count: 420 },
  { model: 'Llama 3', count: 280 },
]

// 处理数据：最多显示5个，剩下的合并为"其他"
const chartData = computed(() => {
  const modelStats = mockData
  if (modelStats.length === 0) {
    return []
  }

  if (modelStats.length <= 5) {
    return modelStats.map((stat) => ({ value: stat.count, name: stat.model }))
  }

  const top5 = modelStats.slice(0, 5)
  const others = modelStats.slice(5)
  const othersTotal = others.reduce((sum, item) => sum + item.count, 0)

  return [
    ...top5.map((stat) => ({ value: stat.count, name: stat.model })),
    { value: othersTotal, name: '其他' },
  ]
})

// 从 CSS 变量获取图表颜色
const getChartColors = () => {
  return [
    getComputedStyleValue('--color-primary') || '#2378ff',
    getComputedStyleValue('--toast-success') || '#22c55e',
    getComputedStyleValue('--toast-info') || '#3b82f6',
    getComputedStyleValue('--toast-warning') || '#f59e0b',
    getComputedStyleValue('--toast-error') || '#ef4444',
    getComputedStyleValue('--color-text-tertiary') || '#adb5bd', // 灰色用于"其他"
  ]
}

const total = computed(() => chartData.value.reduce((acc, item) => acc + item.value, 0))

const getComputedStyleValue = (variable: string) => {
  return getComputedStyle(document.documentElement).getPropertyValue(variable).trim()
}

const updateChart = () => {
  if (!chartInstance) return

  const textColor = getComputedStyleValue('--color-text-secondary') || '#666666'
  const textPrimary = getComputedStyleValue('--color-text') || '#1a1a1a'
  const tooltipBg = getComputedStyleValue('--color-tooltip-bg') || '#ffffff'
  const tooltipBorderColor = getComputedStyleValue('--color-tooltip-border') || '#d4d4d4'
  const tooltipTextColor = getComputedStyleValue('--color-tooltip-text') || '#333333'
  const chartColors = getChartColors()

  const option = {
    color: chartColors,
    tooltip: {
      trigger: 'item',
      formatter: '{b}: {c} ({d}%)',
      backgroundColor: tooltipBg,
      borderColor: tooltipBorderColor,
      borderWidth: 1,
      textStyle: { color: tooltipTextColor },
    },
    legend: { show: false },
    title: {
      text: total.value.toString(),
      subtext: '总调用',
      left: 'center',
      top: '38%',
      textStyle: { fontSize: 32, fontWeight: 'bold', color: textPrimary },
      subtextStyle: { fontSize: 14, color: textColor },
    },
    series: [
      {
        name: '模型分布',
        type: 'pie',
        radius: ['65%', '85%'],
        center: ['50%', '50%'],
        avoidLabelOverlap: false,
        padAngle: 5,
        itemStyle: { borderRadius: 10 },
        label: { show: false },
        data: chartData.value,
        animationType: 'expansion',
        animationDuration: 1000,
        animationEasing: 'exponentialOut',
      },
    ],
  }

  chartInstance.setOption(option)
}

onMounted(() => {
  if (chartRef.value) {
    // 使用 canvas 渲染器并启用脏矩形优化，减少重绘区域
    chartInstance = echarts.init(chartRef.value, null, { renderer: 'canvas', useDirtyRect: true })

    let isFirstRender = true
    // 使用 ResizeObserver 监听容器大小变化
    resizeObserver = new ResizeObserver(() => {
      // 使用 requestAnimationFrame 避免 ResizeObserver loop limit exceeded 错误
      // 并确保在下一帧进行 resize，此时 DOM 尺寸已更新
      window.requestAnimationFrame(() => {
        if (!chartInstance || !chartRef.value) return

        // 检查容器是否可见且有尺寸
        const { clientWidth, clientHeight } = chartRef.value
        if (clientWidth > 0 && clientHeight > 0) {
          chartInstance.resize()

          // 关键：在容器尺寸确定后进行初次渲染，避免动画被初始的 resize 打断
          if (isFirstRender) {
            updateChart()
            isFirstRender = false
          }
        }
      })
    })
    resizeObserver.observe(chartRef.value)
  }
})

watch(isDark, () => {
  setTimeout(updateChart, 50)
})

onUnmounted(() => {
  if (chartInstance) {
    chartInstance.dispose()
    chartInstance = null
  }
  if (resizeObserver) {
    resizeObserver.disconnect()
    resizeObserver = null
  }
})
</script>

<style scoped>
.content-card {
  padding: 1.5rem;
  background: var(--color-background-elevated);
  border: 1px solid var(--color-border);
  border-radius: 12px;
  display: flex;
  flex-direction: column;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
  gap: 1rem;
}

.card-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--color-text);
  margin: 0;
  white-space: nowrap;
}

.chart-wrapper {
  position: relative;
  flex: 1;
  min-height: 240px;
  width: 100%;
}

.chart-container {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
}

.custom-legend {
  margin-top: 1rem;
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 0.75rem;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  min-width: 0;
}

.legend-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  flex-shrink: 0;
}

.legend-text {
  font-size: 0.75rem;
  color: var(--color-text-secondary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
