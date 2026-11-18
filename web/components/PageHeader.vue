<template>
  <div>
    <h2 class="text-2xl font-extrabold tracking-tight text-ink">{{ title }}</h2>
    <p v-if="subtitle" class="text-sm text-ink-secondary mt-1">{{ subtitle }}</p>
    <div v-if="$slots.actions" class="mt-4 flex flex-wrap gap-2">
      <slot name="actions" />
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{ title: string; subtitle?: string }>()
</script>
