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
    <div class="flex flex-wrap gap-2 border-b border-slate-800 pb-2 text-xs font-semibold">
      <button
        v-for="cat in categoryTabs"
        :key="cat.value"
        class="px-4 py-2 rounded-lg transition"
        :class="activeCat === cat.value ? 'bg-sky-600 text-white shadow' : 'text-slate-400 hover:text-white bg-slate-900'"
        @click="activeCat = cat.value; loadRankings()"
      >
        {{ cat.label }}
      </button>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="rounded-xl border border-slate-800 overflow-hidden" aria-busy="true" aria-label="Loading rankings">
      <div v-for="n in 6" :key="n" class="p-4 bg-slate-900 border-b border-slate-800 animate-pulse">
        <div class="h-4 w-2/3 rounded bg-slate-800" />
      </div>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="text-center py-16 bg-red-950/20 rounded-xl border border-red-900 space-y-3">
      <div class="text-3xl">⚠️</div>
      <h4 class="text-base font-bold text-white">{{ error.title }}</h4>
      <p class="text-xs text-slate-400 max-w-sm mx-auto">{{ error.message }}</p>
      <button
        class="px-4 py-2 rounded bg-sky-600 hover:bg-sky-500 text-white text-xs font-semibold transition"
        @click="loadRankings"
      >
        Retry
      </button>
    </div>

    <!-- Leaderboard Table -->
    <div v-else-if="leaderboard.length > 0" class="overflow-x-auto rounded-xl border border-slate-800">
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
          <tr v-for="(item, idx) in leaderboard" :key="item.id" class="hover:bg-slate-800/40 transition">
            <td class="py-4 px-4 text-center font-extrabold text-sm">
              <span v-if="idx === 0">🥇</span>
              <span v-else-if="idx === 1">🥈</span>
              <span v-else-if="idx === 2">🥉</span>
              <span v-else class="text-slate-500">#{{ idx + 1 }}</span>
            </td>
            <td class="py-4 px-4 font-bold text-white">
              <div class="flex items-center space-x-2">
                <span>{{ item.provider_name }}</span>
                <span class="text-slate-400 font-normal">({{ item.plan_code }})</span>
              </div>
            </td>
            <td class="py-4 px-4">
              {{ item.vcpu }} vCPU • {{ (item.ram_mb / 1024).toFixed(1) }} GB RAM • {{ item.disk_gb }} GB {{ item.disk_type }}
            </td>
            <td class="py-4 px-4 font-bold text-white">
              ${{ item.tco_result.subtotal_usd.toFixed(2) }}/mo
            </td>
            <td class="py-4 px-4 text-center">
              <span class="px-2.5 py-1 rounded-full font-bold bg-sky-950 text-sky-400 border border-sky-800 text-xs">
                {{ item.best_value_score.toFixed(2) }}
              </span>
            </td>
            <td class="py-4 px-4 text-center">
              <a
                :href="dealUrl(item.provider_slug)"
                class="px-3 py-1.5 rounded bg-sky-600 hover:bg-sky-500 text-white font-semibold text-[11px] transition shadow"
              >
                Inspect
              </a>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Empty -->
    <div v-else class="text-center py-16 bg-slate-900/40 rounded-xl border border-slate-800 space-y-3">
      <div class="text-3xl">🏆</div>
      <h4 class="text-base font-bold text-white">No ranked instances in this category yet</h4>
      <p class="text-xs text-slate-400 max-w-sm mx-auto">The catalog sync may still be running. Check back shortly.</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ApiError, useApi } from '~/composables/useApi'
import type { ListEnvelope, SearchResultItem } from '~/types/api'

definePageMeta({
  layout: 'public',
})

useHead({
  title: 'Best Value Cloud Rankings',
  meta: [
    { name: 'description', content: 'Algorithmic best-value cloud server rankings. Mathematical TCO scoring with zero pay-for-placement.' },
  ],
})

const categoryTabs = [
  { value: '', label: 'All Categories' },
  { value: 'VPS', label: 'VPS' },
  { value: 'DEDICATED_CPU_VM', label: 'Dedicated CPU' },
  { value: 'GPU_SERVER', label: 'GPU' },
]

const activeCat = ref('')
const loading = ref(false)
const error = ref<ApiError | null>(null)
const leaderboard = ref<SearchResultItem[]>([])

const { request, baseURL } = useApi()

function dealUrl(slug: string): string {
  return `${baseURL}/go/${slug}?src=ranking_table`
}

async function loadRankings() {
  loading.value = true
  error.value = null
  try {
    const res = await request<ListEnvelope<SearchResultItem>>('/api/v1/rankings/best-value', {
      query: { category: activeCat.value || undefined, limit: 20 },
    })
    leaderboard.value = res.data ?? []
  } catch (e) {
    error.value = e instanceof ApiError ? e : ApiError.fromUnknown(e)
    leaderboard.value = []
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadRankings()
})
</script>
