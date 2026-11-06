<template>
  <div class="flex flex-col min-h-screen">
    <!-- Navigation Header -->
    <header class="border-b border-slate-800 bg-slate-900/80 backdrop-blur sticky top-0 z-50">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
        <div class="flex items-center space-x-6">
          <NuxtLink to="/" class="flex items-center space-x-2 text-xl font-bold text-sky-400">
            <span class="text-2xl">⚡</span>
            <span>CloudScanr</span>
          </NuxtLink>
          <nav class="hidden md:flex space-x-4 text-sm font-medium text-slate-300">
            <NuxtLink to="/search" class="hover:text-white transition">Search & Compare</NuxtLink>
            <NuxtLink to="/calculator" class="hover:text-white transition">TCO Calculator</NuxtLink>
            <NuxtLink to="/providers" class="hover:text-white transition">Providers</NuxtLink>
            <NuxtLink to="/rankings" class="hover:text-white transition">Best Value</NuxtLink>
          </nav>
        </div>
        <div class="flex items-center space-x-3">
          <template v-if="auth.isAuthenticated">
            <NuxtLink to="/dashboard" class="px-3 py-1.5 text-sm font-medium text-slate-300 hover:text-white">Console</NuxtLink>
            <span class="hidden sm:inline text-xs text-slate-500 truncate max-w-40">{{ auth.user?.email }}</span>
            <button class="px-4 py-1.5 text-sm font-medium bg-slate-800 hover:bg-slate-700 text-slate-200 rounded-md transition" @click="auth.logout()">Sign Out</button>
          </template>
          <template v-else>
            <NuxtLink to="/dashboard" class="px-3 py-1.5 text-sm font-medium text-slate-300 hover:text-white">Console</NuxtLink>
            <NuxtLink to="/auth/login" class="px-4 py-1.5 text-sm font-medium bg-sky-600 hover:bg-sky-500 text-white rounded-md transition shadow">Sign In</NuxtLink>
          </template>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <main class="flex-1 max-w-7xl w-full mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <slot />
    </main>

    <!-- Footer with FTC / Fair-Trade Disclosures -->
    <footer class="border-t border-slate-800 bg-slate-900 py-10 text-xs text-slate-400">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 space-y-4">
        <div class="flex flex-col md:flex-row justify-between items-center gap-4">
          <p>© 2026 CloudScanr. Independent Global Cloud Price Comparison & BYOC Platform.</p>
          <div class="flex space-x-6">
            <NuxtLink to="/terms" class="hover:underline">Terms</NuxtLink>
            <NuxtLink to="/privacy" class="hover:underline">Privacy</NuxtLink>
            <NuxtLink to="/affiliate-disclosure" class="hover:underline">Affiliate Disclosure</NuxtLink>
          </div>
        </div>
        <p class="text-slate-500 leading-relaxed">
          * FT/KFTC Disclosure: CloudScanr is an independent comparison service. We may receive affiliate commissions from certain cloud providers when you click our outbound links or sign up, at no extra cost to you. Pricing and specifications are normalized across 36 dimensions for transparent evaluation.
        </p>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
const auth = useAuthStore()
</script>

