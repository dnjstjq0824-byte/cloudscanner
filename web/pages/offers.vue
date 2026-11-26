<template>
  <div class="max-w-5xl mx-auto space-y-8 py-4">
    <div class="space-y-2">
      <h1 class="text-3xl font-extrabold text-ink tracking-tight">{{ $t('offers.h1') }}</h1>
      <p class="text-sm text-ink-secondary">
        {{ $t('offers.lede') }}
        <span class="ml-1 inline-block text-[10px] font-bold px-2 py-0.5 rounded bg-amber-500/10 text-amber-700 dark:text-amber-300 border border-amber-500/30">{{ $t('common.demoData') }}</span>
      </p>
    </div>

    <div class="grid gap-4 md:grid-cols-2">
      <div v-for="promo in promos" :key="promo.id" class="card p-5 space-y-2">
        <div class="flex items-center justify-between gap-2">
          <span class="font-bold text-ink">{{ promo.vendor }}</span>
          <span class="text-xs font-extrabold text-emerald-600 dark:text-emerald-400 tabular-nums">−{{ promo.discountPct }}%</span>
        </div>
        <p class="text-sm text-ink-secondary">{{ promo.title }}</p>
        <div class="flex items-center justify-between text-xs">
          <code class="font-mono bg-elevated px-2 py-1 rounded">{{ promo.code }}</code>
          <span class="text-ink-secondary">{{ $t('offers.expires', { date: promo.expiresAt }) }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { mockPromotions } from '~/utils/mock-data'

definePageMeta({
  layout: 'public',
})

const { t } = useI18n()
useHead(() => ({
  title: t('offers.title'),
  meta: [{ name: 'description', content: t('offers.meta') }],
}))

// TODO(api): replace with GET /api/v1/promotions when implemented.
const promos = mockPromotions()
</script>
