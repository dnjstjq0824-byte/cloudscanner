import { applyVat, convertAmount, formatMoney } from '~/utils/money'

/** Locale-aware money formatting bound to the market store. */
export function useMoney() {
  const { locale, t } = useI18n()
  const market = useMarketStore()

  function money(amount: number, currency?: string): string {
    const target = currency || market.currency
    const value = market.taxIncluded ? applyVat(amount, locale.value) : amount
    return formatMoney(value, target, locale.value)
  }

  /** Dual display: converted target + native billing currency with FX tooltip. */
  function dual(amount: number, from: string, cycle?: 'month' | 'year'): {
    converted: string
    native: string
    tip: string
  } {
    const to = market.currency
    const convertedAmount = market.taxIncluded
      ? applyVat(convertAmount(amount, from, to, market.rates), locale.value)
      : convertAmount(amount, from, to, market.rates)
    const suffix = cycle === 'month' ? t('money.perMonth') : cycle === 'year' ? t('money.perYear') : ''
    const rate = market.rates[to] && market.rates[from]
      ? (market.rates[to] / market.rates[from]).toPrecision(4)
      : '—'
    return {
      converted: formatMoney(convertedAmount, to, locale.value) + suffix,
      native: `${t('money.billed')} ${formatMoney(amount, from, locale.value)}${suffix}`,
      tip: t('selector.fxTip', { rate, time: market.ratesAt || '—' }),
    }
  }

  return { money, dual }
}
