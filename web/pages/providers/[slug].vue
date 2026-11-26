<template>
  <div class="max-w-3xl mx-auto space-y-8 py-4">
    <LoadingSkeleton v-if="loading" :rows="3" :label="$t('vendorDetail.loading')" />
    <ErrorState v-else-if="error" :title="error.title" :message="error.message" @retry="load" />
    <EmptyState
      v-else-if="!provider"
      :title="$t('vendorDetail.notFoundTitle')"
      :message="$t('vendorDetail.notFoundBody')"
      :cta-to="$localePath('/providers')"
      :cta-label="$t('vendorDetail.backToVendors')"
    />

    <template v-else-if="provider">
      <div class="space-y-2">
        <NuxtLink :to="$localePath('/providers')" class="text-xs font-semibold text-brand hover:underline">
          {{ $t('vendorDetail.backToVendors') }}
        </NuxtLink>
        <h1 class="text-3xl font-extrabold text-ink tracking-tight">{{ provider.name }}</h1>
        <div class="flex flex-wrap gap-2">
          <StatusBadge :status="provider.affiliate_status" />
          <StatusBadge :status="provider.pricing_api_status" />
        </div>
      </div>

      <div class="card p-6">
        <dl class="grid grid-cols-1 sm:grid-cols-2 gap-x-6 gap-y-3 text-sm">
          <div>
            <dt class="text-xs text-ink-secondary">{{ $t('vendorDetail.tier') }}</dt>
            <dd class="font-semibold text-ink">{{ provider.tier_priority }}</dd>
          </div>
          <div>
            <dt class="text-xs text-ink-secondary">{{ $t('vendorDetail.hq') }}</dt>
            <dd class="font-semibold text-ink">{{ provider.hq_country_code }}</dd>
          </div>
          <div>
            <dt class="text-xs text-ink-secondary">{{ $t('vendorDetail.website') }}</dt>
            <dd><a :href="provider.official_website_url" target="_blank" rel="noopener" class="text-brand hover:underline font-semibold">{{ provider.official_website_url }}</a></dd>
          </div>
          <div>
            <dt class="text-xs text-ink-secondary">{{ $t('vendorDetail.apiDocs') }}</dt>
            <dd><a :href="provider.official_api_docs_url" target="_blank" rel="noopener" class="text-brand hover:underline font-semibold">{{ provider.official_api_docs_url }}</a></dd>
          </div>
        </dl>
      </div>

      <div class="card p-6 space-y-3">
        <h2 class="text-sm font-bold text-ink">{{ $t('vendorDetail.plansTitle') }}</h2>
        <div v-if="plans.length === 0" class="text-xs text-ink-secondary">{{ $t('vendorDetail.noPlans') }}</div>
        <ul v-else class="divide-y divide-subtle">
          <li v-for="sku in plans" :key="sku.id" class="py-3 flex items-center justify-between gap-3">
            <NuxtLink :to="$localePath(`/products/${sku.id}`)" class="text-sm font-semibold text-ink hover:text-brand">
              {{ sku.name }}
            </NuxtLink>
            <span class="text-sm font-bold text-brand tabular-nums">
              <PriceDisplay :amount="num(sku.raw_monthly_price_usd)" currency="USD" cycle="month" />
            </span>
          </li>
        </ul>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ApiError, useApi } from '~/composables/useApi'
import { num } from '~/composables/useNum'
import type { ListEnvelope, Provider, SKUDetail } from '~/types/api'

definePageMeta({
  layout: 'public',
})

const { t } = useI18n()
const route = useRoute()
const slug = typeof route.params.slug === 'string' ? route.params.slug : ''
useHead(() => ({ title: provider.value ? provider.value.name : t('vendors.title') }))

const provider = ref<Provider | null>(null)
const plans = ref<SKUDetail[]>([])
const loading = ref(false)
const error = ref<ApiError | null>(null)
const { request } = useApi()

async function load() {
  loading.value = true
  error.value = null
  try {
    const [provs, skus] = await Promise.all([
      request<ListEnvelope<Provider>>('/api/v1/providers'),
      request<ListEnvelope<SKUDetail>>('/api/v1/catalog/skus'),
    ])
    provider.value = (provs.data ?? []).find((p) => p.slug === slug) || null
    plans.value = (skus.data ?? []).filter((s) => s.provider_slug === slug).slice(0, 12)
  } catch (e) {
    error.value = e instanceof ApiError ? e : ApiError.fromUnknown(e)
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>
