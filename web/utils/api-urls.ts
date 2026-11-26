/** Pure URL/header builders for the Go backend client (unit-tested, no Nuxt deps). */

export const DEFAULT_DEV_API_BASE = 'http://localhost:8080'

/** Absolute API origin for relative /api/* paths when calling direct (non-dev-proxy). */
export function resolveApiUrl(apiBase: string, path: string): string {
  if (path.startsWith('/api/') && apiBase !== DEFAULT_DEV_API_BASE) {
    return `${apiBase}${path}`
  }
  return path
}

/** Outbound deal link (/go/* is served by the API origin). */
export function buildDealUrl(apiBase: string, slug: string, src: string): string {
  return `${apiBase}/go/${slug}?src=${src}`
}

export interface SessionCookies {
  session?: string | null
  org?: string | null
}

/** Bearer + org headers from session cookies; anonymous callers get none. */
export function buildSessionHeaders(
  cookies: SessionCookies,
  opts: { anonymous?: boolean; extra?: Record<string, string> } = {},
): Record<string, string> {
  const headers: Record<string, string> = { ...(opts.extra ?? {}) }
  if (!opts.anonymous) {
    if (cookies.session) {
      headers.Authorization = `Bearer ${cookies.session}`
    }
    if (cookies.org) {
      headers['X-Organization-Id'] = cookies.org
    }
  }
  return headers
}
