<template>
  <div class="space-y-8">
    <PageHeader :title="$t('syncAdmin.h1')" :subtitle="$t('syncAdmin.lede')" />

    <div class="card p-0 overflow-hidden">
      <table class="table-shell">
        <thead>
          <tr>
            <th class="py-3 px-4 text-left">{{ $t('syncAdmin.colVendor') }}</th>
            <th class="py-3 px-4 text-left">{{ $t('syncAdmin.colLastSync') }}</th>
            <th class="py-3 px-4 text-left">{{ $t('syncAdmin.colStatus') }}</th>
            <th class="py-3 px-4 text-right">{{ $t('syncAdmin.colAction') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="job in jobs" :key="job.id" class="border-t border-subtle text-sm">
            <td class="py-3 px-4 font-semibold text-ink">{{ job.name }}</td>
            <td class="py-3 px-4 text-ink-secondary tabular-nums">{{ job.updated }}</td>
            <td class="py-3 px-4"><StatusBadge :status="job.status" /></td>
            <td class="py-3 px-4 text-right">
              <button class="btn-ghost !py-1 !px-2.5 !text-[11px]" @click="trigger(job.name)">
                {{ $t('syncAdmin.trigger') }}
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <p class="text-[11px] text-ink-secondary">
      {{ $t('syncAdmin.note') }}
      <span class="ml-1 inline-block text-[10px] font-bold px-2 py-0.5 rounded bg-amber-500/10 text-amber-700 dark:text-amber-300 border border-amber-500/30">{{ $t('common.demoData') }}</span>
    </p>
  </div>
</template>

<script setup lang="ts">
import { mockRows } from '~/utils/mock-data'
import { useToast } from '~/composables/useToast'

definePageMeta({
  layout: 'admin',
  middleware: 'admin',
})

const { t } = useI18n()
useHead(() => ({ title: t('syncAdmin.title') }))
const toast = useToast()

// TODO(api): replace with pricing sync endpoints when implemented.
const jobs = mockRows('sync', ['Notion', 'Figma', 'Linear', 'Vercel'], ['SYNCED', 'SYNCED', 'SYNCING', 'FAILED'])

function trigger(name: string) {
  toast.success(t('syncAdmin.toastTriggered', { name }))
}
</script>
