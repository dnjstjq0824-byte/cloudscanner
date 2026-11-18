<template>
  <div class="card p-6">
    <div class="text-3xl font-bold tabular-nums text-brand">{{ value }}</div>
    <div class="text-xs text-ink-secondary mt-1 uppercase tracking-wider">{{ label }}</div>
    <NuxtLink v-if="to" :to="to" class="inline-block mt-4 text-xs font-semibold text-brand hover:underline">
      {{ linkLabel }} &rarr;
    </NuxtLink>
  </div>
</template>

<script setup lang="ts">
defineProps<{ value: string | number; label: string; to?: string; linkLabel?: string }>()
</script>
