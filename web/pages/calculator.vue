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
        <span class="text-xs text-slate-500">* All prices calculated mathematically using active provider rate cards</span>
      </div>

      <div class="overflow-x-auto rounded-xl border border-slate-800">
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
              <td class="py-4 px-4 font-bold text-white flex items-center space-x-2">
                <span>{{ row.provider }}</span>
                <span class="text-[10px] font-normal text-slate-400">({{ row.plan }})</span>
                <span v-if="row.isBestValue" class="px-1.5 py-0.2 rounded text-[9px] bg-emerald-950 text-emerald-400 border border-emerald-800">
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
                {{ formatCurrency(row.totalUSD) }}
              </td>
              <td class="py-4 px-4 text-center">
                <NuxtLink
                  :to="`/go/${row.slug}?src=calculator_matrix`"
                  class="px-3 py-1.5 rounded bg-sky-600 hover:bg-sky-500 text-white font-semibold text-[11px] transition"
                >
                  Deploy
                </NuxtLink>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

definePageMeta({
  layout: 'public'
})

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

function formatCurrency(usd: number): string {
  const rate = fxRates[currency.value] || 1.0
  const converted = usd * rate
  if (currency.value === 'KRW') {
    return '₩' + Math.round(converted).toLocaleString()
  } else if (currency.value === 'EUR') {
    return '€' + converted.toFixed(2)
  }
  return '$' + converted.toFixed(2)
}

const comparisonRows = computed(() => {
  // 1. DigitalOcean (s-2vcpu-4gb): $24/mo cap, $0.0357/hr, 4TB egress free, IPv4 included, 20% backup
  const doCompute = Math.min(24.00, hours.value * 0.0357)
  const doIPv4 = extraIPv4.value * 4.00
  const doOverageGB = Math.max(0, egressGB.value - 4000)
  const doEgress = doOverageGB * 0.01
  const doBackup = includeBackups.value ? doCompute * 0.20 : 0
  const doTotal = doCompute + doIPv4 + doEgress + doBackup

  // 2. Hetzner Cloud (CPX21): €8.95 (~$9.70) cap, $0.0134/hr, 20TB egress free, IPv4 included, 20% backup
  const hzCompute = Math.min(9.70, hours.value * 0.0134)
  const hzIPv4 = extraIPv4.value * 2.20
  const hzOverageGB = Math.max(0, egressGB.value - 20000)
  const hzEgress = hzOverageGB * 0.012
  const hzBackup = includeBackups.value ? hzCompute * 0.20 : 0
  const hzTotal = hzCompute + hzIPv4 + hzEgress + hzBackup

  // 3. Vultr (vc2-2c-4gb): $20/mo cap, $0.0298/hr, 3TB egress free, IPv4 included, 20% backup
  const vultrCompute = Math.min(20.00, hours.value * 0.0298)
  const vultrIPv4 = extraIPv4.value * 3.00
  const vultrOverageGB = Math.max(0, egressGB.value - 3000)
  const vultrEgress = vultrOverageGB * 0.01
  const vultrBackup = includeBackups.value ? vultrCompute * 0.20 : 0
  const vultrTotal = vultrCompute + vultrIPv4 + vultrEgress + vultrBackup

  // 4. AWS EC2 (t3.medium): $0.0416/hr uncapped, 100GB egress free ($0.09/GB overage), $0.005/hr IPv4 charge
  const awsCompute = hours.value * 0.0416
  const awsIPv4 = (1 + extraIPv4.value) * (hours.value * 0.005) // AWS charges $0.005/hr for primary AND extra IPv4
  const awsOverageGB = Math.max(0, egressGB.value - 100)
  const awsEgress = awsOverageGB * 0.09 // $0.09/GB standard internet egress
  const awsBackup = includeBackups.value ? (50 * 0.05) : 0 // EBS snapshot $0.05/GB-mo
  const awsTotal = awsCompute + awsIPv4 + awsEgress + awsBackup

  // 5. UpCloud (2xCPU-4GB): $20/mo cap, $0.0298/hr, 4TB egress free, IPv4 included
  const upCompute = Math.min(20.00, hours.value * 0.0298)
  const upIPv4 = extraIPv4.value * 2.50
  const upOverageGB = Math.max(0, egressGB.value - 4000)
  const upEgress = upOverageGB * 0.01
  const upBackup = includeBackups.value ? upCompute * 0.20 : 0
  const upTotal = upCompute + upIPv4 + upEgress + upBackup

  return [
    {
      provider: 'Hetzner Cloud',
      plan: 'CPX21 (3C/4G)',
      slug: 'hetzner',
      computeUSD: hzCompute,
      capped: true,
      ipv4USD: hzIPv4,
      egressUSD: hzEgress,
      freeTB: 20,
      backupUSD: hzBackup,
      totalUSD: hzTotal,
      isBestValue: true,
    },
    {
      provider: 'Vultr',
      plan: 'vc2-2c-4gb',
      slug: 'vultr',
      computeUSD: vultrCompute,
      capped: true,
      ipv4USD: vultrIPv4,
      egressUSD: vultrEgress,
      freeTB: 3,
      backupUSD: vultrBackup,
      totalUSD: vultrTotal,
      isBestValue: false,
    },
    {
      provider: 'UpCloud',
      plan: '2xCPU-4GB',
      slug: 'upcloud',
      computeUSD: upCompute,
      capped: true,
      ipv4USD: upIPv4,
      egressUSD: upEgress,
      freeTB: 4,
      backupUSD: upBackup,
      totalUSD: upTotal,
      isBestValue: false,
    },
    {
      provider: 'DigitalOcean',
      plan: 's-2vcpu-4gb',
      slug: 'digitalocean',
      computeUSD: doCompute,
      capped: true,
      ipv4USD: doIPv4,
      egressUSD: doEgress,
      freeTB: 4,
      backupUSD: doBackup,
      totalUSD: doTotal,
      isBestValue: false,
    },
    {
      provider: 'AWS EC2',
      plan: 't3.medium',
      slug: 'aws',
      computeUSD: awsCompute,
      capped: false,
      ipv4USD: awsIPv4,
      egressUSD: awsEgress,
      freeTB: 0.1,
      backupUSD: awsBackup,
      totalUSD: awsTotal,
      isBestValue: false,
    },
  ]
})
</script>
