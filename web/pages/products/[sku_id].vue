<template>
  <div class="max-w-3xl mx-auto space-y-8 py-4">
    <LoadingSkeleton v-if="loading" :rows="4" :label="$t('product.loading')" />
    <ErrorState v-else-if="error" :title="error.title" :message="error.message" @retry="load" />
    <EmptyState
      v-else-if="!sku"
      :title="$t('product.notFoundTitle')"
      :message="$t('product.notFoundBody')"
      :cta-to="$localePath('/search')"
      :cta-label="$t('product.backToSearch')"
    />

    <template v-else-if="sku">
      <div class="space-y-2">
        <div class="text-xs font-semibold text-brand">{{ sku.provider_name }} · {{ sku.category }}</div>
        <h1 class="text-3xl font-extrabold text-ink tracking-tight">{{ sku.name }}</h1>
        <div class="flex items-center gap-2">
          <StatusBadge :status="sku.freshness_status" />
          <span class="text-lg font-extrabold text-brand tabular-nums">
            <PriceDisplay :amount="num(sku.raw_monthly_price_usd)" currency="USD" cycle="month" />
          </span>
        </div>
      </div>

      <div class="card p-6">
        <h2 class="text-sm font-bold text-ink mb-3">{{ $t('product.factsTitle') }}</h2>
        <dl class="grid grid-cols-1 sm:grid-cols-2 gap-x-6 gap-y-3 text-sm">
          <div>
            <dt class="text-xs text-ink-secondary">{{ $t('product.planCode') }}</dt>
            <dd class="font-semibold text-ink font-mono text-xs">{{ sku.plan_code }}</dd>
          </div>
          <div>
            <dt class="text-xs text-ink-secondary">{{ $t('product.provider') }}</dt>
            <dd>
              <NuxtLink :to="$localePath(`/providers/${sku.provider_slug}`)" class="font-semibold text-brand hover:underline">
                {{ sku.provider_name }}
              </NuxtLink>
            </dd>
          </div>
        </dl>
      </div>

      <div class="flex flex-wrap gap-3">
        <NuxtLink :to="$localePath({ path: '/compare', query: { ids: sku.id } })" class="btn-primary !px-5 !py-2.5">
          {{ $t('product.compareCta') }}
        </NuxtLink>
        <button class="btn-ghost !px-5 !py-2.5" @click="saveFavorite">
          {{ saved ? $t('product.saved') : $t('product.saveFavorite') }}
        </button>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ApiError, useApi } from '~/composables/useApi'
import { num } from '~/composables/useNum'
import type { ItemEnvelope, SKUDetail } from '~/types/api'

definePageMeta({
  layout: 'public',
})

const { t } = useI18n()
const route = useRoute()
const id = typeof route.params.sku_id === 'string' ? route.params.sku_id : ''
useHead(() => ({ title: sku.value ? sku.value.name : t('product.title') }))

const sku = ref<SKUDetail | null>(null)
const saved = ref(false)
const loading = ref(false)
const error = ref<ApiError | null>(null)
const { request } = useApi()

function saveFavorite() {
  // TODO(api): replace with POST /api/v1/users/me/favorites when implemented.
  try {
    const ids = JSON.parse(localStorage.getItem('tirp-favorites') || '[]') as string[]
    if (sku.value && !ids.includes(sku.value.id)) {
      ids.push(sku.value.id)
      localStorage.setItem('tirp-favorites', JSON.stringify(ids))
    }
    saved.value = true
  } catch {
    saved.value = false
  }
}

async function load() {
  loading.value = true
  error.value = null
  try {
    const res = await request<ItemEnvelope<SKUDetail>>(`/api/v1/catalog/skus/${id}`)
    sku.value = res.data ?? null
  } catch (e) {
    error.value = e instanceof ApiError ? e : ApiError.fromUnknown(e)
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>
