<template>
  <div class="space-y-8">
    <PageHeader :title="$t('secEventsAdmin.h1')" :subtitle="$t('secEventsAdmin.lede')" />

    <div class="card p-0 overflow-hidden">
      <table class="table-shell">
        <thead>
          <tr>
            <th class="py-3 px-4 text-left">{{ $t('secEventsAdmin.colTime') }}</th>
            <th class="py-3 px-4 text-left">{{ $t('secEventsAdmin.colEvent') }}</th>
            <th class="py-3 px-4 text-left">{{ $t('secEventsAdmin.colSeverity') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="row in rows" :key="row.id" class="border-t border-subtle text-sm">
            <td class="py-3 px-4 text-ink-secondary tabular-nums">{{ row.at }}</td>
            <td class="py-3 px-4 font-mono text-xs text-ink">{{ row.event }}</td>
            <td class="py-3 px-4"><StatusBadge :status="row.severity" /></td>
          </tr>
        </tbody>
      </table>
    </div>
    <p class="text-[11px] text-ink-secondary">
      {{ $t('secEventsAdmin.note') }}
      <span class="ml-1 inline-block text-[10px] font-bold px-2 py-0.5 rounded bg-amber-500/10 text-amber-700 dark:text-amber-300 border border-amber-500/30">{{ $t('common.demoData') }}</span>
    </p>
  </div>
</template>

<script setup lang="ts">
import { mockRows } from '~/utils/mock-data'

definePageMeta({
  layout: 'admin',
  middleware: 'admin',
})

const { t } = useI18n()
useHead(() => ({ title: t('secEventsAdmin.title') }))

// TODO(api): replace with security event endpoints when implemented.
const rows = mockRows('se', ['idor.blocked', 'rate-limit.tripped', 'mfa.failed', 'csrf.rejected'], ['P1', 'P2', 'P2', 'P3']).map((r, i) => ({
  id: r.id,
  at: `2026-11-${26 - i} 0${8 + i}:1${i}`,
  event: r.name,
  severity: r.status,
}))
</script>
