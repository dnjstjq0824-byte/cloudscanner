<template>
  <div class="space-y-8">
    <PageHeader
      title="Provision Server"
      subtitle="Every billable action passes an explicit cost-approval gate: estimate first, sign the approval, then execute."
    />

    <div class="card p-6 space-y-5">
      <h3 class="text-sm font-bold text-ink">1. Target connection & plan</h3>
      <div v-if="bootLoading" class="text-xs text-ink-secondary">Loading connections and catalog…</div>
      <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div class="space-y-1.5">
          <label for="prov-conn" class="label">Connection (must be read-write)</label>
          <select id="prov-conn" v-model="form.connectionId" class="input">
            <option value="">Select a connection…</option>
            <option v-for="c in writableConnections" :key="c.id" :value="c.id">
              {{ c.connection_name }} ({{ c.provider_slug }})
            </option>
          </select>
          <p v-if="connections.length > 0 && writableConnections.length === 0" class="text-[11px] text-amber-600 dark:text-amber-400">
            All connections are read-only. Reconnect with read-write access to provision.
          </p>
        </div>
        <div class="space-y-1.5">
          <label for="prov-sku" class="label">Plan (SKU)</label>
          <select id="prov-sku" v-model="form.skuId" class="input" @change="estimate">
            <option value="">Select a plan…</option>
            <option v-for="s in skus" :key="s.id" :value="s.id">
              {{ s.provider_name }} {{ s.plan_code }} · {{ s.vcpu }}C/{{ (s.ram_mb / 1024).toFixed(0) }}G · ${{ s.raw_monthly_price_usd }}/mo
            </option>
          </select>
        </div>
        <div class="space-y-1.5">
          <label for="prov-host" class="label">Hostname</label>
          <input id="prov-host" v-model="form.hostname" type="text" placeholder="web-01" autocomplete="off" class="input" />
        </div>
        <div class="grid grid-cols-3 gap-3">
          <div class="space-y-1.5">
            <label for="prov-vcpu" class="label">vCPU</label>
            <input id="prov-vcpu" v-model.number="form.vcpu" type="number" min="1" class="input" />
          </div>
          <div class="space-y-1.5">
            <label for="prov-ram" class="label">RAM MB</label>
            <input id="prov-ram" v-model.number="form.ramMb" type="number" min="512" step="256" class="input" />
          </div>
          <div class="space-y-1.5">
            <label for="prov-disk" class="label">Disk GB</label>
            <input id="prov-disk" v-model.number="form.diskGb" type="number" min="10" class="input" />
          </div>
        </div>
      </div>
    </div>

    <div class="card p-6 space-y-4">
      <h3 class="text-sm font-bold text-ink">2. Cost estimate</h3>
      <div v-if="estimating" class="text-xs text-ink-secondary">Calculating transparent TCO…</div>
      <div v-else-if="tco" class="grid grid-cols-2 md:grid-cols-4 gap-4 text-xs">
        <div class="p-4 rounded-xl bg-elevated border border-subtle">
          <span class="text-ink-secondary block">Compute</span>
          <span class="text-ink font-bold text-base tabular-nums">${{ money(tco.base_compute_usd) }}</span>
        </div>
        <div class="p-4 rounded-xl bg-elevated border border-subtle">
          <span class="text-ink-secondary block">IPv4 + Egress</span>
          <span class="text-ink font-bold text-base tabular-nums">${{ money(num(tco.ipv4_cost_usd) + num(tco.egress_cost_usd)) }}</span>
        </div>
        <div class="p-4 rounded-xl bg-elevated border border-subtle">
          <span class="text-ink-secondary block">Backups</span>
          <span class="text-ink font-bold text-base tabular-nums">${{ money(tco.backup_cost_usd) }}</span>
        </div>
        <div class="p-4 rounded-xl bg-sky-500/10 border border-sky-500/30">
          <span class="text-sky-700 dark:text-sky-300 block">Monthly total</span>
          <span class="text-ink font-bold text-base tabular-nums">${{ money(tco.subtotal_usd) }}</span>
        </div>
      </div>
      <p v-else class="text-xs text-ink-secondary">Select a plan to see the itemized estimate.</p>
    </div>

    <div class="card p-6 space-y-4">
      <h3 class="text-sm font-bold text-ink">3. Sign cost approval</h3>
      <label class="flex items-start gap-2 text-xs text-ink-secondary cursor-pointer">
        <input v-model="form.mfaConfirmed" type="checkbox" class="mt-0.5 rounded accent-sky-600" />
        <span>I confirm step-up verification for this billable operation and accept the estimated charges above.</span>
      </label>

      <div v-if="approvalError" class="px-4 py-3 rounded-xl bg-rose-500/10 border border-rose-500/30 text-xs text-rose-600 dark:text-rose-300" role="alert">
        {{ approvalError }}
      </div>

      <div v-if="approval" class="p-4 rounded-xl bg-emerald-500/10 border border-emerald-500/30 text-xs text-emerald-700 dark:text-emerald-200 space-y-1">
        <p class="font-bold">Approval signed and valid until {{ formatDate(approval.expires_at) }}</p>
        <p class="font-mono break-all text-[11px]">HMAC {{ approval.approval_signature_hmac.slice(0, 32) }}…</p>
      </div>

      <button class="btn-primary" :disabled="!canApprove || approving" @click="signApproval">
        {{ approving ? 'Signing…' : 'Sign Approval' }}
      </button>
    </div>

    <div class="card p-6 space-y-4">
      <h3 class="text-sm font-bold text-ink">4. Execute provisioning</h3>
      <p class="text-xs text-ink-secondary">
        Execution consumes the signed approval exactly once (idempotency-keyed) and enforces budget caps before creating the managed server.
      </p>
      <div v-if="executeError" class="px-4 py-3 rounded-xl bg-rose-500/10 border border-rose-500/30 text-xs text-rose-600 dark:text-rose-300" role="alert">
        {{ executeError }}
      </div>
      <div v-if="provisioned" class="p-4 rounded-xl bg-emerald-500/10 border border-emerald-500/30 text-xs text-emerald-700 dark:text-emerald-200 space-y-1">
        <p class="font-bold">Server '{{ provisioned.hostname }}' created ({{ provisioned.power_status }}).</p>
        <p>{{ syncNote }}</p>
        <NuxtLink to="/dashboard/servers" class="inline-block font-semibold underline underline-offset-2">Open inventory →</NuxtLink>
      </div>
      <button class="btn-primary" :disabled="!canExecute || executing" @click="execute">
        {{ executing ? 'Provisioning…' : 'Execute Provisioning' }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { ApiError, useApi } from '~/composables/useApi'
import { useToast } from '~/composables/useToast'
import type { CloudConnection, CloudResource, ItemEnvelope, ListEnvelope, OperationApproval, ProvisionEnvelope, SKUDetail, TCOResult } from '~/types/api'

definePageMeta({
  layout: 'dashboard',
  middleware: 'auth',
})

useHead({ title: 'Provision Server' })

const bootLoading = ref(true)
const connections = ref<CloudConnection[]>([])
const skus = ref<SKUDetail[]>([])
const estimating = ref(false)
const tco = ref<TCOResult | null>(null)
const approving = ref(false)
const approval = ref<OperationApproval | null>(null)
const approvalError = ref('')
const executing = ref(false)
const executeError = ref('')
const provisioned = ref<CloudResource | null>(null)
const syncNote = ref('')

const form = reactive({
  connectionId: '',
  skuId: '',
  hostname: 'web-01',
  vcpu: 2,
  ramMb: 4096,
  diskGb: 80,
  mfaConfirmed: false,
})

const { request } = useApi()
const toast = useToast()

const writableConnections = computed(() => connections.value.filter((c) => c.access_level === 'READ_WRITE' && c.connection_status === 'ACTIVE'))

const canApprove = computed(() => form.connectionId !== '' && form.skuId !== '' && tco.value !== null && form.mfaConfirmed && approval.value === null)

const canExecute = computed(() => approval.value !== null && provisioned.value === null && form.hostname.trim() !== '')

function formatDate(iso: string): string {
  try {
    return new Date(iso).toLocaleString()
  } catch {
    return iso
  }
}

async function estimate() {
  tco.value = null
  approval.value = null
  provisioned.value = null
  if (!form.skuId) {
    return
  }
  estimating.value = true
  try {
    const res = await request<{ data: TCOResult }>('/api/v1/pricing/estimate', {
      method: 'POST',
      body: {
        sku_id: form.skuId,
        scenario: {
          usage_hours: 730,
          required_ipv4_count: 1,
          estimated_egress_gb: 0,
          enable_backups: false,
          snapshot_storage_gb: 0,
          target_currency: 'USD',
          fx_rate_to_target: 1,
          tax_percentage: 0,
        },
      },
    })
    tco.value = res.data
  } catch {
    tco.value = null
  } finally {
    estimating.value = false
  }
}

async function signApproval() {
  if (!tco.value) {
    return
  }
  approving.value = true
  approvalError.value = ''
  try {
    const res = await request<ItemEnvelope<OperationApproval>>('/api/v1/cloud/approvals', {
      method: 'POST',
      body: {
        connection_id: form.connectionId,
        operation_type: 'CREATE_SERVER',
        target_sku_id: form.skuId,
        estimated_hourly_cost_usd: tco.value.effective_hourly,
        estimated_monthly_cost_usd: tco.value.subtotal_usd,
        step_up_mfa_verified: form.mfaConfirmed,
        duration_minutes: 15,
      },
    })
    approval.value = res.data
  } catch (e) {
    const apiErr = e instanceof ApiError ? e : ApiError.fromUnknown(e)
    approvalError.value = apiErr.message
  } finally {
    approving.value = false
  }
}

async function execute() {
  if (!approval.value) {
    return
  }
  executing.value = true
  executeError.value = ''
  try {
    const res = await request<ProvisionEnvelope>('/api/v1/cloud/servers', {
      method: 'POST',
      body: {
        connection_id: form.connectionId,
        approval_id: approval.value.id,
        hostname: form.hostname.trim().toLowerCase(),
        vcpu_count: form.vcpu,
        ram_mb: form.ramMb,
        disk_gb: form.diskGb,
      },
    })
    provisioned.value = res.data
    syncNote.value = res.provider_sync?.detail || ''
    toast.success(`Server '${res.data.hostname}' provisioned.`)
  } catch (e) {
    const apiErr = e instanceof ApiError ? e : ApiError.fromUnknown(e)
    executeError.value = apiErr.message
  } finally {
    executing.value = false
  }
}

onMounted(async () => {
  try {
    const [conns, skuList] = await Promise.all([
      request<ListEnvelope<CloudConnection>>('/api/v1/cloud/connections'),
      request<ListEnvelope<SKUDetail>>('/api/v1/catalog/skus', { query: { category: 'VPS' } }),
    ])
    connections.value = conns.data ?? []
    skus.value = (skuList.data ?? []).slice(0, 60)
  } catch {
    connections.value = []
    skus.value = []
  } finally {
    bootLoading.value = false
  }
})
</script>
