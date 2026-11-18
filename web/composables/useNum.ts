/** Wire-decimal coercion. The Go API serializes decimal money/score fields as
 * JSON strings; templates must route every such value through these helpers
 * instead of calling .toFixed() on raw payloads. Auto-imported by Nuxt. */

export function num(v: number | string | null | undefined): number {
  if (typeof v === 'number') {
    return v
  }
  if (typeof v === 'string' && v.trim() !== '') {
    return Number(v)
  }
  return NaN
}

export function money(v: number | string | null | undefined, digits = 2): string {
  const n = num(v)
  if (!Number.isFinite(n)) {
    return '—'
  }
  return n.toFixed(digits)
}
