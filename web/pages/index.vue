<template>
  <div class="space-y-12">
    <div v-if="maintenance" class="px-4 py-3 rounded-xl bg-amber-500/10 border border-amber-500/30 text-xs text-amber-700 dark:text-amber-300" role="alert">
      {{ $t('home.maintenance') }}
    </div>

    <section class="py-10 sm:py-16">
      <div class="max-w-3xl space-y-6">
        <p class="text-xs font-bold uppercase tracking-[0.2em] text-brand">{{ $t('home.kicker') }}</p>
        <h1 class="text-4xl sm:text-5xl font-extrabold tracking-tight text-ink leading-[1.1]">
          {{ $t('home.h1a') }}<br />{{ $t('home.h1b') }}
        </h1>
        <p class="text-base sm:text-lg text-ink-secondary max-w-2xl">
          {{ $t('home.lede') }}
        </p>
        <div class="flex flex-wrap gap-3 pt-2">
          <NuxtLink :to="$localePath('/dashboard')" class="btn-primary !px-6 !py-3">{{ $t('home.ctaConsole') }}</NuxtLink>
          <NuxtLink :to="$localePath('/affiliate-disclosure')" class="btn-ghost !px-6 !py-3">{{ $t('home.ctaDisclosure') }}</NuxtLink>
        </div>
      </div>
    </section>

    <section :aria-label="$t('home.statsLabel')" class="grid grid-cols-2 md:grid-cols-4 gap-4 sm:gap-6">
      <div class="card p-6 text-center">
        <div class="text-3xl font-bold text-brand tabular-nums">{{ stats.providers }}</div>
        <div class="text-xs text-ink-secondary mt-1 uppercase tracking-wider">{{ $t('home.statVendors') }}</div>
      </div>
      <div class="card p-6 text-center">
        <div class="text-3xl font-bold text-brand tabular-nums">{{ stats.plans }}</div>
        <div class="text-xs text-ink-secondary mt-1 uppercase tracking-wider">{{ $t('home.statPlans') }}</div>
      </div>
      <div class="card p-6 text-center">
        <div class="text-3xl font-bold text-brand tabular-nums">{{ stats.latency }}</div>
        <div class="text-xs text-ink-secondary mt-1 uppercase tracking-wider">{{ $t('home.statLatency') }}</div>
      </div>
      <div class="card p-6 text-center">
        <div class="text-3xl font-bold text-brand tabular-nums">{{ stats.categories }}</div>
        <div class="text-xs text-ink-secondary mt-1 uppercase tracking-wider">{{ $t('home.statCategories') }}</div>
      </div>
    </section>

    <section class="card p-6 sm:p-8 space-y-5" :aria-label="$t('home.howLabel')">
      <h2 class="text-lg font-bold text-ink">{{ $t('home.howTitle') }}</h2>
      <ol class="space-y-4 text-sm">
        <li class="flex gap-3">
          <span class="font-bold text-brand tabular-nums">01</span>
          <p class="text-ink-secondary"><strong class="text-ink">{{ $t('home.step1title') }}</strong> {{ $t('home.step1body') }}</p>
        </li>
        <li class="flex gap-3">
          <span class="font-bold text-brand tabular-nums">02</span>
          <p class="text-ink-secondary"><strong class="text-ink">{{ $t('home.step2title') }}</strong> {{ $t('home.step2body') }}</p>
        </li>
        <li class="flex gap-3">
          <span class="font-bold text-brand tabular-nums">03</span>
          <p class="text-ink-secondary"><strong class="text-ink">{{ $t('home.step3title') }}</strong> {{ $t('home.step3body') }}</p>
        </li>
      </ol>
      <div class="flex flex-wrap gap-3 pt-1">
        <NuxtLink :to="$localePath('/dashboard')" class="text-xs font-semibold text-brand hover:underline">{{ $t('home.howCta') }}</NuxtLink>
      </div>
    </section>

    <section class="text-xs text-ink-secondary leading-relaxed border-t border-subtle pt-6" :aria-label="$t('home.indepLabel')">
      {{ $t('home.indep1') }} {{ $t('home.indep2pre') }}
      <NuxtLink :to="$localePath('/affiliate-disclosure')" class="text-brand hover:underline font-semibold">{{ $t('home.indep2link') }}</NuxtLink>{{ $t('home.indep2post') }}
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
