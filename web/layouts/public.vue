<template>
  <div class="flex flex-col min-h-screen">
    <header class="border-b border-subtle bg-surface/90 backdrop-blur sticky top-0 z-50">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between gap-2">
        <div class="flex items-center gap-3 min-w-0">
          <button
            class="icon-btn w-10 h-10 text-ink-secondary md:hidden shrink-0"
            :aria-label="$t('a11y.openMenu')"
            :aria-expanded="menuOpen"
            @click="menuOpen = !menuOpen"
          >
            <span aria-hidden="true" class="text-xl">{{ menuOpen ? '✕' : '☰' }}</span>
          </button>
          <NuxtLink :to="$localePath('/')" class="flex items-center gap-2 text-xl font-bold text-brand shrink-0" @click="menuOpen = false">
            <span aria-hidden="true">⚡</span>
            <span>TirpScanner</span>
          </NuxtLink>
          <nav class="hidden md:flex gap-4 text-sm font-medium text-ink-secondary ml-4" :aria-label="$t('a11y.primary')">
            <NuxtLink :to="$localePath('/offers')" class="hover:text-ink transition">{{ $t('nav.deals') }}</NuxtLink>
            <NuxtLink :to="$localePath('/best-value-rankings')" class="hover:text-ink transition">{{ $t('nav.rankings') }}</NuxtLink>
            <NuxtLink :to="$localePath('/providers')" class="hover:text-ink transition">{{ $t('nav.vendors') }}</NuxtLink>
            <NuxtLink :to="$localePath('/compare')" class="hover:text-ink transition">{{ $t('nav.compare') }}</NuxtLink>
            <NuxtLink :to="$localePath('/dashboard')" class="hover:text-ink transition">{{ $t('nav.console') }}</NuxtLink>
            <NuxtLink :to="$localePath('/affiliate-disclosure')" class="hover:text-ink transition">{{ $t('nav.disclosure') }}</NuxtLink>
          </nav>
        </div>
        <div class="flex items-center gap-2 shrink-0">
          <LocaleSwitcher class="hidden lg:flex" />
          <ThemeToggle />
          <template v-if="auth.isAuthenticated">
            <NuxtLink :to="$localePath('/dashboard')" class="hidden sm:inline px-3 py-1.5 text-sm font-medium text-ink-secondary hover:text-ink">{{ $t('nav.console') }}</NuxtLink>
            <span class="hidden lg:inline text-xs text-ink-secondary truncate max-w-40">{{ auth.user?.email }}</span>
            <button class="btn-ghost !py-1.5 !px-3 !text-xs" @click="auth.logout()">{{ $t('auth.signOut') }}</button>
          </template>
          <template v-else>
            <NuxtLink :to="$localePath('/dashboard')" class="hidden sm:inline px-3 py-1.5 text-sm font-medium text-ink-secondary hover:text-ink">{{ $t('nav.console') }}</NuxtLink>
            <NuxtLink :to="$localePath('/auth/login')" class="btn-primary !py-1.5 !px-4 !text-sm">{{ $t('auth.signIn') }}</NuxtLink>
          </template>
        </div>
      </div>
      <nav v-if="menuOpen" class="md:hidden border-t border-subtle bg-surface px-4 py-3 space-y-1 text-sm font-medium" :aria-label="$t('a11y.mobile')">
        <div class="px-3 py-2">
          <LocaleSwitcher />
        </div>
        <NuxtLink :to="$localePath('/offers')" class="block px-3 py-2.5 rounded-lg hover:bg-elevated text-ink" @click="menuOpen = false">{{ $t('nav.deals') }}</NuxtLink>
        <NuxtLink :to="$localePath('/best-value-rankings')" class="block px-3 py-2.5 rounded-lg hover:bg-elevated text-ink" @click="menuOpen = false">{{ $t('nav.rankings') }}</NuxtLink>
        <NuxtLink :to="$localePath('/providers')" class="block px-3 py-2.5 rounded-lg hover:bg-elevated text-ink" @click="menuOpen = false">{{ $t('nav.vendors') }}</NuxtLink>
        <NuxtLink :to="$localePath('/compare')" class="block px-3 py-2.5 rounded-lg hover:bg-elevated text-ink" @click="menuOpen = false">{{ $t('nav.compare') }}</NuxtLink>
        <NuxtLink :to="$localePath('/affiliate-disclosure')" class="block px-3 py-2.5 rounded-lg hover:bg-elevated text-ink" @click="menuOpen = false">{{ $t('nav.disclosure') }}</NuxtLink>
        <NuxtLink :to="$localePath('/dashboard')" class="block px-3 py-2.5 rounded-lg hover:bg-elevated text-brand" @click="menuOpen = false">{{ $t('nav.openConsole') }}</NuxtLink>
      </nav>
    </header>

    <main class="flex-1 max-w-7xl w-full mx-auto px-4 sm:px-6 lg:px-8 py-6 sm:py-8 space-y-6 min-w-0">
      <NoticeBanner />
      <slot />
    </main>

    <footer class="border-t border-subtle bg-surface py-10 text-xs text-ink-secondary">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 space-y-4">
        <div class="flex flex-col md:flex-row justify-between md:items-center gap-4">
          <p>{{ $t('footer.copyright') }}</p>
          <nav class="flex flex-wrap gap-x-4 gap-y-2" :aria-label="$t('footer.sitemap')">
            <NuxtLink :to="$localePath('/search')" class="hover:underline">{{ $t('nav.search') }}</NuxtLink>
            <NuxtLink :to="$localePath('/compare')" class="hover:underline">{{ $t('nav.compare') }}</NuxtLink>
            <NuxtLink :to="$localePath('/best-value-rankings')" class="hover:underline">{{ $t('nav.rankings') }}</NuxtLink>
            <NuxtLink :to="$localePath('/providers')" class="hover:underline">{{ $t('nav.vendors') }}</NuxtLink>
            <NuxtLink :to="$localePath('/offers')" class="hover:underline">{{ $t('nav.deals') }}</NuxtLink>
            <NuxtLink :to="$localePath('/calculator')" class="hover:underline">{{ $t('nav.calculator') }}</NuxtLink>
            <NuxtLink :to="$localePath('/cost-estimation')" class="hover:underline">{{ $t('nav.estimation') }}</NuxtLink>
            <NuxtLink :to="$localePath('/deal-trends')" class="hover:underline">{{ $t('nav.trends') }}</NuxtLink>
            <NuxtLink :to="$localePath('/tutorials')" class="hover:underline">{{ $t('nav.tutorials') }}</NuxtLink>
            <NuxtLink :to="$localePath('/reviews')" class="hover:underline">{{ $t('nav.reviews') }}</NuxtLink>
          </nav>
          <div class="flex flex-wrap gap-x-6 gap-y-2">
            <NuxtLink :to="$localePath('/terms')" class="hover:underline">{{ $t('nav.terms') }}</NuxtLink>
            <NuxtLink :to="$localePath('/privacy')" class="hover:underline">{{ $t('nav.privacy') }}</NuxtLink>
            <NuxtLink :to="$localePath('/affiliate-disclosure')" class="hover:underline">{{ $t('nav.affiliateDisclosure') }}</NuxtLink>
          </div>
        </div>
        <p class="leading-relaxed opacity-80">
          {{ $t('footer.disclosureNote') }}
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
