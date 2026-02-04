import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import type { LogEntry } from '@/utils/logs'

interface UseLogStreamOptions {
  streamUrl?: string
  maxLogs?: number
  onLog?: (log: LogEntry) => void
}

export function useLogStream(options: UseLogStreamOptions = {}) {
  const router = useRouter()
  const authStore = useAuthStore()
  const logs = ref<LogEntry[]>([])
  const streamUrl = options.streamUrl ?? '/api/logs/stream'
  const maxLogs = options.maxLogs ?? 500
  let eventSource: EventSource | null = null
  let reconnectTimer: ReturnType<typeof setTimeout> | null = null

  function connect() {
    if (eventSource) return
    eventSource = new EventSource(streamUrl)

    eventSource.onmessage = (event) => {
      try {
        const logEntry: LogEntry = JSON.parse(event.data)
        logs.value.push(logEntry)

        if (logs.value.length > maxLogs) {
          logs.value.shift()
        }

        options.onLog?.(logEntry)
      } catch (error) {
        console.error('解析日志数据失败:', error)
      }
    }

    eventSource.onerror = (error) => {
      console.error('SSE 连接错误:', error)

      if (eventSource?.readyState === EventSource.CLOSED) {
        authStore.validateSession().then((isValid) => {
          if (!isValid) {
            disconnect()
            router.replace('/login')
            return
          }

          disconnect()
          scheduleReconnect()
        })
      }
    }
  }

  function scheduleReconnect() {
    if (reconnectTimer) return
    reconnectTimer = setTimeout(() => {
      reconnectTimer = null
      connect()
    }, 5000)
  }

  function disconnect() {
    if (reconnectTimer) {
      clearTimeout(reconnectTimer)
      reconnectTimer = null
    }
    eventSource?.close()
    eventSource = null
  }

  function clearLogs() {
    logs.value = []
  }

  return { logs, connect, disconnect, clearLogs }
}
