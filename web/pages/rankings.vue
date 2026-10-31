<template>
  <div class="space-y-8">
    <div class="space-y-2">
      <div class="inline-flex items-center space-x-2 px-3 py-1 rounded-full bg-amber-950 border border-amber-800 text-amber-400 text-xs font-semibold">
        <span>🏆 Algorithmic Leaderboard</span>
      </div>
      <h1 class="text-3xl font-extrabold text-white tracking-tight">Best Value Cloud Rankings</h1>
      <p class="text-sm text-slate-400 max-w-2xl">
        Ranked strictly by mathematical value score: (vCPU × 10 + RAM_GB × 2 + Disk_GB × 0.1) ÷ Monthly TCO. Zero pay-for-placement.
      </p>
    </div>

    <!-- Category Tabs -->
    <div class="flex space-x-2 border-b border-slate-800 pb-2 text-xs font-semibold">
      <button
        v-for="cat in ['VPS', 'Dedicated CPU', 'High RAM']"
        :key="cat"
        class="px-4 py-2 rounded-lg transition"
        :class="activeCat === cat ? 'bg-sky-600 text-white shadow' : 'text-slate-400 hover:text-white bg-slate-900'"
        @click="activeCat = cat"
      >
        {{ cat }}
      </button>
    </div>

    <!-- Leaderboard Table -->
    <div class="overflow-x-auto rounded-xl border border-slate-800">
      <table class="w-full text-left text-xs bg-slate-900 divide-y divide-slate-800">
        <thead class="bg-slate-950 text-slate-400 font-semibold uppercase tracking-wider">
          <tr>
            <th class="py-3 px-4 text-center">Rank</th>
            <th class="py-3 px-4">Provider & Plan</th>
            <th class="py-3 px-4">Specifications</th>
            <th class="py-3 px-4">Monthly TCO</th>
            <th class="py-3 px-4 text-center">Value Score</th>
            <th class="py-3 px-4 text-center">Action</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-800 text-slate-300 tabular-nums">
          <tr v-for="(item, idx) in leaderboard" :key="item.plan" class="hover:bg-slate-800/40 transition">
            <td class="py-4 px-4 text-center font-extrabold text-sm">
              <span v-if="idx === 0">🥇</span>
              <span v-else-if="idx === 1">🥈</span>
              <span v-else-if="idx === 2">🥉</span>
              <span v-else class="text-slate-500">#{{ idx + 1 }}</span>
            </td>
            <td class="py-4 px-4 font-bold text-white">
              <div class="flex items-center space-x-2">
                <span>{{ item.provider }}</span>
                <span class="text-slate-400 font-normal">({{ item.plan }})</span>
              </div>
            </td>
            <td class="py-4 px-4">
              {{ item.vcpu }} vCPU • {{ item.ram }} GB RAM • {{ item.disk }} GB {{ item.diskType }}
            </td>
            <td class="py-4 px-4 font-bold text-white">
              ${{ item.price.toFixed(2) }}/mo
            </td>
            <td class="py-4 px-4 text-center">
              <span class="px-2.5 py-1 rounded-full font-bold bg-sky-950 text-sky-400 border border-sky-800 text-xs">
                {{ item.score.toFixed(2) }}
              </span>
            </td>
            <td class="py-4 px-4 text-center">
              <NuxtLink
                :to="`/go/${item.slug}?src=ranking_table`"
                class="px-3 py-1.5 rounded bg-sky-600 hover:bg-sky-500 text-white font-semibold text-[11px] transition shadow"
              >
                Inspect
              </NuxtLink>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

definePageMeta({
  layout: 'public'
})

const activeCat = ref('VPS')

const leaderboard = computed(() => {
  return [
    { rank: 1, provider: 'Hetzner', plan: 'CPX21', slug: 'hetzner', vcpu: 3, ram: 4, disk: 80, diskType: 'NVMe', price: 8.95, score: 9.45 },
    { rank: 2, provider: 'Hetzner', plan: 'CPX11', slug: 'hetzner', vcpu: 2, ram: 2, disk: 40, diskType: 'NVMe', price: 4.95, score: 8.92 },
    { rank: 3, provider: 'Vultr', plan: 'vc2-2c-4gb', slug: 'vultr', vcpu: 2, ram: 4, disk: 80, diskType: 'NVMe', price: 20.00, score: 5.80 },
    { rank: 4, provider: 'UpCloud', plan: '2xCPU-4GB', slug: 'upcloud', vcpu: 2, ram: 4, disk: 80, diskType: 'MaxIOPS', price: 20.00, score: 5.80 },
    { rank: 5, provider: 'DigitalOcean', plan: 's-2vcpu-4gb', slug: 'digitalocean', vcpu: 2, ram: 4, disk: 80, diskType: 'NVMe', price: 24.00, score: 5.67 },
    { rank: 6, provider: 'DigitalOcean', plan: 's-1vcpu-1gb', slug: 'digitalocean', vcpu: 1, ram: 1, disk: 25, diskType: 'NVMe', price: 6.00, score: 3.42 },
    { rank: 7, provider: 'AWS EC2', plan: 't4g.small', slug: 'aws', vcpu: 2, ram: 2, disk: 40, diskType: 'EBS', price: 15.91, score: 3.10 },
  ]
})
</script>
