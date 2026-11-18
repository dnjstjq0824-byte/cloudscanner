<template>
  <div class="flex flex-col min-h-screen">
    <header class="border-b border-subtle bg-surface/90 backdrop-blur sticky top-0 z-50">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between gap-2">
        <div class="flex items-center gap-3 min-w-0">
          <button
            class="icon-btn w-10 h-10 text-ink-secondary md:hidden shrink-0"
            aria-label="Open navigation menu"
            :aria-expanded="menuOpen"
            @click="menuOpen = !menuOpen"
          >
            <span aria-hidden="true" class="text-xl">{{ menuOpen ? '✕' : '☰' }}</span>
          </button>
          <NuxtLink to="/" class="flex items-center gap-2 text-xl font-bold text-brand shrink-0" @click="menuOpen = false">
            <span aria-hidden="true">⚡</span>
            <span>CloudScanr</span>
          </NuxtLink>
          <nav class="hidden md:flex gap-4 text-sm font-medium text-ink-secondary ml-4" aria-label="Primary">
            <NuxtLink to="/search" class="hover:text-ink transition">Search & Compare</NuxtLink>
            <NuxtLink to="/calculator" class="hover:text-ink transition">TCO Calculator</NuxtLink>
            <NuxtLink to="/providers" class="hover:text-ink transition">Providers</NuxtLink>
            <NuxtLink to="/rankings" class="hover:text-ink transition">Best Value</NuxtLink>
          </nav>
        </div>
        <div class="flex items-center gap-2 shrink-0">
          <ThemeToggle />
          <template v-if="auth.isAuthenticated">
            <NuxtLink to="/dashboard" class="hidden sm:inline px-3 py-1.5 text-sm font-medium text-ink-secondary hover:text-ink">Console</NuxtLink>
            <span class="hidden lg:inline text-xs text-ink-secondary truncate max-w-40">{{ auth.user?.email }}</span>
            <button class="btn-ghost !py-1.5 !px-3 !text-xs" @click="auth.logout()">Sign Out</button>
          </template>
          <template v-else>
            <NuxtLink to="/dashboard" class="hidden sm:inline px-3 py-1.5 text-sm font-medium text-ink-secondary hover:text-ink">Console</NuxtLink>
            <NuxtLink to="/auth/login" class="btn-primary !py-1.5 !px-4 !text-sm">Sign In</NuxtLink>
          </template>
        </div>
      </div>
      <nav v-if="menuOpen" class="md:hidden border-t border-subtle bg-surface px-4 py-3 space-y-1 text-sm font-medium" aria-label="Mobile">
        <NuxtLink to="/search" class="block px-3 py-2.5 rounded-lg hover:bg-elevated text-ink" @click="menuOpen = false">Search & Compare</NuxtLink>
        <NuxtLink to="/calculator" class="block px-3 py-2.5 rounded-lg hover:bg-elevated text-ink" @click="menuOpen = false">TCO Calculator</NuxtLink>
        <NuxtLink to="/providers" class="block px-3 py-2.5 rounded-lg hover:bg-elevated text-ink" @click="menuOpen = false">Providers</NuxtLink>
        <NuxtLink to="/rankings" class="block px-3 py-2.5 rounded-lg hover:bg-elevated text-ink" @click="menuOpen = false">Best Value</NuxtLink>
        <NuxtLink to="/dashboard" class="block px-3 py-2.5 rounded-lg hover:bg-elevated text-brand" @click="menuOpen = false">Open Console →</NuxtLink>
      </nav>
    </header>

    <main class="flex-1 max-w-7xl w-full mx-auto px-4 sm:px-6 lg:px-8 py-6 sm:py-8 space-y-6 min-w-0">
      <NoticeBanner />
      <slot />
    </main>

    <footer class="border-t border-subtle bg-surface py-10 text-xs text-ink-secondary">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 space-y-4">
        <div class="flex flex-col md:flex-row justify-between md:items-center gap-4">
          <p>© 2026 CloudScanr. Independent Global Cloud Price Comparison & BYOC Platform.</p>
          <div class="flex flex-wrap gap-x-6 gap-y-2">
            <NuxtLink to="/terms" class="hover:underline">Terms</NuxtLink>
            <NuxtLink to="/privacy" class="hover:underline">Privacy</NuxtLink>
            <NuxtLink to="/affiliate-disclosure" class="hover:underline">Affiliate Disclosure</NuxtLink>
          </div>
        </div>
        <p class="leading-relaxed opacity-80">
          * FT/KFTC Disclosure: CloudScanr is an independent comparison service. We may receive affiliate commissions from certain cloud providers when you click our outbound links or sign up, at no extra cost to you. Pricing and specifications are normalized for transparent evaluation.
        </p>
      </div>
    </footer>
    <ToastHost />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const auth = useAuthStore()
const menuOpen = ref(false)
</script>
