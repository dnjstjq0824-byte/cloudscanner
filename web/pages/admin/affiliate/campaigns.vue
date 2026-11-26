<template>
  <div class="space-y-8">
    <PageHeader :title="$t('affCampaigns.h1')" :subtitle="$t('affCampaigns.lede')" />

    <div class="card p-0 overflow-hidden">
      <table class="table-shell">
        <thead>
          <tr>
            <th class="py-3 px-4 text-left">{{ $t('affCampaigns.colCampaign') }}</th>
            <th class="py-3 px-4 text-left">{{ $t('affCampaigns.colNetwork') }}</th>
            <th class="py-3 px-4 text-left">{{ $t('affCampaigns.colStatus') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="row in rows" :key="row.id" class="border-t border-subtle text-sm">
            <td class="py-3 px-4 font-semibold text-ink">{{ row.name }}</td>
            <td class="py-3 px-4 text-ink-secondary">{{ row.network }}</td>
            <td class="py-3 px-4"><StatusBadge :status="row.status" /></td>
          </tr>
        </tbody>
      </table>
    </div>
    <p class="text-[11px] text-ink-secondary">
      {{ $t('affCampaigns.note') }}
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
useHead(() => ({ title: t('affCampaigns.title') }))

// TODO(api): replace with affiliate campaign endpoints when implemented.
const rows = [
  { id: 'ac-1', name: 'Q4 SaaS push', network: 'Impact Radius', status: 'ACTIVE' },
  { id: 'ac-2', name: 'AI tools launch', network: 'PartnerStack', status: 'ACTIVE' },
  { id: 'ac-3', name: 'DevTools winter', network: 'CJ Affiliate', status: 'PAUSED' },
]
</script>
