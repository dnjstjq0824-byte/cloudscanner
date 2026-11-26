import { defineStore } from 'pinia'
import { FALLBACK_FX_RATES } from '~/utils/money'

interface FxRateRow {
  base_currency: string
  target_currency: string
  rate: number
  effective_at: string
}

/** Display currency + tax toggle + FX table (live API with static fallback). */
export const useMarketStore = defineStore('market', {
  state: () => ({
    currency: (typeof useCookie === 'function' ? useCookie('tirp-currency').value : null) || 'USD',
    taxIncluded: (typeof useCookie === 'function' ? useCookie('tirp-tax').value : null) === 'incl',
    rates: { ...FALLBACK_FX_RATES } as Record<string, number>,
    ratesAt: '',
    ratesLive: false,
  }),

  actions: {
    setCurrency(code: string) {
      this.currency = code
      useCookie('tirp-currency').value = code
    },

    setTaxIncluded(on: boolean) {
      this.taxIncluded = on
      useCookie('tirp-tax').value = on ? 'incl' : 'excl'
    },

    async fetchFx() {
      try {
        const { request } = useApi()
        const res = await request<{ data: FxRateRow[] }>('/api/v1/pricing/fx-rates')
        const table: Record<string, number> = { USD: 1 }
        let latest = ''
        for (const row of res.data ?? []) {
          if (row.base_currency === 'USD' && row.rate > 0) {
            table[row.target_currency] = row.rate
            if (row.effective_at > latest) {
              latest = row.effective_at
            }
          }
        }
        if (Object.keys(table).length > 1) {
          this.rates = table
          this.ratesAt = latest
          this.ratesLive = true
        }
      } catch {
        // Keep indicative fallback rates; rendering must never break (G §10).
      }
    },
  },
})
