import { describe, expect, it } from 'vitest'
import { buildDealUrl, buildSessionHeaders, resolveApiUrl } from '../utils/api-urls'

describe('resolveApiUrl', () => {
  it('keeps /api/* relative on the dev proxy base', () => {
    expect(resolveApiUrl('http://localhost:8080', '/api/v1/search/skus')).toBe('/api/v1/search/skus')
  })

  it('absolutizes /api/* on an explicit production base', () => {
    expect(resolveApiUrl('https://api.tirpscanr.com', '/api/v1/search/skus')).toBe(
      'https://api.tirpscanr.com/api/v1/search/skus',
    )
  })

  it('never touches non-API paths', () => {
    expect(resolveApiUrl('https://api.tirpscanr.com', '/dashboard')).toBe('/dashboard')
  })
})

describe('buildDealUrl', () => {
  it('builds outbound deal links with src attribution', () => {
    expect(buildDealUrl('https://api.tirpscanr.com', 'notion-pro', 'rankings')).toBe(
      'https://api.tirpscanr.com/go/notion-pro?src=rankings',
    )
  })
})

describe('buildSessionHeaders', () => {
  it('attaches Bearer + org headers from cookies', () => {
    expect(buildSessionHeaders({ session: 'tok', org: 'org-1' })).toEqual({
      Authorization: 'Bearer tok',
      'X-Organization-Id': 'org-1',
    })
  })

  it('omits missing cookies and respects anonymous mode', () => {
    expect(buildSessionHeaders({})).toEqual({})
    expect(buildSessionHeaders({ session: 'tok' }, { anonymous: true })).toEqual({})
  })

  it('preserves caller-supplied extra headers', () => {
    expect(buildSessionHeaders({ session: 'tok' }, { extra: { 'Idempotency-Key': 'k' } })).toEqual({
      'Idempotency-Key': 'k',
      Authorization: 'Bearer tok',
    })
  })
})
