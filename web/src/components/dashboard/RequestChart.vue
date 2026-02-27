<template>
  <section
    class="box-border flex w-full min-h-[400px] flex-col overflow-hidden rounded-xl border border-border bg-bg-surface p-6 max-md:min-h-[350px] max-md:p-4"
  >
    <div class="mb-6 flex items-center justify-between max-md:mb-4">
      <h2 class="m-0 text-xl font-semibold text-text-primary max-md:text-[1.1rem]">{{ title }}</h2>
    </div>
    <div ref="chartRef" class="w-full min-h-[300px] flex-1 max-md:min-h-[250px]"></div>
  </section>
</template>

<script setup lang="ts">
import {
  computed,
  nextTick,
  onMounted,
  onUnmounted,
  onWatcherCleanup,
  useTemplateRef,
  watch,
} from 'vue'
import * as echarts from 'echarts/core'
import { LineChart } from 'echarts/charts'
import { TooltipComponent, GridComponent, LegendComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import { useTheme } from '@/composables/useTheme'

echarts.use([LineChart, TooltipComponent, GridComponent, LegendComponent, CanvasRenderer])

interface HourlyRequestStat {
  timestamp: number
  success: number
  failed: number
}

interface Props {
  title?: string
  stats?: HourlyRequestStat[]
}

const { title = '请求趋势', stats } = defineProps<Props>()

const { isDark } = useTheme()
const chartRef = useTemplateRef<HTMLElement>('chartRef')
let chartInstance: echarts.ECharts | null = null
let resizeObserver: ResizeObserver | null = null

const LINE_SERIES_COLORS = { success: '#22c55e', failed: '#ef4444' } as const

const getComputedStyleValue = (variable: string) => {
  return getComputedStyle(document.documentElement).getPropertyValue(variable).trim()
}

const createMockHourlyStats = (hours = 24): HourlyRequestStat[] => {
  const now = Date.now()
  const stats: HourlyRequestStat[] = []

  for (let index = hours - 1; index >= 0; index -= 1) {
    const timestamp = Math.floor((now - index * 60 * 60 * 1000) / 1000)
    const wave = Math.sin(index / 3) * 25 + Math.cos(index / 5) * 10
    const success = Math.max(20, Math.round(140 + wave))
    const failed = Math.max(1, Math.round(6 + Math.cos(index / 2) * 4))

    stats.push({ timestamp, success, failed })
  }

  return stats
}

const mockHourlyStats = createMockHourlyStats()

const resolvedStats = computed(() => {
  return stats === undefined ? mockHourlyStats : stats
})

// 计算图表数据
const chartData = computed(() => {
  const hourlyStats = resolvedStats.value
  const xAxisData: string[] = []
  const successData: number[] = []
  const failedData: number[] = []

  // 遍历24小时数据（已经是按时间戳排序的滚动窗口）
  for (const stat of hourlyStats) {
    const date = new Date(stat.timestamp * 1000)
    const month = (date.getMonth() + 1).toString()
    const day = date.getDate().toString()
    const hour = date.getHours().toString().padStart(2, '0')
    xAxisData.push(`${month}月${day}日 ${hour}时`)

    successData.push(stat.success)
    failedData.push(stat.failed)
  }

  return { xAxisData, successData, failedData }
})

const updateChart = () => {
  if (!chartInstance) return

  const successColor = LINE_SERIES_COLORS.success
  const errorColor = LINE_SERIES_COLORS.failed
  const textColor = getComputedStyleValue('--sys-color-text-secondary') || '#6c757d'
  const tooltipBg = getComputedStyleValue('--sys-color-tooltip-bg') || '#ffffff'
  const tooltipBorderColor = getComputedStyleValue('--sys-color-tooltip-border') || '#d4d3cc'
  const tooltipTextColor = getComputedStyleValue('--sys-color-tooltip-text') || '#333333'
  const borderColor = getComputedStyleValue('--sys-color-border') || '#d4d3cc'

  const { xAxisData, successData, failedData } = chartData.value

  const option = {
    color: [successColor, errorColor],
    tooltip: {
      trigger: 'axis',
      backgroundColor: tooltipBg,
      borderColor: tooltipBorderColor,
      borderWidth: 1,
      textStyle: { color: tooltipTextColor },
      axisPointer: { type: 'line', lineStyle: { color: borderColor, type: 'dashed' } },
    },
    legend: {
      data: ['成功请求', '失败请求'],
      bottom: 0,
      padding: [5, 0],
      textStyle: { color: textColor, fontSize: 12 },
      itemWidth: 14,
      itemHeight: 14,
    },
    grid: { top: '15%', left: '2%', right: '2%', bottom: '15%', containLabel: true },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: xAxisData,
      axisLine: { lineStyle: { color: borderColor } },
      axisLabel: { color: textColor, fontSize: 11, interval: 'auto' },
    },
    yAxis: {
      type: 'value',
      minInterval: 1,
      axisLine: { show: false },
      axisLabel: { color: textColor, fontSize: 11 },
      splitLine: { lineStyle: { color: borderColor, type: 'dashed' } },
    },
    series: [
      {
        name: '成功请求',
        type: 'line',
        smooth: true,
        showSymbol: false,
        data: successData,
        animationDuration: 1000,
        animationEasing: 'cubicOut',
        lineStyle: { width: 3 },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: `${successColor}4D` },
            { offset: 1, color: `${successColor}00` },
          ]),
        },
      },
      {
        name: '失败请求',
        type: 'line',
        smooth: true,
        showSymbol: false,
        data: failedData,
        animationDuration: 1000,
        animationEasing: 'cubicOut',
        lineStyle: { width: 2 },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: `${errorColor}33` },
            { offset: 1, color: `${errorColor}00` },
          ]),
        },
      },
    ],
  }

  chartInstance.setOption(option)
}

onMounted(async () => {
  await nextTick()
  if (chartRef.value) {
    chartInstance = echarts.init(chartRef.value, null, { renderer: 'canvas', useDirtyRect: true })

    let isFirstRender = true
    // 使用 ResizeObserver 监听容器大小变化
    resizeObserver = new ResizeObserver(() => {
      window.requestAnimationFrame(() => {
        if (!chartInstance || !chartRef.value) return

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

// 监听主题变化
watch(isDark, () => {
  const timer = setTimeout(updateChart, 50)
  onWatcherCleanup(() => clearTimeout(timer))
})

// 监听数据变化
watch(
  resolvedStats,
  () => {
    updateChart()
  },
  { deep: 2 },
)

onUnmounted(() => {
  if (resizeObserver) {
    resizeObserver.disconnect()
    resizeObserver = null
  }
  if (chartInstance) {
    chartInstance.dispose()
    chartInstance = null
  }
})
</script>
