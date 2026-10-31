<template>
  <div class="space-y-8">
    <!-- Provider Header -->
    <div class="p-8 rounded-2xl bg-slate-900 border border-slate-800 space-y-4">
      <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <div class="flex items-center space-x-3">
            <h1 class="text-3xl font-extrabold text-white tracking-tight">{{ provider.name }}</h1>
            <span class="px-2.5 py-0.5 rounded text-xs font-semibold bg-sky-950 text-sky-400 border border-sky-800">
              Verified Provider
            </span>
          </div>
          <p class="text-sm text-slate-400 mt-2 max-w-2xl leading-relaxed">
            {{ provider.description }}
          </p>
        </div>

        <NuxtLink
          :to="`/go/${slug}?src=provider_profile`"
          class="px-6 py-3 rounded-lg bg-sky-600 hover:bg-sky-500 font-semibold text-white text-sm transition shadow-lg shrink-0"
        >
          Deploy on {{ provider.name }}
        </NuxtLink>
      </div>

      <!-- Quick Metrics -->
      <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 pt-4 border-t border-slate-800 text-xs">
        <div>
          <span class="text-slate-500 block">Pricing Model</span>
          <span class="text-white font-semibold">{{ provider.pricingModel }}</span>
        </div>
        <div>
          <span class="text-slate-500 block">Primary Storage</span>
          <span class="text-white font-semibold">{{ provider.storageType }}</span>
        </div>
        <div>
          <span class="text-slate-500 block">Free Egress Quota</span>
          <span class="text-white font-semibold">{{ provider.freeEgress }}</span>
        </div>
        <div>
          <span class="text-slate-500 block">IPv4 Policy</span>
          <span class="text-white font-semibold">{{ provider.ipv4Policy }}</span>
        </div>
      </div>
    </div>

    <!-- SKUs Catalog Table -->
    <div class="space-y-4">
      <h2 class="text-xl font-bold text-white">Popular Instance Types & Rates</h2>
      <div class="overflow-x-auto rounded-xl border border-slate-800">
        <table class="w-full text-left text-xs bg-slate-900 divide-y divide-slate-800">
          <thead class="bg-slate-950 text-slate-400 font-semibold uppercase tracking-wider">
            <tr>
              <th class="py-3 px-4">Plan Code</th>
              <th class="py-3 px-4">vCPU</th>
              <th class="py-3 px-4">Memory</th>
              <th class="py-3 px-4">NVMe Disk</th>
              <th class="py-3 px-4">Bandwidth</th>
              <th class="py-3 px-4">Hourly Rate</th>
              <th class="py-3 px-4 text-right">Monthly Cap</th>
              <th class="py-3 px-4 text-center">Action</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-800 text-slate-300 tabular-nums">
            <tr v-for="sku in provider.skus" :key="sku.code" class="hover:bg-slate-800/40 transition">
              <td class="py-3.5 px-4 font-bold text-white">{{ sku.code }}</td>
              <td class="py-3.5 px-4">{{ sku.vcpu }} vCPU</td>
              <td class="py-3.5 px-4">{{ sku.ram }} GB</td>
              <td class="py-3.5 px-4">{{ sku.disk }} GB NVMe</td>
              <td class="py-3.5 px-4">{{ sku.bandwidth }} TB</td>
              <td class="py-3.5 px-4">${{ sku.hourly.toFixed(4) }}/hr</td>
              <td class="py-3.5 px-4 text-right font-extrabold text-white">${{ sku.monthly.toFixed(2) }}/mo</td>
              <td class="py-3.5 px-4 text-center">
                <NuxtLink
                  :to="`/go/${slug}?src=sku_table`"
                  class="px-3 py-1 rounded bg-sky-600 hover:bg-sky-500 text-white font-semibold text-[11px] transition"
                >
                  Choose
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
import { computed } from 'vue'

definePageMeta({
  layout: 'public'
})

const route = useRoute()
const slug = computed(() => (route.params.slug as string) || 'digitalocean')

const providerMap: Record<string, any> = {
  digitalocean: {
    name: 'DigitalOcean',
    description: 'Cloud hosting with predictable monthly caps, fast NVMe SSD storage, and generous multi-terabyte outbound transfer pools.',
    pricingModel: '672h Monthly Cap',
    storageType: 'Local NVMe SSD',
    freeEgress: '1TB - 10TB per droplet',
    ipv4Policy: 'Primary IPv4 Included',
    skus: [
      { code: 's-1vcpu-1gb', vcpu: 1, ram: 1, disk: 25, bandwidth: 1, hourly: 0.0089, monthly: 6.00 },
      { code: 's-1vcpu-2gb', vcpu: 1, ram: 2, disk: 50, bandwidth: 2, hourly: 0.0179, monthly: 12.00 },
      { code: 's-2vcpu-4gb', vcpu: 2, ram: 4, disk: 80, bandwidth: 4, hourly: 0.0357, monthly: 24.00 },
      { code: 's-4vcpu-8gb', vcpu: 4, ram: 8, disk: 160, bandwidth: 5, hourly: 0.0714, monthly: 48.00 },
    ],
  },
  hetzner: {
    name: 'Hetzner Cloud',
    description: 'European market value benchmark with cutting-edge AMD EPYC & Intel hardware in Germany, Finland, and USA locations.',
    pricingModel: 'Hourly with Low Monthly Cap',
    storageType: 'Fast NVMe RAID',
    freeEgress: '20TB Free per server',
    ipv4Policy: 'IPv4 Included in CPX series',
    skus: [
      { code: 'cpx11', vcpu: 2, ram: 2, disk: 40, bandwidth: 20, hourly: 0.0074, monthly: 4.95 },
      { code: 'cpx21', vcpu: 3, ram: 4, disk: 80, bandwidth: 20, hourly: 0.0134, monthly: 8.95 },
      { code: 'cpx31', vcpu: 4, ram: 8, disk: 160, bandwidth: 20, hourly: 0.0246, monthly: 16.45 },
    ],
  },
}

const provider = computed(() => {
  return providerMap[slug.value] || providerMap.digitalocean
})
</script>
