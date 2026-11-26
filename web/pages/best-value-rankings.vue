<template>
  <div class="max-w-4xl mx-auto space-y-8 py-4">
    <div class="space-y-2">
      <h1 class="text-3xl font-extrabold text-ink tracking-tight">{{ $t('rankings.h1') }}</h1>
      <p class="text-sm text-ink-secondary">{{ $t('rankings.lede') }}</p>
    </div>

    <LoadingSkeleton v-if="loading" :rows="5" :label="$t('rankings.loading')" />
    <ErrorState v-else-if="error" :title="error.title" :message="error.message" @retry="load" />

    <ol v-else class="space-y-3">
      <li v-for="(item, i) in items" :key="item.id" class="card p-5 flex items-center gap-4">
        <span class="text-2xl font-extrabold text-brand tabular-nums w-10 shrink-0">{{ i + 1 }}</span>
        <div class="min-w-0 flex-1">
          <NuxtLink :to="$localePath(`/products/${item.id}`)" class="font-bold text-ink hover:text-brand">
            {{ item.name }}
          </NuxtLink>
          <div class="text-xs text-ink-secondary mt-0.5">{{ item.provider_name }} · {{ item.category }}</div>
        </div>
        <div class="text-right shrink-0">
          <div class="text-sm font-bold text-brand">
            <PriceDisplay :amount="num(item.raw_monthly_price_usd)" currency="USD" cycle="month" />
          </div>
          <div class="text-[11px] text-ink-secondary tabular-nums">{{ $t('rankings.score', { v: score(item) }) }}</div>
        </div>
      </li>
    </ol>
    <EmptyState
      v-if="!loading && !error && items.length === 0"
      :title="$t('rankings.emptyTitle')"
      :message="$t('rankings.emptyBody')"
    />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ApiError, useApi } from '~/composables/useApi'
import { num } from '~/composables/useNum'
import type { ListEnvelope, SearchResultItem } from '~/types/api'

definePageMeta({
  layout: 'public',
})

const { t } = useI18n()
useHead(() => ({
  title: t('rankings.title'),
  meta: [{ name: 'description', content: t('rankings.meta') }],
}))

const items = ref<SearchResultItem[]>([])
const loading = ref(false)
const error = ref<ApiError | null>(null)
const { request } = useApi()

function score(item: SearchResultItem): string {
  const v = num(item.best_value_score)
  return Number.isNaN(v) ? '—' : v.toFixed(1)
}

async function load() {
  loading.value = true
  error.value = null
  try {
    const res = await request<ListEnvelope<SearchResultItem>>('/api/v1/rankings/best-value', { query: { limit: 20 } })
    items.value = res.data ?? []
  } catch (e) {
    error.value = e instanceof ApiError ? e : ApiError.fromUnknown(e)
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>
