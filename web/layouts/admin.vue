<template>
  <div class="flex h-screen bg-canvas text-ink overflow-hidden">
    <aside class="hidden lg:flex w-64 shrink-0 border-r border-rose-500/20 bg-surface flex-col justify-between">
      <div>
        <div class="h-16 px-6 border-b border-rose-500/20 flex items-center justify-between">
          <NuxtLink to="/admin" class="font-bold text-rose-500">🛡️ Ops Console</NuxtLink>
          <span class="text-[10px] px-2 py-0.5 rounded bg-rose-500/10 text-rose-600 dark:text-rose-300 border border-rose-500/30 font-semibold">ADMIN</span>
        </div>
        <nav class="p-4 space-y-1 text-sm font-medium" aria-label="Admin">
          <NuxtLink to="/admin" class="flex items-center px-3 py-2.5 rounded-lg hover:bg-elevated text-ink-secondary hover:text-ink">Ops Health</NuxtLink>
          <NuxtLink to="/admin/users" class="flex items-center px-3 py-2.5 rounded-lg hover:bg-elevated text-ink-secondary hover:text-ink">Members & Roles</NuxtLink>
          <NuxtLink to="/admin/providers" class="flex items-center px-3 py-2.5 rounded-lg hover:bg-elevated text-ink-secondary hover:text-ink">Provider Catalog</NuxtLink>
          <NuxtLink to="/admin/audit-logs" class="flex items-center px-3 py-2.5 rounded-lg hover:bg-elevated text-ink-secondary hover:text-ink">Audit Logs</NuxtLink>
          <NuxtLink to="/admin/notices" class="flex items-center px-3 py-2.5 rounded-lg hover:bg-elevated text-ink-secondary hover:text-ink">Notices</NuxtLink>
          <NuxtLink to="/admin/settings" class="flex items-center px-3 py-2.5 rounded-lg hover:bg-elevated text-ink-secondary hover:text-ink">Site Settings</NuxtLink>
          <NuxtLink to="/dashboard" class="flex items-center px-3 py-2.5 rounded-lg hover:bg-elevated text-ink-secondary hover:text-ink">← Back to Console</NuxtLink>
        </nav>
      </div>
      <div class="p-4 border-t border-rose-500/20 text-xs text-rose-600 dark:text-rose-300">
        <p class="font-bold">Step-Up Gates Enforced</p>
        <p>All mutations audited</p>
      </div>
    </aside>

    <div
      v-if="drawerOpen"
      class="fixed inset-0 z-[90] lg:hidden"
      role="dialog"
      aria-modal="true"
      aria-label="Admin navigation"
    >
      <div class="absolute inset-0 bg-black/60" @click="drawerOpen = false" />
      <aside class="absolute left-0 top-0 bottom-0 w-72 max-w-[85vw] bg-surface border-r border-subtle flex flex-col">
        <div class="h-16 px-5 border-b border-subtle flex items-center justify-between shrink-0">
          <NuxtLink to="/admin" class="font-bold text-rose-500" @click="drawerOpen = false">🛡️ Ops Console</NuxtLink>
          <button class="icon-btn w-9 h-9 text-ink-secondary" aria-label="Close navigation" @click="drawerOpen = false">✕</button>
        </div>
        <nav class="p-4 space-y-1 text-sm font-medium overflow-y-auto" aria-label="Admin mobile">
          <NuxtLink to="/admin" class="flex items-center px-3 py-3 rounded-lg hover:bg-elevated text-ink" @click="drawerOpen = false">Ops Health</NuxtLink>
          <NuxtLink to="/admin/users" class="flex items-center px-3 py-3 rounded-lg hover:bg-elevated text-ink" @click="drawerOpen = false">Members & Roles</NuxtLink>
          <NuxtLink to="/admin/providers" class="flex items-center px-3 py-3 rounded-lg hover:bg-elevated text-ink" @click="drawerOpen = false">Provider Catalog</NuxtLink>
          <NuxtLink to="/admin/audit-logs" class="flex items-center px-3 py-3 rounded-lg hover:bg-elevated text-ink" @click="drawerOpen = false">Audit Logs</NuxtLink>
          <NuxtLink to="/admin/notices" class="flex items-center px-3 py-3 rounded-lg hover:bg-elevated text-ink" @click="drawerOpen = false">Notices</NuxtLink>
          <NuxtLink to="/admin/settings" class="flex items-center px-3 py-3 rounded-lg hover:bg-elevated text-ink" @click="drawerOpen = false">Site Settings</NuxtLink>
          <NuxtLink to="/dashboard" class="flex items-center px-3 py-3 rounded-lg hover:bg-elevated text-ink" @click="drawerOpen = false">← Back to Console</NuxtLink>
        </nav>
      </aside>
    </div>

    <div class="flex-1 flex flex-col overflow-y-auto min-w-0">
      <header class="min-h-16 border-b border-subtle bg-surface/80 backdrop-blur px-4 sm:px-6 lg:px-8 py-3 flex items-center justify-between gap-3 sticky top-0 z-40">
        <div class="flex items-center gap-2 min-w-0">
          <button class="icon-btn w-10 h-10 text-ink-secondary lg:hidden shrink-0" aria-label="Open admin navigation" @click="drawerOpen = true">
            <span aria-hidden="true" class="text-xl">☰</span>
          </button>
          <h1 class="text-base sm:text-lg font-semibold text-ink truncate">System Administration</h1>
        </div>
        <div class="flex items-center gap-2 shrink-0">
          <StatusBadge :status="storeLabel" />
          <ThemeToggle />
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
import { onMounted, ref } from 'vue'

const drawerOpen = ref(false)
const storeLabel = ref('…')

onMounted(async () => {
  try {
    const { request } = useApi()
    const ready = await request<{ database: string }>('/api/v1/readiness')
    storeLabel.value = ready.database === 'connected' ? 'PostgreSQL' : 'In-Memory'
  } catch {
    storeLabel.value = 'API Offline'
  }
})
</script>
