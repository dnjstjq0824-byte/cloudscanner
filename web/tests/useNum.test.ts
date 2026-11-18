import { describe, expect, it } from 'vitest'
import { money, num } from '../composables/useNum'

// Regression: the Go API serializes decimal money/score fields as JSON
// strings (shopspring decimal default). Templates must never call .toFixed()
// on raw wire values — num()/money() coerce both representations.
describe('useNum wire-decimal coercion', () => {
  it('num() accepts decimal strings and numbers', () => {
    expect(num('24')).toBe(24)
    expect(num('0.0329')).toBeCloseTo(0.0329, 4)
    expect(num(24)).toBe(24)
    expect(num('10.50')).toBeCloseTo(10.5, 2)
  })

  it('num() surfaces unusable values as NaN, never throws', () => {
    expect(num('not-a-number')).toBeNaN()
    expect(num(null)).toBeNaN()
    expect(num(undefined)).toBeNaN()
  })

  it('money() formats both representations with fixed digits', () => {
    expect(money('24')).toBe('24.00')
    expect(money(24)).toBe('24.00')
    // String and number inputs render identically for the same value.
    expect(money('0.03285', 4)).toBe(money(0.03285, 4))
    expect(money('100.50')).toBe('100.50')
  })

  it('money() renders a placeholder for unusable values', () => {
    expect(money('nope')).toBe('—')
    expect(money(null)).toBe('—')
    expect(money(undefined)).toBe('—')
  })

  it('string decimals sort numerically after coercion (calculator best-value)', () => {
    const wire = ['100.00', '24.00', '9.50']
    const sorted = [...wire].sort((a, b) => num(a) - num(b))
    expect(sorted).toEqual(['9.50', '24.00', '100.00'])
  })
})
