/** Pure money math (unit-tested). Rates are per-USD multipliers (USD base). */

/** Indicative fallback rates, used only when the FX API is unreachable. */
export const FALLBACK_FX_RATES: Record<string, number> = {
  USD: 1,
  EUR: 0.92,
  KRW: 1380,
  JPY: 149,
  GBP: 0.79,
}

/** Static VAT/GST rates by locale for the tax toggle (display only). */
export const VAT_RATES: Record<string, number> = {
  en: 0,
  ko: 0.1,
  ja: 0.1,
  de: 0.19,
  es: 0.21,
}

export function convertAmount(amount: number, from: string, to: string, rates: Record<string, number>): number {
  const fromRate = rates[from]
  const toRate = rates[to]
  if (!fromRate || !toRate) {
    return amount
  }
  return (amount / fromRate) * toRate
}

export function applyVat(amount: number, locale: string): number {
  return amount * (1 + (VAT_RATES[locale] ?? 0))
}

export function formatMoney(amount: number, currency: string, locale: string): string {
  try {
    return new Intl.NumberFormat(locale, {
      style: 'currency',
      currency,
      maximumFractionDigits: currency === 'KRW' || currency === 'JPY' ? 0 : 2,
    }).format(amount)
  } catch {
    return `${amount} ${currency}`
  }
}

/** Dual-currency display pair: converted target + native billing currency. */
export function dualDisplay(
  amount: number,
  from: string,
  to: string,
  rates: Record<string, number>,
  locale: string,
  withTax: boolean,
): { converted: string; native: string } {
  const convertedAmount = withTax
    ? applyVat(convertAmount(amount, from, to, rates), locale)
    : convertAmount(amount, from, to, rates)
  return {
    converted: formatMoney(convertedAmount, to, locale),
    native: formatMoney(amount, from, locale),
  }
}
