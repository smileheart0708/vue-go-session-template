export type ApiMode = 'real' | 'mock'

function normalizeFlag(value: string | undefined): boolean {
  return value?.trim().toLowerCase() === 'true'
}

function resolveApiMode(value: string | undefined): ApiMode {
  const normalized = value?.trim().toLowerCase()
  if (normalized === 'mock') return 'mock'
  return 'real'
}

export const isMockAuthEnabled = normalizeFlag(import.meta.env.VITE_MOCK_AUTH)
export const apiMode = resolveApiMode(import.meta.env.VITE_API_MODE)
export const isMockApiEnabled = apiMode === 'mock'

