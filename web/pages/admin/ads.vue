<template>
  <div class="space-y-8">
    <PageHeader :title="$t('adsAdmin.h1')" :subtitle="$t('adsAdmin.lede')" />

    <div class="card p-0 overflow-hidden">
      <table class="table-shell">
        <thead>
          <tr>
            <th class="py-3 px-4 text-left">{{ $t('adsAdmin.colSlot') }}</th>
            <th class="py-3 px-4 text-left">{{ $t('adsAdmin.colSponsor') }}</th>
            <th class="py-3 px-4 text-left">{{ $t('adsAdmin.colStatus') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="row in rows" :key="row.id" class="border-t border-subtle text-sm">
            <td class="py-3 px-4 font-semibold text-ink">{{ row.name }}</td>
            <td class="py-3 px-4 text-ink-secondary">{{ row.sponsor }}</td>
            <td class="py-3 px-4"><StatusBadge :status="row.status" /></td>
          </tr>
        </tbody>
      </table>
    </div>
    <p class="text-[11px] text-ink-secondary">
      {{ $t('adsAdmin.note') }}
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
useHead(() => ({ title: t('adsAdmin.title') }))

// TODO(api): replace with ad-slot endpoints when implemented.
const rows = [
  { id: 'ad-1', name: 'Homepage hero', sponsor: '—', status: 'EMPTY' },
  { id: 'ad-2', name: 'Rankings sidebar', sponsor: 'Linear', status: 'LIVE' },
  { id: 'ad-3', name: 'Newsletter footer', sponsor: '—', status: 'EMPTY' },
]
</script>
