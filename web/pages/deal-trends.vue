<template>
  <div class="max-w-4xl mx-auto space-y-8 py-4">
    <div class="space-y-2">
      <h1 class="text-3xl font-extrabold text-ink tracking-tight">{{ $t('trends.h1') }}</h1>
      <p class="text-sm text-ink-secondary">
        {{ $t('trends.lede') }}
        <span class="ml-1 inline-block text-[10px] font-bold px-2 py-0.5 rounded bg-amber-500/10 text-amber-700 dark:text-amber-300 border border-amber-500/30">{{ $t('common.demoData') }}</span>
      </p>
    </div>

    <div class="grid gap-4 md:grid-cols-2">
      <div v-for="trend in trends" :key="trend.category" class="card p-5 space-y-3">
        <div class="flex items-baseline justify-between gap-2">
          <h2 class="font-bold text-ink">{{ trend.category }}</h2>
          <span class="text-sm font-extrabold text-emerald-600 dark:text-emerald-400 tabular-nums">
            {{ $t('trends.avgDrop', { v: trend.avgDropPct, n: trend.dealCount }) }}
          </span>
        </div>
        <div class="flex items-end gap-1 h-12" aria-hidden="true">
          <div
            v-for="(v, i) in trend.spark"
            :key="i"
            class="flex-1 rounded-sm bg-brand/70"
            :style="{ height: `${Math.max(8, (v / 24) * 100)}%` }"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { mockTrends } from '~/utils/mock-data'

definePageMeta({
  layout: 'public',
})

const { t } = useI18n()
useHead(() => ({ title: t('trends.title') }))

// TODO(api): replace with GET /api/v1/analytics/deal-performance when implemented.
const trends = mockTrends()
</script>
