<template>
  <div class="space-y-12">
    <!-- Hero Section -->
    <section class="text-center py-16 space-y-6">
      <div class="inline-flex items-center space-x-2 px-3 py-1 rounded-full bg-sky-950 border border-sky-800 text-sky-400 text-xs font-semibold">
        <span>✨ 18+ Global Cloud Providers Cataloged</span>
      </div>
      <h1 class="text-4xl sm:text-6xl font-extrabold tracking-tight text-white max-w-4xl mx-auto leading-tight">
        Transparent Cloud Price Comparison & Bring-Your-Own-Cloud Management
      </h1>
      <p class="text-lg sm:text-xl text-slate-400 max-w-2xl mx-auto">
        Stop overpaying. Calculate real total-cost-of-ownership across CPU, RAM, NVMe, IPv4, and egress bandwidth with 36-dimension normalization.
      </p>
      <div class="flex flex-wrap justify-center gap-4 pt-4">
        <NuxtLink to="/search" class="px-6 py-3 rounded-lg bg-sky-600 hover:bg-sky-500 font-semibold text-white shadow-lg transition">
          Compare Cloud Instances
        </NuxtLink>
        <NuxtLink to="/calculator" class="px-6 py-3 rounded-lg bg-slate-800 hover:bg-slate-700 font-semibold text-slate-200 border border-slate-700 transition">
          Total TCO Calculator
        </NuxtLink>
      </div>
    </section>

    <!-- Quick Stats Grid -->
    <section class="grid grid-cols-2 md:grid-cols-4 gap-6">
      <div class="p-6 rounded-xl bg-slate-900 border border-slate-800 text-center">
        <div class="text-3xl font-bold text-sky-400 tabular-nums">18</div>
        <div class="text-xs text-slate-400 mt-1 uppercase tracking-wider">Cloud Providers</div>
      </div>
      <div class="p-6 rounded-xl bg-slate-900 border border-slate-800 text-center">
        <div class="text-3xl font-bold text-sky-400 tabular-nums">36</div>
        <div class="text-xs text-slate-400 mt-1 uppercase tracking-wider">Price Dimensions</div>
      </div>
      <div class="p-6 rounded-xl bg-slate-900 border border-slate-800 text-center">
        <div class="text-3xl font-bold text-sky-400 tabular-nums">&lt;50ms</div>
        <div class="text-xs text-slate-400 mt-1 uppercase tracking-wider">Search Latency</div>
      </div>
      <div class="p-6 rounded-xl bg-slate-900 border border-slate-800 text-center">
        <div class="text-3xl font-bold text-sky-400 tabular-nums">100%</div>
        <div class="text-xs text-slate-400 mt-1 uppercase tracking-wider">BYOC Safe Gated</div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'public'
})
</script>

