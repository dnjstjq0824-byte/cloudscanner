<template>
  <div class="space-y-8">
    <div>
      <h2 class="text-2xl font-extrabold text-white tracking-tight">Overview</h2>
      <p class="text-sm text-slate-400 mt-1">
        {{ auth.org?.name || 'Your organization' }} · Role: {{ auth.role || '—' }}
      </p>
    </div>

    <div v-if="loading" class="grid grid-cols-1 md:grid-cols-3 gap-6" aria-busy="true" aria-label="Loading overview">
      <div v-for="n in 3" :key="n" class="p-6 rounded-xl bg-slate-900 border border-slate-800 animate-pulse">
        <div class="h-8 w-16 rounded bg-slate-800" />
        <div class="mt-3 h-4 w-2/3 rounded bg-slate-800" />
      </div>
    </div>

    <div v-else-if="error" class="p-8 rounded-xl bg-red-950/20 border border-red-900 text-center space-y-3">
      <h3 class="text-base font-bold text-white">{{ error.title }}</h3>
      <p class="text-xs text-slate-400">{{ error.message }}</p>
      <button class="px-4 py-2 rounded bg-sky-600 hover:bg-sky-500 text-white text-xs font-semibold" @click="load">
        Retry
      </button>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <div class="p-6 rounded-xl bg-slate-900 border border-slate-800">
        <div class="text-3xl font-bold text-sky-400 tabular-nums">{{ connections.length }}</div>
        <div class="text-xs text-slate-400 mt-1 uppercase tracking-wider">Connected Clouds</div>
        <NuxtLink to="/dashboard/connections" class="inline-block mt-4 text-xs font-semibold text-sky-400 hover:text-sky-300">
          Manage connections &rarr;
        </NuxtLink>
      </div>
      <div class="p-6 rounded-xl bg-slate-900 border border-slate-800">
        <div class="text-3xl font-bold text-sky-400 tabular-nums">{{ resources.length }}</div>
        <div class="text-xs text-slate-400 mt-1 uppercase tracking-wider">Managed Servers</div>
        <NuxtLink to="/dashboard/servers" class="inline-block mt-4 text-xs font-semibold text-sky-400 hover:text-sky-300">
          View inventory &rarr;
        </NuxtLink>
      </div>
      <div class="p-6 rounded-xl bg-slate-900 border border-slate-800">
        <div class="text-3xl font-bold text-sky-400 tabular-nums">
          {{ budget ? '$' + Number(budget.monthly_limit_usd).toFixed(0) : '—' }}
        </div>
        <div class="text-xs text-slate-400 mt-1 uppercase tracking-wider">Monthly Budget Cap</div>
        <NuxtLink to="/dashboard/budgets" class="inline-block mt-4 text-xs font-semibold text-sky-400 hover:text-sky-300">
          Configure budgets &rarr;
        </NuxtLink>
      </div>
    </div>

    <div v-if="!loading && !error" class="p-6 rounded-xl bg-slate-900 border border-slate-800">
      <h3 class="text-sm font-bold text-white">Getting started</h3>
      <ol class="mt-3 space-y-2 text-xs text-slate-400 list-decimal list-inside">
        <li>Connect a cloud account with a scoped API token (stored with KMS envelope encryption).</li>
        <li>Set a monthly budget cap so provisioning is blocked before overspend.</li>
        <li>Provision servers through the signed cost-approval gate — every action is audited.</li>
      </ol>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ApiError, useApi } from '~/composables/useApi'
import type { BudgetPolicy, CloudConnection, CloudResource, ItemEnvelope, ListEnvelope } from '~/types/api'

definePageMeta({
  layout: 'dashboard',
  middleware: 'auth',
})

useHead({ title: 'Console Overview' })

const auth = useAuthStore()
const loading = ref(false)
const error = ref<ApiError | null>(null)
const connections = ref<CloudConnection[]>([])
const resources = ref<CloudResource[]>([])
const budget = ref<BudgetPolicy | null>(null)

const { request } = useApi()

async function load() {
  loading.value = true
  error.value = null
  try {
    const [conns, res] = await Promise.all([
      request<ListEnvelope<CloudConnection>>('/api/v1/cloud/connections'),
      request<ListEnvelope<CloudResource>>('/api/v1/cloud/resources'),
    ])
    connections.value = conns.data ?? []
    resources.value = res.data ?? []
    try {
      const b = await request<ItemEnvelope<BudgetPolicy>>('/api/v1/cloud/budgets')
      budget.value = b.data
    } catch {
      budget.value = null
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
