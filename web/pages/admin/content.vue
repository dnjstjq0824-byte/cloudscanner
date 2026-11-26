<template>
  <div class="space-y-8">
    <PageHeader :title="$t('contentAdmin.h1')" :subtitle="$t('contentAdmin.lede')" />

    <div class="card p-0 overflow-hidden">
      <table class="table-shell">
        <thead>
          <tr>
            <th class="py-3 px-4 text-left">{{ $t('contentAdmin.colTitle') }}</th>
            <th class="py-3 px-4 text-left">{{ $t('contentAdmin.colType') }}</th>
            <th class="py-3 px-4 text-left">{{ $t('contentAdmin.colStatus') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="row in rows" :key="row.id" class="border-t border-subtle text-sm">
            <td class="py-3 px-4 font-semibold text-ink">{{ row.name }}</td>
            <td class="py-3 px-4 text-ink-secondary">{{ row.kind }}</td>
            <td class="py-3 px-4"><StatusBadge :status="row.status" /></td>
          </tr>
        </tbody>
      </table>
    </div>
    <p class="text-[11px] text-ink-secondary">
      {{ $t('contentAdmin.note') }}
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
useHead(() => ({ title: t('contentAdmin.title') }))

// TODO(api): replace with content endpoints when implemented.
const rows = [
  { id: 'ct-1', name: 'How to compare CRM plans', kind: 'Tutorial', status: 'PUBLISHED' },
  { id: 'ct-2', name: 'AI video editor deals, ranked', kind: 'Review', status: 'PUBLISHED' },
  { id: 'ct-3', name: 'Annual vs monthly billing', kind: 'Tutorial', status: 'DRAFT' },
]
</script>
