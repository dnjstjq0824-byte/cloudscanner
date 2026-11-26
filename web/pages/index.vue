<template>
  <div class="space-y-12">
    <div v-if="maintenance" class="px-4 py-3 rounded-xl bg-amber-500/10 border border-amber-500/30 text-xs text-amber-700 dark:text-amber-300" role="alert">
      Scheduled maintenance is in progress. Search and console actions may return 503 until it completes.
    </div>

    <section class="py-10 sm:py-16">
      <div class="max-w-3xl space-y-6">
        <p class="text-xs font-bold uppercase tracking-[0.2em] text-brand">Independent SaaS deal intelligence</p>
        <h1 class="text-4xl sm:text-5xl font-extrabold tracking-tight text-ink leading-[1.1]">
          Find real SaaS deals.<br />Never overpay.
        </h1>
        <p class="text-base sm:text-lg text-ink-secondary max-w-2xl">
          TirpScanner indexes subscription plans across vendors, tracks price
          drops and promos, and discloses every affiliate link. Cloud catalog
          and server management were retired with the 2026-11-26 pivot.
        </p>
        <div class="flex flex-wrap gap-3 pt-2">
          <NuxtLink to="/dashboard" class="btn-primary !px-6 !py-3">Open the Console</NuxtLink>
          <NuxtLink to="/affiliate-disclosure" class="btn-ghost !px-6 !py-3">How We Disclose</NuxtLink>
        </div>
      </div>
    </section>

    <section aria-label="Live catalog statistics" class="grid grid-cols-2 md:grid-cols-4 gap-4 sm:gap-6">
      <div class="card p-6 text-center">
        <div class="text-3xl font-bold text-brand tabular-nums">{{ stats.providers }}</div>
        <div class="text-xs text-ink-secondary mt-1 uppercase tracking-wider">Vendors Tracked</div>
      </div>
      <div class="card p-6 text-center">
        <div class="text-3xl font-bold text-brand tabular-nums">{{ stats.plans }}</div>
        <div class="text-xs text-ink-secondary mt-1 uppercase tracking-wider">Plans Indexed</div>
      </div>
      <div class="card p-6 text-center">
        <div class="text-3xl font-bold text-brand tabular-nums">{{ stats.latency }}</div>
        <div class="text-xs text-ink-secondary mt-1 uppercase tracking-wider">Search Latency</div>
      </div>
      <div class="card p-6 text-center">
        <div class="text-3xl font-bold text-brand tabular-nums">{{ stats.categories }}</div>
        <div class="text-xs text-ink-secondary mt-1 uppercase tracking-wider">Categories</div>
      </div>
    </section>

    <section class="card p-6 sm:p-8 space-y-5" aria-label="How TirpScanner works">
      <h2 class="text-lg font-bold text-ink">One flow: scan → compare → save</h2>
      <ol class="space-y-4 text-sm">
        <li class="flex gap-3">
          <span class="font-bold text-brand tabular-nums">01</span>
          <p class="text-ink-secondary"><strong class="text-ink">Scan live deals</strong> across SaaS vendors: plan prices, trials, coupons, and renewals, normalized per seat and billing cycle.</p>
        </li>
        <li class="flex gap-3">
          <span class="font-bold text-brand tabular-nums">02</span>
          <p class="text-ink-secondary"><strong class="text-ink">Compare plans side by side</strong> with transparent effective pricing — no sponsored ranking boosts, ever.</p>
        </li>
        <li class="flex gap-3">
          <span class="font-bold text-brand tabular-nums">03</span>
          <p class="text-ink-secondary"><strong class="text-ink">Get price-drop alerts</strong> when tracked plans go on sale or cheaper tiers appear.</p>
        </li>
      </ol>
      <div class="flex flex-wrap gap-3 pt-1">
        <NuxtLink to="/dashboard" class="text-xs font-semibold text-brand hover:underline">Open the console →</NuxtLink>
      </div>
    </section>

    <section class="text-xs text-ink-secondary leading-relaxed border-t border-subtle pt-6" aria-label="Independence note">
      TirpScanner ranks by measured value, never by commission. Outbound deal links are labeled and governed by our
      <NuxtLink to="/affiliate-disclosure" class="text-brand hover:underline font-semibold">affiliate disclosure</NuxtLink>.
    </section>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import type { Category, ListEnvelope, Provider, SearchResponse, SKUDetail } from '~/types/api'

definePageMeta({
  layout: 'public',
})

const stats = reactive({ providers: '—', plans: '—', latency: '—', categories: '—' })
const maintenance = ref(false)

onMounted(async () => {
  try {
    const { request } = useApi()
    const started = performance.now()
    const [provs, skus, cats, probe] = await Promise.all([
      request<ListEnvelope<Provider>>('/api/v1/providers'),
      request<ListEnvelope<SKUDetail>>('/api/v1/catalog/skus'),
      request<ListEnvelope<Category>>('/api/v1/catalog/categories'),
      request<SearchResponse>('/api/v1/search/skus', { query: { q: 'crm', limit: 1 } }),
    ])
    stats.providers = String(provs.meta?.count ?? provs.data?.length ?? '—')
    stats.plans = String(skus.meta?.count ?? skus.data?.length ?? '—')
    stats.categories = String(cats.meta?.count ?? cats.data?.length ?? '—')
    const measured = probe.execution_ms > 0 ? probe.execution_ms : Math.round(performance.now() - started)
    stats.latency = `${measured}ms`
    try {
      const cfg = await request<{ data: { maintenance_mode: boolean } }>('/api/v1/site-config')
      maintenance.value = cfg.data?.maintenance_mode === true
    } catch {
      maintenance.value = false
    }
  } catch {
    // Homepage stays readable when the API is unreachable; stats show placeholders.
  }
})
</script>
