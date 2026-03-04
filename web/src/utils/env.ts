import { z } from 'zod'
import { parseWithFallback } from '@/types/zod'

const apiModeSchema = z.enum(['real', 'mock'])
export type ApiMode = z.infer<typeof apiModeSchema>

const booleanFlagSchema = z.preprocess(
  (value) => {
    if (typeof value !== 'string') {
      return value
    }
    return value.trim().toLowerCase()
  },
  z.union([z.literal('true'), z.literal('false')]).transform((value) => value === 'true'),
)

export const isMockAuthEnabled = parseWithFallback(
  booleanFlagSchema,
  import.meta.env.VITE_MOCK_AUTH,
  false,
)
export const apiMode = parseWithFallback<ApiMode>(apiModeSchema, import.meta.env.VITE_API_MODE, 'real')
export const isMockApiEnabled = apiMode === 'mock'
