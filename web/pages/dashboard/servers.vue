<template>
  <div class="space-y-8">
    <PageHeader title="Server Inventory" subtitle="Unified view of servers synced from your connected clouds." />

    <div class="flex items-center gap-2 text-xs">
      <span class="text-ink-secondary">Connection:</span>
      <select v-model="connectionFilter" class="input !w-auto !py-1.5" @change="load">
        <option value="">All connections</option>
        <option v-for="c in connections" :key="c.id" :value="c.id">{{ c.connection_name }}</option>
      </select>
    </div>

    <LoadingSkeleton v-if="loading" :rows="4" label="Loading servers" />

    <ErrorState v-else-if="error" :title="error.title" :message="error.message" @retry="load" />

    <div v-else-if="resources.length > 0" class="card overflow-hidden">
      <div class="overflow-x-auto">
        <table class="table-shell min-w-[760px]">
          <thead>
            <tr>
              <th class="py-3 px-4">Hostname</th>
              <th class="py-3 px-4">Provider</th>
              <th class="py-3 px-4">Specs</th>
              <th class="py-3 px-4">Status</th>
              <th class="py-3 px-4 text-right">Est. Monthly</th>
              <th class="py-3 px-4">Last Sync</th>
              <th class="py-3 px-4"><span class="sr-only">Actions</span></th>
            </tr>
          </thead>
          <tbody class="tabular-nums">
            <tr v-for="r in resources" :key="r.id" class="hover:bg-elevated/60 transition">
              <td class="py-3.5 px-4 font-bold">{{ r.hostname }}</td>
              <td class="py-3.5 px-4">{{ r.provider_slug }}</td>
              <td class="py-3.5 px-4 whitespace-nowrap">{{ r.vcpu_count }} vCPU · {{ (r.ram_mb / 1024).toFixed(1) }} GB · {{ r.disk_gb }} GB</td>
              <td class="py-3.5 px-4">
                <StatusBadge :status="r.power_status" />
              </td>
              <td class="py-3.5 px-4 text-right font-bold">${{ Number(r.estimated_monthly_cost_usd).toFixed(2) }}</td>
              <td class="py-3.5 px-4 text-ink-secondary whitespace-nowrap">{{ formatDate(r.last_synced_at) }}</td>
              <td class="py-3.5 px-4 text-right">
                <button class="btn-danger !py-1 !px-2.5 !text-[11px]" @click="askDelete(r)">Delete</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <EmptyState
      v-else
      title="No servers synced yet"
      message="Inventory appears after you connect a cloud account and provision or sync servers. Deletion gates operate on managed servers only."
      cta-to="/dashboard/connections"
      cta-label="Connect a Cloud"
    />

    <ConfirmModal
      :open="deleteOpen"
      title="Delete this server?"
      :message="`Type the hostname to confirm deletion of '${target?.hostname}'. Step-up verification is required.`"
      confirm-label="Delete Server"
      danger
      :busy="deleting"
      :error="deleteError"
      @confirm="remove"
      @cancel="closeDelete"
    >
      <div class="space-y-3">
        <input v-model="confirmHostname" type="text" placeholder="hostname" autocomplete="off" class="input" />
        <label class="flex items-start gap-2 text-xs text-ink-secondary cursor-pointer">
          <input v-model="mfaConfirmed" type="checkbox" class="mt-0.5 rounded accent-sky-600" />
          <span>I completed step-up verification for this destructive operation.</span>
        </label>
      </div>
    </ConfirmModal>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ApiError, useApi } from '~/composables/useApi'
import { useToast } from '~/composables/useToast'
import type { CloudConnection, CloudResource, ListEnvelope } from '~/types/api'

definePageMeta({
  layout: 'dashboard',
  middleware: 'auth',
})

useHead({ title: 'Server Inventory' })

const loading = ref(false)
const error = ref<ApiError | null>(null)
const resources = ref<CloudResource[]>([])
const connections = ref<CloudConnection[]>([])
const connectionFilter = ref('')
const deleteOpen = ref(false)
const deleting = ref(false)
const deleteError = ref('')
const target = ref<CloudResource | null>(null)
const confirmHostname = ref('')
const mfaConfirmed = ref(false)

const { request } = useApi()
const toast = useToast()

function formatDate(iso: string): string {
  try {
    return new Date(iso).toLocaleString()
  } catch {
    return iso
  }
}

async function load() {
  loading.value = true
  error.value = null
  try {
    const [res, conns] = await Promise.all([
      request<ListEnvelope<CloudResource>>('/api/v1/cloud/resources', {
        query: { connection_id: connectionFilter.value || undefined },
      }),
      connections.value.length > 0
        ? Promise.resolve(null)
        : request<ListEnvelope<CloudConnection>>('/api/v1/cloud/connections'),
    ])
    resources.value = res.data ?? []
    if (conns) {
      connections.value = conns.data ?? []
    }
  } catch (e) {
    error.value = e instanceof ApiError ? e : ApiError.fromUnknown(e)
  } finally {
    loading.value = false
  }
}

function askDelete(r: CloudResource) {
  target.value = r
  confirmHostname.value = ''
  mfaConfirmed.value = false
  deleteError.value = ''
  deleteOpen.value = true
}

function closeDelete() {
  deleteOpen.value = false
  target.value = null
}

async function remove() {
  if (!target.value) {
    return
  }
  deleting.value = true
  deleteError.value = ''
  try {
    await request(`/api/v1/cloud/servers/${target.value.id}`, {
      method: 'DELETE',
      body: { confirmation_hostname: confirmHostname.value.trim(), mfa_verified: mfaConfirmed.value },
    })
    toast.success(`Server '${target.value.hostname}' deleted.`)
    closeDelete()
    await load()
  } catch (e) {
    const apiErr = e instanceof ApiError ? e : ApiError.fromUnknown(e)
    deleteError.value = apiErr.message
  } finally {
    deleting.value = false
  }
}

onMounted(() => {
  load()
})
</script>
