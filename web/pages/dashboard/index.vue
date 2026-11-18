<template>
  <div class="space-y-8">
    <PageHeader title="Overview" :subtitle="`${auth.org?.name || 'Your organization'} · Role: ${auth.role || '—'}`" />

    <LoadingSkeleton v-if="loading" :rows="3" :lines="2" label="Loading overview" />

    <ErrorState v-else-if="error" :title="error.title" :message="error.message" @retry="load" />

    <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <StatCard :value="connections.length" label="Connected Clouds" to="/dashboard/connections" link-label="Manage connections" />
      <StatCard :value="resources.length" label="Managed Servers" to="/dashboard/servers" link-label="View inventory" />
      <StatCard
        :value="budget ? '$' + Number(budget.monthly_limit_usd).toFixed(0) : '—'"
        label="Monthly Budget Cap"
        to="/dashboard/budgets"
        link-label="Configure budgets"
      />
    </div>

    <div v-if="!loading && !error" class="card p-6">
      <h3 class="text-sm font-bold text-ink">Getting started</h3>
      <ol class="mt-3 space-y-2 text-xs text-ink-secondary list-decimal list-inside">
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
