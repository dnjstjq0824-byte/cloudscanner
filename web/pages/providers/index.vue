<template>
  <div class="space-y-8">
    <div class="space-y-2">
      <h1 class="text-3xl font-extrabold text-white tracking-tight">Cloud Providers Directory</h1>
      <p class="text-sm text-slate-400">
        {{ providers.length }} global and regional cloud infrastructure providers evaluated and normalized across 36 price dimensions.
      </p>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6" aria-busy="true" aria-label="Loading providers">
      <div v-for="n in 6" :key="n" class="p-6 rounded-xl bg-slate-900 border border-slate-800 animate-pulse space-y-3">
        <div class="h-5 w-1/2 rounded bg-slate-800" />
        <div class="h-3 w-full rounded bg-slate-800" />
        <div class="h-3 w-2/3 rounded bg-slate-800" />
      </div>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="text-center py-16 bg-red-950/20 rounded-xl border border-red-900 space-y-3">
      <div class="text-3xl">⚠️</div>
      <h4 class="text-base font-bold text-white">{{ error.title }}</h4>
      <p class="text-xs text-slate-400 max-w-sm mx-auto">{{ error.message }}</p>
      <button
        class="px-4 py-2 rounded bg-sky-600 hover:bg-sky-500 text-white text-xs font-semibold transition"
        @click="loadProviders"
      >
        Retry
      </button>
    </div>

    <!-- Provider Cards Grid -->
    <div v-else-if="providers.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div
        v-for="prov in providers"
        :key="prov.slug"
        class="p-6 rounded-xl bg-slate-900 border border-slate-800 hover:border-slate-700 transition flex flex-col justify-between space-y-4"
      >
        <div class="space-y-3">
          <div class="flex items-center justify-between">
            <span class="text-lg font-bold text-white">{{ prov.name }}</span>
            <span class="px-2 py-0.5 rounded text-[10px] font-semibold bg-sky-950 text-sky-400 border border-sky-800">
              {{ prov.hq_country_code }}
            </span>
          </div>

          <div class="flex flex-wrap gap-1.5 pt-1">
            <span
              class="px-2 py-0.5 rounded text-[10px] border"
              :class="prov.pricing_api_status === 'VERIFIED' ? 'bg-emerald-950 text-emerald-400 border-emerald-800' : 'bg-slate-800 text-slate-300 border-slate-700'"
            >
              Pricing: {{ prov.pricing_api_status }}
            </span>
            <span
              class="px-2 py-0.5 rounded text-[10px] border"
              :class="prov.management_api_status === 'VERIFIED' ? 'bg-emerald-950 text-emerald-400 border-emerald-800' : 'bg-slate-800 text-slate-300 border-slate-700'"
            >
              BYOC: {{ prov.management_api_status }}
            </span>
            <span class="px-2 py-0.5 rounded text-[10px] bg-slate-800 text-slate-300 border border-slate-700">
              {{ prov.tier_priority }}
            </span>
          </div>

          <div class="flex flex-wrap gap-x-4 gap-y-1 text-[11px] text-slate-400">
            <span v-if="prov.capabilities.supports_provisioning">✓ Provisioning</span>
            <span v-if="prov.capabilities.supports_snapshots">✓ Snapshots</span>
            <span v-if="prov.capabilities.supports_firewalls">✓ Firewalls</span>
            <span v-if="prov.supports_byoc_api_token">✓ API Token</span>
            <span v-if="prov.supports_byoc_oauth">✓ OAuth</span>
          </div>
        </div>

        <div class="pt-4 border-t border-slate-800 flex items-center justify-between">
          <span class="text-xs text-slate-500 tabular-nums">{{ prov.regions?.length ?? 0 }} Regions</span>
          <NuxtLink
            :to="`/providers/${prov.slug}`"
            class="text-xs font-semibold text-sky-400 hover:text-sky-300"
          >
            View SKUs &rarr;
          </NuxtLink>
        </div>
      </div>
    </div>

    <!-- Empty -->
    <div v-else class="text-center py-16 bg-slate-900/40 rounded-xl border border-slate-800 space-y-3">
      <div class="text-3xl">☁️</div>
      <h4 class="text-base font-bold text-white">No providers available</h4>
      <p class="text-xs text-slate-400 max-w-sm mx-auto">The provider catalog is empty. Sync may still be running.</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ApiError, useApi } from '~/composables/useApi'
import type { ListEnvelope, Provider } from '~/types/api'

definePageMeta({
  layout: 'public',
})

useHead({
  title: 'Cloud Providers Directory',
  meta: [
    { name: 'description', content: 'Directory of 18+ verified cloud providers with pricing API status, BYOC support, and region coverage.' },
  ],
})

const loading = ref(false)
const error = ref<ApiError | null>(null)
const providers = ref<Provider[]>([])

const { request } = useApi()

async function loadProviders() {
  loading.value = true
  error.value = null
  try {
    const res = await request<ListEnvelope<Provider>>('/api/v1/providers')
    providers.value = res.data ?? []
  } catch (e) {
    error.value = e instanceof ApiError ? e : ApiError.fromUnknown(e)
    providers.value = []
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadProviders()
})
</script>
