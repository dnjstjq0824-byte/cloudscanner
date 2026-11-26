import type { ProblemDetails } from '~/types/api'
import { buildDealUrl, buildSessionHeaders, DEFAULT_DEV_API_BASE, resolveApiUrl } from '~/utils/api-urls'

/** Normalized API error carrying problem details when present. */
export class ApiError extends Error {
  status: number
  title: string
  traceId?: string

  constructor(status: number, title: string, detail: string, traceId?: string) {
    super(detail || title)
    this.name = 'ApiError'
    this.status = status
    this.title = title
    this.traceId = traceId
  }

  static fromUnknown(err: unknown): ApiError {
    const anyErr = err as { response?: { status?: number; _data?: ProblemDetails }; message?: string }
    const data = anyErr?.response?._data
    if (data && typeof data.status === 'number') {
      return new ApiError(data.status, data.title || 'Request Failed', data.detail || data.title || 'Request failed', data.trace_id)
    }
    const status = anyErr?.response?.status ?? 0
    return new ApiError(status, 'Network Error', anyErr?.message || 'The API did not respond.')
  }
}

export interface ApiOptions {
  method?: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE'
  query?: Record<string, string | number | undefined>
  body?: unknown
  /** Skip session headers (login/register/public bootstrap). */
  anonymous?: boolean
  /** Override the auto-generated idempotency key for mutating requests. */
  idempotencyKey?: string
  headers?: Record<string, string>
}

function newIdempotencyKey(): string {
  try {
    return crypto.randomUUID()
  } catch {
    return `key-${Date.now().toString(36)}-${Math.random().toString(36).slice(2, 18)}`
  }
}

/**
 * API client for the Go backend.
 *
 * Development (default apiBase) issues same-origin calls through the Nitro
 * `/api/**` proxy. Any explicit NUXT_PUBLIC_API_BASE (e.g. production
 * https://api.tirpscanr.com) issues direct absolute calls with Bearer auth;
 * the API authorizes those origins via exact-match CORS.
 */
export function useApi() {
  const apiBase = (useRuntimeConfig().public.apiBase as string) || DEFAULT_DEV_API_BASE

  function resolveUrl(path: string): string {
    return resolveApiUrl(apiBase, path)
  }

  async function request<T>(path: string, opts: ApiOptions = {}): Promise<T> {
    const cleanQuery: Record<string, string | number> = {}
    for (const [k, v] of Object.entries(opts.query ?? {})) {
      if (v !== undefined && v !== '') cleanQuery[k] = v
    }
    const method = opts.method ?? 'GET'
    const headers = buildSessionHeaders(
      { session: useCookie('cs_session').value, org: useCookie('cs_org').value },
      { anonymous: opts.anonymous, extra: opts.headers },
    )
    if (method !== 'GET') {
      headers['Idempotency-Key'] = opts.idempotencyKey || newIdempotencyKey()
    }
    try {
      return await $fetch<T>(resolveUrl(path), {
        method,
        query: cleanQuery,
        body: opts.body as Record<string, unknown> | undefined,
        headers,
      })
    } catch (err) {
      throw ApiError.fromUnknown(err)
    }
  }

  /** Absolute API origin for outbound deal links (/go/* is served by the API). */
  const baseURL = apiBase

  function dealUrl(slug: string, src: string): string {
    return buildDealUrl(apiBase, slug, src)
  }

  return { request, baseURL, dealUrl }
}
