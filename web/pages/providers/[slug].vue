<template>
  <div class="space-y-8">
    <!-- Loading -->
    <div v-if="loading" class="space-y-4" aria-busy="true" aria-label="Loading provider">
      <div class="p-8 rounded-2xl bg-slate-900 border border-slate-800 animate-pulse">
        <div class="h-8 w-1/3 rounded bg-slate-800" />
        <div class="mt-4 h-4 w-2/3 rounded bg-slate-800" />
      </div>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="text-center py-16 bg-red-950/20 rounded-xl border border-red-900 space-y-3">
      <div class="text-3xl">⚠️</div>
      <h4 class="text-base font-bold text-white">{{ error.title }}</h4>
      <p class="text-xs text-slate-400 max-w-sm mx-auto">{{ error.message }}</p>
      <div class="flex justify-center gap-2">
        <button
          class="px-4 py-2 rounded bg-sky-600 hover:bg-sky-500 text-white text-xs font-semibold transition"
          @click="loadProvider"
        >
          Retry
        </button>
        <NuxtLink
          to="/providers"
          class="px-4 py-2 rounded bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs font-semibold transition"
        >
          All Providers
        </NuxtLink>
      </div>
    </div>

    <template v-else-if="provider">
      <!-- Provider Header -->
      <div class="p-8 rounded-2xl bg-slate-900 border border-slate-800 space-y-4">
        <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
          <div>
            <div class="flex items-center space-x-3 flex-wrap">
              <h1 class="text-3xl font-extrabold text-white tracking-tight">{{ provider.name }}</h1>
              <span
                class="px-2.5 py-0.5 rounded text-xs font-semibold border"
                :class="provider.pricing_api_status === 'VERIFIED' ? 'bg-emerald-950 text-emerald-400 border-emerald-800' : 'bg-slate-800 text-slate-300 border-slate-700'"
              >
                {{ provider.pricing_api_status === 'VERIFIED' ? 'Verified Provider' : provider.pricing_api_status }}
              </span>
            </div>
            <p class="text-sm text-slate-400 mt-2 max-w-2xl leading-relaxed">
              {{ provider.tier_priority }} · HQ {{ provider.hq_country_code }} ·
              <a :href="provider.official_website_url" target="_blank" rel="noopener sponsored" class="text-sky-400 hover:text-sky-300">Official Site</a>
              ·
              <a :href="provider.official_api_docs_url" target="_blank" rel="noopener" class="text-sky-400 hover:text-sky-300">API Docs</a>
            </p>
          </div>

          <a
            :href="dealUrl"
            class="px-6 py-3 rounded-lg bg-sky-600 hover:bg-sky-500 font-semibold text-white text-sm transition shadow-lg shrink-0"
          >
            Deploy on {{ provider.name }}
          </a>
        </div>

        <!-- Quick Metrics -->
        <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 pt-4 border-t border-slate-800 text-xs">
          <div>
            <span class="text-slate-500 block">Pricing API</span>
            <span class="text-white font-semibold">{{ provider.pricing_api_status }}</span>
          </div>
          <div>
            <span class="text-slate-500 block">Management API</span>
            <span class="text-white font-semibold">{{ provider.management_api_status }}</span>
          </div>
          <div>
            <span class="text-slate-500 block">Regions</span>
            <span class="text-white font-semibold tabular-nums">{{ provider.regions?.length ?? 0 }} locations</span>
          </div>
          <div>
            <span class="text-slate-500 block">BYOC</span>
            <span class="text-white font-semibold">{{ byocLabel }}</span>
          </div>
        </div>

        <!-- Regions -->
        <div v-if="provider.regions && provider.regions.length > 0" class="flex flex-wrap gap-1.5 pt-1">
          <span
            v-for="r in provider.regions"
            :key="r.region_code"
            class="px-2 py-0.5 rounded text-[10px] bg-slate-800 text-slate-300 border border-slate-700"
            :title="`${r.name}, ${r.city} (${r.country_code})`"
          >
            {{ r.region_code }} · {{ r.city }}
          </span>
        </div>
      </div>

      <!-- SKUs Catalog Table -->
      <div class="space-y-4">
        <h2 class="text-xl font-bold text-white">Popular Instance Types & Rates</h2>
        <div v-if="skusLoading" class="rounded-xl border border-slate-800 overflow-hidden" aria-busy="true" aria-label="Loading SKUs">
          <div v-for="n in 4" :key="n" class="p-4 bg-slate-900 border-b border-slate-800 animate-pulse">
            <div class="h-4 w-3/4 rounded bg-slate-800" />
          </div>
        </div>
        <div v-else-if="skus.length > 0" class="overflow-x-auto rounded-xl border border-slate-800">
          <table class="w-full text-left text-xs bg-slate-900 divide-y divide-slate-800">
            <thead class="bg-slate-950 text-slate-400 font-semibold uppercase tracking-wider">
              <tr>
                <th class="py-3 px-4">Plan Code</th>
                <th class="py-3 px-4">vCPU</th>
                <th class="py-3 px-4">Memory</th>
                <th class="py-3 px-4">Disk</th>
                <th class="py-3 px-4">Bandwidth</th>
                <th class="py-3 px-4">Hourly Rate</th>
                <th class="py-3 px-4 text-right">Monthly Cap</th>
                <th class="py-3 px-4 text-center">Action</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-800 text-slate-300 tabular-nums">
              <tr v-for="sku in skus" :key="sku.id" class="hover:bg-slate-800/40 transition">
                <td class="py-3.5 px-4 font-bold text-white">{{ sku.plan_code }}</td>
                <td class="py-3.5 px-4">{{ sku.vcpu }} vCPU</td>
                <td class="py-3.5 px-4">{{ (sku.ram_mb / 1024).toFixed(1) }} GB</td>
                <td class="py-3.5 px-4">{{ sku.disk_gb }} GB {{ sku.disk_type }}</td>
                <td class="py-3.5 px-4">{{ (sku.included_bandwidth_gb / 1024).toFixed(1) }} TB</td>
                <td class="py-3.5 px-4">${{ sku.raw_hourly_price_usd.toFixed(4) }}/hr</td>
                <td class="py-3.5 px-4 text-right font-extrabold text-white">${{ sku.raw_monthly_price_usd.toFixed(2) }}/mo</td>
                <td class="py-3.5 px-4 text-center">
                  <a
                    :href="dealUrl"
                    class="px-3 py-1 rounded bg-sky-600 hover:bg-sky-500 text-white font-semibold text-[11px] transition"
                  >
                    Choose
                  </a>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        <div v-else class="text-center py-12 bg-slate-900/40 rounded-xl border border-slate-800">
          <p class="text-sm text-slate-400">No instance types synced for this provider yet.</p>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { ApiError, useApi } from '~/composables/useApi'
import type { ItemEnvelope, ListEnvelope, Provider, SKUDetail } from '~/types/api'

definePageMeta({
  layout: 'public',
})

const route = useRoute()
const slug = computed(() => (route.params.slug as string) || 'digitalocean')

const loading = ref(false)
const skusLoading = ref(false)
const error = ref<ApiError | null>(null)
const provider = ref<Provider | null>(null)
const skus = ref<SKUDetail[]>([])

const { request, baseURL } = useApi()

const dealUrl = computed(() => `${baseURL}/go/${slug.value}?src=provider_profile`)

const byocLabel = computed(() => {
  if (!provider.value) {
    return '—'
  }
  const modes: string[] = []
  if (provider.value.supports_byoc_oauth) {
    modes.push('OAuth')
  }
  if (provider.value.supports_byoc_api_token) {
    modes.push('API Token')
  }
  return modes.length > 0 ? modes.join(' + ') : 'Not supported'
})

useHead(() => ({
  title: provider.value ? `${provider.value.name} Instances & Pricing` : 'Provider',
  meta: [
    { name: 'description', content: provider.value ? `Compare ${provider.value.name} instance types, hourly rates, and monthly caps with transparent TCO.` : 'Cloud provider profile.' },
  ],
}))

async function loadProvider() {
  loading.value = true
  error.value = null
  try {
    const res = await request<ItemEnvelope<Provider>>(`/api/v1/providers/${slug.value}`)
    provider.value = res.data
    await loadSkus()
  } catch (e) {
    error.value = e instanceof ApiError ? e : ApiError.fromUnknown(e)
    provider.value = null
  } finally {
    loading.value = false
  }
}

async function loadSkus() {
  skusLoading.value = true
  try {
    const res = await request<ListEnvelope<SKUDetail>>('/api/v1/catalog/skus', {
      query: { provider: slug.value },
    })
    skus.value = res.data ?? []
  } catch {
    skus.value = []
  } finally {
    skusLoading.value = false
  }
}

watch(slug, () => loadProvider())

onMounted(() => {
  loadProvider()
})
</script>
