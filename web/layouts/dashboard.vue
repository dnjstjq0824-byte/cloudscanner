<template>
  <div class="flex h-screen bg-canvas text-ink overflow-hidden">
    <aside class="hidden lg:flex w-64 shrink-0 border-r border-subtle bg-surface flex-col justify-between">
      <div>
        <div class="h-16 px-6 border-b border-subtle flex items-center justify-between">
          <NuxtLink to="/" class="font-bold text-brand">⚡ CloudScanr</NuxtLink>
          <span class="text-[10px] px-2 py-0.5 rounded bg-sky-500/10 text-sky-700 dark:text-sky-300 border border-sky-500/30 font-semibold">BYOC</span>
        </div>
        <nav class="p-4 space-y-1 text-sm font-medium" aria-label="Console">
          <NuxtLink to="/dashboard" class="flex items-center px-3 py-2.5 rounded-lg hover:bg-elevated text-ink-secondary hover:text-ink">Overview</NuxtLink>
          <NuxtLink to="/dashboard/connections" class="flex items-center px-3 py-2.5 rounded-lg hover:bg-elevated text-ink-secondary hover:text-ink">Cloud Connections</NuxtLink>
          <NuxtLink to="/dashboard/servers" class="flex items-center px-3 py-2.5 rounded-lg hover:bg-elevated text-ink-secondary hover:text-ink">Server Inventory</NuxtLink>
          <NuxtLink to="/dashboard/provision" class="flex items-center px-3 py-2.5 rounded-lg hover:bg-elevated text-ink-secondary hover:text-ink">Provision Server</NuxtLink>
          <NuxtLink to="/dashboard/budgets" class="flex items-center px-3 py-2.5 rounded-lg hover:bg-elevated text-ink-secondary hover:text-ink">Budgets & Alerts</NuxtLink>
        </nav>
      </div>
      <div class="p-4 border-t border-subtle text-xs text-ink-secondary">
        <p class="font-medium text-ink">Tenant Isolation Active</p>
        <p class="truncate">KMS Envelope Encrypted</p>
      </div>
    </aside>

    <div
      v-if="drawerOpen"
      class="fixed inset-0 z-[90] lg:hidden"
      role="dialog"
      aria-modal="true"
      aria-label="Console navigation"
    >
      <div class="absolute inset-0 bg-black/60" @click="drawerOpen = false" />
      <aside class="absolute left-0 top-0 bottom-0 w-72 max-w-[85vw] bg-surface border-r border-subtle flex flex-col">
        <div class="h-16 px-5 border-b border-subtle flex items-center justify-between shrink-0">
          <NuxtLink to="/" class="font-bold text-brand" @click="drawerOpen = false">⚡ CloudScanr</NuxtLink>
          <button class="icon-btn w-9 h-9 text-ink-secondary" aria-label="Close navigation" @click="drawerOpen = false">✕</button>
        </div>
        <nav class="p-4 space-y-1 text-sm font-medium overflow-y-auto" aria-label="Console mobile">
          <NuxtLink to="/dashboard" class="flex items-center px-3 py-3 rounded-lg hover:bg-elevated text-ink" @click="drawerOpen = false">Overview</NuxtLink>
          <NuxtLink to="/dashboard/connections" class="flex items-center px-3 py-3 rounded-lg hover:bg-elevated text-ink" @click="drawerOpen = false">Cloud Connections</NuxtLink>
          <NuxtLink to="/dashboard/servers" class="flex items-center px-3 py-3 rounded-lg hover:bg-elevated text-ink" @click="drawerOpen = false">Server Inventory</NuxtLink>
          <NuxtLink to="/dashboard/provision" class="flex items-center px-3 py-3 rounded-lg hover:bg-elevated text-ink" @click="drawerOpen = false">Provision Server</NuxtLink>
          <NuxtLink to="/dashboard/budgets" class="flex items-center px-3 py-3 rounded-lg hover:bg-elevated text-ink" @click="drawerOpen = false">Budgets & Alerts</NuxtLink>
        </nav>
      </aside>
    </div>

    <div class="flex-1 flex flex-col overflow-y-auto min-w-0">
      <header class="min-h-16 border-b border-subtle bg-surface/80 backdrop-blur px-4 sm:px-6 lg:px-8 py-3 flex items-center justify-between gap-3 sticky top-0 z-40">
        <div class="flex items-center gap-2 min-w-0">
          <button class="icon-btn w-10 h-10 text-ink-secondary lg:hidden shrink-0" aria-label="Open console navigation" @click="drawerOpen = true">
            <span aria-hidden="true" class="text-xl">☰</span>
          </button>
          <h1 class="text-base sm:text-lg font-semibold text-ink truncate">Multi-Cloud Management Console</h1>
        </div>
        <div class="flex items-center gap-2 sm:gap-3 text-sm shrink-0">
          <span class="hidden md:inline text-ink-secondary text-xs truncate max-w-48">Org: <strong class="text-ink">{{ auth.org?.name || '—' }}</strong></span>
          <ThemeToggle />
          <button class="btn-ghost !py-1.5 !px-3 !text-xs" @click="auth.logout()">Sign Out</button>
        </div>
      </header>
      <main class="flex-1 p-4 sm:p-6 lg:p-8 min-w-0">
        <slot />
      </main>
    </div>
    <ToastHost />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const auth = useAuthStore()
const drawerOpen = ref(false)
</script>
