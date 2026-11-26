<template>
  <div class="space-y-8">
    <PageHeader :title="$t('admin.auditTitle')" :subtitle="$t('admin.auditSubtitle')" />

    <div class="flex flex-wrap items-center gap-2 text-xs">
      <span class="text-ink-secondary">{{ $t('admin.actionFilter') }}</span>
      <select v-model="actionFilter" class="input !w-auto !py-1.5" @change="reload">
        <option value="">{{ $t('admin.optAllActions') }}</option>
        <option v-for="a in actions" :key="a" :value="a">{{ a }}</option>
      </select>
      <span class="text-ink-secondary ml-auto">{{ $t('admin.entriesCount', { n: total }) }}</span>
    </div>

    <LoadingSkeleton v-if="loading" :rows="5" :label="$t('admin.loadingAudit')" />

    <ErrorState v-else-if="error" :title="error.title" :message="error.message" @retry="load" />

    <div v-else class="card overflow-hidden">
      <div class="overflow-x-auto">
        <table class="table-shell min-w-[760px]">
          <thead>
            <tr>
              <th class="py-3 px-4">{{ $t('admin.colTime') }}</th>
              <th class="py-3 px-4">{{ $t('admin.colAction') }}</th>
              <th class="py-3 px-4">{{ $t('admin.colActor') }}</th>
              <th class="py-3 px-4">{{ $t('admin.colTarget') }}</th>
              <th class="py-3 px-4">{{ $t('admin.colDetail') }}</th>
            </tr>
          </thead>
          <tbody class="tabular-nums">
            <tr v-for="e in entries" :key="e.id" class="hover:bg-elevated/60 transition">
              <td class="py-2.5 px-4 text-ink-secondary whitespace-nowrap">{{ formatDate(e.created_at) }}</td>
              <td class="py-2.5 px-4 font-mono text-[11px] font-semibold">{{ e.action }}</td>
              <td class="py-2.5 px-4 font-mono text-[11px] text-ink-secondary">{{ short(e.actor_user_id) }}</td>
              <td class="py-2.5 px-4 font-mono text-[11px] text-ink-secondary">{{ e.target_type }}:{{ short(e.target_id) }}</td>
              <td class="py-2.5 px-4 text-[11px] text-ink-secondary max-w-64 truncate" :title="e.detail">{{ e.detail || '—' }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <EmptyState
      v-if="!loading && !error && entries.length === 0"
      :title="$t('admin.emptyAuditTitle')"
      :message="$t('admin.emptyAuditBody')"
    />

    <div v-if="!loading && !error && entries.length > 0" class="flex items-center justify-between text-xs text-ink-secondary">
      <button class="btn-ghost !py-1.5 !text-xs" :disabled="offset === 0" @click="page(-1)">{{ $t('common.prev') }}</button>
      <span>{{ $t('admin.showingOf', { shown: entries.length, total }) }}</span>
      <button class="btn-ghost !py-1.5 !text-xs" :disabled="offset + entries.length >= total" @click="page(1)">{{ $t('common.next') }}</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ApiError, useApi } from '~/composables/useApi'
import type { AuditEntry } from '~/types/api'

definePageMeta({
  layout: 'admin',
  middleware: 'admin',
})

const { t, locale } = useI18n()
useHead(() => ({ title: t('admin.auditTitle') }))

const actions = [
  'server.provision', 'server.delete', 'connection.create', 'connection.delete',
  'approval.sign', 'budget.save', 'admin.user.status', 'admin.user.role',
  'admin.provider.status', 'admin.notice.create', 'admin.notice.update',
  'admin.notice.delete', 'admin.settings.update',
]

const loading = ref(false)
const error = ref<ApiError | null>(null)
const entries = ref<AuditEntry[]>([])
const total = ref(0)
const offset = ref(0)
const limit = 50
const actionFilter = ref('')

const { request } = useApi()

function short(id: string): string {
  return id && id.length > 8 ? `${id.slice(0, 8)}…` : id || '—'
}

function formatDate(iso: string): string {
  try {
    return new Date(iso).toLocaleString(locale.value)
  } catch {
    return iso
  }
}

async function load() {
  loading.value = true
  error.value = null
  try {
    const res = await request<{ data: AuditEntry[]; meta: { total: number } }>('/api/v1/admin/audit-logs', {
      query: { limit, offset: offset.value, action: actionFilter.value || undefined },
    })
    entries.value = res.data ?? []
    total.value = res.meta?.total ?? 0
  } catch (e) {
    error.value = e instanceof ApiError ? e : ApiError.fromUnknown(e)
  } finally {
    loading.value = false
  }
}

function reload() {
  offset.value = 0
  load()
}

function page(dir: number) {
  offset.value = Math.max(0, offset.value + dir * limit)
  load()
}

onMounted(() => {
  load()
})
</script>
