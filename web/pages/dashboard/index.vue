<template>
  <div class="space-y-8">
    <PageHeader :title="$t('nav.overview')" :subtitle="$t('dashboard.subtitle', { org: auth.org?.name || '—', role: auth.role || '—' })" />

    <div class="card p-6 space-y-3">
      <h3 class="text-sm font-bold text-ink">{{ $t('dashboard.retiredTitle') }}</h3>
      <p class="text-xs text-ink-secondary leading-relaxed">
        {{ $t('dashboard.retiredBody') }}
      </p>
      <NuxtLink :to="$localePath('/')" class="inline-block text-xs font-semibold text-brand hover:underline">
        {{ $t('dashboard.backHome') }}
      </NuxtLink>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'dashboard',
  middleware: 'auth',
})

const { t } = useI18n()
useHead(() => ({ title: t('dashboard.title') }))

const auth = useAuthStore()
</script>
