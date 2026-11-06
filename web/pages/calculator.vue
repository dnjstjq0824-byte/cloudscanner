<template>
  <div class="space-y-10">
    <!-- Header -->
    <div class="space-y-2">
      <div class="inline-flex items-center space-x-2 px-3 py-1 rounded-full bg-sky-950 border border-sky-800 text-sky-400 text-xs font-semibold">
        <span>Deterministic 36-Dimension Engine</span>
      </div>
      <h1 class="text-3xl font-extrabold text-white tracking-tight">Total Cost of Ownership (TCO) Calculator</h1>
      <p class="text-sm text-slate-400 max-w-3xl">
        Calculate actual, transparent cloud bills. Account for monthly hourly caps, separate IPv4 charges ($0.005/hr on AWS), egress bandwidth overages, and backup storage fees.
      </p>
    </div>

    <!-- Interactive Parameter Controls -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6 bg-slate-900/80 p-6 rounded-xl border border-slate-800">
      <!-- Usage Hours -->
      <div class="space-y-3">
        <div class="flex justify-between items-center text-xs">
          <label class="font-semibold text-slate-300">Monthly Usage Hours</label>
          <span class="font-bold text-sky-400 tabular-nums">{{ hours }} hrs / mo</span>
        </div>
        <input
          v-model.number="hours"
          type="range"
          min="1"
          max="744"
          step="1"
          class="w-full accent-sky-500 cursor-pointer"
        />
        <div class="flex justify-between text-[10px] text-slate-500">
          <button class="hover:text-slate-300" @click="hours = 168">1 Wk (168h)</button>
          <button class="hover:text-slate-300" @click="hours = 672">DO Cap (672h)</button>
          <button class="hover:text-slate-300" @click="hours = 730">Full Mo (730h)</button>
        </div>
      </div>

      <!-- Outbound Bandwidth -->
      <div class="space-y-3">
        <div class="flex justify-between items-center text-xs">
          <label class="font-semibold text-slate-300">Outbound Egress Bandwidth</label>
          <span class="font-bold text-sky-400 tabular-nums">{{ egressGB }} GB</span>
        </div>
        <input
          v-model.number="egressGB"
          type="range"
          min="0"
          max="10000"
          step="100"
          class="w-full accent-sky-500 cursor-pointer"
        />
        <div class="flex justify-between text-[10px] text-slate-500">
          <button class="hover:text-slate-300" @click="egressGB = 500">500 GB</button>
          <button class="hover:text-slate-300" @click="egressGB = 2000">2,000 GB</button>
          <button class="hover:text-slate-300" @click="egressGB = 5000">5,000 GB</button>
        </div>
      </div>

      <!-- Add-ons & Currency -->
      <div class="space-y-4">
        <div class="flex items-center justify-between">
          <label class="text-xs font-semibold text-slate-300">Include Daily Backups (+20%)</label>
          <input
            v-model="includeBackups"
            type="checkbox"
            class="rounded text-sky-600 bg-slate-800 border-slate-700 h-4 w-4 focus:ring-0"
          />
        </div>

        <div class="flex items-center justify-between">
          <label class="text-xs font-semibold text-slate-300">Extra IPv4 Addresses</label>
          <select
            v-model.number="extraIPv4"
            class="bg-slate-800 border border-slate-700 rounded px-2 py-1 text-xs text-slate-200"
          >
            <option :value="0">0 (Base only)</option>
            <option :value="1">+1 Extra IPv4</option>
            <option :value="2">+2 Extra IPv4</option>
          </select>
        </div>

        <div class="flex items-center justify-between pt-1">
          <label class="text-xs font-semibold text-slate-300">Display Currency</label>
          <div class="flex space-x-1">
            <button
              v-for="curr in ['USD', 'KRW', 'EUR']"
              :key="curr"
              class="px-2.5 py-0.5 rounded text-xs font-bold transition"
              :class="currency === curr ? 'bg-sky-600 text-white' : 'bg-slate-800 text-slate-400 hover:text-white'"
              @click="currency = curr"
            >
              {{ curr }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Comparison Matrix Table -->
    <div class="space-y-4">
      <div class="flex justify-between items-center">
        <h2 class="text-lg font-bold text-white">Scenario: 2 vCPU, 4GB RAM Workload Comparison</h2>
        <span class="text-xs text-slate-500">* Server-side TCO using active provider rate cards</span>
      </div>

      <div v-if="loading && comparisonRows.length === 0" class="rounded-xl border border-slate-800 overflow-hidden" aria-busy="true" aria-label="Calculating TCO">
        <div v-for="n in 5" :key="n" class="p-4 bg-slate-900 border-b border-slate-800 animate-pulse">
          <div class="h-4 w-3/4 rounded bg-slate-800" />
        </div>
      </div>

      <div v-else-if="error" class="text-center py-12 bg-red-950/20 rounded-xl border border-red-900 space-y-3">
        <h4 class="text-base font-bold text-white">{{ error.title }}</h4>
        <p class="text-xs text-slate-400 max-w-sm mx-auto">{{ error.message }}</p>
        <button
          class="px-4 py-2 rounded bg-sky-600 hover:bg-sky-500 text-white text-xs font-semibold transition"
          @click="recalculate"
        >
          Retry Calculation
        </button>
      </div>

      <div v-else class="overflow-x-auto rounded-xl border border-slate-800">
        <table class="w-full text-left text-xs bg-slate-900 divide-y divide-slate-800">
          <thead class="bg-slate-950/80 text-slate-400 font-semibold uppercase tracking-wider">
            <tr>
              <th class="py-3 px-4">Provider & Plan</th>
              <th class="py-3 px-4">Compute Base</th>
              <th class="py-3 px-4">IPv4 Charge</th>
              <th class="py-3 px-4">Egress Overage</th>
              <th class="py-3 px-4">Backup Snapshots</th>
              <th class="py-3 px-4 text-right">Effective Total ({{ currency }})</th>
              <th class="py-3 px-4 text-center">Action</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-800 text-slate-300 tabular-nums">
            <tr
              v-for="row in comparisonRows"
              :key="row.provider"
              class="hover:bg-slate-800/40 transition"
              :class="{ 'bg-sky-950/20': row.isBestValue }"
            >
              <td class="py-4 px-4 font-bold text-white">
                <span>{{ row.provider }}</span>
                <span class="text-[10px] font-normal text-slate-400"> ({{ row.plan }})</span>
                <span v-if="row.isBestValue" class="ml-2 px-1.5 py-0.5 rounded text-[9px] bg-emerald-950 text-emerald-400 border border-emerald-800">
                  Best Value
                </span>
              </td>
              <td class="py-4 px-4">
                ${{ row.computeUSD.toFixed(2) }}
                <span class="block text-[10px] text-slate-500 font-normal">
                  {{ row.capped ? 'Capped at monthly max' : 'Uncapped hourly' }}
                </span>
              </td>
              <td class="py-4 px-4">
                <span :class="row.ipv4USD > 0 ? 'text-amber-400 font-semibold' : 'text-slate-400'">
                  ${{ row.ipv4USD.toFixed(2) }}
                </span>
                <span v-if="row.ipv4USD > 0" class="block text-[10px] text-slate-500">
                  $0.005/hr surcharge
                </span>
              </td>
              <td class="py-4 px-4">
                <span :class="row.egressUSD > 0 ? 'text-rose-400 font-semibold' : 'text-slate-400'">
                  ${{ row.egressUSD.toFixed(2) }}
                </span>
                <span class="block text-[10px] text-slate-500">
                  {{ row.freeTB }} TB Included
                </span>
              </td>
              <td class="py-4 px-4">
                ${{ row.backupUSD.toFixed(2) }}
              </td>
              <td class="py-4 px-4 text-right font-extrabold text-sm text-white">
                {{ formatTotal(row.totalTarget) }}
              </td>
              <td class="py-4 px-4 text-center">
                <a
                  :href="dealUrl(row.slug)"
                  class="px-3 py-1.5 rounded bg-sky-600 hover:bg-sky-500 text-white font-semibold text-[11px] transition"
                >
                  Deploy
                </a>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { ApiError, useApi } from '~/composables/useApi'
import type { ListEnvelope, SKUDetail, TCOResult } from '~/types/api'

definePageMeta({
  layout: 'public',
})

useHead({
  title: 'TCO Calculator',
  meta: [
    { name: 'description', content: 'Transparent total-cost-of-ownership calculator: hourly caps, IPv4 surcharges, egress overages, and backups across 5 providers.' },
  ],
})

interface CalcRow {
  provider: string
  plan: string
  slug: string
  computeUSD: number
  capped: boolean
  ipv4USD: number
  egressUSD: number
  freeTB: number
  backupUSD: number
  totalUSD: number
  totalTarget: number
  isBestValue: boolean
}

const hours = ref(730)
const egressGB = ref(1000)
const includeBackups = ref(true)
const extraIPv4 = ref(0)
const currency = ref('USD')

const fxRates: Record<string, number> = {
  USD: 1.0,
  KRW: 1380.0,
  EUR: 0.92,
}

// Canonical 2C/4G comparison set resolved against the live catalog.
const comparisonPlans = [
  { planCode: 's-2vcpu-4gb', provider: 'DigitalOcean' },
  { planCode: 'cpx21', provider: 'Hetzner Cloud' },
  { planCode: 'vc2-2c-4gb', provider: 'Vultr' },
  { planCode: 't3.medium', provider: 'AWS EC2' },
  { planCode: '2xCPU-4GB', provider: 'UpCloud' },
]

const loading = ref(false)
const error = ref<ApiError | null>(null)
const comparisonRows = ref<CalcRow[]>([])
const catalogSkus = ref<SKUDetail[]>([])

const { request, baseURL } = useApi()

function dealUrl(slug: string): string {
  return `${baseURL}/go/${slug}?src=calculator_matrix`
}

function formatTotal(target: number): string {
  if (currency.value === 'KRW') {
    return '₩' + Math.round(target).toLocaleString()
  } else if (currency.value === 'EUR') {
    return '€' + target.toFixed(2)
  }
  return '$' + target.toFixed(2)
}

let debounce: ReturnType<typeof setTimeout> | null = null

async function recalculate() {
  loading.value = true
  error.value = null
  try {
    if (catalogSkus.value.length === 0) {
      const list = await request<ListEnvelope<SKUDetail>>('/api/v1/catalog/skus', { query: { category: 'VPS' } })
      catalogSkus.value = list.data ?? []
    }
    const byPlan = new Map(catalogSkus.value.map((s) => [s.plan_code, s]))
    const scenario = {
      usage_hours: hours.value,
      required_ipv4_count: 1 + extraIPv4.value,
      estimated_egress_gb: egressGB.value,
      enable_backups: includeBackups.value,
      snapshot_storage_gb: 0,
      target_currency: currency.value,
      fx_rate_to_target: fxRates[currency.value] || 1.0,
      tax_percentage: 0,
    }
    const rows: CalcRow[] = []
    for (const plan of comparisonPlans) {
      const sku = byPlan.get(plan.planCode)
      if (!sku) {
        continue
      }
      const res = await request<{ data: TCOResult }>('/api/v1/estimate/tco', {
        method: 'POST',
        body: { sku_id: sku.id, scenario },
      })
      const t = res.data
      rows.push({
        provider: plan.provider,
        plan: sku.plan_code,
        slug: sku.provider_slug,
        computeUSD: t.base_compute_usd,
        capped: sku.has_monthly_cap,
        ipv4USD: t.ipv4_cost_usd,
        egressUSD: t.egress_cost_usd,
        freeTB: Math.round((sku.included_bandwidth_gb / 1024) * 10) / 10,
        backupUSD: t.backup_cost_usd,
        totalUSD: t.subtotal_usd,
        totalTarget: t.total_target,
        isBestValue: false,
      })
    }
    if (rows.length === 0) {
      throw new ApiError(0, 'No Plans Available', 'The comparison plans are not in the catalog yet.')
    }
    rows.sort((a, b) => a.totalUSD - b.totalUSD)
    rows[0].isBestValue = true
    comparisonRows.value = rows
  } catch (e) {
    error.value = e instanceof ApiError ? e : ApiError.fromUnknown(e)
  } finally {
    loading.value = false
  }
}

function scheduleRecalc() {
  if (debounce) {
    clearTimeout(debounce)
  }
  debounce = setTimeout(() => recalculate(), 350)
}

watch([hours, egressGB, includeBackups, extraIPv4, currency], scheduleRecalc)

onMounted(() => {
  recalculate()
})
</script>
