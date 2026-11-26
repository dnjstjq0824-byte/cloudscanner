<template>
  <div class="max-w-5xl mx-auto space-y-8 py-4">
    <div class="space-y-2">
      <h1 class="text-3xl font-extrabold text-ink tracking-tight">{{ $t('compare.h1') }}</h1>
      <p class="text-sm text-ink-secondary">{{ $t('compare.lede') }}</p>
    </div>

    <LoadingSkeleton v-if="loading" :rows="3" :label="$t('compare.loading')" />
    <ErrorState v-else-if="error" :title="error.title" :message="error.message" @retry="load" />
    <EmptyState
      v-else-if="rows.length === 0"
      :title="$t('compare.emptyTitle')"
      :message="$t('compare.emptyBody')"
      :cta-to="$localePath('/search')"
      :cta-label="$t('compare.emptyCta')"
    />

    <div v-else class="overflow-x-auto -mx-4 px-4">
      <table class="table-shell min-w-[640px]">
        <thead>
          <tr>
            <th class="py-3 px-4 text-left">{{ $t('compare.colPlan') }}</th>
            <th class="py-3 px-4 text-left">{{ $t('compare.colVendor') }}</th>
            <th class="py-3 px-4 text-right">{{ $t('compare.colMonthly') }}</th>
            <th class="py-3 px-4 text-left">{{ $t('compare.colFreshness') }}</th>
            <th class="py-3 px-4"></th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="sku in rows" :key="sku.id" class="border-t border-subtle">
            <td class="py-3 px-4 font-semibold text-ink">{{ sku.name }}</td>
            <td class="py-3 px-4 text-ink-secondary">{{ sku.provider_name }}</td>
            <td class="py-3 px-4 text-right font-bold text-brand tabular-nums">
              <PriceDisplay :amount="num(sku.raw_monthly_price_usd)" currency="USD" cycle="month" />
            </td>
            <td class="py-3 px-4"><StatusBadge :status="sku.freshness_status" /></td>
            <td class="py-3 px-4 text-right">
              <NuxtLink :to="$localePath(`/products/${sku.id}`)" class="text-xs font-semibold text-brand hover:underline">
                {{ $t('compare.view') }}
              </NuxtLink>
            </td>
          </tr>
        </tbody>
      </table>
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
  title: t('compare.title'),
  meta: [{ name: 'description', content: t('compare.meta') }],
}))

const route = useRoute()
const ids = computed(() => {
  const q = route.query.ids
  if (typeof q === 'string' && q.length > 0) {
    return q.split(',').slice(0, 6)
  }
  return [] as string[]
})
const rows = ref<SKUDetail[]>([])
const loading = ref(false)
const error = ref<ApiError | null>(null)
const { request } = useApi()

async function load() {
  if (ids.value.length === 0) {
    return
  }
  loading.value = true
  error.value = null
  try {
    const res = await request<ListEnvelope<SKUDetail>>('/api/v1/catalog/skus')
    const wanted = new Set(ids.value)
    rows.value = (res.data ?? []).filter((s) => wanted.has(s.id)).slice(0, 6)
  } catch (e) {
    error.value = e instanceof ApiError ? e : ApiError.fromUnknown(e)
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>
