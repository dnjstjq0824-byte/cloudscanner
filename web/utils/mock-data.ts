/** Demo fixtures for pages whose APIs are PLANNED (see openapi.yaml).
 * Each consumer is marked with common.demoData until wired to live endpoints. */

export interface MockPromotion {
  id: string
  vendor: string
  title: string
  code: string
  discountPct: number
  expiresAt: string
}

export interface MockArticle {
  slug: string
  title: string
  excerpt: string
  kind: 'tutorial' | 'review'
  rating: number
  updatedAt: string
}

export interface MockTrend {
  category: string
  avgDropPct: number
  dealCount: number
  spark: number[]
}

export interface MockPriceAlert {
  id: string
  label: string
  targetPrice: number
  currency: string
  active: boolean
}

export interface MockSecurityLog {
  id: string
  at: string
  event: string
  ip: string
  status: 'OK' | 'BLOCKED'
}

export interface MockUserNotice {
  id: string
  title: string
  body: string
  at: string
  read: boolean
}

export interface MockJob {
  id: string
  queue: string
  state: 'queued' | 'running' | 'dead'
  attempts: number
  updatedAt: string
}

export function mockPromotions(): MockPromotion[] {
  return [
    { id: 'promo-1', vendor: 'Notion', title: 'Plus annual −20% for new teams', code: 'TEAM20', discountPct: 20, expiresAt: '2026-12-31' },
    { id: 'promo-2', vendor: 'Figma', title: 'Professional 3 months free', code: 'DESIGN3', discountPct: 25, expiresAt: '2026-11-30' },
    { id: 'promo-3', vendor: 'Linear', title: 'Startup plan free for 12 months', code: 'STARTUP', discountPct: 100, expiresAt: '2027-01-15' },
    { id: 'promo-4', vendor: 'Vercel', title: 'Pro trial credit $20', code: 'SHIP20', discountPct: 15, expiresAt: '2026-12-15' },
  ]
}

export function mockArticles(kind: 'tutorial' | 'review'): MockArticle[] {
  const all: MockArticle[] = [
    { slug: 'compare-crm-plans', title: 'How to compare CRM plans without a spreadsheet', excerpt: 'Normalize per-seat prices, trials, and renewal jumps before you buy.', kind: 'tutorial', rating: 5, updatedAt: '2026-11-20' },
    { slug: 'video-editor-deals', title: 'AI video editor deals, ranked by value', excerpt: 'Seven editors scored on price-per-seat, trial generosity, and renewal fairness.', kind: 'review', rating: 4, updatedAt: '2026-11-18' },
    { slug: 'annual-vs-monthly', title: 'Annual vs monthly billing: the real math', excerpt: 'When the annual discount beats monthly flexibility, with seat-count tables.', kind: 'tutorial', rating: 5, updatedAt: '2026-11-10' },
    { slug: 'pm-tools-2026', title: 'Project tools 2026: five plans compared', excerpt: 'Notion, Asana, Linear, ClickUp, and Monday measured side by side.', kind: 'review', rating: 4, updatedAt: '2026-11-05' },
  ]
  return all.filter((a) => a.kind === kind)
}

export function mockTrends(): MockTrend[] {
  return [
    { category: 'AI video', avgDropPct: 18, dealCount: 12, spark: [4, 6, 5, 9, 12, 10, 12] },
    { category: 'CRM', avgDropPct: 11, dealCount: 23, spark: [20, 18, 22, 19, 23, 21, 23] },
    { category: 'DevOps', avgDropPct: 7, dealCount: 9, spark: [6, 7, 6, 8, 7, 9, 9] },
    { category: 'Design', avgDropPct: 22, dealCount: 7, spark: [2, 3, 5, 4, 6, 5, 7] },
  ]
}

export function mockPriceAlerts(): MockPriceAlert[] {
  return [
    { id: 'al-1', label: 'Notion Plus annual', targetPrice: 96, currency: 'USD', active: true },
    { id: 'al-2', label: 'Figma Professional monthly', targetPrice: 12, currency: 'USD', active: true },
    { id: 'al-3', label: 'Linear Business', targetPrice: 8, currency: 'USD', active: false },
  ]
}

export function mockSecurityLogs(): MockSecurityLog[] {
  return [
    { id: 'sl-1', at: '2026-11-25 09:14', event: 'login.success', ip: '***.***.12.34', status: 'OK' },
    { id: 'sl-2', at: '2026-11-24 22:03', event: 'login.failed', ip: '***.***.56.78', status: 'BLOCKED' },
    { id: 'sl-3', at: '2026-11-24 22:02', event: 'login.failed', ip: '***.***.56.78', status: 'BLOCKED' },
    { id: 'sl-4', at: '2026-11-20 14:55', event: 'token.refresh', ip: '***.***.12.34', status: 'OK' },
  ]
}

export function mockUserNotices(): MockUserNotice[] {
  return [
    { id: 'un-1', title: 'Price drop: Notion Plus annual', body: 'Now $96/yr, below your $100 target.', at: '2026-11-25', read: false },
    { id: 'un-2', title: 'New coupon: DESIGN3', body: 'Figma Professional 3 months free for new teams.', at: '2026-11-22', read: false },
    { id: 'un-3', title: 'Trial ending soon', body: 'Your Vercel Pro trial ends in 3 days.', at: '2026-11-18', read: true },
  ]
}

export function mockJobs(): MockJob[] {
  return [
    { id: 'job-8812', queue: 'deal-ingest', state: 'running', attempts: 1, updatedAt: '2026-11-26 13:40' },
    { id: 'job-8809', queue: 'deal-ingest', state: 'queued', attempts: 0, updatedAt: '2026-11-26 13:39' },
    { id: 'job-8790', queue: 'fx-sync', state: 'dead', attempts: 5, updatedAt: '2026-11-26 11:02' },
  ]
}

export interface MockRow {
  id: string
  name: string
  status: string
  updated: string
}

/** Generic small-table fixture for admin consoles awaiting live endpoints. */
export function mockRows(prefix: string, names: string[], statuses: string[]): MockRow[] {
  return names.map((name, i) => ({
    id: `${prefix}-${i + 1}`,
    name,
    status: statuses[i % statuses.length],
    updated: '2026-11-26',
  }))
}
