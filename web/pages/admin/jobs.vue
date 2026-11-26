<template>
  <div class="space-y-8">
    <PageHeader :title="$t('jobsAdmin.h1')" :subtitle="$t('jobsAdmin.lede')" />

    <div class="card p-0 overflow-hidden">
      <table class="table-shell">
        <thead>
          <tr>
            <th class="py-3 px-4 text-left">{{ $t('jobsAdmin.colJob') }}</th>
            <th class="py-3 px-4 text-left">{{ $t('jobsAdmin.colQueue') }}</th>
            <th class="py-3 px-4 text-left">{{ $t('jobsAdmin.colState') }}</th>
            <th class="py-3 px-4 text-right">{{ $t('jobsAdmin.colAction') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="job in jobs" :key="job.id" class="border-t border-subtle text-sm">
            <td class="py-3 px-4 font-mono text-xs text-ink">{{ job.id }}</td>
            <td class="py-3 px-4 text-ink-secondary">{{ job.queue }}</td>
            <td class="py-3 px-4"><StatusBadge :status="job.state" /></td>
            <td class="py-3 px-4 text-right">
              <button
                v-if="job.state === 'dead'"
                class="btn-ghost !py-1 !px-2.5 !text-[11px]"
                @click="retry(job.id)"
              >
                {{ $t('common.retry') }}
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <p class="text-[11px] text-ink-secondary">
      {{ $t('jobsAdmin.note') }}
      <span class="ml-1 inline-block text-[10px] font-bold px-2 py-0.5 rounded bg-amber-500/10 text-amber-700 dark:text-amber-300 border border-amber-500/30">{{ $t('common.demoData') }}</span>
    </p>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { mockJobs } from '~/utils/mock-data'
import { useToast } from '~/composables/useToast'

definePageMeta({
  layout: 'admin',
  middleware: 'admin',
})

const { t } = useI18n()
useHead(() => ({ title: t('jobsAdmin.title') }))
const toast = useToast()

// TODO(api): replace with job queue endpoints when implemented.
const jobs = ref(mockJobs())

function retry(id: string) {
  const found = jobs.value.find((j) => j.id === id)
  if (found) {
    found.state = 'queued'
  }
  toast.success(t('jobsAdmin.toastRetried'))
}
</script>
