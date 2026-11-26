<template>
  <div class="space-y-8">
    <PageHeader :title="$t('auditVerifyAdmin.h1')" :subtitle="$t('auditVerifyAdmin.lede')" />

    <div class="card p-6 space-y-3">
      <div v-for="check in checks" :key="check.id" class="flex items-center justify-between gap-3 text-sm">
        <span class="text-ink">{{ check.name }}</span>
        <StatusBadge :status="check.status" />
      </div>
      <div class="pt-2">
        <NuxtLink :to="$localePath('/admin/audit-logs')" class="text-xs font-semibold text-brand hover:underline">
          {{ $t('auditVerifyAdmin.openLogs') }}
        </NuxtLink>
      </div>
    </div>
    <p class="text-[11px] text-ink-secondary">
      {{ $t('auditVerifyAdmin.note') }}
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
useHead(() => ({ title: t('auditVerifyAdmin.title') }))

// TODO(api): replace with audit verification endpoints when implemented.
const checks = [
  { id: 'av-1', name: 'Hash chain intact (last 1,000 entries)', status: 'PASS' },
  { id: 'av-2', name: 'Payout batches dual-approved', status: 'PASS' },
  { id: 'av-3', name: 'Log retention ≥ 365 days', status: 'WARN' },
]
</script>
