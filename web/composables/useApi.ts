import type { ProblemDetails } from '~/types/api'

/** Normalized API error carrying RFC 9457 problem details when present. */
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
    return new ApiError(status, 'Network Error', anyErr?.message || 'The API did not respond. Check that the backend is running.')
  }
}

export interface ApiOptions {
  method?: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE'
  query?: Record<string, string | number | undefined>
  body?: unknown
  /** Skip attaching the session bearer [REDACTED] */
  anonymous?: boolean
}

/**
 * Authenticated API client. Attaches the session token, normalizes
 * RFC 9457 failures into ApiError, and clears the session on 401.
 */
export function useApi() {
  const config = useRuntimeConfig()
  const baseURL = (config.public.apiBase as string) || 'http://localhost:8080'

  async function request<T>(path: string, opts: ApiOptions = {}): Promise<T> {
    const headers: Record<string, string> = { 'Content-Type': 'application/json' }
    if (!opts.anonymous) {
      const token = useCookie('cs_session').value
      if (token) {
        headers.Authorization = `Bearer ${token}`
      }
      const orgId = useCookie('cs_org').value
      if (orgId) {
        headers['X-Organization-Id'] = orgId
      }
    }
    const cleanQuery: Record<string, string | number> = {}
    for (const [k, v] of Object.entries(opts.query ?? {})) {
      if (v !== undefined && v !== '') {
        cleanQuery[k] = v
      }
    }
    try {
      return await $fetch<T>(path, {
        baseURL,
        method: opts.method ?? 'GET',
        query: cleanQuery,
        body: opts.body as Record<string, unknown> | undefined,
        headers,
      })
    } catch (err) {
      const apiErr = ApiError.fromUnknown(err)
      if (apiErr.status === 401 && !opts.anonymous && import.meta.client) {
        useCookie('cs_session').value = null
      }
      throw apiErr
    }
  }

  return { request, baseURL }
}
