<template>
  <div class="flex h-screen bg-slate-950 text-slate-100 overflow-hidden">
    <!-- BYOC Sidebar -->
    <aside class="w-64 border-r border-slate-800 bg-slate-900 flex flex-col justify-between">
      <div>
        <div class="h-16 px-6 border-b border-slate-800 flex items-center justify-between">
          <NuxtLink to="/" class="font-bold text-sky-400 flex items-center space-x-2">
            <span>⚡ CloudScanr</span>
          </NuxtLink>
          <span class="text-xs px-2 py-0.5 rounded bg-sky-950 text-sky-400 border border-sky-800">BYOC</span>
        </div>
        <nav class="p-4 space-y-1 text-sm font-medium">
          <NuxtLink to="/dashboard" class="flex items-center px-3 py-2 rounded-md hover:bg-slate-800 text-slate-300 hover:text-white">Overview</NuxtLink>
          <NuxtLink to="/dashboard/connections" class="flex items-center px-3 py-2 rounded-md hover:bg-slate-800 text-slate-300 hover:text-white">Cloud Connections</NuxtLink>
          <NuxtLink to="/dashboard/servers" class="flex items-center px-3 py-2 rounded-md hover:bg-slate-800 text-slate-300 hover:text-white">Server Inventory</NuxtLink>
          <NuxtLink to="/dashboard/provision" class="flex items-center px-3 py-2 rounded-md hover:bg-slate-800 text-slate-300 hover:text-white">Provision Server</NuxtLink>
          <NuxtLink to="/dashboard/budgets" class="flex items-center px-3 py-2 rounded-md hover:bg-slate-800 text-slate-300 hover:text-white">Budgets & Alerts</NuxtLink>
        </nav>
      </div>
      <div class="p-4 border-t border-slate-800 text-xs text-slate-400">
        <p class="font-medium text-slate-300">Tenant Isolation Active</p>
        <p class="truncate">KMS Envelope Encrypted</p>
      </div>
    </aside>

    <!-- Main Content Area -->
    <div class="flex-1 flex flex-col overflow-y-auto">
      <header class="h-16 border-b border-slate-800 bg-slate-900/50 px-8 flex items-center justify-between">
        <h1 class="text-lg font-semibold text-white">Multi-Cloud Management Console</h1>
        <div class="flex items-center space-x-4 text-sm">
          <span class="text-slate-400">Org: <strong class="text-white">{{ auth.org?.name || '—' }}</strong></span>
          <span class="hidden md:inline text-slate-500 text-xs">{{ auth.user?.email }}</span>
          <button class="px-3 py-1 rounded bg-slate-800 hover:bg-slate-700 text-xs" @click="auth.logout()">Sign Out</button>
        </div>
      </header>
      <main class="flex-1 p-8">
        <slot />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
const auth = useAuthStore()
</script>

