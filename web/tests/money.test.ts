import { describe, expect, it } from 'vitest'
import { applyVat, convertAmount, dualDisplay, FALLBACK_FX_RATES, formatMoney } from '../utils/money'

describe('convertAmount', () => {
  it('converts through USD-base rates', () => {
    expect(convertAmount(24, 'USD', 'KRW', FALLBACK_FX_RATES)).toBeCloseTo(24 * 1380, 0)
    expect(convertAmount(100, 'EUR', 'USD', FALLBACK_FX_RATES)).toBeCloseTo(100 / 0.92, 2)
  })

  it('returns the input when a rate is unknown', () => {
    expect(convertAmount(10, 'USD', 'XXX', FALLBACK_FX_RATES)).toBe(10)
  })
})

describe('applyVat', () => {
  it('applies the locale VAT rate', () => {
    expect(applyVat(100, 'ko')).toBeCloseTo(110, 6)
    expect(applyVat(100, 'de')).toBeCloseTo(119, 6)
    expect(applyVat(100, 'en')).toBe(100)
  })
})

describe('formatMoney', () => {
  it('formats per locale with zero decimals for KRW/JPY', () => {
    expect(formatMoney(33200, 'KRW', 'ko')).toContain('33,200')
    expect(formatMoney(24, 'USD', 'en')).toContain('24.00')
  })

  it('never throws on unknown currency', () => {
    expect(formatMoney(10, 'XXX', 'en')).toContain('10')
  })
})

describe('dualDisplay', () => {
  it('pairs converted target with native billing currency', () => {
    const d = dualDisplay(24, 'USD', 'KRW', FALLBACK_FX_RATES, 'en', false)
    expect(d.converted).toContain('33,120')
    expect(d.native).toContain('24.00')
  })
})
