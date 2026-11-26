<template>
  <div class="space-y-8">
    <PageHeader :title="$t('failuresAdmin.h1')" :subtitle="$t('failuresAdmin.lede')" />

    <div class="card p-0 overflow-hidden">
      <table class="table-shell">
        <thead>
          <tr>
            <th class="py-3 px-4 text-left">{{ $t('failuresAdmin.colFeed') }}</th>
            <th class="py-3 px-4 text-left">{{ $t('failuresAdmin.colError') }}</th>
            <th class="py-3 px-4 text-left">{{ $t('failuresAdmin.colStatus') }}</th>
            <th class="py-3 px-4 text-right">{{ $t('failuresAdmin.colAction') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="row in rows" :key="row.id" class="border-t border-subtle text-sm">
            <td class="py-3 px-4 font-semibold text-ink">{{ row.name }}</td>
            <td class="py-3 px-4 font-mono text-xs text-ink-secondary">{{ row.error }}</td>
            <td class="py-3 px-4"><StatusBadge :status="row.status" /></td>
            <td class="py-3 px-4 text-right">
              <button class="btn-ghost !py-1 !px-2.5 !text-[11px]" @click="resolve(row.id)">
                {{ $t('failuresAdmin.resolve') }}
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <p class="text-[11px] text-ink-secondary">
      {{ $t('failuresAdmin.note') }}
      <span class="ml-1 inline-block text-[10px] font-bold px-2 py-0.5 rounded bg-amber-500/10 text-amber-700 dark:text-amber-300 border border-amber-500/30">{{ $t('common.demoData') }}</span>
    </p>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useToast } from '~/composables/useToast'

definePageMeta({
  layout: 'admin',
  middleware: 'admin',
})

const { t } = useI18n()
useHead(() => ({ title: t('failuresAdmin.title') }))
const toast = useToast()

// TODO(api): replace with pricing quarantine endpoints when implemented.
const rows = ref([
  { id: 'qf-1', name: 'Acme CRM pricing page', error: 'HTTP 429 rate limited', status: 'QUARANTINED' },
  { id: 'qf-2', name: 'Beta PM changelog feed', error: 'Schema drift: missing field', status: 'QUARANTINED' },
  { id: 'qf-3', name: 'Gamma CRM API', error: 'Auth expired', status: 'FAILED' },
])

function resolve(id: string) {
  rows.value = rows.value.filter((r) => r.id !== id)
  toast.success(t('failuresAdmin.toastResolved'))
}
</script>
