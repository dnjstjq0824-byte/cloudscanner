<template>
  <div class="space-y-8">
    <PageHeader :title="$t('deployAdmin.h1')" :subtitle="$t('deployAdmin.lede')" />

    <div class="card p-0 overflow-hidden">
      <table class="table-shell">
        <thead>
          <tr>
            <th class="py-3 px-4 text-left">{{ $t('deployAdmin.colRelease') }}</th>
            <th class="py-3 px-4 text-left">{{ $t('deployAdmin.colTarget') }}</th>
            <th class="py-3 px-4 text-left">{{ $t('deployAdmin.colStatus') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="row in rows" :key="row.id" class="border-t border-subtle text-sm">
            <td class="py-3 px-4 font-mono text-xs text-ink">{{ row.name }}</td>
            <td class="py-3 px-4 text-ink-secondary">{{ row.target }}</td>
            <td class="py-3 px-4"><StatusBadge :status="row.status" /></td>
          </tr>
        </tbody>
      </table>
    </div>
    <p class="text-[11px] text-ink-secondary">
      {{ $t('deployAdmin.note') }}
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
useHead(() => ({ title: t('deployAdmin.title') }))

// TODO(api): replace with deployment endpoints when implemented.
const rows = [
  { id: 'dp-1', name: 'web@0a91223', target: 'Pages production', status: 'LIVE' },
  { id: 'dp-2', name: 'api@0a91223', target: 'API container', status: 'LIVE' },
  { id: 'dp-3', name: 'web@ebdcf81', target: 'Pages preview', status: 'SUPERSEDED' },
]
</script>
