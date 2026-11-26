<template>
  <div class="space-y-8">
    <PageHeader :title="$t('permsAdmin.h1')" :subtitle="$t('permsAdmin.lede')" />

    <div class="card p-0 overflow-hidden">
      <table class="table-shell min-w-[560px]">
        <thead>
          <tr>
            <th class="py-3 px-4 text-left">{{ $t('permsAdmin.colRole') }}</th>
            <th class="py-3 px-4 text-left">{{ $t('permsAdmin.colCurate') }}</th>
            <th class="py-3 px-4 text-left">{{ $t('permsAdmin.colPayout') }}</th>
            <th class="py-3 px-4 text-left">{{ $t('permsAdmin.colMembers') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="row in matrix" :key="row.role" class="border-t border-subtle text-sm">
            <td class="py-3 px-4 font-mono text-xs font-bold text-ink">{{ row.role }}</td>
            <td class="py-3 px-4"><StatusBadge :status="row.curate" /></td>
            <td class="py-3 px-4"><StatusBadge :status="row.payout" /></td>
            <td class="py-3 px-4"><StatusBadge :status="row.members" /></td>
          </tr>
        </tbody>
      </table>
    </div>
    <p class="text-[11px] text-ink-secondary">
      {{ $t('permsAdmin.note') }}
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
useHead(() => ({ title: t('permsAdmin.title') }))

// TODO(api): replace with IAM policy endpoints when implemented. Read-only matrix.
const matrix = [
  { role: 'ORG_OWNER', curate: 'ALLOW', payout: 'APPROVE', members: 'MANAGE' },
  { role: 'ORG_ADMIN', curate: 'ALLOW', payout: 'REQUEST', members: 'MANAGE' },
  { role: 'CATALOG_OPERATOR', curate: 'ALLOW', payout: 'DENY', members: 'VIEW' },
  { role: 'PAYOUT_APPROVER', curate: 'VIEW', payout: 'APPROVE', members: 'VIEW' },
  { role: 'ORG_VIEWER', curate: 'VIEW', payout: 'DENY', members: 'VIEW' },
]
</script>
