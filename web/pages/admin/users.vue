<template>
  <div class="space-y-8">
    <PageHeader title="Members & Roles" subtitle="Suspend accounts or reassign organization roles. Both actions are audited." />

    <div class="flex flex-wrap items-center gap-2 text-xs">
      <span class="text-ink-secondary">Status:</span>
      <select v-model="statusFilter" class="input !w-auto !py-1.5" @change="load">
        <option value="">All</option>
        <option value="ACTIVE">Active</option>
        <option value="SUSPENDED">Suspended</option>
        <option value="PENDING_VERIFICATION">Pending</option>
      </select>
      <span class="text-ink-secondary ml-auto">{{ total }} members</span>
    </div>

    <LoadingSkeleton v-if="loading" :rows="4" label="Loading members" />

    <ErrorState v-else-if="error" :title="error.title" :message="error.message" @retry="load" />

    <div v-else class="card overflow-hidden">
      <div class="overflow-x-auto">
        <table class="table-shell min-w-[720px]">
          <thead>
            <tr>
              <th class="py-3 px-4">Member</th>
              <th class="py-3 px-4">Organization</th>
              <th class="py-3 px-4">Role</th>
              <th class="py-3 px-4">Status</th>
              <th class="py-3 px-4">Joined</th>
              <th class="py-3 px-4 text-right">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="u in users" :key="u.id" class="hover:bg-elevated/60 transition">
              <td class="py-3 px-4">
                <div class="font-bold">{{ u.full_name || '—' }}</div>
                <div class="text-ink-secondary text-[11px]">{{ u.email }}</div>
              </td>
              <td class="py-3 px-4 text-ink-secondary">{{ u.organization_name || '—' }}</td>
              <td class="py-3 px-4">
                <select
                  :value="u.role"
                  class="input !w-auto !py-1 !text-[11px]"
                  :disabled="saving === u.id"
                  aria-label="Change role"
                  @change="changeRole(u, ($event.target as HTMLSelectElement).value)"
                >
                  <option>OWNER</option>
                  <option>ADMIN</option>
                  <option>MEMBER</option>
                  <option>VIEWER</option>
                </select>
              </td>
              <td class="py-3 px-4"><StatusBadge :status="u.status" /></td>
              <td class="py-3 px-4 text-ink-secondary whitespace-nowrap">{{ formatDate(u.created_at) }}</td>
              <td class="py-3 px-4 text-right whitespace-nowrap">
                <button
                  v-if="u.status !== 'SUSPENDED'"
                  class="btn-danger !py-1 !px-2.5 !text-[11px]"
                  :disabled="saving === u.id"
                  @click="changeStatus(u, 'SUSPENDED')"
                >
                  Suspend
                </button>
                <button v-else class="btn-ghost !py-1 !px-2.5 !text-[11px]" :disabled="saving === u.id" @click="changeStatus(u, 'ACTIVE')">
                  Reactivate
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div v-if="!loading && !error" class="flex items-center justify-between text-xs text-ink-secondary">
      <button class="btn-ghost !py-1.5 !text-xs" :disabled="offset === 0" @click="page(-1)">← Prev</button>
      <span>Showing {{ users.length }} of {{ total }}</span>
      <button class="btn-ghost !py-1.5 !text-xs" :disabled="offset + users.length >= total" @click="page(1)">Next →</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ApiError, useApi } from '~/composables/useApi'
import { useToast } from '~/composables/useToast'
import type { AdminUser } from '~/types/api'

definePageMeta({
  layout: 'admin',
  middleware: 'admin',
})

useHead({ title: 'Members & Roles' })

const loading = ref(false)
const error = ref<ApiError | null>(null)
const users = ref<AdminUser[]>([])
const total = ref(0)
const offset = ref(0)
const limit = 25
const statusFilter = ref('')
const saving = ref('')

const { request } = useApi()
const toast = useToast()

function formatDate(iso: string): string {
  try {
    return new Date(iso).toLocaleDateString()
  } catch {
    return iso
  }
}

async function load() {
  loading.value = true
  error.value = null
  try {
    const res = await request<{ data: AdminUser[]; meta: { total: number } }>('/api/v1/admin/users', {
      query: { limit, offset: offset.value, status: statusFilter.value || undefined },
    })
    users.value = res.data ?? []
    total.value = res.meta?.total ?? 0
  } catch (e) {
    error.value = e instanceof ApiError ? e : ApiError.fromUnknown(e)
  } finally {
    loading.value = false
  }
}

function page(dir: number) {
  offset.value = Math.max(0, offset.value + dir * limit)
  load()
}

async function patch(u: AdminUser, body: Record<string, string>) {
  saving.value = u.id
  try {
    await request(`/api/v1/admin/users/${u.id}`, { method: 'PATCH', body })
    toast.success(`Member ${u.email} updated.`)
    await load()
  } catch (e) {
    const apiErr = e instanceof ApiError ? e : ApiError.fromUnknown(e)
    toast.error(apiErr.message)
  } finally {
    saving.value = ''
  }
}

function changeRole(u: AdminUser, role: string) {
  if (role !== u.role) {
    patch(u, { role })
  }
}

function changeStatus(u: AdminUser, status: string) {
  patch(u, { status })
}

onMounted(() => {
  load()
})
</script>
