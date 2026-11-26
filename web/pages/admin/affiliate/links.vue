<template>
  <div class="space-y-8">
    <PageHeader :title="$t('affLinks.h1')" :subtitle="$t('affLinks.lede')" />

    <div class="card p-0 overflow-hidden">
      <table class="table-shell">
        <thead>
          <tr>
            <th class="py-3 px-4 text-left">{{ $t('affLinks.colSlug') }}</th>
            <th class="py-3 px-4 text-left">{{ $t('affLinks.colTarget') }}</th>
            <th class="py-3 px-4 text-right">{{ $t('affLinks.colClicks') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="row in rows" :key="row.id" class="border-t border-subtle text-sm">
            <td class="py-3 px-4 font-mono text-xs text-ink">/go/{{ row.slug }}</td>
            <td class="py-3 px-4 text-ink-secondary">{{ row.target }}</td>
            <td class="py-3 px-4 text-right tabular-nums">{{ row.clicks }}</td>
          </tr>
        </tbody>
      </table>
    </div>
    <p class="text-[11px] text-ink-secondary">
      {{ $t('affLinks.note') }}
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
useHead(() => ({ title: t('affLinks.title') }))

// TODO(api): replace with affiliate link endpoints when implemented.
const rows = [
  { id: 'al-1', slug: 'notion-pro', target: 'Notion Plus annual', clicks: 1842 },
  { id: 'al-2', slug: 'figma-team', target: 'Figma Professional', clicks: 967 },
  { id: 'al-3', slug: 'linear-biz', target: 'Linear Business', clicks: 423 },
]
</script>
