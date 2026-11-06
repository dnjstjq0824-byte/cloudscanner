<template>
  <div class="space-y-8">
    <div>
      <h2 class="text-2xl font-extrabold text-white tracking-tight">Provision Server</h2>
      <p class="text-sm text-slate-400 mt-1">
        Every billable action passes an explicit cost-approval gate: estimate first, sign the approval, then execute.
      </p>
    </div>

    <!-- Step 1: pick connection + SKU -->
    <div class="p-6 rounded-xl bg-slate-900 border border-slate-800 space-y-5">
      <h3 class="text-sm font-bold text-white">1. Target connection & plan</h3>
      <div v-if="bootLoading" class="text-xs text-slate-500">Loading connections and catalog…</div>
      <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div class="space-y-1.5">
          <label for="prov-conn" class="text-xs font-semibold text-slate-300">Connection (must be read-write)</label>
          <select
            id="prov-conn"
            v-model="form.connectionId"
            class="w-full px-3 py-2.5 bg-slate-950 border border-slate-700 rounded-lg text-white text-sm focus:outline-none focus:border-sky-500"
          >
            <option value="">Select a connection…</option>
            <option v-for="c in writableConnections" :key="c.id" :value="c.id">
              {{ c.connection_name }} ({{ c.provider_slug }})
            </option>
          </select>
          <p v-if="connections.length > 0 && writableConnections.length === 0" class="text-[11px] text-amber-400">
            All connections are read-only. Reconnect with read-write access to provision.
          </p>
        </div>
        <div class="space-y-1.5">
          <label for="prov-sku" class="text-xs font-semibold text-slate-300">Plan (SKU)</label>
          <select
            id="prov-sku"
            v-model="form.skuId"
            class="w-full px-3 py-2.5 bg-slate-950 border border-slate-700 rounded-lg text-white text-sm focus:outline-none focus:border-sky-500"
            @change="estimate"
          >
            <option value="">Select a plan…</option>
            <option v-for="s in skus" :key="s.id" :value="s.id">
              {{ s.provider_name }} {{ s.plan_code }} · {{ s.vcpu }}C/{{ (s.ram_mb / 1024).toFixed(0) }}G · ${{ s.raw_monthly_price_usd }}/mo
            </option>
          </select>
        </div>
      </div>
    </div>

    <!-- Step 2: estimate -->
    <div class="p-6 rounded-xl bg-slate-900 border border-slate-800 space-y-4">
      <h3 class="text-sm font-bold text-white">2. Cost estimate</h3>
      <div v-if="estimating" class="text-xs text-slate-500">Calculating transparent TCO…</div>
      <div v-else-if="tco" class="grid grid-cols-2 md:grid-cols-4 gap-4 text-xs">
        <div class="p-4 rounded-lg bg-slate-950 border border-slate-800">
          <span class="text-slate-500 block">Compute</span>
          <span class="text-white font-bold text-base tabular-nums">${{ tco.base_compute_usd.toFixed(2) }}</span>
        </div>
        <div class="p-4 rounded-lg bg-slate-950 border border-slate-800">
          <span class="text-slate-500 block">IPv4 + Egress</span>
          <span class="text-white font-bold text-base tabular-nums">${{ (tco.ipv4_cost_usd + tco.egress_cost_usd).toFixed(2) }}</span>
        </div>
        <div class="p-4 rounded-lg bg-slate-950 border border-slate-800">
          <span class="text-slate-500 block">Backups</span>
          <span class="text-white font-bold text-base tabular-nums">${{ tco.backup_cost_usd.toFixed(2) }}</span>
        </div>
        <div class="p-4 rounded-lg bg-sky-950 border border-sky-800">
          <span class="text-sky-400 block">Monthly total</span>
          <span class="text-white font-bold text-base tabular-nums">${{ tco.subtotal_usd.toFixed(2) }}</span>
        </div>
      </div>
      <p v-else class="text-xs text-slate-500">Select a plan to see the itemized estimate.</p>
    </div>

    <!-- Step 3: approval gate -->
    <div class="p-6 rounded-xl bg-slate-900 border border-slate-800 space-y-4">
      <h3 class="text-sm font-bold text-white">3. Sign cost approval</h3>
      <label class="flex items-start gap-2 text-xs text-slate-300 cursor-pointer">
        <input v-model="form.mfaConfirmed" type="checkbox" class="mt-0.5 rounded text-sky-600 bg-slate-800 border-slate-700" />
        <span>I confirm step-up verification for this billable operation and accept the estimated charges above.</span>
      </label>

      <div v-if="approvalError" class="px-4 py-3 rounded-lg bg-red-950/40 border border-red-900 text-xs text-red-300" role="alert">
        {{ approvalError }}
      </div>

      <div v-if="approval" class="p-4 rounded-lg bg-emerald-950/40 border border-emerald-900 text-xs text-emerald-200 space-y-1">
        <p class="font-bold">Approval signed and valid until {{ formatDate(approval.expires_at) }}</p>
        <p class="font-mono break-all text-[11px]">HMAC {{ approval.approval_signature_hmac.slice(0, 32) }}…</p>
        <p class="text-emerald-300/80">Execution hands the signed approval to the provider adapter with an idempotency key. Live execution requires an authorized sandbox budget.</p>
      </div>

      <button
        class="px-5 py-2.5 rounded-lg bg-sky-600 hover:bg-sky-500 disabled:opacity-60 text-white text-sm font-semibold transition"
        :disabled="!canApprove || approving"
        @click="signApproval"
      >
        {{ approving ? 'Signing…' : 'Sign Approval' }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { ApiError, useApi } from '~/composables/useApi'
import type { CloudConnection, ItemEnvelope, ListEnvelope, OperationApproval, SKUDetail, TCOResult } from '~/types/api'

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

const form = reactive({
  connectionId: '',
  skuId: '',
  mfaConfirmed: false,
})

const { request } = useApi()

const writableConnections = computed(() => connections.value.filter((c) => c.access_level === 'READ_WRITE' && c.connection_status === 'ACTIVE'))

const canApprove = computed(() => form.connectionId !== '' && form.skuId !== '' && tco.value !== null && form.mfaConfirmed && approval.value === null)

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
  if (!form.skuId) {
    return
  }
  estimating.value = true
  try {
    const res = await request<{ data: TCOResult }>('/api/v1/estimate/tco', {
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
