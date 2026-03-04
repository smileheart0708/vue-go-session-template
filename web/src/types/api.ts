import type { HttpResponseSchema } from '@/utils/http'
import type { LogEntry } from '@/utils/logs'
import {
  ensureJsonValue,
  expectBooleanField,
  expectFiniteNumberField,
  expectIntegerField,
  expectObjectRecord,
  expectOptionalStringField,
  expectStringField,
  isObjectRecord,
} from './schema-helpers'

export interface LoginResponse {
  success: boolean
  message: string
  session_id?: string
}

export interface ValidateSessionResponse {
  valid: boolean
}

export interface LogoutResponse {
  success: boolean
  message: string
}

export interface DashboardStatsResponse {
  memory_used: number
  memory_total: number
  memory_percent: number
  start_time: number
}

export interface LogsHistoryResponse {
  logs: LogEntry[]
  count: number
}

function parseLogAttrs(value: unknown): Record<string, unknown> | undefined {
  if (value === undefined) {
    return undefined
  }

  if (!isObjectRecord(value)) {
    throw new Error('LogEntry.attrs must be an object when provided')
  }

  const attrs: Record<string, unknown> = {}
  for (const [key, rawValue] of Object.entries(value)) {
    attrs[key] = ensureJsonValue(rawValue, `LogEntry.attrs.${key}`)
  }

  return attrs
}

function parseLoginResponse(value: unknown): LoginResponse {
  const payload = expectObjectRecord(value, 'LoginResponse')
  const success = expectBooleanField(payload, 'success', 'LoginResponse')
  const message = expectStringField(payload, 'message', 'LoginResponse')
  const sessionId = expectOptionalStringField(payload, 'session_id', 'LoginResponse')

  if (success && !sessionId) {
    throw new Error('LoginResponse.session_id is required when success is true')
  }

  const response: LoginResponse = { success, message }
  if (sessionId !== undefined) {
    response.session_id = sessionId
  }
  return response
}

function parseValidateSessionResponse(value: unknown): ValidateSessionResponse {
  const payload = expectObjectRecord(value, 'ValidateSessionResponse')
  return { valid: expectBooleanField(payload, 'valid', 'ValidateSessionResponse') }
}

function parseLogoutResponse(value: unknown): LogoutResponse {
  const payload = expectObjectRecord(value, 'LogoutResponse')
  return {
    success: expectBooleanField(payload, 'success', 'LogoutResponse'),
    message: expectStringField(payload, 'message', 'LogoutResponse'),
  }
}

function parseDashboardStatsResponse(value: unknown): DashboardStatsResponse {
  const payload = expectObjectRecord(value, 'DashboardStatsResponse')
  const memoryUsed = expectFiniteNumberField(payload, 'memory_used', 'DashboardStatsResponse')
  const memoryTotal = expectFiniteNumberField(payload, 'memory_total', 'DashboardStatsResponse')
  const memoryPercent = expectFiniteNumberField(payload, 'memory_percent', 'DashboardStatsResponse')
  const startTime = expectIntegerField(payload, 'start_time', 'DashboardStatsResponse')

  if (memoryUsed < 0) {
    throw new Error('DashboardStatsResponse.memory_used must be >= 0')
  }
  if (memoryTotal <= 0) {
    throw new Error('DashboardStatsResponse.memory_total must be > 0')
  }
  if (memoryUsed > memoryTotal) {
    throw new Error('DashboardStatsResponse.memory_used must be <= memory_total')
  }
  if (memoryPercent < 0 || memoryPercent > 100) {
    throw new Error('DashboardStatsResponse.memory_percent must be between 0 and 100')
  }
  if (startTime <= 0) {
    throw new Error('DashboardStatsResponse.start_time must be a positive integer')
  }

  return {
    memory_used: memoryUsed,
    memory_total: memoryTotal,
    memory_percent: memoryPercent,
    start_time: startTime,
  }
}

export function parseLogEntry(value: unknown): LogEntry {
  const payload = expectObjectRecord(value, 'LogEntry')
  const time = expectStringField(payload, 'time', 'LogEntry').trim()
  const level = expectStringField(payload, 'level', 'LogEntry').trim()
  const msg = expectStringField(payload, 'msg', 'LogEntry')
  const attrs = parseLogAttrs(payload['attrs'])

  if (!time) {
    throw new Error('LogEntry.time must not be empty')
  }
  if (!level) {
    throw new Error('LogEntry.level must not be empty')
  }

  const entry: LogEntry = { time, level, msg }
  if (attrs && Object.keys(attrs).length > 0) {
    entry.attrs = attrs
  }
  return entry
}

function parseLogsHistoryResponse(value: unknown): LogsHistoryResponse {
  const payload = expectObjectRecord(value, 'LogsHistoryResponse')
  const rawLogs = payload['logs']
  if (!Array.isArray(rawLogs)) {
    throw new Error('LogsHistoryResponse.logs must be an array')
  }

  const logs = rawLogs.map((item, index) => {
    try {
      return parseLogEntry(item)
    } catch (error) {
      const detail = error instanceof Error ? error.message : 'invalid log entry'
      throw new Error(`LogsHistoryResponse.logs[${index}] invalid: ${detail}`)
    }
  })
  const count = expectIntegerField(payload, 'count', 'LogsHistoryResponse')
  if (count < 0) {
    throw new Error('LogsHistoryResponse.count must be >= 0')
  }
  if (count !== logs.length) {
    throw new Error('LogsHistoryResponse.count must equal logs length')
  }

  return { logs, count }
}

export const loginResponseSchema: HttpResponseSchema<LoginResponse> = {
  name: 'LoginResponse',
  parse: parseLoginResponse,
}

export const validateSessionResponseSchema: HttpResponseSchema<ValidateSessionResponse> = {
  name: 'ValidateSessionResponse',
  parse: parseValidateSessionResponse,
}

export const logoutResponseSchema: HttpResponseSchema<LogoutResponse> = {
  name: 'LogoutResponse',
  parse: parseLogoutResponse,
}

export const dashboardStatsResponseSchema: HttpResponseSchema<DashboardStatsResponse> = {
  name: 'DashboardStatsResponse',
  parse: parseDashboardStatsResponse,
}

export const logEntrySchema: HttpResponseSchema<LogEntry> = {
  name: 'LogEntry',
  parse: parseLogEntry,
}

export const logsHistoryResponseSchema: HttpResponseSchema<LogsHistoryResponse> = {
  name: 'LogsHistoryResponse',
  parse: parseLogsHistoryResponse,
}
