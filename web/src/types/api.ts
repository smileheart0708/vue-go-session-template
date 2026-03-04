import { z } from 'zod'
import { jsonValueSchema } from './zod'

const nonEmptyTrimmedStringSchema = z.string().trim().min(1)

export const loginResponseSchema = z
  .object({ success: z.boolean(), message: z.string() })
  .describe('LoginResponse')

export type LoginResponse = z.infer<typeof loginResponseSchema>

export const sessionStatusResponseSchema = z
  .object({ authenticated: z.boolean(), message: z.string().optional() })
  .describe('SessionStatusResponse')

export type SessionStatusResponse = z.infer<typeof sessionStatusResponseSchema>

export const logoutResponseSchema = z
  .object({ success: z.boolean(), message: z.string() })
  .describe('LogoutResponse')

export type LogoutResponse = z.infer<typeof logoutResponseSchema>

export const dashboardStatsResponseSchema = z
  .object({
    memory_used: z.number().nonnegative(),
    memory_total: z.number().positive(),
    memory_percent: z.number().nonnegative().lte(100),
    start_time: z.int().positive(),
  })
  .refine((value) => value.memory_used <= value.memory_total, {
    path: ['memory_used'],
    message: 'memory_used must be <= memory_total',
  })
  .describe('DashboardStatsResponse')

export type DashboardStatsResponse = z.infer<typeof dashboardStatsResponseSchema>

const logEntryBaseSchema = z.object({
  time: nonEmptyTrimmedStringSchema,
  level: nonEmptyTrimmedStringSchema,
  msg: z.string(),
  attrs: z.record(z.string(), jsonValueSchema).optional(),
})

export const logEntrySchema = logEntryBaseSchema
  .transform((value) => {
    if (value.attrs && Object.keys(value.attrs).length === 0) {
      return { time: value.time, level: value.level, msg: value.msg }
    }
    return value
  })
  .describe('LogEntry')

export type LogEntry = z.infer<typeof logEntrySchema>

export const logsHistoryResponseSchema = z
  .object({ logs: z.array(logEntrySchema), count: z.int().nonnegative() })
  .refine((value) => value.count === value.logs.length, {
    path: ['count'],
    message: 'count must equal logs length',
  })
  .describe('LogsHistoryResponse')

export type LogsHistoryResponse = z.infer<typeof logsHistoryResponseSchema>
