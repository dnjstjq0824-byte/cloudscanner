<template>
  <div class="flex flex-col min-h-screen bg-canvas text-ink">
    <header class="border-b border-subtle bg-surface/90">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center">
        <NuxtLink :to="$localePath('/')" class="text-xl font-bold text-brand">
          <span>TirpScanner</span>
        </NuxtLink>
      </div>
    </header>

    <main class="flex-1 flex items-center justify-center px-4 py-16">
      <div class="max-w-md w-full text-center space-y-5">
        <div class="text-6xl font-extrabold text-ink tabular-nums">{{ error.statusCode }}</div>
        <h1 class="text-xl font-bold text-ink">{{ headline }}</h1>
        <p class="text-sm text-ink-secondary">{{ detail }}</p>
        <div class="flex justify-center gap-3 pt-2">
          <button class="btn-ghost" @click="handleError">{{ $t('error.goBack') }}</button>
          <NuxtLink :to="$localePath('/')" class="btn-primary">{{ $t('error.home') }}</NuxtLink>
        </div>
        <p v-if="traceId" class="text-[11px] text-ink-secondary font-mono">{{ $t('error.trace', { id: traceId }) }}</p>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { NuxtError } from '#app'

const props = defineProps<{ error: NuxtError }>()

const { t } = useI18n()
const localePath = useLocalePath()

const headline = computed(() => {
  if (props.error.statusCode === 404) {
    return t('error.notFoundTitle')
  }
  return t('error.genericTitle')
})

const detail = computed(() => {
  if (props.error.statusCode === 404) {
    return t('error.notFoundBody')
  }
  return t('error.genericBody')
})

const traceId = computed(() => {
  const data = props.error.data as Record<string, unknown> | undefined
  const id = data?.trace_id
  return typeof id === 'string' ? id : ''
})

function handleError() {
  clearError({ redirect: localePath('/') })
}
</script>
