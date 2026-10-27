<template>
  <div class="flex h-screen bg-slate-950 text-slate-100 overflow-hidden">
    <!-- Admin Sidebar -->
    <aside class="w-64 border-r border-red-900/40 bg-slate-900 flex flex-col justify-between">
      <div>
        <div class="h-16 px-6 border-b border-red-900/40 flex items-center justify-between">
          <NuxtLink to="/admin" class="font-bold text-red-400 flex items-center space-x-2">
            <span>🛡️ Ops Console</span>
          </NuxtLink>
          <span class="text-xs px-2 py-0.5 rounded bg-red-950 text-red-400 border border-red-800">SUPERADMIN</span>
        </div>
        <nav class="p-4 space-y-1 text-sm font-medium">
          <NuxtLink to="/admin" class="flex items-center px-3 py-2 rounded-md hover:bg-slate-800 text-slate-300 hover:text-white">Ops Health</NuxtLink>
          <NuxtLink to="/admin/pricing/quarantine" class="flex items-center px-3 py-2 rounded-md hover:bg-slate-800 text-slate-300 hover:text-white">Price Quarantine</NuxtLink>
          <NuxtLink to="/admin/affiliate/reconcile" class="flex items-center px-3 py-2 rounded-md hover:bg-slate-800 text-slate-300 hover:text-white">Affiliate Reconciliation</NuxtLink>
          <NuxtLink to="/admin/cloud/audit" class="flex items-center px-3 py-2 rounded-md hover:bg-slate-800 text-slate-300 hover:text-white">BYOC Audit Logs</NuxtLink>
          <NuxtLink to="/admin/flags" class="flex items-center px-3 py-2 rounded-md hover:bg-slate-800 text-slate-300 hover:text-white">Feature Flags</NuxtLink>
        </nav>
      </div>
      <div class="p-4 border-t border-red-900/40 text-xs text-red-400">
        <p class="font-bold">MFA Step-Up Enforced</p>
        <p>All mutations audited</p>
      </div>
    </aside>

    <!-- Main Content Area -->
    <div class="flex-1 flex flex-col overflow-y-auto">
      <header class="h-16 border-b border-slate-800 bg-slate-900/50 px-8 flex items-center justify-between">
        <h1 class="text-lg font-semibold text-red-200">System Administration & Operations</h1>
        <div class="flex items-center space-x-4 text-sm">
          <span class="px-2 py-1 rounded bg-green-950 text-green-400 border border-green-800 text-xs">PostgreSQL 16 HA Online</span>
          <span class="px-2 py-1 rounded bg-sky-950 text-sky-400 border border-sky-800 text-xs">OpenSearch 2.17 Online</span>
        </div>
      </header>
      <main class="flex-1 p-8">
        <slot />
      </main>
    </div>
  </div>
</template>
