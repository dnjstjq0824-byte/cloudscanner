<template>
  <div class="space-y-8">
    <!-- Header & Natural Language Search Bar -->
    <div class="space-y-4">
      <div>
        <h1 class="text-3xl font-extrabold text-white tracking-tight">Compare Cloud Instances</h1>
        <p class="text-sm text-slate-400 mt-1">
          Search across 18+ verified cloud providers with 36-dimension TCO and Best Value scoring.
          <span v-if="lastExecutionMs !== null" class="text-slate-500">· {{ lastExecutionMs }}ms</span>
        </p>
      </div>

      <div class="flex flex-col sm:flex-row gap-3">
        <div class="relative flex-1">
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Try: '2 vcpu 4gb ram under $30 in frankfurt' or 'nvme vps'"
            class="w-full pl-4 pr-10 py-3 bg-slate-900 border border-slate-700 rounded-lg text-white placeholder-slate-500 focus:outline-none focus:border-sky-500 transition text-sm"
            @keyup.enter="performSearch"
          />
          <button
            v-if="searchQuery"
            class="absolute right-3 top-3 text-slate-500 hover:text-slate-300 text-sm"
            @click="searchQuery = ''; performSearch()"
          >
            ✕
          </button>
        </div>
        <button
          class="px-6 py-3 bg-sky-600 hover:bg-sky-500 text-white font-semibold rounded-lg text-sm transition shadow-md flex items-center justify-center space-x-2 disabled:opacity-60"
          :disabled="loading"
          @click="performSearch"
        >
          <span>{{ loading ? 'Searching…' : 'Search' }}</span>
        </button>
      </div>

      <!-- Parsed NL query hints -->
      <div v-if="parsedQuery" class="flex flex-wrap gap-2 text-[11px]">
        <span class="text-slate-500 self-center">Detected:</span>
        <span v-if="parsedQuery.detected_vcpu" class="px-2 py-0.5 rounded bg-sky-950 text-sky-400 border border-sky-800">{{ parsedQuery.detected_vcpu }} vCPU</span>
        <span v-if="parsedQuery.detected_ram_gb" class="px-2 py-0.5 rounded bg-sky-950 text-sky-400 border border-sky-800">{{ parsedQuery.detected_ram_gb }} GB RAM</span>
        <span v-if="parsedQuery.detected_disk_type" class="px-2 py-0.5 rounded bg-sky-950 text-sky-400 border border-sky-800">{{ parsedQuery.detected_disk_type }}</span>
        <span v-if="parsedQuery.detected_region" class="px-2 py-0.5 rounded bg-sky-950 text-sky-400 border border-sky-800">{{ parsedQuery.detected_region }}</span>
        <span v-if="parsedQuery.detected_max_price" class="px-2 py-0.5 rounded bg-sky-950 text-sky-400 border border-sky-800">≤ ${{ parsedQuery.detected_max_price }}</span>
      </div>

      <!-- Quick Preset Filter Chips -->
      <div class="flex flex-wrap gap-2 text-xs">
        <span class="text-slate-500 self-center">Popular:</span>
        <button
          v-for="chip in quickChips"
          :key="chip.label"
          class="px-2.5 py-1 rounded bg-slate-800 text-slate-300 hover:bg-slate-700 border border-slate-700 transition"
          @click="applyChip(chip.query)"
        >
          {{ chip.label }}
        </button>
      </div>
    </div>

    <!-- Main Content Layout (Sidebar + Results) -->
    <div class="grid grid-cols-1 lg:grid-cols-4 gap-8">
      <!-- Filter Sidebar -->
      <aside class="space-y-6 bg-slate-900/60 p-5 rounded-xl border border-slate-800 h-fit">
        <div>
          <h3 class="text-xs font-bold uppercase tracking-wider text-slate-400 mb-3">Category</h3>
          <div class="space-y-1.5 text-sm">
            <label
              v-for="cat in categoryOptions"
              :key="cat.value"
              class="flex items-center space-x-2 text-slate-300 hover:text-white cursor-pointer"
            >
              <input
                v-model="selectedCategory"
                type="radio"
                :value="cat.value"
                name="category"
                class="text-sky-600 bg-slate-800 border-slate-700 focus:ring-0"
                @change="performSearch"
              />
              <span>{{ cat.label }}</span>
              <span v-if="facets && cat.value !== 'ALL'" class="text-[11px] text-slate-500 tabular-nums">({{ facets.categories[cat.value] || 0 }})</span>
            </label>
          </div>
        </div>

        <div class="border-t border-slate-800 pt-4">
          <h3 class="text-xs font-bold uppercase tracking-wider text-slate-400 mb-3">Provider</h3>
          <div class="space-y-1.5 text-sm max-h-48 overflow-y-auto">
            <label
              v-for="prov in providersList"
              :key="prov.slug"
              class="flex items-center space-x-2 text-slate-300 hover:text-white cursor-pointer"
            >
              <input
                v-model="selectedProviders"
                type="checkbox"
                :value="prov.slug"
                class="rounded text-sky-600 bg-slate-800 border-slate-700 focus:ring-0"
                @change="performSearch"
              />
              <span>{{ prov.name }}</span>
            </label>
          </div>
        </div>

        <div class="border-t border-slate-800 pt-4">
          <h3 class="text-xs font-bold uppercase tracking-wider text-slate-400 mb-3">Min vCPU</h3>
          <div class="grid grid-cols-4 gap-1.5">
            <button
              v-for="c in [0, 1, 2, 4, 8]"
              :key="c"
              class="py-1 text-xs rounded border transition"
              :class="minVCPU === c ? 'bg-sky-600 border-sky-500 text-white font-bold' : 'bg-slate-800 border-slate-700 text-slate-300 hover:bg-slate-700'"
              @click="minVCPU = c; performSearch()"
            >
              {{ c === 0 ? 'Any' : c + 'C' }}
            </button>
          </div>
        </div>

        <div class="border-t border-slate-800 pt-4">
          <h3 class="text-xs font-bold uppercase tracking-wider text-slate-400 mb-3">Min RAM (GB)</h3>
          <div class="grid grid-cols-4 gap-1.5">
            <button
              v-for="r in [0, 1, 2, 4, 8, 16]"
              :key="r"
              class="py-1 text-xs rounded border transition"
              :class="minRAMGB === r ? 'bg-sky-600 border-sky-500 text-white font-bold' : 'bg-slate-800 border-slate-700 text-slate-300 hover:bg-slate-700'"
              @click="minRAMGB = r; performSearch()"
            >
              {{ r === 0 ? 'Any' : r + 'G' }}
            </button>
          </div>
        </div>

        <div class="border-t border-slate-800 pt-4">
          <button
            class="w-full py-2 text-xs font-semibold rounded bg-slate-800 hover:bg-slate-700 text-slate-400 hover:text-white border border-slate-700 transition"
            @click="resetFilters"
          >
            Reset All Filters
          </button>
        </div>
      </aside>

      <!-- Results Area -->
      <main class="lg:col-span-3 space-y-4">
        <!-- Results Header Controls -->
        <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-2 pb-2 border-b border-slate-800 text-sm">
          <div class="text-slate-400">
            Showing <span class="font-bold text-white tabular-nums">{{ filteredResults.length }}</span> instance matches
          </div>
          <div class="flex items-center space-x-2 text-xs">
            <span class="text-slate-500">Sort by:</span>
            <select
              v-model="sortBy"
              class="bg-slate-900 border border-slate-700 rounded px-2.5 py-1 text-slate-200 focus:outline-none"
              @change="performSearch"
            >
              <option value="best_value">Best Value Score (Calculated)</option>
              <option value="price_asc">Price: Low to High</option>
              <option value="price_desc">Price: High to Low</option>
              <option value="vcpu_desc">vCPU: Most to Least</option>
              <option value="ram_desc">RAM: Most to Least</option>
            </select>
          </div>
        </div>

        <!-- Loading Skeleton -->
        <div v-if="loading" class="space-y-3" aria-busy="true" aria-label="Loading results">
          <div v-for="n in 4" :key="n" class="p-5 rounded-xl bg-slate-900 border border-slate-800 animate-pulse">
            <div class="h-4 w-1/3 rounded bg-slate-800" />
            <div class="mt-3 h-3 w-2/3 rounded bg-slate-800" />
          </div>
        </div>

        <!-- Error State -->
        <div v-else-if="error" class="text-center py-16 bg-red-950/20 rounded-xl border border-red-900 space-y-3">
          <div class="text-3xl">⚠️</div>
          <h4 class="text-base font-bold text-white">{{ error.title }}</h4>
          <p class="text-xs text-slate-400 max-w-sm mx-auto">{{ error.message }}</p>
          <button
            class="px-4 py-2 rounded bg-sky-600 hover:bg-sky-500 text-white text-xs font-semibold transition"
            @click="performSearch"
          >
            Retry Search
          </button>
        </div>

        <!-- Populated Results Grid -->
        <div v-else-if="filteredResults.length > 0" class="space-y-3">
          <div
            v-for="item in pagedResults"
            :key="item.id"
            class="p-5 rounded-xl bg-slate-900 border border-slate-800 hover:border-slate-700 transition flex flex-col md:flex-row justify-between items-start md:items-center gap-4"
          >
            <!-- Left Info -->
            <div class="space-y-1.5 flex-1">
              <div class="flex items-center space-x-2.5 flex-wrap">
                <span class="px-2 py-0.5 rounded text-[11px] font-semibold uppercase bg-slate-800 text-sky-400 border border-slate-700">
                  {{ item.provider_name }}
                </span>
                <span class="text-base font-bold text-white">{{ item.plan_code }}</span>
                <span v-if="item.best_value_score > 5" class="px-2 py-0.5 rounded-full text-[10px] font-bold bg-amber-950 text-amber-400 border border-amber-800">
                  ⚡ High Value
                </span>
                <span
                  class="px-2 py-0.5 rounded-full text-[10px] font-semibold border"
                  :class="item.freshness_status === 'STALE' ? 'bg-amber-950 text-amber-400 border-amber-800' : 'bg-emerald-950 text-emerald-400 border-emerald-800'"
                >
                  {{ item.freshness_status || 'FRESH' }}
                </span>
              </div>

              <!-- Hardware Specs Grid -->
              <div class="flex flex-wrap gap-x-4 gap-y-1 text-xs text-slate-400 pt-1">
                <span class="tabular-nums"><strong>{{ item.vcpu }}</strong> vCPU</span>
                <span>•</span>
                <span class="tabular-nums"><strong>{{ (item.ram_mb / 1024).toFixed(1) }}</strong> GB RAM</span>
                <span>•</span>
                <span class="tabular-nums"><strong>{{ item.disk_gb }}</strong> GB {{ item.disk_type }}</span>
                <span>•</span>
                <span class="tabular-nums"><strong>{{ (item.included_bandwidth_gb / 1024).toFixed(1) }}</strong> TB Bandwidth</span>
                <span>•</span>
                <span :class="item.ipv4_included_count > 0 ? 'text-emerald-400' : 'text-amber-400'">
                  {{ item.ipv4_included_count > 0 ? 'IPv4 Included' : '+IPv4 Surcharge' }}
                </span>
              </div>
              <div class="text-[11px] text-slate-500 tabular-nums">
                TCO ${{ item.tco_result.subtotal_usd.toFixed(2) }}/mo · score {{ item.best_value_score.toFixed(2) }}
              </div>
            </div>

            <!-- Right Pricing & CTA -->
            <div class="flex items-center space-x-6 w-full md:w-auto justify-between md:justify-end border-t md:border-t-0 pt-3 md:pt-0 border-slate-800">
              <div class="text-right">
                <div class="text-2xl font-extrabold text-white tabular-nums">
                  ${{ item.tco_result.subtotal_usd.toFixed(2) }}
                  <span class="text-xs font-normal text-slate-400">/mo</span>
                </div>
                <div class="text-[11px] text-slate-500 tabular-nums">
                  ${{ item.raw_hourly_price_usd.toFixed(4) }}/hr
                </div>
              </div>

              <a
                :href="dealUrl(item.provider_slug)"
                class="px-4 py-2.5 rounded-lg bg-sky-600 hover:bg-sky-500 font-semibold text-white text-xs transition shadow shrink-0"
              >
                Deploy Deal
              </a>
            </div>
          </div>

          <!-- Pagination -->
          <div v-if="totalPages > 1" class="flex items-center justify-center gap-2 pt-2 text-xs">
            <button
              class="px-3 py-1.5 rounded bg-slate-800 border border-slate-700 text-slate-300 disabled:opacity-40"
              :disabled="page <= 1"
              @click="page--; scrollTop()"
            >
              ← Prev
            </button>
            <span class="text-slate-400 tabular-nums">Page {{ page }} / {{ totalPages }}</span>
            <button
              class="px-3 py-1.5 rounded bg-slate-800 border border-slate-700 text-slate-300 disabled:opacity-40"
              :disabled="page >= totalPages"
              @click="page++; scrollTop()"
            >
              Next →
            </button>
          </div>
        </div>

        <!-- Empty State -->
        <div v-else class="text-center py-16 bg-slate-900/40 rounded-xl border border-slate-800 space-y-3">
          <div class="text-3xl">🔍</div>
          <h4 class="text-base font-bold text-white">No cloud instances matched your filters</h4>
          <p class="text-xs text-slate-400 max-w-sm mx-auto">
            Try adjusting your query, widening the CPU/RAM constraints, or clearing specific provider selections.
          </p>
          <button
            class="px-4 py-2 rounded bg-slate-800 hover:bg-slate-700 text-sky-400 text-xs font-semibold transition"
            @click="resetFilters"
          >
            Clear Filters
          </button>
        </div>
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { ApiError, useApi } from '~/composables/useApi'
import type { FacetCounts, ParsedQuery, SearchResponse, SearchResultItem } from '~/types/api'

definePageMeta({
  layout: 'public',
})

useHead({
  title: 'Compare Cloud Instances',
  meta: [
    { name: 'description', content: 'Search and compare cloud VPS instances across 18+ providers with transparent 36-dimension TCO and best-value scoring.' },
  ],
})

const searchQuery = ref('')
const selectedCategory = ref('ALL')
const selectedProviders = ref<string[]>([])
const minVCPU = ref(0)
const minRAMGB = ref(0)
const sortBy = ref('best_value')
const page = ref(1)
const pageSize = 20

const loading = ref(false)
const error = ref<ApiError | null>(null)
const results = ref<SearchResultItem[]>([])
const facets = ref<FacetCounts | null>(null)
const parsedQuery = ref<ParsedQuery | null>(null)
const lastExecutionMs = ref<number | null>(null)

const quickChips = [
  { label: 'Budget VPS (<$10)', query: 'vps under $10' },
  { label: '4 vCPU Dev Server', query: '4 vcpu 8gb ram' },
  { label: 'NVMe Fast IO', query: 'nvme' },
  { label: 'High Bandwidth', query: 'vps 4tb bandwidth' },
]

const categoryOptions = [
  { value: 'ALL', label: 'All Categories' },
  { value: 'VPS', label: 'VPS' },
  { value: 'DEDICATED_CPU_VM', label: 'DEDICATED CPU VM' },
  { value: 'GPU_SERVER', label: 'GPU SERVER' },
]

const providersList = [
  { slug: 'digitalocean', name: 'DigitalOcean' },
  { slug: 'hetzner', name: 'Hetzner Cloud' },
  { slug: 'vultr', name: 'Vultr' },
  { slug: 'upcloud', name: 'UpCloud' },
  { slug: 'aws', name: 'AWS EC2' },
  { slug: 'gcp', name: 'Google Cloud' },
  { slug: 'azure', name: 'Microsoft Azure' },
  { slug: 'linode', name: 'Linode (Akamai)' },
  { slug: 'ovhcloud', name: 'OVHcloud' },
  { slug: 'scaleway', name: 'Scaleway' },
]

const { request, baseURL } = useApi()

function dealUrl(slug: string): string {
  return `${baseURL}/go/${slug}?src=search_table`
}

function applyChip(q: string) {
  searchQuery.value = q
  performSearch()
}

function resetFilters() {
  searchQuery.value = ''
  selectedCategory.value = 'ALL'
  selectedProviders.value = []
  minVCPU.value = 0
  minRAMGB.value = 0
  sortBy.value = 'best_value'
  page.value = 1
  performSearch()
}

function scrollTop() {
  if (import.meta.client) {
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }
}

async function performSearch() {
  loading.value = true
  error.value = null
  page.value = 1
  try {
    const res = await request<SearchResponse>('/api/v1/search/skus', {
      query: {
        q: searchQuery.value.trim() || undefined,
        category: selectedCategory.value === 'ALL' ? undefined : selectedCategory.value,
        provider: selectedProviders.value.length === 1 ? selectedProviders.value[0] : undefined,
        min_vcpu: minVCPU.value || undefined,
        min_ram_mb: minRAMGB.value ? minRAMGB.value * 1024 : undefined,
        sort_by: sortBy.value,
        limit: 100,
      },
    })
    results.value = res.data ?? []
    facets.value = res.facets
    parsedQuery.value = res.parsed_query && Object.keys(res.parsed_query).length > 0 ? res.parsed_query : null
    lastExecutionMs.value = res.execution_ms
  } catch (e) {
    error.value = e instanceof ApiError ? e : ApiError.fromUnknown(e)
    results.value = []
  } finally {
    loading.value = false
  }
}

/** Multi-provider selections are intersected client-side over the server result set. */
const filteredResults = computed(() => {
  if (selectedProviders.value.length <= 1) {
    return results.value
  }
  const wanted = new Set(selectedProviders.value)
  return results.value.filter((item) => wanted.has(item.provider_slug))
})

const totalPages = computed(() => Math.max(1, Math.ceil(filteredResults.value.length / pageSize)))

const pagedResults = computed(() => {
  const start = (page.value - 1) * pageSize
  return filteredResults.value.slice(start, start + pageSize)
})

onMounted(() => {
  performSearch()
})
</script>
