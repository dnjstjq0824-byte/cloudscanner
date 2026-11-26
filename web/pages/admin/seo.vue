<template>
  <div class="space-y-8">
    <PageHeader :title="$t('seoAdmin.h1')" :subtitle="$t('seoAdmin.lede')" />

    <div class="card p-0 overflow-hidden">
      <table class="table-shell">
        <thead>
          <tr>
            <th class="py-3 px-4 text-left">{{ $t('seoAdmin.colCheck') }}</th>
            <th class="py-3 px-4 text-left">{{ $t('seoAdmin.colStatus') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="row in rows" :key="row.id" class="border-t border-subtle text-sm">
            <td class="py-3 px-4 font-semibold text-ink">{{ row.name }}</td>
            <td class="py-3 px-4"><StatusBadge :status="row.status" /></td>
          </tr>
        </tbody>
      </table>
    </div>
    <p class="text-[11px] text-ink-secondary">
      {{ $t('seoAdmin.note') }}
      <span class="ml-1 inline-block text-[10px] font-bold px-2 py-0.5 rounded bg-amber-500/10 text-amber-700 dark:text-amber-300 border border-amber-500/30">{{ $t('common.demoData') }}</span>
    </p>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'admin',
  middleware: 'admin',
})

const { t } = useI18n()
useHead(() => ({ title: t('seoAdmin.title') }))

// TODO(api): replace with SEO audit endpoints when implemented.
const rows = [
  { id: 'seo-1', name: 'hreflang x-default present', status: 'PASS' },
  { id: 'seo-2', name: 'Canonical per locale', status: 'PASS' },
  { id: 'seo-3', name: 'Sitemap includes /compare', status: 'FAIL' },
  { id: 'seo-4', name: 'OG images 1200×630', status: 'WARN' },
]
</script>
