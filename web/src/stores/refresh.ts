import { computed, ref, watch } from 'vue'
import { defineStore } from 'pinia'
import { useDocumentVisibility, useIntervalFn, useLocalStorage, useNetwork } from '@vueuse/core'

type RefreshHandler = () => Promise<void>

interface RegisterOptions {
  immediate?: boolean
}

const MIN_INTERVAL_SECONDS = 5
const MAX_INTERVAL_SECONDS = 60
const DEFAULT_INTERVAL_SECONDS = 30

const normalizeIntervalSeconds = (value: number): number => {
  if (!Number.isFinite(value)) return DEFAULT_INTERVAL_SECONDS

  let nextValue = Math.round(value)

  // 兼容旧的毫秒存储格式
  if (nextValue >= MIN_INTERVAL_SECONDS * 1000) {
    nextValue = Math.round(nextValue / 1000)
  }

  if (nextValue < MIN_INTERVAL_SECONDS) return MIN_INTERVAL_SECONDS
  if (nextValue > MAX_INTERVAL_SECONDS) return MAX_INTERVAL_SECONDS
  return nextValue
}

export const useRefreshStore = defineStore('refresh', () => {
  // ==========================
  // 1. 配置状态 (持久化)
  // ==========================
  const isEnabled = useLocalStorage('settings.refresh_enabled', true)
  const intervalSeconds = useLocalStorage('settings.refresh_interval_seconds', DEFAULT_INTERVAL_SECONDS)

  // ==========================
  // 2. 运行时环境检测 (VueUse)
  // ==========================
  const { isOnline } = useNetwork()
  const visibility = useDocumentVisibility() // 'visible' | 'hidden'

  // 判断是否应该进行倒计时：开启 + 在线 + 页面可见
  const isActive = computed(
    () => isEnabled.value && isOnline.value && visibility.value === 'visible',
  )

  // ==========================
  // 3. 刷新注册表
  // ==========================
  // 存储所有需要刷新的回调函数
  const subscribers = ref(new Map<string, RefreshHandler>())
  const isRefreshing = ref(false)
  const lastRefreshedAt = ref<number | null>(null)

  // ==========================
  // 4. 核心执行逻辑
  // ==========================
  const executeHandler = async (key: string, handler: RefreshHandler) => {
    try {
      await handler()
    } catch (error) {
      console.error(`[AutoRefresh] ${key} failed:`, error)
    }
  }

  const triggerRefresh = async () => {
    if (isRefreshing.value || subscribers.value.size === 0) return

    isRefreshing.value = true
    try {
      // 并行执行所有注册的刷新函数
      await Promise.all(
        Array.from(subscribers.value.entries()).map(([key, handler]) =>
          executeHandler(key, handler),
        ),
      )
      lastRefreshedAt.value = Date.now()
    } catch (error) {
      console.error('[AutoRefresh] Error during refresh cycle:', error)
    } finally {
      isRefreshing.value = false
    }
  }

  // ==========================
  // 5. 定时器引擎
  // ==========================
  const intervalMs = computed(() => intervalSeconds.value * 1000)
  const hasSubscribers = computed(() => subscribers.value.size > 0)
  const shouldRun = computed(() => isActive.value && hasSubscribers.value)

  const { pause, resume } = useIntervalFn(triggerRefresh, intervalMs, {
    immediate: false,
    immediateCallback: false,
  })

  // 监听状态变化来控制定时器
  watch(
    [shouldRun, intervalMs],
    ([active]) => {
      if (!active) {
        pause()
        return
      }
      pause()
      resume()
    },
    { immediate: true },
  )

  // ==========================
  // 6. 暴露给 Composable 的接口
  // ==========================
  const setIntervalSeconds = (value: number) => {
    intervalSeconds.value = normalizeIntervalSeconds(value)
  }

  watch(
    intervalSeconds,
    (value) => {
      const normalized = normalizeIntervalSeconds(value)
      if (normalized !== value) {
        intervalSeconds.value = normalized
      }
    },
    { immediate: true },
  )

  function register(key: string, handler: RefreshHandler, options?: RegisterOptions) {
    if (!key.trim()) return
    subscribers.value.set(key, handler)
    if (options?.immediate) {
      void executeHandler(key, handler)
    }
    return () => unregister(key)
  }

  function unregister(key: string) {
    subscribers.value.delete(key)
  }

  return {
    isEnabled,
    intervalSeconds,
    minIntervalSeconds: MIN_INTERVAL_SECONDS,
    maxIntervalSeconds: MAX_INTERVAL_SECONDS,
    isActive,
    isRefreshing,
    lastRefreshedAt,
    triggerRefresh,
    register,
    unregister,
    setIntervalSeconds,
  }
})
