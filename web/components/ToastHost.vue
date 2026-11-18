<template>
  <div class="fixed bottom-4 right-4 z-[110] space-y-2 w-[calc(100%-2rem)] max-w-sm" aria-live="polite">
    <div
      v-for="t in toasts"
      :key="t.id"
      class="card px-4 py-3 text-xs font-medium shadow-pop flex items-start justify-between gap-3"
      :class="tone(t.kind)"
      role="status"
    >
      <span>{{ t.message }}</span>
      <button class="opacity-60 hover:opacity-100 shrink-0" :aria-label="'Dismiss notification'" @click="dismiss(t.id)">✕</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useToast } from '~/composables/useToast'

const { toasts, dismiss } = useToast()

function tone(kind: string): string {
  if (kind === 'success') return 'text-emerald-700 dark:text-emerald-300'
  if (kind === 'error') return 'text-rose-700 dark:text-rose-300'
  return 'text-ink'
}
</script>
