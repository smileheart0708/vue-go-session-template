import ky, { HTTPError, type Options } from 'ky'
import { ZodError, z, type ZodType } from 'zod'

const DEFAULT_API_BASE_URL = '/api'
const DEFAULT_TIMEOUT_MS = 30_000

type UnauthorizedHandler = () => void | Promise<void>

let unauthorizedHandler: UnauthorizedHandler | null = null
let unauthorizedTask: Promise<void> | null = null
let redirectingToLogin = false

const API_BASE_URL = resolveApiBaseUrl(import.meta.env.VITE_API_BASE_URL)

const contextSchema = z
  .object({
    skipUnauthorizedHandler: z.boolean().optional(),
  })
  .catchall(z.unknown())

const redirectQuerySchema = z.union([z.string(), z.array(z.string())])

function isAbsoluteUrl(value: string): boolean {
  return /^[a-z][a-z\d+\-.]*:\/\//i.test(value)
}

function resolveApiBaseUrl(rawValue: string | undefined): string {
  const raw = rawValue?.trim()
  if (!raw) {
    return DEFAULT_API_BASE_URL
  }

  if (isAbsoluteUrl(raw)) {
    const url = new URL(raw)
    const normalizedPath =
      url.pathname === '/' ? DEFAULT_API_BASE_URL : url.pathname.replace(/\/+$/, '')
    url.pathname = normalizedPath
    return url.toString().replace(/\/+$/, '')
  }

  const normalized = raw.startsWith('/') ? raw : `/${raw}`
  return normalized.replace(/\/+$/, '') || DEFAULT_API_BASE_URL
}

function shouldSkipUnauthorizedHandler(options: { context?: unknown }): boolean {
  const parsedContext = contextSchema.safeParse(options.context)
  if (!parsedContext.success) {
    return false
  }

  return parsedContext.data.skipUnauthorizedHandler === true
}

function ensureRedirectToLogin(currentPath: string = getCurrentPath()): void {
  if (typeof window === 'undefined' || redirectingToLogin) {
    return
  }

  const targetPath = buildLoginRedirectPath(currentPath)
  if (targetPath === getCurrentPath()) {
    return
  }

  redirectingToLogin = true
  window.location.replace(targetPath)
}

async function runUnauthorizedHandler(): Promise<void> {
  if (!unauthorizedHandler) {
    return
  }

  if (!unauthorizedTask) {
    unauthorizedTask = Promise.resolve(unauthorizedHandler()).finally(() => {
      unauthorizedTask = null
    })
  }

  await unauthorizedTask
}

function formatIssuePath(path: ReadonlyArray<PropertyKey>): string {
  if (path.length === 0) {
    return 'root'
  }

  return path
    .map((segment, index) => {
      if (typeof segment === 'number') {
        return `[${segment}]`
      }

      const segmentText = String(segment)
      if (index === 0) {
        return segmentText
      }
      return `.${segmentText}`
    })
    .join('')
}

function summarizeZodError(error: ZodError): string {
  const issueMessages = error.issues.slice(0, 3).map((issue) => {
    const path = formatIssuePath(issue.path)
    return `${path}: ${issue.message}`
  })

  if (error.issues.length > 3) {
    issueMessages.push(`...and ${error.issues.length - 3} more issue(s)`)
  }

  return issueMessages.join('; ')
}

function resolveSchemaName(schema: ZodType<unknown>): string {
  const description = schema.description?.trim()
  if (description) {
    return description
  }
  return 'AnonymousSchema'
}

export const api = ky.create({
  prefixUrl: API_BASE_URL,
  credentials: 'include',
  retry: 0,
  timeout: DEFAULT_TIMEOUT_MS,
  hooks: {
    afterResponse: [
      async (_request, options, response) => {
        if (shouldSkipUnauthorizedHandler(options) || response.status !== 401) {
          return response
        }

        await runUnauthorizedHandler()
        ensureRedirectToLogin()
        return response
      },
    ],
  },
})

export { HTTPError }

export class ApiResponseValidationError extends Error {
  readonly requestUrl: string
  readonly schemaName: string
  readonly data: unknown

  constructor(requestUrl: string, schemaName: string, detail: string, data: unknown) {
    super(`Invalid response payload for schema "${schemaName}" at "${requestUrl}": ${detail}`)
    this.name = 'ApiResponseValidationError'
    this.requestUrl = requestUrl
    this.schemaName = schemaName
    this.data = data
  }
}

export function parseWithSchema<T>(
  value: unknown,
  schema: ZodType<T>,
  requestUrl: string,
): T {
  try {
    return schema.parse(value)
  } catch (error) {
    if (error instanceof ZodError) {
      throw new ApiResponseValidationError(
        requestUrl,
        resolveSchemaName(schema),
        summarizeZodError(error),
        value,
      )
    }

    const detail = error instanceof Error ? error.message : 'Unknown schema parse error'
    throw new ApiResponseValidationError(requestUrl, resolveSchemaName(schema), detail, value)
  }
}

export async function readHttpErrorData(error: HTTPError): Promise<unknown | null> {
  const response = error.response.clone()
  if (response.status === 204 || response.status === 205) {
    return null
  }

  const contentType = response.headers.get('content-type')?.toLowerCase() ?? ''
  if (contentType.includes('application/json')) {
    try {
      return await response.json()
    } catch {
      return null
    }
  }

  try {
    const text = await response.text()
    return text || null
  } catch {
    return null
  }
}

export function setUnauthorizedHandler(handler: UnauthorizedHandler | null): void {
  unauthorizedHandler = handler
}

export function withUnauthorizedHandlerSkipped(options: Options = {}): Options {
  const parsedContext = contextSchema.safeParse(options.context)
  const context = parsedContext.success ? parsedContext.data : {}

  return {
    ...options,
    context: {
      ...context,
      skipUnauthorizedHandler: true,
    },
  }
}

export function normalizeApiEndpoint(endpoint: string): string {
  const trimmedEndpoint = endpoint.trim()
  if (!trimmedEndpoint) {
    throw new Error('API endpoint is required')
  }

  if (trimmedEndpoint.startsWith('/api/')) {
    return trimmedEndpoint.slice('/api/'.length)
  }
  if (trimmedEndpoint.startsWith('api/')) {
    return trimmedEndpoint.slice('api/'.length)
  }
  if (trimmedEndpoint.startsWith('/')) {
    return trimmedEndpoint.slice(1)
  }
  return trimmedEndpoint
}

export function getCurrentPath(): string {
  if (typeof window === 'undefined') {
    return '/'
  }
  return `${window.location.pathname}${window.location.search}${window.location.hash}`
}

export function resolveRedirectPath(value: unknown): string | null {
  const parsedQueryValue = redirectQuerySchema.safeParse(value)
  if (!parsedQueryValue.success) {
    return null
  }

  const rawValue = Array.isArray(parsedQueryValue.data)
    ? parsedQueryValue.data[0]
    : parsedQueryValue.data
  if (!rawValue) {
    return null
  }

  const trimmed = rawValue.trim()
  if (!trimmed.startsWith('/') || trimmed.startsWith('//') || trimmed.startsWith('/login')) {
    return null
  }

  return trimmed
}

export function buildLoginRedirectPath(currentPath: string = getCurrentPath()): string {
  const redirectPath = resolveRedirectPath(currentPath)
  if (!redirectPath) {
    return '/login'
  }
  return `/login?redirect=${encodeURIComponent(redirectPath)}`
}
