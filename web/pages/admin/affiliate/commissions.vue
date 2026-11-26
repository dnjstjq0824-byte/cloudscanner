<template>
  <div class="space-y-8">
    <PageHeader :title="$t('affCommissions.h1')" :subtitle="$t('affCommissions.lede')" />

    <div class="card p-0 overflow-hidden">
      <table class="table-shell">
        <thead>
          <tr>
            <th class="py-3 px-4 text-left">{{ $t('affCommissions.colPeriod') }}</th>
            <th class="py-3 px-4 text-right">{{ $t('affCommissions.colEarned') }}</th>
            <th class="py-3 px-4 text-left">{{ $t('affCommissions.colStatus') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="row in rows" :key="row.id" class="border-t border-subtle text-sm">
            <td class="py-3 px-4 font-semibold text-ink">{{ row.period }}</td>
            <td class="py-3 px-4 text-right font-bold text-brand tabular-nums">
              <PriceDisplay :amount="row.amount" currency="USD" />
            </td>
            <td class="py-3 px-4"><StatusBadge :status="row.status" /></td>
          </tr>
        </tbody>
      </table>
    </div>
    <p class="text-[11px] text-ink-secondary">
      {{ $t('affCommissions.note') }}
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
useHead(() => ({ title: t('affCommissions.title') }))

// TODO(api): replace with commission ledger endpoints when implemented.
const rows = [
  { id: 'cm-1', period: '2026-11', amount: 1842.5, status: 'ACCRUED' },
  { id: 'cm-2', period: '2026-10', amount: 2310.0, status: 'PAID' },
  { id: 'cm-3', period: '2026-09', amount: 1975.25, status: 'PAID' },
]
</script>
