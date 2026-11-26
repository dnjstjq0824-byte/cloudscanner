<template>
  <div class="space-y-8">
    <PageHeader :title="$t('favorites.h1')" :subtitle="$t('favorites.lede')" />

    <LoadingSkeleton v-if="loading" :rows="3" :label="$t('favorites.loading')" />
    <ErrorState v-else-if="error" :title="error.title" :message="error.message" @retry="load" />
    <EmptyState
      v-else-if="items.length === 0"
      :title="$t('favorites.emptyTitle')"
      :message="$t('favorites.emptyBody')"
      :cta-to="$localePath('/search')"
      :cta-label="$t('favorites.emptyCta')"
    />

    <div v-else class="grid gap-4 md:grid-cols-2">
      <div v-for="item in items" :key="item.id" class="card p-5">
        <div class="flex items-start justify-between gap-2">
          <NuxtLink :to="$localePath(`/products/${item.id}`)" class="font-bold text-ink hover:text-brand">
            {{ item.name }}
          </NuxtLink>
          <button class="text-xs font-semibold text-rose-500 hover:underline shrink-0" @click="remove(item.id)">
            {{ $t('favorites.remove') }}
          </button>
        </div>
        <div class="text-xs text-ink-secondary mt-1">{{ item.provider_name }} · {{ item.category }}</div>
        <div class="text-sm font-bold text-brand mt-2">
          <PriceDisplay :amount="num(item.raw_monthly_price_usd)" currency="USD" cycle="month" />
        </div>
      </div>
    </div>
    <p class="text-[11px] text-ink-secondary">{{ $t('favorites.localNote') }}</p>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ApiError, useApi } from '~/composables/useApi'
import { num } from '~/composables/useNum'
import type { ListEnvelope, SKUDetail } from '~/types/api'

definePageMeta({
  layout: 'dashboard',
  middleware: 'auth',
})

const { t } = useI18n()
useHead(() => ({ title: t('favorites.title') }))

const items = ref<SKUDetail[]>([])
const loading = ref(false)
const error = ref<ApiError | null>(null)
const { request } = useApi()

function readIds(): string[] {
  try {
    return JSON.parse(localStorage.getItem('tirp-favorites') || '[]') as string[]
  } catch {
    return []
  }
}

async function load() {
  loading.value = true
  error.value = null
  try {
    const res = await request<ListEnvelope<SKUDetail>>('/api/v1/catalog/skus')
    const wanted = new Set(readIds())
    items.value = (res.data ?? []).filter((s) => wanted.has(s.id))
  } catch (e) {
    error.value = e instanceof ApiError ? e : ApiError.fromUnknown(e)
  } finally {
    loading.value = false
  }
}

function remove(id: string) {
  // TODO(api): replace with DELETE /api/v1/users/me/favorites when implemented.
  localStorage.setItem('tirp-favorites', JSON.stringify(readIds().filter((x) => x !== id)))
  items.value = items.value.filter((s) => s.id !== id)
}

onMounted(load)
</script>
