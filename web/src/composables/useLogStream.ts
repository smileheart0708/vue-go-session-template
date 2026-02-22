import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { logEntrySchema, logsHistoryResponseSchema } from '@/types/api'
import { useAuthStore } from '@/stores/auth'
import { HttpError, HttpResponseValidationError, buildLoginRedirectPath, http } from '@/utils'
import { isMockApiEnabled } from '@/utils/env'
import type { LogEntry } from '@/utils/logs'

export type LogStreamStatus = 'connecting' | 'connected' | 'disconnected'

interface UseLogStreamOptions {
  streamUrl?: string
  historyEndpoint?: string
  maxLogs?: number
  onLog?: (log: LogEntry) => void
}

function createLogEntryKey(entry: LogEntry): string {
  const attrs = entry.attrs ? JSON.stringify(entry.attrs) : ''
  return `${entry.time}|${entry.level}|${entry.msg}|${attrs}`
}

export function useLogStream(options: UseLogStreamOptions = {}) {
  const router = useRouter()
  const authStore = useAuthStore()
  const logs = ref<LogEntry[]>([])
  const status = ref<LogStreamStatus>('connecting')
  const streamUrl = options.streamUrl ?? '/api/logs/stream?history=0'
  const historyEndpoint = options.historyEndpoint ?? '/logs/history'
  const maxLogs = options.maxLogs ?? 500
  let eventSource: EventSource | null = null
  let reconnectTimer: ReturnType<typeof setTimeout> | null = null
  let historyLoaded = false

  function normalizeLogs(entries: LogEntry[]): LogEntry[] {
    const seen = new Set<string>()
    const deduped: LogEntry[] = []

    for (const entry of entries) {
      const key = createLogEntryKey(entry)
      if (seen.has(key)) {
        continue
      }
      seen.add(key)
      deduped.push(entry)
    }

    if (deduped.length <= maxLogs) {
      return deduped
    }
    return deduped.slice(-maxLogs)
  }

  function appendLog(logEntry: LogEntry): void {
    logs.value = normalizeLogs([...logs.value, logEntry])
    options.onLog?.(logEntry)
  }

  function mergeHistory(historyLogs: LogEntry[]): void {
    logs.value = normalizeLogs([...historyLogs, ...logs.value])
  }

  async function loadHistory(): Promise<void> {
    if (historyLoaded) {
      return
    }

    try {
      const data = await http(historyEndpoint, {
        schema: logsHistoryResponseSchema,
      })
      mergeHistory(data.logs)
      historyLoaded = true
    } catch (error) {
      if (error instanceof HttpResponseValidationError) {
        console.error('历史日志响应格式异常:', error)
        return
      }
      if (error instanceof HttpError && error.status === 401) {
        return
      }
      console.error('加载历史日志失败:', error)
    }
  }

  function connect(): void {
    if (eventSource) return
    status.value = 'connecting'
    void loadHistory()
    eventSource = new EventSource(streamUrl)

    eventSource.onopen = () => {
      status.value = 'connected'
    }

    eventSource.onmessage = (event) => {
      try {
        const rawData: unknown = JSON.parse(event.data)
        const logEntry = logEntrySchema.parse(rawData)
        appendLog(logEntry)
      } catch (error) {
        console.error('解析日志数据失败:', error)
      }
    }

    eventSource.onerror = (error) => {
      console.error('SSE 连接错误:', error)

      if (eventSource?.readyState === EventSource.CLOSED) {
        status.value = 'disconnected'
        if (isMockApiEnabled) {
          disconnect()
          scheduleReconnect()
          return
        }

        authStore.validateSession().then((isValid) => {
          if (!isValid) {
            disconnect()
            router.replace(buildLoginRedirectPath(router.currentRoute.value.fullPath))
            return
          }

          disconnect()
          scheduleReconnect()
        })
      } else {
        status.value = 'connecting'
      }
    }
  }

  function scheduleReconnect() {
    if (reconnectTimer) return
    status.value = 'connecting'
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
    status.value = 'disconnected'
  }

  function clearLogs() {
    logs.value = []
  }

  return { logs, status, connect, disconnect, clearLogs }
}
