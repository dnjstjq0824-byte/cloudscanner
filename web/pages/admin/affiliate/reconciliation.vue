<template>
  <div class="space-y-8">
    <PageHeader :title="$t('affRecon.h1')" :subtitle="$t('affRecon.lede')" />

    <div class="card p-0 overflow-hidden">
      <table class="table-shell">
        <thead>
          <tr>
            <th class="py-3 px-4 text-left">{{ $t('affRecon.colBatch') }}</th>
            <th class="py-3 px-4 text-right">{{ $t('affRecon.colMatched') }}</th>
            <th class="py-3 px-4 text-right">{{ $t('affRecon.colOrphan') }}</th>
            <th class="py-3 px-4 text-left">{{ $t('affRecon.colStatus') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="row in rows" :key="row.id" class="border-t border-subtle text-sm">
            <td class="py-3 px-4 font-semibold text-ink">{{ row.batch }}</td>
            <td class="py-3 px-4 text-right tabular-nums">{{ row.matched }}</td>
            <td class="py-3 px-4 text-right tabular-nums">{{ row.orphan }}</td>
            <td class="py-3 px-4"><StatusBadge :status="row.status" /></td>
          </tr>
        </tbody>
      </table>
    </div>
    <p class="text-[11px] text-ink-secondary">
      {{ $t('affRecon.note') }}
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
useHead(() => ({ title: t('affRecon.title') }))

// TODO(api): replace with reconciliation endpoints when implemented.
const rows = [
  { id: 'rc-1', batch: '2026-11-25 Impact', matched: 412, orphan: 3, status: 'CLEAR' },
  { id: 'rc-2', batch: '2026-11-24 PartnerStack', matched: 188, orphan: 11, status: 'REVIEW' },
  { id: 'rc-3', batch: '2026-11-23 CJ', matched: 96, orphan: 0, status: 'CLEAR' },
]
</script>
