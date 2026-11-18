<template>
  <div class="space-y-8">
    <PageHeader title="Operations Health" subtitle="Live platform signals. Every value below is read from the running API." />

    <LoadingSkeleton v-if="loading" :rows="3" :lines="2" label="Loading ops health" />

    <ErrorState v-else-if="error" :title="error.title" :message="error.message" @retry="load" />

    <template v-else>
      <div class="grid grid-cols-2 lg:grid-cols-4 gap-4 sm:gap-6">
        <StatCard :value="overview?.users_total ?? '—'" label="Registered Users" to="/admin/users" link-label="Manage members" />
        <StatCard :value="overview?.servers_managed ?? '—'" label="Managed Servers" />
        <StatCard :value="overview?.connections_active ?? '—'" label="Active Connections" />
        <StatCard :value="overview?.approvals_signed ?? '—'" label="Approvals Signed" />
        <StatCard :value="providerCount" label="Providers in Catalog" to="/admin/providers" link-label="Moderate catalog" />
        <StatCard :value="verifiedCount" label="Verified Pricing APIs" />
        <StatCard :value="overview?.skus_indexed ?? '—'" label="Plans Indexed" />
        <div class="card p-6">
          <div class="flex items-center gap-2">
            <span class="w-2.5 h-2.5 rounded-full live-dot" :class="readiness?.status === 'ready' ? 'bg-emerald-400' : 'bg-rose-400'" />
            <span class="text-lg font-bold text-ink">{{ readiness?.status === 'ready' ? 'API Ready' : 'API Not Ready' }}</span>
          </div>
          <div class="text-xs text-ink-secondary mt-2">
            Identity: <strong class="text-ink">{{ overview?.identity_store || '—' }}</strong>
            · BYOC store: <strong class="text-ink">{{ overview?.management_store || '—' }}</strong>
          </div>
        </div>
      </div>

      <div v-if="overview?.maintenance_mode" class="px-4 py-3 rounded-xl bg-amber-500/10 border border-amber-500/30 text-xs text-amber-700 dark:text-amber-300" role="alert">
        Maintenance mode is ON: user traffic receives 503 while admin, health, notices, and site-config stay reachable.
        <NuxtLink to="/admin/settings" class="font-bold underline underline-offset-2">Open settings</NuxtLink>
      </div>

      <div class="card p-6 space-y-3">
        <h3 class="text-sm font-bold text-ink">Affiliate compliance</h3>
        <div class="text-xs text-ink-secondary space-y-1">
          <p>Disclosure policy version: <strong class="text-ink">{{ disclosure?.policy_version || '—' }}</strong></p>
          <p>Independent comparison: <strong class="text-ink">{{ disclosure?.independent_comparison ? 'enforced' : '—' }}</strong></p>
          <p>Pay-for-placement: <strong class="text-ink">{{ disclosure && !disclosure.pay_for_placement_allowed ? 'blocked' : '—' }}</strong></p>
        </div>
        <NuxtLink to="/affiliate-disclosure" class="inline-block text-xs font-semibold text-brand hover:underline">
          View public disclosure &rarr;
        </NuxtLink>
      </div>

      <div class="card p-6 space-y-3">
        <h3 class="text-sm font-bold text-ink">Sync status by provider</h3>
        <div class="overflow-x-auto -mx-6 px-6">
          <table class="table-shell min-w-[560px]">
            <thead>
              <tr>
                <th class="py-2 pr-4">Provider</th>
                <th class="py-2 pr-4">Pricing API</th>
                <th class="py-2 pr-4">Management API</th>
                <th class="py-2">Affiliate</th>
              </tr>
            </thead>
            <tbody class="text-ink-secondary">
              <tr v-for="p in providers" :key="p.slug">
                <td class="py-2 pr-4 font-semibold text-ink">{{ p.name }}</td>
                <td class="py-2 pr-4"><StatusBadge :status="p.pricing_api_status" /></td>
                <td class="py-2 pr-4"><StatusBadge :status="p.management_api_status" /></td>
                <td class="py-2"><StatusBadge :status="p.affiliate_status" /></td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { ApiError, useApi } from '~/composables/useApi'
import type { AdminOverview, ListEnvelope, Provider } from '~/types/api'

definePageMeta({
  layout: 'admin',
  middleware: 'admin',
})

useHead({ title: 'Ops Health' })

interface Readiness {
  status: string
  service: string
  database: string
}

interface Disclosure {
  policy_version: string
  independent_comparison: boolean
  pay_for_placement_allowed: boolean
}

const loading = ref(false)
const error = ref<ApiError | null>(null)
const readiness = ref<Readiness | null>(null)
const providers = ref<Provider[]>([])
const disclosure = ref<Disclosure | null>(null)
const overview = ref<AdminOverview | null>(null)

const providerCount = computed(() => providers.value.length)
const verifiedCount = computed(() => providers.value.filter((p) => p.pricing_api_status === 'VERIFIED').length)

const { request } = useApi()

async function load() {
  loading.value = true
  error.value = null
  try {
    const [ready, provs, disc, over] = await Promise.all([
      request<Readiness>('/api/v1/readiness'),
      request<ListEnvelope<Provider>>('/api/v1/providers'),
      request<{ data: Disclosure }>('/api/v1/affiliate/disclosure'),
      request<AdminOverview>('/api/v1/admin/overview'),
    ])
    readiness.value = ready
    providers.value = provs.data ?? []
    disclosure.value = disc.data
    overview.value = over
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
