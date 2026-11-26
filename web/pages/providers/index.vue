<template>
  <div class="max-w-5xl mx-auto space-y-8 py-4">
    <div class="space-y-2">
      <h1 class="text-3xl font-extrabold text-ink tracking-tight">{{ $t('vendors.h1') }}</h1>
      <p class="text-sm text-ink-secondary">{{ $t('vendors.lede') }}</p>
    </div>

    <LoadingSkeleton v-if="loading" :rows="5" :label="$t('vendors.loading')" />
    <ErrorState v-else-if="error" :title="error.title" :message="error.message" @retry="load" />

    <div v-else class="grid gap-4 md:grid-cols-2">
      <NuxtLink
        v-for="p in providers"
        :key="p.slug"
        :to="$localePath(`/providers/${p.slug}`)"
        class="card p-5 hover:border-brand/60 transition block"
      >
        <div class="flex items-center justify-between gap-2">
          <span class="font-bold text-ink">{{ p.name }}</span>
          <StatusBadge :status="p.affiliate_status" />
        </div>
        <div class="text-xs text-ink-secondary mt-1">{{ $t('vendors.tier', { v: p.tier_priority }) }}</div>
      </NuxtLink>
    </div>
    <EmptyState
      v-if="!loading && !error && providers.length === 0"
      :title="$t('vendors.emptyTitle')"
      :message="$t('vendors.emptyBody')"
    />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ApiError, useApi } from '~/composables/useApi'
import type { ListEnvelope, Provider } from '~/types/api'

definePageMeta({
  layout: 'public',
})

const { t } = useI18n()
useHead(() => ({
  title: t('vendors.title'),
  meta: [{ name: 'description', content: t('vendors.meta') }],
}))

const providers = ref<Provider[]>([])
const loading = ref(false)
const error = ref<ApiError | null>(null)
const { request } = useApi()

async function load() {
  loading.value = true
  error.value = null
  try {
    const res = await request<ListEnvelope<Provider>>('/api/v1/providers')
    providers.value = res.data ?? []
  } catch (e) {
    error.value = e instanceof ApiError ? e : ApiError.fromUnknown(e)
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>
