<template>
  <div class="max-w-5xl mx-auto space-y-8 py-4">
    <div class="space-y-3">
      <h1 class="text-3xl font-extrabold text-ink tracking-tight">{{ $t('search.h1') }}</h1>
      <p class="text-sm text-ink-secondary">{{ $t('search.lede') }}</p>
      <form class="flex flex-col sm:flex-row gap-2" @submit.prevent="run">
        <input
          v-model="q"
          type="search"
          required
          minlength="2"
          :placeholder="$t('search.placeholder')"
          class="input flex-1"
          :aria-label="$t('search.h1')"
        />
        <button type="submit" class="btn-primary" :disabled="loading">
          {{ loading ? $t('common.working') : $t('search.submit') }}
        </button>
      </form>
    </div>

    <LoadingSkeleton v-if="loading" :rows="4" :label="$t('search.loading')" />
    <ErrorState v-else-if="error" :title="error.title" :message="error.message" @retry="run" />

    <div v-else-if="searched" class="space-y-4">
      <p class="text-xs text-ink-secondary">{{ $t('search.results', { n: items.length, ms: ms }) }}</p>
      <EmptyState
        v-if="items.length === 0"
        :title="$t('search.emptyTitle')"
        :message="$t('search.emptyBody')"
      />
      <div class="grid gap-4 md:grid-cols-2">
        <NuxtLink
          v-for="item in items"
          :key="item.id"
          :to="$localePath(`/products/${item.id}`)"
          class="card p-5 hover:border-brand/60 transition block"
        >
          <div class="flex items-start justify-between gap-2">
            <div class="font-bold text-ink">{{ item.name }}</div>
            <StatusBadge :status="item.freshness_status" />
          </div>
          <div class="text-xs text-ink-secondary mt-1">{{ item.provider_name }} · {{ item.category }}</div>
          <div class="text-sm font-bold text-brand mt-2">
            <PriceDisplay :amount="num(item.raw_monthly_price_usd)" currency="USD" cycle="month" />
          </div>
        </NuxtLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ApiError, useApi } from '~/composables/useApi'
import { num } from '~/composables/useNum'
import type { SearchResponse, SearchResultItem } from '~/types/api'

definePageMeta({
  layout: 'public',
})

const { t } = useI18n()
useHead(() => ({
  title: t('search.title'),
  meta: [{ name: 'description', content: t('search.meta') }],
}))

const route = useRoute()
const q = ref(typeof route.query.q === 'string' ? route.query.q : '')
const items = ref<SearchResultItem[]>([])
const ms = ref(0)
const loading = ref(false)
const searched = ref(false)
const error = ref<ApiError | null>(null)
const { request } = useApi()

async function run() {
  if (q.value.trim().length < 2) {
    return
  }
  loading.value = true
  error.value = null
  try {
    const res = await request<SearchResponse>('/api/v1/search/skus', { query: { q: q.value.trim(), limit: 20 } })
    items.value = res.data ?? []
    ms.value = res.execution_ms ?? 0
    searched.value = true
  } catch (e) {
    error.value = e instanceof ApiError ? e : ApiError.fromUnknown(e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  if (q.value.trim().length >= 2) {
    run()
  }
})
</script>
