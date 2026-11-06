<template>
  <div class="flex flex-col min-h-screen bg-slate-950 text-slate-100">
    <header class="border-b border-slate-800 bg-slate-900/80">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center">
        <NuxtLink to="/" class="flex items-center space-x-2 text-xl font-bold text-sky-400">
          <span>CloudScanr</span>
        </NuxtLink>
      </div>
    </header>

    <main class="flex-1 flex items-center justify-center px-4 py-16">
      <div class="max-w-md w-full text-center space-y-5">
        <div class="text-6xl font-extrabold text-white tabular-nums">{{ error.statusCode }}</div>
        <h1 class="text-xl font-bold text-white">{{ headline }}</h1>
        <p class="text-sm text-slate-400">{{ detail }}</p>
        <div class="flex justify-center gap-3 pt-2">
          <button
            class="px-5 py-2.5 rounded-lg bg-slate-800 hover:bg-slate-700 border border-slate-700 text-slate-200 text-sm font-semibold transition"
            @click="handleError"
          >
            Go Back
          </button>
          <NuxtLink
            to="/"
            class="px-5 py-2.5 rounded-lg bg-sky-600 hover:bg-sky-500 text-white text-sm font-semibold transition"
          >
            Home
          </NuxtLink>
        </div>
        <p v-if="traceId" class="text-[11px] text-slate-600 font-mono">trace {{ traceId }}</p>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { NuxtError } from '#app'

const props = defineProps<{ error: NuxtError }>()

const headline = computed(() => {
  if (props.error.statusCode === 404) {
    return 'This page drifted off the cloud map'
  }
  return 'Something broke on our side'
})

const detail = computed(() => {
  if (props.error.statusCode === 404) {
    return 'The address may be mistyped, or the resource moved. Rankings, search, and your console are one click away.'
  }
  return 'Our team has been notified with a trace ID. Try again in a moment, or head back home.'
})

const traceId = computed(() => {
  const data = props.error.data as Record<string, unknown> | undefined
  const id = data?.trace_id
  return typeof id === 'string' ? id : ''
})

function handleError() {
  clearError({ redirect: '/' })
}
</script>
