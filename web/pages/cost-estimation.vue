<template>
  <div class="max-w-3xl mx-auto space-y-8 py-4">
    <div class="space-y-2">
      <h1 class="text-3xl font-extrabold text-ink tracking-tight">{{ $t('estimation.h1') }}</h1>
      <p class="text-sm text-ink-secondary">{{ $t('estimation.lede') }}</p>
    </div>

    <form class="card p-6 grid grid-cols-1 sm:grid-cols-2 gap-4" @submit.prevent="run">
      <div class="space-y-1.5 sm:col-span-2">
        <label for="est-sku" class="label">{{ $t('estimation.plan') }}</label>
        <select id="est-sku" v-model="skuId" class="input" required>
          <option v-for="s in skus" :key="s.id" :value="s.id">{{ s.name }} — {{ s.provider_name }}</option>
        </select>
      </div>
      <div class="space-y-1.5">
        <label for="est-currency" class="label">{{ $t('estimation.currency') }}</label>
        <select id="est-currency" v-model="currency" class="input">
          <option>USD</option>
          <option>EUR</option>
          <option>KRW</option>
        </select>
      </div>
      <div class="space-y-1.5">
        <label for="est-tax" class="label">{{ $t('estimation.taxPct') }}</label>
        <input id="est-tax" v-model.number="taxPct" type="number" min="0" max="50" step="0.5" class="input" />
      </div>
      <div class="sm:col-span-2">
        <button type="submit" class="btn-primary" :disabled="running || !skuId">
          {{ running ? $t('common.working') : $t('estimation.submit') }}
        </button>
      </div>
    </form>

    <LoadingSkeleton v-if="loadingSkus" :rows="2" :label="$t('estimation.loading')" />
    <ErrorState v-else-if="error" :title="error.title" :message="error.message" @retry="run" />

    <div v-if="result" class="card p-6 space-y-2">
      <div class="flex items-baseline justify-between text-sm">
        <span class="text-ink-secondary">{{ $t('estimation.subtotal') }}</span>
        <strong class="tabular-nums"><PriceDisplay :amount="num(result.subtotal_target)" :currency="result.currency" /></strong>
      </div>
      <div class="flex items-baseline justify-between text-sm">
        <span class="text-ink-secondary">{{ $t('estimation.tax') }}</span>
        <strong class="tabular-nums"><PriceDisplay :amount="num(result.tax_target)" :currency="result.currency" /></strong>
      </div>
      <div class="flex items-baseline justify-between border-t border-subtle pt-2">
        <span class="font-bold text-ink">{{ $t('estimation.total') }}</span>
        <span class="text-xl font-extrabold text-brand tabular-nums">
          <PriceDisplay :amount="num(result.total_target)" :currency="result.currency" />
        </span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ApiError, useApi } from '~/composables/useApi'
import { num } from '~/composables/useNum'
import type { ListEnvelope, SKUDetail } from '~/types/api'

definePageMeta({
  layout: 'public',
})

const { t } = useI18n()
useHead(() => ({
  title: t('estimation.title'),
  meta: [{ name: 'description', content: t('estimation.meta') }],
}))

interface TCOResult {
  subtotal_target: number | string
  tax_target: number | string
  total_target: number | string
  currency: string
}

const skus = ref<SKUDetail[]>([])
const skuId = ref('')
const currency = ref('USD')
const taxPct = ref(0)
const result = ref<TCOResult | null>(null)
const loadingSkus = ref(false)
const running = ref(false)
const error = ref<ApiError | null>(null)
const { request } = useApi()

async function run() {
  if (!skuId.value) {
    return
  }
  running.value = true
  error.value = null
  try {
    const res = await request<{ data: TCOResult }>('/api/v1/pricing/estimate', {
      method: 'POST',
      body: {
        sku_id: skuId.value,
        scenario: { target_currency: currency.value, tax_percentage: taxPct.value },
      },
    })
    result.value = res.data
  } catch (e) {
    error.value = e instanceof ApiError ? e : ApiError.fromUnknown(e)
  } finally {
    running.value = false
  }
}

onMounted(async () => {
  loadingSkus.value = true
  try {
    const res = await request<ListEnvelope<SKUDetail>>('/api/v1/catalog/skus')
    skus.value = res.data ?? []
    if (skus.value.length > 0) {
      skuId.value = skus.value[0].id
    }
  } catch (e) {
    error.value = e instanceof ApiError ? e : ApiError.fromUnknown(e)
  } finally {
    loadingSkus.value = false
  }
})
</script>
