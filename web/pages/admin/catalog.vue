<template>
  <div class="space-y-8">
    <PageHeader :title="$t('catalogAdmin.h1')" :subtitle="$t('catalogAdmin.lede')" />

    <div class="card p-0 overflow-hidden">
      <table class="table-shell">
        <thead>
          <tr>
            <th class="py-3 px-4 text-left">{{ $t('catalogAdmin.colTool') }}</th>
            <th class="py-3 px-4 text-left">{{ $t('catalogAdmin.colAction') }}</th>
            <th class="py-3 px-4 text-left">{{ $t('catalogAdmin.colStatus') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="row in rows" :key="row.id" class="border-t border-subtle text-sm">
            <td class="py-3 px-4 font-semibold text-ink">{{ row.name }}</td>
            <td class="py-3 px-4 text-ink-secondary">{{ row.action }}</td>
            <td class="py-3 px-4"><StatusBadge :status="row.status" /></td>
          </tr>
        </tbody>
      </table>
    </div>
    <p class="text-[11px] text-ink-secondary">
      {{ $t('catalogAdmin.note') }}
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
useHead(() => ({ title: t('catalogAdmin.title') }))

// TODO(api): replace with catalog moderation endpoints when implemented.
const rows = [
  { id: 'cm-1', name: 'Notion', action: 'New tool submission', status: 'PENDING' },
  { id: 'cm-2', name: 'Linear', status: 'VERIFIED', action: 'Plan refresh verified' },
  { id: 'cm-3', name: 'Acme CRM', status: 'REJECTED', action: 'Duplicate of existing listing' },
]
</script>
