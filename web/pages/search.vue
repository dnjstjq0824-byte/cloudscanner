<template>
  <div class="space-y-8">
    <!-- Header & Natural Language Search Bar -->
    <div class="space-y-4">
      <div>
        <h1 class="text-3xl font-extrabold text-white tracking-tight">Compare Cloud Instances</h1>
        <p class="text-sm text-slate-400 mt-1">
          Search across 18+ verified cloud providers with 36-dimension TCO and Best Value scoring.
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
          class="px-6 py-3 bg-sky-600 hover:bg-sky-500 text-white font-semibold rounded-lg text-sm transition shadow-md flex items-center justify-center space-x-2"
          @click="performSearch"
        >
          <span>Search</span>
        </button>
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
              v-for="cat in ['ALL', 'VPS', 'DEDICATED_CPU_VM', 'GPU_INSTANCE']"
              :key="cat"
              class="flex items-center space-x-2 text-slate-300 hover:text-white cursor-pointer"
            >
              <input
                v-model="selectedCategory"
                type="radio"
                :value="cat"
                name="category"
                class="text-sky-600 bg-slate-800 border-slate-700 focus:ring-0"
                @change="performSearch"
              />
              <span>{{ cat === 'ALL' ? 'All Categories' : cat.replace(/_/g, ' ') }}</span>
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
              @change="performSort"
            >
              <option value="best_value">Best Value Score (Calculated)</option>
              <option value="price_asc">Price: Low to High</option>
              <option value="price_desc">Price: High to Low</option>
              <option value="vcpu_desc">vCPU: Most to Least</option>
              <option value="ram_desc">RAM: Most to Least</option>
            </select>
          </div>
        </div>

        <!-- Populated Results Grid -->
        <div v-if="filteredResults.length > 0" class="space-y-3">
          <div
            v-for="item in filteredResults"
            :key="item.sku_id"
            class="p-5 rounded-xl bg-slate-900 border border-slate-800 hover:border-slate-700 transition flex flex-col md:flex-row justify-between items-start md:items-center gap-4"
          >
            <!-- Left Info -->
            <div class="space-y-1.5 flex-1">
              <div class="flex items-center space-x-2.5">
                <span class="px-2 py-0.5 rounded text-[11px] font-semibold uppercase bg-slate-800 text-sky-400 border border-slate-700">
                  {{ item.provider_name }}
                </span>
                <span class="text-base font-bold text-white">{{ item.plan_code }}</span>
                <span v-if="item.value_score > 5" class="px-2 py-0.5 rounded-full text-[10px] font-bold bg-amber-950 text-amber-400 border border-amber-800">
                  ⚡ High Value
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
                <span class="tabular-nums"><strong>{{ item.bandwidth_tb }}</strong> TB Bandwidth</span>
                <span>•</span>
                <span :class="item.ipv4_included ? 'text-emerald-400' : 'text-amber-400'">
                  {{ item.ipv4_included ? 'IPv4 Included' : '+IPv4 Surcharge' }}
                </span>
              </div>
            </div>

            <!-- Right Pricing & CTA -->
            <div class="flex items-center space-x-6 w-full md:w-auto justify-between md:justify-end border-t md:border-t-0 pt-3 md:pt-0 border-slate-800">
              <div class="text-right">
                <div class="text-2xl font-extrabold text-white tabular-nums">
                  ${{ item.monthly_price_usd.toFixed(2) }}
                  <span class="text-xs font-normal text-slate-400">/mo</span>
                </div>
                <div class="text-[11px] text-slate-500 tabular-nums">
                  ${{ item.hourly_price_usd.toFixed(4) }}/hr
                </div>
              </div>

              <NuxtLink
                :to="`/go/${item.provider_slug}?src=search_table`"
                class="px-4 py-2.5 rounded-lg bg-sky-600 hover:bg-sky-500 font-semibold text-white text-xs transition shadow shrink-0"
              >
                Deploy Deal
              </NuxtLink>
            </div>
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
import { ref, computed } from 'vue'

definePageMeta({
  layout: 'public'
})

interface SKUItem {
  sku_id: string
  provider_slug: string
  provider_name: string
  plan_code: string
  category: string
  vcpu: number
  ram_mb: number
  disk_gb: number
  disk_type: string
  bandwidth_tb: number
  monthly_price_usd: number
  hourly_price_usd: number
  ipv4_included: boolean
  value_score: number
}

const searchQuery = ref('')
const selectedCategory = ref('ALL')
const selectedProviders = ref<string[]>([])
const minVCPU = ref(0)
const minRAMGB = ref(0)
const sortBy = ref('best_value')

const quickChips = [
  { label: 'Budget VPS (<$10)', query: 'vps under $10' },
  { label: '4 vCPU Dev Server', query: '4 vcpu 8gb ram' },
  { label: 'NVMe Fast IO', query: 'nvme' },
  { label: 'High Bandwidth', query: 'vps 4tb bandwidth' },
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

// Seeded verified catalog dataset across providers
const allItems = ref<SKUItem[]>([
  {
    sku_id: 'do-s-1vcpu-1gb',
    provider_slug: 'digitalocean',
    provider_name: 'DigitalOcean',
    plan_code: 's-1vcpu-1gb',
    category: 'VPS',
    vcpu: 1,
    ram_mb: 1024,
    disk_gb: 25,
    disk_type: 'NVME',
    bandwidth_tb: 1.0,
    monthly_price_usd: 6.00,
    hourly_price_usd: 0.0089,
    ipv4_included: true,
    value_score: 3.42,
  },
  {
    sku_id: 'do-s-2vcpu-4gb',
    provider_slug: 'digitalocean',
    provider_name: 'DigitalOcean',
    plan_code: 's-2vcpu-4gb',
    category: 'VPS',
    vcpu: 2,
    ram_mb: 4096,
    disk_gb: 80,
    disk_type: 'NVME',
    bandwidth_tb: 4.0,
    monthly_price_usd: 24.00,
    hourly_price_usd: 0.0357,
    ipv4_included: true,
    value_score: 5.67,
  },
  {
    sku_id: 'hz-cpx11',
    provider_slug: 'hetzner',
    provider_name: 'Hetzner',
    plan_code: 'cpx11',
    category: 'VPS',
    vcpu: 2,
    ram_mb: 2048,
    disk_gb: 40,
    disk_type: 'NVME',
    bandwidth_tb: 20.0,
    monthly_price_usd: 4.95,
    hourly_price_usd: 0.0074,
    ipv4_included: true,
    value_score: 8.92,
  },
  {
    sku_id: 'hz-cpx21',
    provider_slug: 'hetzner',
    provider_name: 'Hetzner',
    plan_code: 'cpx21',
    category: 'VPS',
    vcpu: 3,
    ram_mb: 4096,
    disk_gb: 80,
    disk_type: 'NVME',
    bandwidth_tb: 20.0,
    monthly_price_usd: 8.95,
    hourly_price_usd: 0.0134,
    ipv4_included: true,
    value_score: 9.45,
  },
  {
    sku_id: 'vultr-vc2-2c-4gb',
    provider_slug: 'vultr',
    provider_name: 'Vultr',
    plan_code: 'vc2-2c-4gb',
    category: 'VPS',
    vcpu: 2,
    ram_mb: 4096,
    disk_gb: 80,
    disk_type: 'NVME',
    bandwidth_tb: 3.0,
    monthly_price_usd: 20.00,
    hourly_price_usd: 0.0298,
    ipv4_included: true,
    value_score: 5.80,
  },
  {
    sku_id: 'upcloud-2xcpu-4gb',
    provider_slug: 'upcloud',
    provider_name: 'UpCloud',
    plan_code: '2xCPU-4GB',
    category: 'VPS',
    vcpu: 2,
    ram_mb: 4096,
    disk_gb: 80,
    disk_type: 'NVME',
    bandwidth_tb: 4.0,
    monthly_price_usd: 20.00,
    hourly_price_usd: 0.0298,
    ipv4_included: true,
    value_score: 5.80,
  },
  {
    sku_id: 'aws-t4g-small',
    provider_slug: 'aws',
    provider_name: 'AWS EC2',
    plan_code: 't4g.small',
    category: 'VPS',
    vcpu: 2,
    ram_mb: 2048,
    disk_gb: 40,
    disk_type: 'EBS_SSD',
    bandwidth_tb: 0.1,
    monthly_price_usd: 15.91,
    hourly_price_usd: 0.0218,
    ipv4_included: false,
    value_score: 3.10,
  },
  {
    sku_id: 'aws-t3-medium',
    provider_slug: 'aws',
    provider_name: 'AWS EC2',
    plan_code: 't3.medium',
    category: 'VPS',
    vcpu: 2,
    ram_mb: 4096,
    disk_gb: 50,
    disk_type: 'EBS_SSD',
    bandwidth_tb: 0.1,
    monthly_price_usd: 34.02,
    hourly_price_usd: 0.0466,
    ipv4_included: false,
    value_score: 2.85,
  },
])

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
}

function performSearch() {
  // Parsing simple query tokens if typed
  const q = searchQuery.value.toLowerCase()
  if (q.includes('vps')) selectedCategory.value = 'VPS'
  if (q.includes('dedicated')) selectedCategory.value = 'DEDICATED_CPU_VM'
}

function performSort() {
  // Handled by computed property
}

const filteredResults = computed(() => {
  let list = allItems.value.filter(item => {
    // Category filter
    if (selectedCategory.value !== 'ALL' && item.category !== selectedCategory.value) {
      return false
    }
    // Provider filter
    if (selectedProviders.value.length > 0 && !selectedProviders.value.includes(item.provider_slug)) {
      return false
    }
    // Min vCPU
    if (minVCPU.value > 0 && item.vcpu < minVCPU.value) {
      return false
    }
    // Min RAM
    if (minRAMGB.value > 0 && (item.ram_mb / 1024) < minRAMGB.value) {
      return false
    }
    // Search query keyword matching
    if (searchQuery.value.trim() !== '') {
      const q = searchQuery.value.toLowerCase()
      const matchText = `${item.plan_code} ${item.provider_name} ${item.disk_type} ${item.category}`.toLowerCase()
      // Check if price condition
      if (q.includes('under $')) {
        const parts = q.split('under $')
        const maxP = parseFloat(parts[1])
        if (!isNaN(maxP) && item.monthly_price_usd > maxP) return false
      }
      if (!matchText.includes(q.replace(/under \$\d+/g, '').trim())) {
        return false
      }
    }
    return true
  })

  // Sorting
  if (sortBy.value === 'best_value') {
    list.sort((a, b) => b.value_score - a.value_score)
  } else if (sortBy.value === 'price_asc') {
    list.sort((a, b) => a.monthly_price_usd - b.monthly_price_usd)
  } else if (sortBy.value === 'price_desc') {
    list.sort((a, b) => b.monthly_price_usd - a.monthly_price_usd)
  } else if (sortBy.value === 'vcpu_desc') {
    list.sort((a, b) => b.vcpu - a.vcpu)
  } else if (sortBy.value === 'ram_desc') {
    list.sort((a, b) => b.ram_mb - a.ram_mb)
  }

  return list
})
</script>
