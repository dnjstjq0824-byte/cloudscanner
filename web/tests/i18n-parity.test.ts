import { describe, expect, it } from 'vitest'
import { readdirSync, readFileSync } from 'node:fs'
import { join } from 'node:path'

function leafKeys(obj: unknown, prefix = ''): string[] {
  const out: string[] = []
  for (const [k, v] of Object.entries(obj as Record<string, unknown>)) {
    if (v !== null && typeof v === 'object' && !Array.isArray(v)) {
      out.push(...leafKeys(v, `${prefix}${k}.`))
    } else {
      out.push(`${prefix}${k}`)
    }
  }
  return out
}

const dir = join(__dirname, '..', 'i18n', 'locales')
const files = readdirSync(dir).filter((f) => f.endsWith('.json')).sort()

describe('i18n key parity (G §12)', () => {
  it('ships exactly the five spec locales', () => {
    expect(files).toEqual(['de.json', 'en.json', 'es.json', 'ja.json', 'ko.json'])
  })

  it('keeps 100% key parity across locales', () => {
    const [base, ...rest] = files.map((f) => ({
      file: f,
      keys: new Set(leafKeys(JSON.parse(readFileSync(join(dir, f), 'utf-8')))),
    }))
    expect(base.keys.size).toBeGreaterThan(0)
    for (const other of rest) {
      const missing = [...base.keys].filter((k) => !other.keys.has(k))
      const extra = [...other.keys].filter((k) => !base.keys.has(k))
      expect({ file: other.file, missing, extra }).toEqual({ file: other.file, missing: [], extra: [] })
    }
  })

  it('keeps interpolation params consistent across locales', () => {
    const paramOf = (s: string) => [...s.matchAll(/\{(\w+)\}/g)].map((m) => m[1]).sort().join(',')
    const bundles = files.map((f) => JSON.parse(readFileSync(join(dir, f), 'utf-8')) as Record<string, unknown>)
    const walk = (node: unknown, path: string, cb: (path: string, value: string) => void): void => {
      for (const [k, v] of Object.entries(node as Record<string, unknown>)) {
        if (v !== null && typeof v === 'object' && !Array.isArray(v)) walk(v, `${path}${k}.`, cb)
        else if (typeof v === 'string') cb(`${path}${k}`, v)
      }
    }
    const byKey = new Map<string, string[]>()
    for (const b of bundles) {
      walk(b, '', (path, value) => {
        byKey.set(path, [...(byKey.get(path) ?? []), paramOf(value)])
      })
    }
    for (const [path, sigs] of byKey) {
      expect({ path, sigs }).toEqual({ path, sigs: sigs.map(() => sigs[0]) })
    }
  })
})
