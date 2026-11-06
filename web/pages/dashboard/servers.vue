<template>
  <div class="space-y-8">
    <div>
      <h2 class="text-2xl font-extrabold text-white tracking-tight">Server Inventory</h2>
      <p class="text-sm text-slate-400 mt-1">Unified view of servers synced from your connected clouds.</p>
    </div>

    <div class="flex items-center gap-2 text-xs">
      <span class="text-slate-500">Connection:</span>
      <select
        v-model="connectionFilter"
        class="bg-slate-900 border border-slate-700 rounded px-2.5 py-1.5 text-slate-200 focus:outline-none"
        @change="load"
      >
        <option value="">All connections</option>
        <option v-for="c in connections" :key="c.id" :value="c.id">{{ c.connection_name }}</option>
      </select>
    </div>

    <div v-if="loading" class="space-y-3" aria-busy="true" aria-label="Loading servers">
      <div v-for="n in 4" :key="n" class="p-5 rounded-xl bg-slate-900 border border-slate-800 animate-pulse">
        <div class="h-4 w-1/2 rounded bg-slate-800" />
      </div>
    </div>

    <div v-else-if="error" class="p-8 rounded-xl bg-red-950/20 border border-red-900 text-center space-y-3">
      <h3 class="text-base font-bold text-white">{{ error.title }}</h3>
      <p class="text-xs text-slate-400">{{ error.message }}</p>
      <button class="px-4 py-2 rounded bg-sky-600 hover:bg-sky-500 text-white text-xs font-semibold" @click="load">
        Retry
      </button>
    </div>

    <div v-else-if="resources.length > 0" class="overflow-x-auto rounded-xl border border-slate-800">
      <table class="w-full text-left text-xs bg-slate-900 divide-y divide-slate-800">
        <thead class="bg-slate-950 text-slate-400 font-semibold uppercase tracking-wider">
          <tr>
            <th class="py-3 px-4">Hostname</th>
            <th class="py-3 px-4">Provider</th>
            <th class="py-3 px-4">Specs</th>
            <th class="py-3 px-4">Status</th>
            <th class="py-3 px-4 text-right">Est. Monthly</th>
            <th class="py-3 px-4">Last Sync</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-800 text-slate-300 tabular-nums">
          <tr v-for="r in resources" :key="r.id" class="hover:bg-slate-800/40 transition">
            <td class="py-3.5 px-4 font-bold text-white">{{ r.hostname }}</td>
            <td class="py-3.5 px-4">{{ r.provider_slug }}</td>
            <td class="py-3.5 px-4">{{ r.vcpu_count }} vCPU · {{ (r.ram_mb / 1024).toFixed(1) }} GB · {{ r.disk_gb }} GB</td>
            <td class="py-3.5 px-4">
              <span
                class="px-2 py-0.5 rounded-full text-[10px] font-semibold border"
                :class="r.power_status === 'RUNNING' ? 'bg-emerald-950 text-emerald-400 border-emerald-800' : 'bg-slate-800 text-slate-300 border-slate-700'"
              >
                {{ r.power_status }}
              </span>
            </td>
            <td class="py-3.5 px-4 text-right font-bold text-white">${{ Number(r.estimated_monthly_cost_usd).toFixed(2) }}</td>
            <td class="py-3.5 px-4 text-slate-500">{{ formatDate(r.last_synced_at) }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-else class="text-center py-16 bg-slate-900/40 rounded-xl border border-slate-800 space-y-3">
      <h4 class="text-base font-bold text-white">No servers synced yet</h4>
      <p class="text-xs text-slate-400 max-w-md mx-auto">
        Inventory appears after you connect a cloud account and the provider sync runs.
        Power actions, resize, and deletion gates operate on synced servers only.
      </p>
      <NuxtLink
        to="/dashboard/connections"
        class="inline-block px-4 py-2 rounded bg-sky-600 hover:bg-sky-500 text-white text-xs font-semibold transition"
      >
        Connect a Cloud
      </NuxtLink>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ApiError, useApi } from '~/composables/useApi'
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

const { request } = useApi()

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

onMounted(() => {
  load()
})
</script>
