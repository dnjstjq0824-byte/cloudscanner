<template>
  <div class="space-y-8">
    <PageHeader :title="$t('affAccounts.h1')" :subtitle="$t('affAccounts.lede')" />

    <div class="card p-0 overflow-hidden">
      <table class="table-shell">
        <thead>
          <tr>
            <th class="py-3 px-4 text-left">{{ $t('affAccounts.colNetwork') }}</th>
            <th class="py-3 px-4 text-left">{{ $t('affAccounts.colAccount') }}</th>
            <th class="py-3 px-4 text-left">{{ $t('affAccounts.colStatus') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="row in rows" :key="row.id" class="border-t border-subtle text-sm">
            <td class="py-3 px-4 font-semibold text-ink">{{ row.name }}</td>
            <td class="py-3 px-4 font-mono text-xs text-ink-secondary">{{ row.account }}</td>
            <td class="py-3 px-4"><StatusBadge :status="row.status" /></td>
          </tr>
        </tbody>
      </table>
    </div>
    <p class="text-[11px] text-ink-secondary">
      {{ $t('affAccounts.note') }}
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
useHead(() => ({ title: t('affAccounts.title') }))

// TODO(api): replace with affiliate account endpoints when implemented.
const rows = [
  { id: 'aa-1', name: 'Impact Radius', account: 'tirp-prod-01', status: 'CONNECTED' },
  { id: 'aa-2', name: 'PartnerStack', account: 'tirp-prod-02', status: 'CONNECTED' },
  { id: 'aa-3', name: 'CJ Affiliate', account: 'tirp-review-01', status: 'PENDING' },
]
</script>
