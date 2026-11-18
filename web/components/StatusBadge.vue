<template>
  <span class="inline-flex items-center px-2 py-0.5 rounded-full text-[10px] font-semibold border whitespace-nowrap" :class="tone">
    {{ status }}
  </span>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{ status: string }>()

const tone = computed(() => {
  const s = (props.status || '').toUpperCase()
  if (['RUNNING', 'ACTIVE', 'VERIFIED', 'ENFORCED', 'READY'].includes(s)) {
    return 'bg-emerald-500/10 text-emerald-700 dark:text-emerald-300 border-emerald-500/30'
  }
  if (['PROVISIONING', 'PENDING', 'PENDING_VERIFICATION', 'WARNING'].includes(s)) {
    return 'bg-amber-500/10 text-amber-700 dark:text-amber-300 border-amber-500/30'
  }
  if (['ERROR', 'FAILED', 'EXPIRED', 'REVOKED', 'SUSPENDED', 'CRITICAL', 'FROZEN_SECURITY'].includes(s)) {
    return 'bg-rose-500/10 text-rose-700 dark:text-rose-300 border-rose-500/30'
  }
  if (['OWNER', 'ADMIN', 'SUPERADMIN'].includes(s)) {
    return 'bg-sky-500/10 text-sky-700 dark:text-sky-300 border-sky-500/30'
  }
  return 'bg-slate-500/10 text-slate-600 dark:text-slate-300 border-slate-500/20'
})
</script>
