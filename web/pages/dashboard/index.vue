<template>
  <div class="space-y-8">
    <PageHeader title="Overview" :subtitle="`${auth.org?.name || 'Your organization'} · Role: ${auth.role || '—'}`" />

    <div class="card p-6 space-y-3">
      <h3 class="text-sm font-bold text-ink">Cloud management console retired</h3>
      <p class="text-xs text-ink-secondary leading-relaxed">
        Connections, server inventory, provisioning, and cloud budgets were
        retired with the TirpScanner pivot on 2026-11-26. The console now
        covers your account and organization only; deal, catalog, and
        affiliate tooling lives on the public site.
      </p>
      <NuxtLink to="/" class="inline-block text-xs font-semibold text-brand hover:underline">
        Back to TirpScanner home →
      </NuxtLink>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'dashboard',
  middleware: 'auth',
})

useHead({ title: 'Console Overview' })

const auth = useAuthStore()
</script>
