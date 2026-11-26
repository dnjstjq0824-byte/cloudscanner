<template>
  <div class="space-y-8">
    <PageHeader :title="$t('seclogs.h1')" :subtitle="$t('seclogs.lede')" />

    <div class="card p-0 overflow-hidden">
      <table class="table-shell">
        <thead>
          <tr>
            <th class="py-3 px-4 text-left">{{ $t('seclogs.colTime') }}</th>
            <th class="py-3 px-4 text-left">{{ $t('seclogs.colEvent') }}</th>
            <th class="py-3 px-4 text-left">{{ $t('seclogs.colIp') }}</th>
            <th class="py-3 px-4 text-left">{{ $t('seclogs.colStatus') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="log in logs" :key="log.id" class="border-t border-subtle text-sm">
            <td class="py-3 px-4 text-ink-secondary tabular-nums">{{ log.at }}</td>
            <td class="py-3 px-4 font-mono text-xs text-ink">{{ log.event }}</td>
            <td class="py-3 px-4 font-mono text-xs text-ink-secondary">{{ log.ip }}</td>
            <td class="py-3 px-4"><StatusBadge :status="log.status" /></td>
          </tr>
        </tbody>
      </table>
    </div>
    <p class="text-[11px] text-ink-secondary">
      {{ $t('seclogs.note') }}
      <span class="ml-1 inline-block text-[10px] font-bold px-2 py-0.5 rounded bg-amber-500/10 text-amber-700 dark:text-amber-300 border border-amber-500/30">{{ $t('common.demoData') }}</span>
    </p>
  </div>
</template>

<script setup lang="ts">
import { mockSecurityLogs } from '~/utils/mock-data'

definePageMeta({
  layout: 'dashboard',
  middleware: 'auth',
})

const { t } = useI18n()
useHead(() => ({ title: t('seclogs.title') }))

// TODO(api): replace with GET /api/v1/users/me/security-logs when implemented.
const logs = mockSecurityLogs()
</script>
