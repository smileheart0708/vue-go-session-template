<template>
  <section class="flex flex-col rounded-xl border border-border bg-bg-surface p-6 max-md:p-4">
    <div class="mb-6 flex items-center justify-between gap-4 max-md:mb-4">
      <h2
        class="m-0 whitespace-nowrap text-xl font-semibold text-text-primary max-md:text-[1.1rem]"
      >
        模型分布
      </h2>
    </div>
    <div class="relative w-full min-h-60 flex-1">
      <div
        ref="chartRef"
        class="absolute left-0 top-0 h-full w-full"
      ></div>
    </div>

    <!-- 自定义两栏图例 -->
    <div class="mt-4 grid grid-cols-2 gap-3">
      <div
        v-for="(item, index) in chartData"
        :key="item.name"
        class="flex min-w-0 items-center gap-2"
      >
        <span
          class="h-2.5 w-2.5 shrink-0 rounded-full"
          :style="{ backgroundColor: PIE_CHART_COLORS[index] }"
        ></span>
        <span
          class="overflow-hidden text-ellipsis whitespace-nowrap text-xs text-text-secondary"
          :title="item.name"
        >
          {{ item.name }}
        </span>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, useTemplateRef, watch } from 'vue'
import * as echarts from 'echarts/core'
import { PieChart } from 'echarts/charts'
import { TitleComponent, TooltipComponent, LegendComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import { useTheme } from '@/composables/useTheme'

echarts.use([PieChart, TitleComponent, TooltipComponent, LegendComponent, CanvasRenderer])

const { isDark } = useTheme()
const chartRef = useTemplateRef<HTMLElement>('chartRef')
let chartInstance: echarts.ECharts | null = null
let resizeObserver: ResizeObserver | null = null

const PIE_CHART_COLORS: string[] = [
  '#0078d4',
  '#22c55e',
  '#3b82f6',
  '#f59e0b',
  '#ef4444',
  '#adb5bd',
]

// 模拟数据 - 后续可替换为真实 API 数据
const mockData = [
  { model: 'GPT-5-3-codex', count: 1250 },
  { model: 'GPT-5.2', count: 890 },
  { model: 'Claude-4.6', count: 640 },
  { model: 'Gemini 3 Pro', count: 420 },
  { model: 'Grok 4.2', count: 280 },
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

const total = computed(() => chartData.value.reduce((acc, item) => acc + item.value, 0))

const getComputedStyleValue = (variable: string) => {
  return getComputedStyle(document.documentElement).getPropertyValue(variable).trim()
}

const updateChart = () => {
  if (!chartInstance) return

  const textColor = getComputedStyleValue('--sys-color-text-secondary') || '#666666'
  const textPrimary = getComputedStyleValue('--sys-color-text-primary') || '#1a1a1a'
  const tooltipBg = getComputedStyleValue('--sys-color-tooltip-bg') || '#ffffff'
  const tooltipBorderColor = getComputedStyleValue('--sys-color-tooltip-border') || '#d4d4d4'
  const tooltipTextColor = getComputedStyleValue('--sys-color-tooltip-text') || '#333333'

  const option = {
    color: PIE_CHART_COLORS,
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
  const timer = setTimeout(updateChart, 50)
  return () => clearTimeout(timer)
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
