import { z } from 'zod'

export const jsonValueSchema = z.json()

function isJsonEqual(left: unknown, right: unknown): boolean {
  if (Object.is(left, right)) {
    return true
  }

  try {
    return JSON.stringify(left) === JSON.stringify(right)
  } catch {
    return false
  }
}

function cloneFallbackValue<T>(value: T): T {
  if (typeof value === 'object' && value !== null) {
    if (typeof structuredClone === 'function') {
      return structuredClone(value)
    }
  }
  return value
}

export function parseWithFallback<T>(schema: z.ZodType<T>, value: unknown, fallback: T): T {
  const parsed = schema.safeParse(value)
  if (parsed.success) {
    return parsed.data
  }
  return cloneFallbackValue(fallback)
}

export function parseAndCheckChanged<T>(
  schema: z.ZodType<T>,
  value: unknown,
  fallback: T,
): { value: T; changed: boolean } {
  const parsed = schema.safeParse(value)
  if (parsed.success) {
    return {
      value: parsed.data,
      changed: !isJsonEqual(value, parsed.data),
    }
  }

  const fallbackValue = cloneFallbackValue(fallback)
  return {
    value: fallbackValue,
    changed: !isJsonEqual(value, fallbackValue),
  }
}
