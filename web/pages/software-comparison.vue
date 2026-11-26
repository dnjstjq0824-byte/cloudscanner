<template>
  <div class="max-w-5xl mx-auto space-y-8 py-4">
    <div class="space-y-2">
      <h1 class="text-3xl font-extrabold text-ink tracking-tight">{{ $t('software.h1') }}</h1>
      <p class="text-sm text-ink-secondary">{{ $t('software.lede') }}</p>
    </div>

    <LoadingSkeleton v-if="loading" :rows="4" :label="$t('software.loading')" />
    <ErrorState v-else-if="error" :title="error.title" :message="error.message" @retry="load" />
    <EmptyState
      v-else-if="items.length === 0"
      :title="$t('software.emptyTitle')"
      :message="$t('software.emptyBody')"
    />

    <div v-else class="grid gap-4 md:grid-cols-2">
      <div v-for="item in items" :key="item.id" class="card p-5">
        <div class="flex items-start justify-between gap-2">
          <NuxtLink :to="$localePath(`/products/${item.id}`)" class="font-bold text-ink hover:text-brand">
            {{ item.name }}
          </NuxtLink>
          <StatusBadge :status="item.freshness_status" />
        </div>
        <div class="text-xs text-ink-secondary mt-1">{{ item.provider_name }} · {{ item.category }}</div>
        <div class="text-sm font-bold text-brand mt-2">
          <PriceDisplay :amount="num(item.raw_monthly_price_usd)" currency="USD" cycle="month" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ApiError, useApi } from '~/composables/useApi'
import { num } from '~/composables/useNum'
import type { SearchResponse, SearchResultItem } from '~/types/api'

definePageMeta({
  layout: 'public',
})

const { t } = useI18n()
useHead(() => ({
  title: t('software.title'),
  meta: [{ name: 'description', content: t('software.meta') }],
}))

const items = ref<SearchResultItem[]>([])
const loading = ref(false)
const error = ref<ApiError | null>(null)
const { request } = useApi()

async function load() {
  loading.value = true
  error.value = null
  try {
    const res = await request<SearchResponse>('/api/v1/search/skus', { query: { category: 'SOFTWARE', limit: 20 } })
    items.value = res.data ?? []
  } catch (e) {
    error.value = e instanceof ApiError ? e : ApiError.fromUnknown(e)
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>
