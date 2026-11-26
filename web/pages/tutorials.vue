<template>
  <div class="max-w-4xl mx-auto space-y-8 py-4">
    <div class="space-y-2">
      <h1 class="text-3xl font-extrabold text-ink tracking-tight">{{ $t('tutorials.h1') }}</h1>
      <p class="text-sm text-ink-secondary">
        {{ $t('tutorials.lede') }}
        <span class="ml-1 inline-block text-[10px] font-bold px-2 py-0.5 rounded bg-amber-500/10 text-amber-700 dark:text-amber-300 border border-amber-500/30">{{ $t('common.demoData') }}</span>
      </p>
    </div>

    <div class="space-y-4">
      <article v-for="a in items" :key="a.slug" class="card p-5 space-y-1">
        <h2 class="font-bold text-ink">{{ a.title }}</h2>
        <p class="text-sm text-ink-secondary">{{ a.excerpt }}</p>
        <p class="text-[11px] text-ink-secondary">{{ $t('tutorials.updated', { date: a.updatedAt }) }}</p>
      </article>
    </div>
  </div>
</template>

<script setup lang="ts">
import { mockArticles } from '~/utils/mock-data'

definePageMeta({
  layout: 'public',
})

const { t } = useI18n()
useHead(() => ({ title: t('tutorials.title') }))

// TODO(api): replace with GET /api/v1/content when implemented.
const items = mockArticles('tutorial')
</script>
