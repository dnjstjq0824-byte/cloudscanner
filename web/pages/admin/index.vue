<template>
  <div class="space-y-8">
    <div>
      <h2 class="text-2xl font-extrabold text-white tracking-tight">Operations Health</h2>
      <p class="text-sm text-slate-400 mt-1">Live platform signals. Every value below is read from the running API.</p>
    </div>

    <div v-if="loading" class="grid grid-cols-1 md:grid-cols-3 gap-6" aria-busy="true" aria-label="Loading ops health">
      <div v-for="n in 3" :key="n" class="p-6 rounded-xl bg-slate-900 border border-slate-800 animate-pulse">
        <div class="h-8 w-16 rounded bg-slate-800" />
      </div>
    </div>

    <div v-else-if="error" class="p-8 rounded-xl bg-red-950/20 border border-red-900 text-center space-y-3">
      <h3 class="text-base font-bold text-white">{{ error.title }}</h3>
      <p class="text-xs text-slate-400">{{ error.message }}</p>
      <button class="px-4 py-2 rounded bg-sky-600 hover:bg-sky-500 text-white text-xs font-semibold" @click="load">
        Retry
      </button>
    </div>

    <template v-else>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div class="p-6 rounded-xl bg-slate-900 border border-slate-800">
          <div class="flex items-center gap-2">
            <span class="w-2.5 h-2.5 rounded-full" :class="readiness?.status === 'ready' ? 'bg-emerald-400' : 'bg-red-400'" />
            <span class="text-lg font-bold text-white">{{ readiness?.status === 'ready' ? 'API Ready' : 'API Not Ready' }}</span>
          </div>
          <div class="text-xs text-slate-400 mt-2">Database: <strong class="text-slate-200">{{ readiness?.database || 'unknown' }}</strong></div>
        </div>
        <div class="p-6 rounded-xl bg-slate-900 border border-slate-800">
          <div class="text-3xl font-bold text-sky-400 tabular-nums">{{ providerCount }}</div>
          <div class="text-xs text-slate-400 mt-1 uppercase tracking-wider">Providers in Catalog</div>
        </div>
        <div class="p-6 rounded-xl bg-slate-900 border border-slate-800">
          <div class="text-3xl font-bold text-sky-400 tabular-nums">{{ verifiedCount }}</div>
          <div class="text-xs text-slate-400 mt-1 uppercase tracking-wider">Verified Pricing APIs</div>
        </div>
      </div>

      <div class="p-6 rounded-xl bg-slate-900 border border-slate-800 space-y-3">
        <h3 class="text-sm font-bold text-white">Affiliate compliance</h3>
        <div class="text-xs text-slate-400 space-y-1">
          <p>Disclosure policy version: <strong class="text-slate-200">{{ disclosure?.policy_version || '—' }}</strong></p>
          <p>Independent comparison: <strong class="text-slate-200">{{ disclosure?.independent_comparison ? 'enforced' : '—' }}</strong></p>
          <p>Pay-for-placement: <strong class="text-slate-200">{{ disclosure && !disclosure.pay_for_placement_allowed ? 'blocked' : '—' }}</strong></p>
        </div>
        <NuxtLink to="/affiliate-disclosure" class="inline-block text-xs font-semibold text-sky-400 hover:text-sky-300">
          View public disclosure &rarr;
        </NuxtLink>
      </div>

      <div class="p-6 rounded-xl bg-slate-900 border border-slate-800 space-y-3">
        <h3 class="text-sm font-bold text-white">Sync status by provider</h3>
        <div class="overflow-x-auto">
          <table class="w-full text-left text-xs">
            <thead class="text-slate-500 uppercase tracking-wider">
              <tr>
                <th class="py-2 pr-4">Provider</th>
                <th class="py-2 pr-4">Pricing API</th>
                <th class="py-2 pr-4">Management API</th>
                <th class="py-2">Affiliate</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-800 text-slate-300">
              <tr v-for="p in providers" :key="p.slug">
                <td class="py-2 pr-4 font-semibold text-white">{{ p.name }}</td>
                <td class="py-2 pr-4">{{ p.pricing_api_status }}</td>
                <td class="py-2 pr-4">{{ p.management_api_status }}</td>
                <td class="py-2">{{ p.affiliate_status }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { ApiError, useApi } from '~/composables/useApi'
import type { ListEnvelope, Provider } from '~/types/api'

definePageMeta({
  layout: 'admin',
  middleware: 'admin',
})

useHead({ title: 'Ops Health' })

interface Readiness {
  status: string
  service: string
  database: string
}

interface Disclosure {
  policy_version: string
  independent_comparison: boolean
  pay_for_placement_allowed: boolean
}

const loading = ref(false)
const error = ref<ApiError | null>(null)
const readiness = ref<Readiness | null>(null)
const providers = ref<Provider[]>([])
const disclosure = ref<Disclosure | null>(null)

const providerCount = computed(() => providers.value.length)
const verifiedCount = computed(() => providers.value.filter((p) => p.pricing_api_status === 'VERIFIED').length)

const { request } = useApi()

async function load() {
  loading.value = true
  error.value = null
  try {
    const [ready, provs, disc] = await Promise.all([
      request<Readiness>('/api/v1/readiness'),
      request<ListEnvelope<Provider>>('/api/v1/providers'),
      request<Disclosure>('/api/v1/affiliate/disclosure'),
    ])
    readiness.value = ready
    providers.value = provs.data ?? []
    disclosure.value = disc
  } catch (e) {
    error.value = e instanceof ApiError ? e : ApiError.fromUnknown(e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  load()
})
</script>
