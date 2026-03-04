import { delay, http, HttpResponse, sse } from 'msw'
import type {
  DashboardStatsResponse,
  LoginResponse,
  LogsHistoryResponse,
  SessionStatusResponse,
} from '@/types/api'
import type { LogEntry } from '@/utils/logs'

interface LoginRequestBody {
  auth_key: string
}

const MOCK_MEMORY_TOTAL = 16 * 1024 * 1024 * 1024
const MOCK_SSE_INTERVAL_MS = 1200
const MOCK_UPTIME_DAYS = 14
const MOCK_START_TIME = Math.floor(Date.now() / 1000) - MOCK_UPTIME_DAYS * 24 * 60 * 60

function parseBooleanEnv(value: unknown, fallback: boolean): boolean {
  if (typeof value !== 'string') {
    return fallback
  }

  const normalized = value.trim().toLowerCase()
  if (normalized === 'true') {
    return true
  }
  if (normalized === 'false') {
    return false
  }
  return fallback
}

const isMockAuthEnabledByEnv = parseBooleanEnv(import.meta.env.VITE_MOCK_AUTH, false)
let isMockAuthenticated = isMockAuthEnabledByEnv

const LOG_LEVELS = ['INFO', 'WARN', 'ERROR', 'DEBUG'] as const
type LogLevel = (typeof LOG_LEVELS)[number]

const LOG_MESSAGES: Record<LogLevel, ReadonlyArray<string>> = {
  INFO: ['请求处理完成', '上游响应成功', '定时任务执行中', '配置热更新生效'],
  WARN: ['请求耗时偏高', '上游重试触发', '缓存命中率下降'],
  ERROR: ['上游请求失败', '会话校验失败', '流式连接中断'],
  DEBUG: ['路由命中: /api/chat', '开始刷新 dashboard 指标', '下发 SSE 消息'],
}

function pickRandom<T>(items: ReadonlyArray<T>): T {
  const index = Math.floor(Math.random() * items.length)
  const item = items[index]
  if (item === undefined) {
    throw new Error('items must not be empty')
  }
  return item
}

function randomBetween(min: number, max: number): number {
  return min + Math.random() * (max - min)
}

function pad(value: number): string {
  return String(value).padStart(2, '0')
}

function formatLogTime(date: Date): string {
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}`
}

function createMockLogEntry(
  forcedLevel?: LogLevel,
  forcedMessage?: string,
  date: Date = new Date(),
): LogEntry {
  const level = forcedLevel ?? pickRandom(LOG_LEVELS)
  const message = forcedMessage ?? pickRandom(LOG_MESSAGES[level])

  const entry: LogEntry = { time: formatLogTime(date), level, msg: message }
  if (Math.random() > 0.5) {
    entry.attrs = { trace_id: Math.random().toString(16).slice(2, 10) }
  }
  return entry
}

function createDashboardStats(): DashboardStatsResponse {
  const memoryPercent = Number(randomBetween(28, 73).toFixed(1))
  const memoryUsed = Math.floor((memoryPercent / 100) * MOCK_MEMORY_TOTAL)

  return {
    memory_used: memoryUsed,
    memory_total: MOCK_MEMORY_TOTAL,
    memory_percent: memoryPercent,
    start_time: MOCK_START_TIME,
  }
}

async function parseJsonBody<T>(request: Request): Promise<T | null> {
  try {
    const payload = await request.json()
    return payload
  } catch {
    return null
  }
}

function buildLoginSuccessResponse(): LoginResponse {
  return { success: true, message: '登录成功（Mock）' }
}

function buildLoginInvalidResponse(message: string): LoginResponse {
  return { success: false, message }
}

function buildSessionUnauthorizedResponse(): SessionStatusResponse {
  return { authenticated: false, message: '未授权，请先登录' }
}

function buildSessionAuthorizedResponse(): SessionStatusResponse {
  return { authenticated: true }
}

function requireAuthenticated() {
  if (isMockAuthenticated) {
    return null
  }

  return HttpResponse.json({ error: '未授权，请先登录' }, { status: 401 })
}

export const handlers = [
  http.post('/api/login', async ({ request }) => {
    await delay(120)

    const body = await parseJsonBody<Partial<LoginRequestBody>>(request)
    const authKey = body?.auth_key?.trim()

    if (!authKey) {
      return HttpResponse.json(buildLoginInvalidResponse('认证令牌不能为空'), { status: 400 })
    }

    isMockAuthenticated = true
    return HttpResponse.json(buildLoginSuccessResponse())
  }),

  http.get('/api/session', async () => {
    await delay(60)

    if (!isMockAuthenticated) {
      return HttpResponse.json(buildSessionUnauthorizedResponse(), { status: 401 })
    }

    return HttpResponse.json(buildSessionAuthorizedResponse())
  }),

  http.post('/api/logout', async () => {
    await delay(80)
    isMockAuthenticated = false
    return HttpResponse.json({ success: true, message: '登出成功（Mock）' })
  }),

  http.get('/api/dashboard/stats', async () => {
    await delay(100)

    const unauthorizedResponse = requireAuthenticated()
    if (unauthorizedResponse) {
      return unauthorizedResponse
    }

    return HttpResponse.json(createDashboardStats())
  }),

  http.get('/api/logs/history', async () => {
    await delay(100)

    const unauthorizedResponse = requireAuthenticated()
    if (unauthorizedResponse) {
      return unauthorizedResponse
    }

    const logs: LogEntry[] = Array.from({ length: 60 }, (_, index) => {
      const offsetMs = (60 - index) * 1500
      return createMockLogEntry(undefined, undefined, new Date(Date.now() - offsetMs))
    })

    const response: LogsHistoryResponse = { logs, count: logs.length }
    return HttpResponse.json(response)
  }),

  sse('/api/logs/stream', ({ client, request }) => {
    if (!isMockAuthenticated) {
      client.close()
      return
    }

    client.send({ retry: 5000 })
    client.send({ data: createMockLogEntry('INFO', 'Mock SSE 连接已建立') })

    const timer = globalThis.setInterval(() => {
      client.send({ data: createMockLogEntry() })
    }, MOCK_SSE_INTERVAL_MS)

    const cleanup = () => {
      globalThis.clearInterval(timer)
      client.close()
    }

    request.signal.addEventListener('abort', cleanup, { once: true })
  }),
]
