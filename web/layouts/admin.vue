<template>
  <div class="flex h-screen bg-canvas text-ink overflow-hidden">
    <aside class="hidden lg:flex w-64 shrink-0 border-r border-rose-500/20 bg-surface flex-col justify-between">
      <div>
        <div class="h-16 px-6 border-b border-rose-500/20 flex items-center justify-between">
          <NuxtLink :to="$localePath('/admin')" class="font-bold text-rose-500">🛡️ {{ $t('nav.opsConsole') }}</NuxtLink>
          <span class="text-[10px] px-2 py-0.5 rounded bg-rose-500/10 text-rose-600 dark:text-rose-300 border border-rose-500/30 font-semibold">{{ $t('admin.badge') }}</span>
        </div>
        <nav class="p-4 space-y-1 text-sm font-medium overflow-y-auto" :aria-label="$t('a11y.admin')">
          <NuxtLink :to="$localePath('/admin')" class="flex items-center px-3 py-2.5 rounded-lg hover:bg-elevated text-ink-secondary hover:text-ink">{{ $t('nav.opsHealth') }}</NuxtLink>
          <NuxtLink :to="$localePath('/admin/users')" class="flex items-center px-3 py-2.5 rounded-lg hover:bg-elevated text-ink-secondary hover:text-ink">{{ $t('nav.membersRoles') }}</NuxtLink>
          <NuxtLink :to="$localePath('/admin/providers')" class="flex items-center px-3 py-2.5 rounded-lg hover:bg-elevated text-ink-secondary hover:text-ink">{{ $t('nav.vendorCatalog') }}</NuxtLink>
          <p class="px-3 pt-3 text-[10px] font-bold uppercase tracking-wider text-ink-secondary">{{ $t('admin.groupCatalog') }}</p>
          <NuxtLink :to="$localePath('/admin/catalog')" class="flex items-center px-3 py-2.5 rounded-lg hover:bg-elevated text-ink-secondary hover:text-ink">{{ $t('nav.catalog') }}</NuxtLink>
          <NuxtLink :to="$localePath('/admin/pricing-sync')" class="flex items-center px-3 py-2.5 rounded-lg hover:bg-elevated text-ink-secondary hover:text-ink">{{ $t('nav.pricingSync') }}</NuxtLink>
          <NuxtLink :to="$localePath('/admin/pricing-failures')" class="flex items-center px-3 py-2.5 rounded-lg hover:bg-elevated text-ink-secondary hover:text-ink">{{ $t('nav.pricingFailures') }}</NuxtLink>
          <p class="px-3 pt-3 text-[10px] font-bold uppercase tracking-wider text-ink-secondary">{{ $t('admin.groupAffiliate') }}</p>
          <NuxtLink :to="$localePath('/admin/affiliate/accounts')" class="flex items-center px-3 py-2.5 rounded-lg hover:bg-elevated text-ink-secondary hover:text-ink">{{ $t('nav.affAccounts') }}</NuxtLink>
          <NuxtLink :to="$localePath('/admin/affiliate/campaigns')" class="flex items-center px-3 py-2.5 rounded-lg hover:bg-elevated text-ink-secondary hover:text-ink">{{ $t('nav.affCampaigns') }}</NuxtLink>
          <NuxtLink :to="$localePath('/admin/affiliate/links')" class="flex items-center px-3 py-2.5 rounded-lg hover:bg-elevated text-ink-secondary hover:text-ink">{{ $t('nav.affLinks') }}</NuxtLink>
          <NuxtLink :to="$localePath('/admin/affiliate/commissions')" class="flex items-center px-3 py-2.5 rounded-lg hover:bg-elevated text-ink-secondary hover:text-ink">{{ $t('nav.affCommissions') }}</NuxtLink>
          <NuxtLink :to="$localePath('/admin/affiliate/reconciliation')" class="flex items-center px-3 py-2.5 rounded-lg hover:bg-elevated text-ink-secondary hover:text-ink">{{ $t('nav.affRecon') }}</NuxtLink>
          <p class="px-3 pt-3 text-[10px] font-bold uppercase tracking-wider text-ink-secondary">{{ $t('admin.groupSecurity') }}</p>
          <NuxtLink :to="$localePath('/admin/audit-logs')" class="flex items-center px-3 py-2.5 rounded-lg hover:bg-elevated text-ink-secondary hover:text-ink">{{ $t('nav.auditLogs') }}</NuxtLink>
          <NuxtLink :to="$localePath('/admin/permissions')" class="flex items-center px-3 py-2.5 rounded-lg hover:bg-elevated text-ink-secondary hover:text-ink">{{ $t('nav.permissions') }}</NuxtLink>
          <NuxtLink :to="$localePath('/admin/security-events')" class="flex items-center px-3 py-2.5 rounded-lg hover:bg-elevated text-ink-secondary hover:text-ink">{{ $t('nav.securityEvents') }}</NuxtLink>
          <NuxtLink :to="$localePath('/admin/audit-verification')" class="flex items-center px-3 py-2.5 rounded-lg hover:bg-elevated text-ink-secondary hover:text-ink">{{ $t('nav.auditVerify') }}</NuxtLink>
          <p class="px-3 pt-3 text-[10px] font-bold uppercase tracking-wider text-ink-secondary">{{ $t('admin.groupContent') }}</p>
          <NuxtLink :to="$localePath('/admin/content')" class="flex items-center px-3 py-2.5 rounded-lg hover:bg-elevated text-ink-secondary hover:text-ink">{{ $t('nav.content') }}</NuxtLink>
          <NuxtLink :to="$localePath('/admin/ads')" class="flex items-center px-3 py-2.5 rounded-lg hover:bg-elevated text-ink-secondary hover:text-ink">{{ $t('nav.ads') }}</NuxtLink>
          <NuxtLink :to="$localePath('/admin/seo')" class="flex items-center px-3 py-2.5 rounded-lg hover:bg-elevated text-ink-secondary hover:text-ink">{{ $t('nav.seo') }}</NuxtLink>
          <NuxtLink :to="$localePath('/admin/notices')" class="flex items-center px-3 py-2.5 rounded-lg hover:bg-elevated text-ink-secondary hover:text-ink">{{ $t('nav.notices') }}</NuxtLink>
          <p class="px-3 pt-3 text-[10px] font-bold uppercase tracking-wider text-ink-secondary">{{ $t('admin.groupSystem') }}</p>
          <NuxtLink :to="$localePath('/admin/integrations')" class="flex items-center px-3 py-2.5 rounded-lg hover:bg-elevated text-ink-secondary hover:text-ink">{{ $t('nav.integrations') }}</NuxtLink>
          <NuxtLink :to="$localePath('/admin/jobs')" class="flex items-center px-3 py-2.5 rounded-lg hover:bg-elevated text-ink-secondary hover:text-ink">{{ $t('nav.jobs') }}</NuxtLink>
          <NuxtLink :to="$localePath('/admin/deployments')" class="flex items-center px-3 py-2.5 rounded-lg hover:bg-elevated text-ink-secondary hover:text-ink">{{ $t('nav.deployments') }}</NuxtLink>
          <NuxtLink :to="$localePath('/admin/settings')" class="flex items-center px-3 py-2.5 rounded-lg hover:bg-elevated text-ink-secondary hover:text-ink">{{ $t('nav.siteSettings') }}</NuxtLink>
          <NuxtLink :to="$localePath('/dashboard')" class="flex items-center px-3 py-2.5 rounded-lg hover:bg-elevated text-ink-secondary hover:text-ink">{{ $t('nav.backToConsole') }}</NuxtLink>
        </nav>
      </div>
      <div class="p-4 border-t border-rose-500/20 text-xs text-rose-600 dark:text-rose-300">
        <p class="font-bold">{{ $t('admin.stepUpTitle') }}</p>
        <p>{{ $t('admin.stepUpSub') }}</p>
      </div>
    </aside>

    <div
      v-if="drawerOpen"
      class="fixed inset-0 z-[90] lg:hidden"
      role="dialog"
      aria-modal="true"
      :aria-label="$t('a11y.adminMobile')"
    >
      <div class="absolute inset-0 bg-black/60" @click="drawerOpen = false" />
      <aside class="absolute left-0 top-0 bottom-0 w-72 max-w-[85vw] bg-surface border-r border-subtle flex flex-col">
        <div class="h-16 px-5 border-b border-subtle flex items-center justify-between shrink-0">
          <NuxtLink :to="$localePath('/admin')" class="font-bold text-rose-500" @click="drawerOpen = false">🛡️ {{ $t('nav.opsConsole') }}</NuxtLink>
          <button class="icon-btn w-9 h-9 text-ink-secondary" :aria-label="$t('a11y.closeNav')" @click="drawerOpen = false">✕</button>
        </div>
        <nav class="p-4 space-y-1 text-sm font-medium overflow-y-auto" :aria-label="$t('a11y.adminMobile')">
          <NuxtLink :to="$localePath('/admin')" class="flex items-center px-3 py-3 rounded-lg hover:bg-elevated text-ink" @click="drawerOpen = false">{{ $t('nav.opsHealth') }}</NuxtLink>
          <NuxtLink :to="$localePath('/admin/users')" class="flex items-center px-3 py-3 rounded-lg hover:bg-elevated text-ink" @click="drawerOpen = false">{{ $t('nav.membersRoles') }}</NuxtLink>
          <NuxtLink :to="$localePath('/admin/providers')" class="flex items-center px-3 py-3 rounded-lg hover:bg-elevated text-ink" @click="drawerOpen = false">{{ $t('nav.vendorCatalog') }}</NuxtLink>
          <p class="px-3 pt-3 text-[10px] font-bold uppercase tracking-wider text-ink-secondary">{{ $t('admin.groupCatalog') }}</p>
          <NuxtLink :to="$localePath('/admin/catalog')" class="flex items-center px-3 py-3 rounded-lg hover:bg-elevated text-ink" @click="drawerOpen = false">{{ $t('nav.catalog') }}</NuxtLink>
          <NuxtLink :to="$localePath('/admin/pricing-sync')" class="flex items-center px-3 py-3 rounded-lg hover:bg-elevated text-ink" @click="drawerOpen = false">{{ $t('nav.pricingSync') }}</NuxtLink>
          <NuxtLink :to="$localePath('/admin/pricing-failures')" class="flex items-center px-3 py-3 rounded-lg hover:bg-elevated text-ink" @click="drawerOpen = false">{{ $t('nav.pricingFailures') }}</NuxtLink>
          <p class="px-3 pt-3 text-[10px] font-bold uppercase tracking-wider text-ink-secondary">{{ $t('admin.groupAffiliate') }}</p>
          <NuxtLink :to="$localePath('/admin/affiliate/accounts')" class="flex items-center px-3 py-3 rounded-lg hover:bg-elevated text-ink" @click="drawerOpen = false">{{ $t('nav.affAccounts') }}</NuxtLink>
          <NuxtLink :to="$localePath('/admin/affiliate/campaigns')" class="flex items-center px-3 py-3 rounded-lg hover:bg-elevated text-ink" @click="drawerOpen = false">{{ $t('nav.affCampaigns') }}</NuxtLink>
          <NuxtLink :to="$localePath('/admin/affiliate/links')" class="flex items-center px-3 py-3 rounded-lg hover:bg-elevated text-ink" @click="drawerOpen = false">{{ $t('nav.affLinks') }}</NuxtLink>
          <NuxtLink :to="$localePath('/admin/affiliate/commissions')" class="flex items-center px-3 py-3 rounded-lg hover:bg-elevated text-ink" @click="drawerOpen = false">{{ $t('nav.affCommissions') }}</NuxtLink>
          <NuxtLink :to="$localePath('/admin/affiliate/reconciliation')" class="flex items-center px-3 py-3 rounded-lg hover:bg-elevated text-ink" @click="drawerOpen = false">{{ $t('nav.affRecon') }}</NuxtLink>
          <p class="px-3 pt-3 text-[10px] font-bold uppercase tracking-wider text-ink-secondary">{{ $t('admin.groupSecurity') }}</p>
          <NuxtLink :to="$localePath('/admin/audit-logs')" class="flex items-center px-3 py-3 rounded-lg hover:bg-elevated text-ink" @click="drawerOpen = false">{{ $t('nav.auditLogs') }}</NuxtLink>
          <NuxtLink :to="$localePath('/admin/permissions')" class="flex items-center px-3 py-3 rounded-lg hover:bg-elevated text-ink" @click="drawerOpen = false">{{ $t('nav.permissions') }}</NuxtLink>
          <NuxtLink :to="$localePath('/admin/security-events')" class="flex items-center px-3 py-3 rounded-lg hover:bg-elevated text-ink" @click="drawerOpen = false">{{ $t('nav.securityEvents') }}</NuxtLink>
          <NuxtLink :to="$localePath('/admin/audit-verification')" class="flex items-center px-3 py-3 rounded-lg hover:bg-elevated text-ink" @click="drawerOpen = false">{{ $t('nav.auditVerify') }}</NuxtLink>
          <p class="px-3 pt-3 text-[10px] font-bold uppercase tracking-wider text-ink-secondary">{{ $t('admin.groupContent') }}</p>
          <NuxtLink :to="$localePath('/admin/content')" class="flex items-center px-3 py-3 rounded-lg hover:bg-elevated text-ink" @click="drawerOpen = false">{{ $t('nav.content') }}</NuxtLink>
          <NuxtLink :to="$localePath('/admin/ads')" class="flex items-center px-3 py-3 rounded-lg hover:bg-elevated text-ink" @click="drawerOpen = false">{{ $t('nav.ads') }}</NuxtLink>
          <NuxtLink :to="$localePath('/admin/seo')" class="flex items-center px-3 py-3 rounded-lg hover:bg-elevated text-ink" @click="drawerOpen = false">{{ $t('nav.seo') }}</NuxtLink>
          <NuxtLink :to="$localePath('/admin/notices')" class="flex items-center px-3 py-3 rounded-lg hover:bg-elevated text-ink" @click="drawerOpen = false">{{ $t('nav.notices') }}</NuxtLink>
          <p class="px-3 pt-3 text-[10px] font-bold uppercase tracking-wider text-ink-secondary">{{ $t('admin.groupSystem') }}</p>
          <NuxtLink :to="$localePath('/admin/integrations')" class="flex items-center px-3 py-3 rounded-lg hover:bg-elevated text-ink" @click="drawerOpen = false">{{ $t('nav.integrations') }}</NuxtLink>
          <NuxtLink :to="$localePath('/admin/jobs')" class="flex items-center px-3 py-3 rounded-lg hover:bg-elevated text-ink" @click="drawerOpen = false">{{ $t('nav.jobs') }}</NuxtLink>
          <NuxtLink :to="$localePath('/admin/deployments')" class="flex items-center px-3 py-3 rounded-lg hover:bg-elevated text-ink" @click="drawerOpen = false">{{ $t('nav.deployments') }}</NuxtLink>
          <NuxtLink :to="$localePath('/admin/settings')" class="flex items-center px-3 py-3 rounded-lg hover:bg-elevated text-ink" @click="drawerOpen = false">{{ $t('nav.siteSettings') }}</NuxtLink>
          <NuxtLink :to="$localePath('/dashboard')" class="flex items-center px-3 py-3 rounded-lg hover:bg-elevated text-ink" @click="drawerOpen = false">{{ $t('nav.backToConsole') }}</NuxtLink>
        </nav>
      </aside>
    </div>

    <div class="flex-1 flex flex-col overflow-y-auto min-w-0">
      <header class="min-h-16 border-b border-subtle bg-surface/80 backdrop-blur px-4 sm:px-6 lg:px-8 py-3 flex items-center justify-between gap-3 sticky top-0 z-40">
        <div class="flex items-center gap-2 min-w-0">
          <button class="icon-btn w-10 h-10 text-ink-secondary lg:hidden shrink-0" :aria-label="$t('a11y.openAdminNav')" @click="drawerOpen = true">
            <span aria-hidden="true" class="text-xl">☰</span>
          </button>
          <h1 class="text-base sm:text-lg font-semibold text-ink truncate">{{ $t('admin.title') }}</h1>
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

const { t } = useI18n()
const drawerOpen = ref(false)
const storeLabel = ref('…')

onMounted(async () => {
  try {
    const { request } = useApi()
    const ready = await request<{ database: string }>('/api/v1/readiness')
    storeLabel.value = ready.database === 'connected' ? t('admin.storePostgres') : t('admin.storeMemory')
  } catch {
    storeLabel.value = t('admin.storeOffline')
  }
})
</script>
