const DEFAULT_API_BASE_URL = '/api'

type QueryValue = string | number | boolean | null | undefined

export type HttpQuery = Record<string, QueryValue | QueryValue[]>
export type HttpResponseType = 'json' | 'text' | 'blob' | 'arrayBuffer' | 'response'

export interface HttpOptions
  extends Omit<RequestInit, 'body' | 'headers' | 'method' | 'credentials'> {
  method?: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE' | 'HEAD' | 'OPTIONS'
  headers?: HeadersInit
  body?: BodyInit | Record<string, unknown> | unknown[] | null
  query?: HttpQuery
  responseType?: HttpResponseType
  credentials?: RequestCredentials
  skipUnauthorizedHandler?: boolean
}

interface HttpContext {
  options: HttpOptions
  requestUrl: string
}

export type ResponseInterceptor = (
  response: Response,
  context: HttpContext,
) => Response | Promise<Response>

type UnauthorizedHandler = () => void | Promise<void>

const isAbsoluteUrl = (value: string): boolean => /^[a-z][a-z\d+\-.]*:\/\//i.test(value)

const API_BASE_URL = resolveApiBaseUrl(import.meta.env.VITE_API_BASE_URL)
const responseInterceptors: ResponseInterceptor[] = []

let unauthorizedHandler: UnauthorizedHandler | null = null
let unauthorizedTask: Promise<void> | null = null
let redirectingToLogin = false

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

function mergePath(basePath: string, endpoint: string): string {
  const normalizedEndpoint = endpoint.startsWith('/') ? endpoint : `/${endpoint}`
  if (normalizedEndpoint.startsWith('/api/')) {
    return normalizedEndpoint
  }

  const normalizedBase = basePath.replace(/\/+$/, '')
  return `${normalizedBase}${normalizedEndpoint}`
}

function appendQuery(url: string, query: HttpQuery): string {
  const [withoutHash = '', hashPart] = url.split('#', 2)
  const [pathname = '', existingQuery] = withoutHash.split('?', 2)
  const params = new URLSearchParams(existingQuery ?? '')

  for (const [key, rawValue] of Object.entries(query)) {
    params.delete(key)
    if (rawValue === undefined || rawValue === null) {
      continue
    }

    const values = Array.isArray(rawValue) ? rawValue : [rawValue]
    for (const value of values) {
      if (value === undefined || value === null) {
        continue
      }
      params.append(key, String(value))
    }
  }

  const queryString = params.toString()
  const nextUrl = queryString ? `${pathname}?${queryString}` : pathname
  return hashPart ? `${nextUrl}#${hashPart}` : nextUrl
}

function buildRequestUrl(endpoint: string, query?: HttpQuery): string {
  const trimmedEndpoint = endpoint.trim()
  if (!trimmedEndpoint) {
    throw new Error('HTTP endpoint is required')
  }

  let requestUrl = trimmedEndpoint
  if (!isAbsoluteUrl(requestUrl)) {
    if (isAbsoluteUrl(API_BASE_URL)) {
      const baseUrl = new URL(API_BASE_URL)
      baseUrl.pathname = mergePath(baseUrl.pathname, requestUrl)
      requestUrl = baseUrl.toString()
    } else {
      requestUrl = mergePath(API_BASE_URL, requestUrl)
    }
  }

  if (!query) {
    return requestUrl
  }

  return appendQuery(requestUrl, query)
}

function isJsonBody(body: unknown): body is Record<string, unknown> | unknown[] {
  if (body === null || body === undefined) return false
  if (typeof body !== 'object') return false
  if (body instanceof FormData) return false
  if (body instanceof URLSearchParams) return false
  if (body instanceof Blob) return false
  if (body instanceof ArrayBuffer) return false
  if (ArrayBuffer.isView(body)) return false
  if (body instanceof ReadableStream) return false
  return true
}

function createRequestInit(options: HttpOptions): RequestInit {
  const {
    method = 'GET',
    headers,
    body,
    credentials = 'include',
    query: _query,
    responseType: _responseType,
    skipUnauthorizedHandler: _skipUnauthorizedHandler,
    ...rest
  } = options

  const requestHeaders = new Headers(headers ?? {})
  const requestInit: RequestInit = {
    ...rest,
    method,
    headers: requestHeaders,
    credentials,
  }

  if (body === undefined || body === null) {
    return requestInit
  }

  if (isJsonBody(body)) {
    if (!requestHeaders.has('Content-Type')) {
      requestHeaders.set('Content-Type', 'application/json')
    }
    requestInit.body = JSON.stringify(body)
    return requestInit
  }

  requestInit.body = body
  return requestInit
}

async function parseErrorData(response: Response): Promise<unknown> {
  if (response.status === 204 || response.status === 205) {
    return null
  }

  const contentType = response.headers.get('content-type')?.toLowerCase() ?? ''
  if (contentType.includes('application/json')) {
    try {
      return (await response.json()) as unknown
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

async function parseSuccessData<T>(response: Response, responseType: HttpResponseType): Promise<T> {
  if (responseType === 'response') {
    return response as T
  }

  if (response.status === 204 || response.status === 205) {
    return undefined as T
  }

  switch (responseType) {
    case 'text':
      return (await response.text()) as T
    case 'blob':
      return (await response.blob()) as T
    case 'arrayBuffer':
      return (await response.arrayBuffer()) as T
    case 'json': {
      const text = await response.text()
      if (!text) {
        return undefined as T
      }
      return JSON.parse(text) as T
    }
    default:
      return undefined as T
  }
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

responseInterceptors.push(async (response, context) => {
  if (context.options.skipUnauthorizedHandler || response.status !== 401) {
    return response
  }

  await runUnauthorizedHandler()
  ensureRedirectToLogin()
  return response
})

export class HttpError<T = unknown> extends Error {
  readonly status: number
  readonly response: Response
  readonly data: T | null

  constructor(response: Response, data: T | null) {
    super(response.statusText || `HTTP ${response.status}`)
    this.name = 'HttpError'
    this.status = response.status
    this.response = response
    this.data = data
  }
}

export function addResponseInterceptor(interceptor: ResponseInterceptor): () => void {
  responseInterceptors.push(interceptor)
  return () => {
    const index = responseInterceptors.indexOf(interceptor)
    if (index > -1) {
      responseInterceptors.splice(index, 1)
    }
  }
}

export function setUnauthorizedHandler(handler: UnauthorizedHandler | null): void {
  unauthorizedHandler = handler
}

export function getCurrentPath(): string {
  if (typeof window === 'undefined') {
    return '/'
  }
  return `${window.location.pathname}${window.location.search}${window.location.hash}`
}

export function resolveRedirectPath(value: unknown): string | null {
  const rawValue = Array.isArray(value) ? value[0] : value
  if (typeof rawValue !== 'string') {
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

export async function http<T = unknown>(endpoint: string, options: HttpOptions = {}): Promise<T> {
  const responseType = options.responseType ?? 'json'
  const requestUrl = buildRequestUrl(endpoint, options.query)
  const requestInit = createRequestInit(options)

  let response = await fetch(requestUrl, requestInit)
  const context: HttpContext = {
    options,
    requestUrl,
  }

  for (const interceptor of responseInterceptors) {
    response = await interceptor(response, context)
  }

  if (!response.ok) {
    const data = await parseErrorData(response)
    throw new HttpError(response, data)
  }

  return parseSuccessData<T>(response, responseType)
}
