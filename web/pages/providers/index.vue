<template>
  <div class="space-y-8">
    <div class="space-y-2">
      <h1 class="text-3xl font-extrabold text-white tracking-tight">Cloud Providers Directory</h1>
      <p class="text-sm text-slate-400">
        18+ Global and regional cloud infrastructure providers evaluated and normalized across 36 price dimensions.
      </p>
    </div>

    <!-- Provider Cards Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div
        v-for="prov in providers"
        :key="prov.slug"
        class="p-6 rounded-xl bg-slate-900 border border-slate-800 hover:border-slate-700 transition flex flex-col justify-between space-y-4"
      >
        <div class="space-y-3">
          <div class="flex items-center justify-between">
            <span class="text-lg font-bold text-white">{{ prov.name }}</span>
            <span class="px-2 py-0.5 rounded text-[10px] font-semibold bg-sky-950 text-sky-400 border border-sky-800">
              {{ prov.country }}
            </span>
          </div>

          <p class="text-xs text-slate-400 leading-relaxed">
            {{ prov.description }}
          </p>

          <div class="flex flex-wrap gap-1.5 pt-1">
            <span
              v-for="cat in prov.categories"
              :key="cat"
              class="px-2 py-0.5 rounded text-[10px] bg-slate-800 text-slate-300 border border-slate-700"
            >
              {{ cat }}
            </span>
          </div>
        </div>

        <div class="pt-4 border-t border-slate-800 flex items-center justify-between">
          <span class="text-xs text-slate-500 tabular-nums">{{ prov.regionsCount }} Global Regions</span>
          <NuxtLink
            :to="`/providers/${prov.slug}`"
            class="text-xs font-semibold text-sky-400 hover:text-sky-300"
          >
            View SKUs &rarr;
          </NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'public'
})

const providers = [
  {
    slug: 'digitalocean',
    name: 'DigitalOcean',
    country: 'United States',
    description: 'Developer-friendly cloud platform offering predictable hourly/monthly capped pricing for droplets and Kubernetes.',
    categories: ['VPS', 'Dedicated CPU', 'Storage', 'Kubernetes'],
    regionsCount: 9,
  },
  {
    slug: 'hetzner',
    name: 'Hetzner Online',
    country: 'Germany',
    description: 'Unmatched European price-to-performance leader with generous 20TB traffic allocations and low-cost NVMe servers.',
    categories: ['VPS', 'Dedicated Root', 'Storage'],
    regionsCount: 5,
  },
  {
    slug: 'vultr',
    name: 'Vultr',
    country: 'United States',
    description: 'High-performance cloud compute with 32+ global datacenter locations, fractional GPU instances, and bare metal.',
    categories: ['VPS', 'Bare Metal', 'GPU Cloud', 'Kubernetes'],
    regionsCount: 32,
  },
  {
    slug: 'upcloud',
    name: 'UpCloud',
    country: 'Finland',
    description: 'Nordic cloud provider renowned for proprietary MaxIOPS storage technology offering exceptional I/O performance.',
    categories: ['VPS', 'Managed DB', 'Storage'],
    regionsCount: 12,
  },
  {
    slug: 'aws',
    name: 'Amazon Web Services (AWS)',
    country: 'United States',
    description: 'Global hyperscale market leader with comprehensive EC2 instance catalogs, separate IPv4 surcharges, and tiered egress.',
    categories: ['VPS', 'Compute Optimized', 'GPU', 'Enterprise'],
    regionsCount: 34,
  },
  {
    slug: 'gcp',
    name: 'Google Cloud Platform (GCP)',
    country: 'United States',
    description: 'Global cloud provider with industry-leading private network backbone, custom machine types, and TPU accelerators.',
    categories: ['Compute Engine', 'TPU/GPU', 'Cloud Storage'],
    regionsCount: 40,
  },
]
</script>
